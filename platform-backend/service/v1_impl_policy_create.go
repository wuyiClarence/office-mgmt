package service

import (
	"context"
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

func (svc *V1ServiceImpl) PolicyCreate(c *gin.Context, req *dto.PolicyCreateReq) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	// check
	var (
		resourceType  = enum.ResourceTypeDevice
		permissionKey = enum.PermissionKeyPowerOff
	)
	if req.ActionType == enum.ActionTypeOff {
		permissionKey = enum.PermissionKeyPowerOff
	}
	if req.ActionType == enum.ActionTypeOn {
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

	policyRepo := repository.NewPolicyRepository(db.MysqlDB.DB())
	err = repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		newPolicy := models.Policy{
			PolicyName:    req.PolicyName,
			ActionType:    req.ActionType,
			AssociateType: req.AssociateType,
			Status:        int8(req.Status),
			ExecuteType:   req.ExecuteType,
			StartDate:     req.StartDate,
			EndDate:       req.EndDate,
			ExecuteTime:   req.ExecuteTime,
			DayOfWeek:     req.DayOfWeek,
			CreateUser:    userID,
		}
		err = policyRepo.Create(ctx, &newPolicy)
		if err != nil {
			return err
		}
		if req.AssociateType == enum.AssociateTypeDevice {
			policyDevices := make([]models.PolicyDevice, 0, len(req.DeviceIDS))
			for _, deviceId := range req.DeviceIDS {
				policyDevices = append(policyDevices, models.PolicyDevice{
					DeviceId: deviceId,
					PolicyId: newPolicy.ID,
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
					PolicyId:      newPolicy.ID,
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
