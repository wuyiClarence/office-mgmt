package service

import (
	"errors"

	"github.com/gin-gonic/gin"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/repository"
	myjwt "platform-backend/utils/jwt"
)

func (svc *V1ServiceImpl) UserLogIn(c *gin.Context, req *dto.UserLoginApiReq) (*dto.UserLoginApiRes, error) {

	repo := repository.NewUserRepository(db.MysqlDB.DB())
	user, err := repo.FindOne(c.Request.Context(), map[string]interface{}{
		"user_name": req.UserName,
	})
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// req.Password, err = password.Decrypt(config.MyConfig.PasswordAuthKey, req.Password)
	// if err != nil {
	// 	return nil, err
	// }

	if !user.VerifyPassword(req.Password) {
		return nil, errors.New("密码错误")
	}

	accessToken, refreshToken, err := myjwt.GenerateTokens(user.ID, user.UserName)
	if err != nil {
		return nil, err
	}
	return &dto.UserLoginApiRes{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
