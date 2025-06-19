package service

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/models"
	"platform-backend/repository"
	sp "platform-backend/service/permission"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) PolicyUpdate(c *gin.Context, req *dto.PolicyUpdateReq) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	err = sp.CheckPermission(c, enum.ResourceTypePolicy, req.ID, enum.PermissionKeyEdit, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	// check
	policyRepo := repository.NewPolicyRepository(db.MysqlDB.DB())
	policy, err := policyRepo.FindOne(c.Request.Context(), map[string]interface{}{
		"id": req.ID,
	})
	if err != nil {
		return errors.New("无效的策略")
	}

	var (
		resourceType  = enum.ResourceTypeDevice
		permissionKey = enum.PermissionKeyPowerOff
	)
	if req.ActionType == enum.ActionTypeOff {
		permissionKey = enum.PermissionKeyPowerOff
	}
	if req.ActionType == enum.ActionTypeOff {
		permissionKey = enum.PermissionKeyPowerOn
	}

	if req.AssociateType == enum.AssociateTypeDevice {
		resourceType = enum.ResourceTypeDevice
		for _, deviceId := range req.DeviceIDS {
			err = sp.CheckPermission(c, resourceType, deviceId, permissionKey, userID, utils.IsSuperAdmin(c))
			if err != nil {
				return err
			}
		}
	} else {
		resourceType = enum.ResourceTypeDeviceGroup
		for _, deviceGroupId := range req.DeviceGroupIDS {
			err = sp.CheckPermission(c, resourceType, deviceGroupId, permissionKey, userID, utils.IsSuperAdmin(c))
			if err != nil {
				return err
			}
		}
	}

	if req.ExecuteType == enum.TypeOnce {
		curtime := time.Now()
		req.StartDate = &curtime
		req.EndDate = &curtime
	}

	// do update
	if len(req.PolicyName) > 0 {
		policy.PolicyName = req.PolicyName
	}
	policy.ActionType = req.ActionType
	policy.AssociateType = req.AssociateType
	policy.Status = req.Status
	policy.ExecuteType = req.ExecuteType
	policy.StartDate = req.StartDate
	policy.EndDate = req.EndDate
	policy.ExecuteTime = req.ExecuteTime
	policy.DayOfWeek = req.DayOfWeek

	policy.Status = req.Status

	// todo 更新已经存在的定时任务以及未执行的任务  特别是状态变更

	err = repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		err = policyRepo.Update(c.Request.Context(), policy)

		if err != nil {
			return err
		}

		// 删除与设备绑定关系
		err = repository.NewPolicyDeviceRelRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"policy_id": policy.ID,
		})
		if err != nil {
			return err
		}

		// 删除与设备组绑定关系
		err = repository.NewPolicyDeviceGroupRelRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"policy_id": policy.ID,
		})
		if err != nil {
			return err
		}

		if req.AssociateType == enum.AssociateTypeDevice {
			policyDevices := make([]models.PolicyDevice, 0, len(req.DeviceIDS))
			for _, deviceId := range req.DeviceIDS {
				policyDevices = append(policyDevices, models.PolicyDevice{
					DeviceId: deviceId,
					PolicyId: policy.ID,
				})
			}

			policyDeviceRepo := repository.NewPolicyDeviceRelRepository(tx)
			if err = policyDeviceRepo.BatchInsert(ctx, policyDevices); err != nil {
				return err
			}
		}

		if req.AssociateType == enum.AssociateTypeDeviceGroup {
			policyDeviceGroups := make([]models.PolicyDeviceGroup, 0, len(req.DeviceIDS))
			for _, deviceGroupId := range req.DeviceGroupIDS {
				policyDeviceGroups = append(policyDeviceGroups, models.PolicyDeviceGroup{
					DeviceGroupId: deviceGroupId,
					PolicyId:      policy.ID,
				})
			}

			policyDeviceGroupRepo := repository.NewPolicyDeviceGroupRelRepository(tx)
			if err = policyDeviceGroupRepo.BatchInsert(ctx, policyDeviceGroups); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
