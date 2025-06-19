package service

import (
	"github.com/gin-gonic/gin"
	"platform-backend/utils/jwt"

	"platform-backend/dto"
)

func (svc *V1ServiceImpl) UserRefreshToken(c *gin.Context, req *dto.UserRefreshTokenReq) (*dto.UserRefreshTokenRes, error) {
	accessToken, err := jwt.RefreshTokens(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.UserRefreshTokenRes{AccessToken: accessToken}, nil
}
