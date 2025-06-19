package service

import (
	"time"

	"github.com/gin-gonic/gin"

	db "platform-backend/db"
	"platform-backend/models"
	"platform-backend/repository"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) UserLogOut(c *gin.Context) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	repo := repository.NewUserRepository(db.MysqlDB.DB())

	changeAt := time.Now()

	err = repo.UpdateByCondition(c.Request.Context(), map[string]interface{}{
		"id": userID,
	}, &models.User{LogOutAt: &changeAt})

	return err
}
