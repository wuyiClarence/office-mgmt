package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/models"
	"platform-backend/repository"
	"platform-backend/service/password"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) UserCreate(c *gin.Context, req *dto.UserCreateReq) error {
	var err error
	if len(req.Password) > 0 {
		err = password.CheckPassWord(req.Password)
		if err != nil {
			return err
		}
	}
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	return repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)

		newUser := models.User{
			UserName:        req.UserName,
			UserDisplayName: req.UserDisplayName,
			Status:          enum.UserStatusOk,
			Email:           req.Email,
			PhoneNumber:     req.PhoneNumber,
			Password: sql.NullString{
				String: req.Password,
				Valid:  true,
			},
			CreateUser: userID,
		}

		if len(req.Password) == 0 {
			req.Password = enum.UserDefaultPassword
		}
		err = newUser.SetPassword(req.Password)
		if err != nil {
			return err
		}

		err = userRepo.Create(ctx, &newUser)
		if err != nil {
			return errors.New(fmt.Sprintf("已存在用户%s", req.UserName))
		}

		roles := make([]models.UserRole, 0, len(req.RoleIDS))
		for _, roleID := range req.RoleIDS {
			roles = append(roles, models.UserRole{
				UserID: newUser.ID,
				RoleID: roleID,
			})
		}

		userRoleRepo := repository.NewUserRoleRepository(tx)
		if err = userRoleRepo.BatchInsert(ctx, roles); err != nil {
			return err
		}

		return nil
	})
}
