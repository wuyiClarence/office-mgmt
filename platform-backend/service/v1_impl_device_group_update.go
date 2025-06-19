package service

import (
	"context"
	"platform-backend/dto/enum"
	sp "platform-backend/service/permission"
	"platform-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/models"
	"platform-backend/repository"
)

func (svc *V1ServiceImpl) DeviceGroupUpdate(c *gin.Context, req *dto.DeviceGroupUpdateReq) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	err = sp.CheckPermission(c, enum.ResourceTypeDeviceGroup, req.ID, enum.PermissionKeyEdit, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	deviceMap := make(map[int64]struct{})
	for _, deviceID := range req.DeviceIDS {
		deviceMap[deviceID] = struct{}{}
	}

	return repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		deviceDeviceGroupRelRepo := repository.NewDeviceDeviceGroupRelRepository(tx)

		existedDevices, err := deviceDeviceGroupRelRepo.FindAll(ctx, map[string]interface{}{
			"group_id": req.ID,
		}, true)
		if err != nil {
			return err
		}

		needDelRecords := make([]int64, 0)
		for _, existedOne := range existedDevices {
			// 需要移除的
			if _, ok := deviceMap[existedOne.DeviceId]; !ok {
				needDelRecords = append(needDelRecords, existedOne.DeviceId)
			} else {
				delete(deviceMap, existedOne.DeviceId)
			}
		}

		needAddRecords := make([]models.DeviceDeviceGroupRel, 0)
		for deviceID, _ := range deviceMap {
			needAddRecords = append(needAddRecords, models.DeviceDeviceGroupRel{
				DeviceId: deviceID,
				GroupId:  req.ID,
			})
		}

		err = deviceDeviceGroupRelRepo.BatchInsert(ctx, needAddRecords)
		if err != nil {
			return err
		}

		return deviceDeviceGroupRelRepo.DeleteByCondition(ctx, map[string]interface{}{
			"device_id": needDelRecords,
			"group_id":  req.ID,
		})
	})
}
