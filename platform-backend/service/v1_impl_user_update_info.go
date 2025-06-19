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
	"platform-backend/service/password"
	sp "platform-backend/service/permission"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) UserUpdateInfo(c *gin.Context, req *dto.UserUpdateInfoReq) error {
	if req.UserID <= 0 {
		return errors.New("无效用户")
	}
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	var userIds []int64
	userIds = append(userIds, req.UserID)

	err = sp.CheckPermission(c, enum.ResourceTypeUser, userIds, enum.PermissionKeyEdit, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	err = repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(db.MysqlDB.DB())
		user, err := userRepo.FindOne(c.Request.Context(), map[string]interface{}{
			"id": req.UserID,
		})
		if err != nil {
			return errors.New("无效用户")
		}

		if len(req.Password) > 0 {
			// 至少包含一个字母、一个数字，且长度大于等于 6 位 且不包含空格
			if err := password.CheckPassWord(req.Password); err != nil {
				return err
			}
			err = user.SetPassword(req.Password)
			if err != nil {
				return err
			}
		}

		if len(req.UserDisplayName) > 0 {
			user.UserDisplayName = req.UserDisplayName
		}

		if len(req.Email) > 0 {
			user.Email = req.Email
		}

		if len(req.PhoneNumber) > 0 {
			user.PhoneNumber = req.PhoneNumber
		}

		err = userRepo.Update(c.Request.Context(), user)
		if err != nil {
			return errors.New("更新失败")
		}

		userRoleRepo := repository.NewUserRoleRepository(tx)

		err = userRoleRepo.DeleteByCondition(ctx, map[string]interface{}{
			"user_id": req.UserID,
		})
		if err != nil {
			return err
		}
		roles := make([]models.UserRole, 0, len(req.RoleIDS))
		for _, roleID := range req.RoleIDS {
			roles = append(roles, models.UserRole{
				UserID: req.UserID,
				RoleID: roleID,
			})
		}

		if err = userRoleRepo.BatchInsert(ctx, roles); err != nil {
			return err
		}
		return nil
	})

	return err
}
