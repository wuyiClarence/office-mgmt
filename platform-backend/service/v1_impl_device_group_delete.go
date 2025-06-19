package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/repository"
	sp "platform-backend/service/permission"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) DeviceGroupDelete(c *gin.Context, req *dto.DeviceGroupDelReq) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	err = sp.CheckPermission(c, enum.ResourceTypeDeviceGroup, req.IDs, enum.PermissionKeyDelete, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	return repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		err := repository.NewUserResourcePermissionRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"resource_id":   req.IDs,
			"resource_type": enum.ResourceTypeDeviceGroup,
		})
		if err != nil {
			return err
		}

		// 删除与角色绑定关系
		err = repository.NewRoleResourcePermissionRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"resource_id":   req.IDs,
			"resource_type": enum.ResourceTypeDeviceGroup,
		})
		if err != nil {
			return err
		}

		// 删除device group
		deviceGroupRepo := repository.NewDeviceGroupRepository(tx)
		err = deviceGroupRepo.DeleteByCondition(ctx, map[string]interface{}{
			"id": req.IDs,
		})
		if err != nil {
			return err
		}

		// 删除策略
		devicePolicyRepo := repository.NewPolicyDeviceGroupRelRepository(tx)
		err = devicePolicyRepo.DeleteByCondition(ctx, map[string]interface{}{
			"device_group_id": req.IDs,
		})
		if err != nil {
			return err
		}

		// 删除与设备绑定关系
		deviceDeviceGroupRelRepo := repository.NewDeviceDeviceGroupRelRepository(tx)
		return deviceDeviceGroupRelRepo.DeleteByCondition(ctx, map[string]interface{}{
			"group_id": req.IDs,
		})
	})
}
