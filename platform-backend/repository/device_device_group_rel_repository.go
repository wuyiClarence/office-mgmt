package repository

import (
	"context"
	"platform-backend/models"

	"gorm.io/gorm"
)

type DeviceDeviceGroupRelRepository struct {
	*BaseRepository[models.DeviceDeviceGroupRel]
}

func NewDeviceDeviceGroupRelRepository(db *gorm.DB) *DeviceDeviceGroupRelRepository {
	return &DeviceDeviceGroupRelRepository{BaseRepository: NewBaseRepository[models.DeviceDeviceGroupRel](db)}
}
func (d *DeviceDeviceGroupRelRepository) GetDeviceIdsByDeviceGroupIds(ctx context.Context, deviceGroupIds []int64) ([]int64, error) {

	var allIDs []int64
	err := d.db.WithContext(ctx).Model(models.DeviceDeviceGroupRel{}).Where("group_id IN (?)", deviceGroupIds).
		Select("device_id").Find(&allIDs).Error

	if err != nil {
		return nil, err
	}
	return allIDs, nil
}
