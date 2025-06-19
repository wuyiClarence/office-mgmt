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

func (svc *V1ServiceImpl) PolicyDelete(c *gin.Context, req *dto.PolicyDelReq) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	err = sp.CheckPermission(c, enum.ResourceTypePolicy, req.IDs, enum.PermissionKeyDelete, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	return repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {

		// 删除与用户绑定关系
		err := repository.NewUserResourcePermissionRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"resource_id":   req.IDs,
			"resource_type": enum.ResourceTypePolicy,
		})
		if err != nil {
			return err
		}

		// 删除与角色绑定关系
		err = repository.NewRoleResourcePermissionRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"resource_id":   req.IDs,
			"resource_type": enum.ResourceTypePolicy,
		})
		if err != nil {
			return err
		}

		err = repository.NewPolicyRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"id": req.IDs,
		})
		if err != nil {
			return err
		}

		// 删除与设备绑定关系
		err = repository.NewPolicyDeviceRelRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"policy_id": req.IDs,
		})
		if err != nil {
			return err
		}

		// 删除与设备组绑定关系
		err = repository.NewPolicyDeviceGroupRelRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"policy_id": req.IDs,
		})
		if err != nil {
			return err
		}
		return nil
	})
}
