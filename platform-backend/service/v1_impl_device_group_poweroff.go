package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	mymqtt "platform-backend/mqtt"
	"platform-backend/repository"
	"platform-backend/service/deviceoprator"
	sp "platform-backend/service/permission"
	"platform-backend/utils"
	"platform-backend/utils/log"
)

func (svc *V1ServiceImpl) DeviceGroupPowerOff(c *gin.Context, req *dto.DeviceGroupPowerOffReq) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	err = sp.CheckPermission(c, enum.ResourceTypeDeviceGroup, req.IDs, enum.PermissionKeyPowerOff, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	mqttClient, err := mymqtt.GetMqttClient()
	if err != nil {
		return err
	}

	deviceGroupRelRep := repository.NewDeviceDeviceGroupRelRepository(db.MysqlDB.DB())

	deviceIds, err := deviceGroupRelRep.GetDeviceIdsByDeviceGroupIds(c.Request.Context(), req.IDs)
	if err != nil {
		return err
	}

	deviceRep := repository.NewDeviceRepository(db.MysqlDB.DB())
	devices, err := deviceRep.GetDevicesByIds(c.Request.Context(), deviceIds)
	if err != nil {
		return err
	}

	for _, device := range devices {
		if device.DeviceType == enum.DeviceTypePhysical {
			deviceoprator.PowerOff(c.Request.Context(), device)
		}
		if device.DeviceType == enum.DeviceTypeKvm {
			err = mqttClient.SendVirHostPowerOff(device)
		}
		if err != nil {
			_, _ = fmt.Fprintf(log.SystemLogger, "%s send MQTT to %s,err %v\n", time.Now().Format("2006-01-02 15:04:05"), device.Mac, err)
		} else {
			_, _ = fmt.Fprintf(log.SystemLogger, "%s send MQTT to %s\n", time.Now().Format("2006-01-02 15:04:05"), device.Mac)
		}
	}

	return nil
}
