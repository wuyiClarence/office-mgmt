package service

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/models"
	"platform-backend/repository"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) RoleCreate(c *gin.Context, req *dto.RoleCreateReq) error {
	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return errors.New("非法用户")
	}

	err = repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		var err error
		// 1.role
		role := models.Role{
			RoleName:    req.RoleName,
			Description: req.Description,
			CreateUser:  curUserID,
		}

		if err = repository.NewRoleRepository(tx).Create(ctx, &role); err != nil {
			return err
		}
		permissionReq := dto.RoleUpdateMenuPermissionReq{
			RoleID: role.ID,
			Menu:   req.Menu,
		}
		err = svc.RoleMenuPermissionUpdate(c, &permissionReq)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return errors.New("创建失败")
	}

	return nil
}
