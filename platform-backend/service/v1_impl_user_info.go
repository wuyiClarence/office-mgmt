package service

import (
	"github.com/gin-gonic/gin"
	db "platform-backend/db"
	"platform-backend/repository"

	"platform-backend/dto"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) UserInfo(c *gin.Context) (*dto.UserInfoRes, error) {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, err
	}

	user, err := repository.NewUserRepository(db.MysqlDB.DB()).FindOne(c.Request.Context(), map[string]interface{}{
		"id": userID,
	})
	if err != nil {
		return nil, err
	}

	res, err := svc.PermissionList(c)
	if err != nil {
		return nil, err
	}

	return &dto.UserInfoRes{Name: user.UserDisplayName, UserId: userID, Permissions: *res}, nil
}
