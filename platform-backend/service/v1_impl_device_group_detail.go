package service

import (
	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/repository"
	sp "platform-backend/service/permission"
	"platform-backend/utils"

	"github.com/gin-gonic/gin"
)

func (svc *V1ServiceImpl) DeviceGroupDetail(c *gin.Context, req *dto.DeviceGroupDetailReq) (*dto.DeviceGroupDetail, error) {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, err
	}

	err = sp.CheckPermission(c, enum.ResourceTypeDeviceGroup, req.ID, enum.PermissionKeyView, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}

	deviceGroup, err := repository.NewDeviceGroupRepository(db.MysqlDB.DB()).FindOne(c.Request.Context(), map[string]interface{}{
		"id": req.ID,
	})
	if err != nil {
		return nil, err
	}

	res := &dto.DeviceGroupDetail{
		DeviceGroupName: deviceGroup.GroupName,
		Devices:         nil,
	}

	deviceList, err := repository.NewDeviceRepository(db.MysqlDB.DB()).GetDevicesByGroupID(c.Request.Context(), deviceGroup.ID)
	if err != nil {
		return nil, err
	}

	for _, device := range deviceList {
		item := &dto.DeviceDetail{
			ID:         device.ID,
			AliasName:  device.AliasName,
			DeviceName: device.DeviceName,
			Mac:        device.Mac,
			Ip:         device.Ip,
			OsType:     device.OsType,
			DeviceType: device.DeviceType,
			Status:     device.Status,
		}
		res.Devices = append(res.Devices, item)
	}
	return res, nil
}
