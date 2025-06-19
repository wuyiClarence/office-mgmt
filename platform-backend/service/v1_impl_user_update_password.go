package service

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"platform-backend/config"
	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/repository"
	"platform-backend/service/password"
)

func (svc *V1ServiceImpl) UserUpdatePassword(c *gin.Context, req *dto.UserUpdatePasswordApiReq) error {
	var err error
	req.OldPassword, err = password.Decrypt(config.MyConfig.PasswordAuthKey, req.OldPassword)
	if err != nil {
		return errors.New("无法解析密码")
	}

	if req.NewPassword != req.ConfirmPassword {
		return errors.New("两次输入的新密码不一致")
	}

	// 至少包含一个字母、一个数字，且长度大于等于 6 位 且不包含空格
	if err := password.CheckPassWord(req.NewPassword); err != nil {
		return err
	}

	if req.OldPassword == req.NewPassword {
		return errors.New("新旧密码不能相同")
	}

	userIDAny, _ := c.Get("userID")
	var userID uint
	if v, ok := userIDAny.(uint); ok {
		userID = v
	}

	if userID == 0 {
		return errors.New("invalid user")
	}

	return repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		repo := repository.NewUserRepository(tx)

		user, err := repo.FindOne(c.Request.Context(), map[string]interface{}{"user_id": userID})
		if err != nil {
			return errors.New("invalid user")
		}

		if !user.VerifyPassword(req.OldPassword) {
			return errors.New("原密码错误")
		}

		changeAt := time.Now()
		user.PasswordUpdatedAt = &changeAt
		err = user.SetPassword(req.NewPassword)
		if err != nil {
			return err
		}

		return repo.Update(ctx, user)
	})
}
