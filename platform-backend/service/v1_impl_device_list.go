package service

import (
	"platform-backend/dto/enum"
	"platform-backend/models"
	sp "platform-backend/service/permission"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/repository"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) DeviceList(c *gin.Context, req *dto.ListReq) (*dto.DeviceListRes, error) {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, err
	}

	allIds, err := sp.GetUserWithResourceIds(c, enum.ResourceTypeDevice, enum.PermissionKeyView, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}
	deviceRep := repository.NewDeviceRepository(db.MysqlDB.DB())
	deviceGroupRep := repository.NewDeviceGroupRepository(db.MysqlDB.DB())

	deviceList, total, err := deviceRep.FindWithPagination(c.Request.Context(), map[string]interface{}{
		"id": allIds,
	}, req.PageIndex, req.PageSize, "id asc")
	if err != nil {
		return nil, err
	}

	res := &dto.DeviceListRes{
		List:  make([]*dto.DeviceDetail, 0, len(deviceList)),
		Total: total,
	}

	resourceIds := lo.Map(deviceList, func(p models.Device, _ int) int64 {
		return p.ID
	})
	permissions, err := sp.GetUserWithResourcePermission(c, enum.ResourceTypeDevice, resourceIds, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}

	for _, device := range deviceList {

		item := &dto.DeviceDetail{
			ID:         device.ID,
			DeviceName: device.DeviceName,
			AliasName:  device.AliasName,
			Mac:        device.Mac,
			Ip:         device.Ip,
			OsType:     device.OsType,
			DeviceType: device.DeviceType,
			Status:     device.Status,
		}
		for _, value := range permissions[item.ID] {
			item.Permissions = append(item.Permissions, value)
		}
		sort.Slice(item.Permissions, func(i, j int) bool {
			return item.Permissions[i].Order < item.Permissions[j].Order
		})
		deviceGroups, err := deviceGroupRep.FindDeviceDeviceGroup(c.Request.Context(), device.ID)
		if err != nil {
			return nil, err
		}
		item.DeviceGroupInfo = make([]dto.DeviceGroupInfo, 0)
		for _, deviceGroup := range deviceGroups {
			itemDeviceGroup := dto.DeviceGroupInfo{
				ID:              deviceGroup.ID,
				DeviceGroupName: deviceGroup.GroupName,
			}
			item.DeviceGroupInfo = append(item.DeviceGroupInfo, itemDeviceGroup)
		}
		res.List = append(res.List, item)
	}

	return res, nil
}
