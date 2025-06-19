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

func (svc *V1ServiceImpl) PolicyList(c *gin.Context, req *dto.ListReq) (*dto.PolicyListRes, error) {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, err
	}

	allIds, err := sp.GetUserWithResourceIds(c, enum.ResourceTypePolicy, enum.PermissionKeyView, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}

	policies, total, err := repository.NewPolicyRepository(db.MysqlDB.DB()).FindWithPagination(c.Request.Context(), map[string]interface{}{
		"id": allIds,
	}, req.PageIndex, req.PageSize, "created_at desc")
	if err != nil {
		return nil, err
	}

	res := &dto.PolicyListRes{
		List:  make([]*dto.PolicyList, 0, len(policies)),
		Total: total,
	}

	policyDeviceRepo := repository.NewPolicyDeviceRelRepository(db.MysqlDB.DB())
	policyDeviceGroupRepo := repository.NewPolicyDeviceGroupRelRepository(db.MysqlDB.DB())
	deviceRepo := repository.NewDeviceRepository(db.MysqlDB.DB())
	deviceGroupRepo := repository.NewDeviceGroupRepository(db.MysqlDB.DB())

	resourceIds := lo.Map(policies, func(p models.Policy, _ int) int64 {
		return p.ID
	})
	permissions, err := sp.GetUserWithResourcePermission(c, enum.ResourceTypePolicy, resourceIds, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}

	for _, policy := range policies {

		policy := dto.PolicyList{
			ID: policy.ID,
			PolicyBaseInfo: dto.PolicyBaseInfo{
				PolicyName:    policy.PolicyName,
				ActionType:    policy.ActionType,
				ExecuteTime:   policy.ExecuteTime,
				ExecuteType:   policy.ExecuteType,
				AssociateType: policy.AssociateType,
				Status:        policy.Status,
				DayOfWeek:     policy.DayOfWeek,
				StartDate:     policy.StartDate,
				EndDate:       policy.EndDate,
			},
			Devices:      make([]*dto.DeviceDetail, 0),
			DeviceGroups: make([]*dto.DeviceGroupDetail, 0),
		}

		if policy.AssociateType == enum.AssociateTypeDevice {
			deviceIds, err := policyDeviceRepo.GetDeviceIDs(c.Request.Context(), policy.ID)
			if err != nil {
				return nil, err
			}

			deviceList, err := deviceRepo.GetDevicesByIds(c.Request.Context(), deviceIds)
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
				policy.Devices = append(policy.Devices, item)
			}
		}

		if policy.AssociateType == enum.AssociateTypeDeviceGroup {
			deviceGroupIds, err := policyDeviceGroupRepo.GetDeviceGroupIDs(c.Request.Context(), policy.ID)
			if err != nil {
				return nil, err
			}

			deviceGroupList, err := deviceGroupRepo.GetDevcieGroupsByIds(c.Request.Context(), deviceGroupIds)
			if err != nil {
				return nil, err
			}

			for _, deviceGroup := range deviceGroupList {
				item := &dto.DeviceGroupDetail{
					ID:              deviceGroup.ID,
					DeviceGroupName: deviceGroup.GroupName,
				}
				policy.DeviceGroups = append(policy.DeviceGroups, item)
			}
		}

		for _, value := range permissions[policy.ID] {
			policy.Permissions = append(policy.Permissions, value)
		}

		sort.Slice(policy.Permissions, func(i, j int) bool {
			return policy.Permissions[i].Order < policy.Permissions[j].Order
		})

		res.List = append(res.List, &policy)
	}

	return res, nil
}
