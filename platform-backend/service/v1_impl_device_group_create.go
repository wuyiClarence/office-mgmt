package service

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/models"
	"platform-backend/repository"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) DeviceGroupCreate(c *gin.Context, req *dto.DeviceGroupCreateReq) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	return repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		deviceGroupRepo := repository.NewDeviceGroupRepository(tx)

		group := &models.DeviceGroup{GroupName: req.DeviceGroupName, CreateUser: userID}
		err := deviceGroupRepo.Create(ctx, group)
		if err != nil {
			if repository.IsUniqueConstraintError(err) {
				return errors.New("已存在同名的设备组")
			}

			return err
		}

		deviceGroups := make([]models.DeviceDeviceGroupRel, 0, len(req.DeviceIDS))
		for _, deviceID := range req.DeviceIDS {
			deviceGroups = append(deviceGroups, models.DeviceDeviceGroupRel{
				DeviceId: deviceID,
				GroupId:  group.ID,
			})
		}

		deviceDeviceGroupRelRepo := repository.NewDeviceDeviceGroupRelRepository(tx)
		return deviceDeviceGroupRelRepo.BatchInsert(ctx, deviceGroups)
	})
}
