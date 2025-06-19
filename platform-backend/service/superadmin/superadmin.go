package superadmin

import (
	"context"
	"errors"

	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto/enum"
	"platform-backend/models"
	"platform-backend/repository"
)

func InitSuperAdmin() error {
	userRepo := repository.NewUserRepository(db.MysqlDB.DB())

	ctx := context.Background()
	superAdmin, err := userRepo.FindOne(ctx, map[string]interface{}{
		"user_name": enum.SuperAdminUserAccName,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return nil
	}

	superAdmin = &models.User{
		UserName:        enum.SuperAdminUserAccName,
		UserDisplayName: "超级管理员",
		Status:          enum.UserStatusOk,
	}

	err = superAdmin.SetPassword(enum.SuperAdminDefaultPassword)
	if err != nil {
		return err
	}

	return userRepo.Create(ctx, superAdmin)
}
