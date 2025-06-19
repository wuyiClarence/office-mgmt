package service

import (
	"context"
	"errors"

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

func (svc *V1ServiceImpl) RoleUpdate(c *gin.Context, req *dto.RoleUpdateReq) error {
	var err error
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	var roleIds []int64
	roleIds = append(roleIds, req.RoleID)
	err = sp.CheckPermission(c, enum.ResourceTypeRole, roleIds, enum.PermissionKeyEdit, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	err = repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		// 1.role
		if len(req.Description) > 0 || len(req.RoleName) > 0 {
			if err := repository.NewRoleRepository(tx).UpdateByCondition(ctx, map[string]interface{}{
				"id": req.RoleID,
			}, &models.Role{Description: req.Description, RoleName: req.RoleName}); err != nil {
				return err
			}
		}

		permissionReq := dto.RoleUpdateMenuPermissionReq{
			RoleID: req.RoleID,
			Menu:   req.Menu,
		}
		err = svc.RoleMenuPermissionUpdate(c, &permissionReq)
		if err != nil {
			return err
		}
		return nil

	})

	if err != nil {
		return errors.New("更新失败")
	}

	return nil
}
