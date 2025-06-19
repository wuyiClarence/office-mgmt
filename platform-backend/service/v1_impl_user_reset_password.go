package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/repository"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) UserResetPassword(c *gin.Context, req *dto.UserResetPasswordReq) error {
	if !utils.IsSuperAdmin(c) {
		return errors.New("无权操作")
	}

	userRepo := repository.NewUserRepository(db.MysqlDB.DB())
	user, err := userRepo.FindOne(c.Request.Context(), map[string]interface{}{
		"id": req.UserID,
	})
	if err != nil {
		return errors.New("无效用户")
	}

	err = user.SetPassword(enum.UserDefaultPassword)
	if err != nil {
		return err
	}

	changeAt := time.Now()
	user.PasswordUpdatedAt = &changeAt

	err = userRepo.Update(c.Request.Context(), user)
	if err != nil {
		return errors.New("更新失败")
	}

	return nil
}
