package service

import (
	"context"

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

func (svc *V1ServiceImpl) UserDelete(c *gin.Context, req *dto.UserDeleteReq) error {
	// 检查是否有权限删除
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	err = sp.CheckPermission(c, enum.ResourceTypeUser, req.UserIDs, enum.PermissionKeyDelete, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	return repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		err := repository.NewUserRepository(tx).UpdateByCondition(ctx, map[string]interface{}{
			"id": req.UserIDs,
		}, &models.User{Status: enum.UserStatusDeleted})
		if err != nil {
			return err
		}

		err = repository.NewUserRoleRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"user_id": req.UserIDs,
		})
		if err != nil {
			return err
		}

		err = repository.NewUserResourcePermissionRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"user_id": req.UserIDs,
		})
		if err != nil {
			return err
		}

		return nil
	})
}
