package service

import (
	"github.com/gin-gonic/gin"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/repository"
	"platform-backend/service/deviceoprator"
	sp "platform-backend/service/permission"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) DevicePowerOn(c *gin.Context, req *dto.DevicePowerOnReq) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	err = sp.CheckPermission(c, enum.ResourceTypeDevice, req.IDs, enum.PermissionKeyPowerOn, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	deviceRep := repository.NewDeviceRepository(db.MysqlDB.DB())

	devices, err := deviceRep.GetDevicesByIds(c.Request.Context(), req.IDs)
	if err != nil {
		return err
	}

	for _, device := range devices {
		deviceoprator.PowerOn(c.Request.Context(), device)
	}

	return nil
}
