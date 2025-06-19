package service

import (
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/models"
	"platform-backend/repository"
	sp "platform-backend/service/permission"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) DeviceGroupList(c *gin.Context, req *dto.ListReq) (*dto.DeviceGroupListRes, error) {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, err
	}

	allIds, err := sp.GetUserWithResourceIds(c, enum.ResourceTypeDeviceGroup, enum.PermissionKeyView, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}

	deviceGroups, total, err := repository.NewDeviceGroupRepository(db.MysqlDB.DB()).FindWithPagination(c.Request.Context(), map[string]interface{}{
		"id": allIds,
	}, req.PageIndex, req.PageSize, "created_at desc")
	if err != nil {
		return nil, err
	}

	res := &dto.DeviceGroupListRes{
		List:  make([]*dto.DeviceGroupDetail, 0, len(deviceGroups)),
		Total: total,
	}

	resourceIds := lo.Map(deviceGroups, func(p models.DeviceGroup, _ int) int64 {
		return p.ID
	})

	permissions, err := sp.GetUserWithResourcePermission(c, enum.ResourceTypeDeviceGroup, resourceIds, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(deviceGroups); i++ {
		detail := dto.DeviceGroupDetail{
			ID:              deviceGroups[i].ID,
			DeviceGroupName: deviceGroups[i].GroupName,
			Devices:         make([]*dto.DeviceDetail, 0),
		}

		deviceRepo := repository.NewDeviceRepository(db.MysqlDB.DB())
		deviceList, err := deviceRepo.GetDevicesByGroupID(c.Request.Context(), deviceGroups[i].ID)
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
			detail.Devices = append(detail.Devices, item)
		}

		for _, value := range permissions[detail.ID] {
			detail.Permissions = append(detail.Permissions, value)
		}
		sort.Slice(detail.Permissions, func(i, j int) bool {
			return detail.Permissions[i].Order < detail.Permissions[j].Order
		})

		res.List = append(res.List, &detail)
	}

	return res, nil
}
