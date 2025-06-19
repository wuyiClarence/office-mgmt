package service

import (
	"errors"

	"github.com/gin-gonic/gin"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/repository"
	sp "platform-backend/service/permission"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) DeviceUpdate(c *gin.Context, req *dto.DeviceUpdateReq) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	err = sp.CheckPermission(c, enum.ResourceTypeDevice, req.ID, enum.PermissionKeyEdit, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	// check
	deviceRepo := repository.NewDeviceRepository(db.MysqlDB.DB())
	device, err := deviceRepo.FindOne(c.Request.Context(), map[string]interface{}{
		"id": req.ID,
	})
	if err != nil {
		return errors.New("无效的设备")
	}

	device.AliasName = req.AliasName
	if device.DeviceType == enum.DeviceTypeKvm {
		device.Mac = req.Mac
		device.Ip = req.Ip
	}

	err = deviceRepo.Update(c.Request.Context(), device)

	if err != nil {
		return err
	}

	return nil
}
