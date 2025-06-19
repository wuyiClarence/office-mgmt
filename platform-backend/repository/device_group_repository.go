package repository

import (
	"context"
	"platform-backend/models"

	"gorm.io/gorm"
)

type DeviceGroupRepository struct {
	*BaseRepository[models.DeviceGroup]
}

func NewDeviceGroupRepository(db *gorm.DB) *DeviceGroupRepository {
	return &DeviceGroupRepository{BaseRepository: NewBaseRepository[models.DeviceGroup](db)}
}

// FindDeviceDeviceGroup 查找设备归属到了哪些服务器组
func (s *DeviceGroupRepository) FindDeviceDeviceGroup(ctx context.Context, deviceID int64) ([]*models.DeviceGroup, error) {
	var results []*models.DeviceGroup

	db := s.db.WithContext(ctx).Model(models.Device{})
	err := db.Joins("JOIN device_device_group_rel ON device.id = device_device_group_rel.device_id").
		Joins("JOIN device_group ON device_device_group_rel.group_id = device_group.id").
		Where("device.id = ? ", deviceID).Select("device_group.*").Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *DeviceGroupRepository) GetUserDeviceGroupIDs(ctx context.Context, userID int64, isSuperAdmin bool) ([]int64, error) {
	var allIDs []int64

	if isSuperAdmin {
		err := s.db.WithContext(ctx).Model(models.DeviceGroup{}).Select("id").Find(&allIDs).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := s.db.WithContext(ctx).Model(models.DeviceGroup{}).Where("create_user = ?", userID).Select("id").Find(&allIDs).Error
		if err != nil {
			return nil, err
		}
	}

	return allIDs, nil
}

func (r *DeviceGroupRepository) GetDevcieGroups(ctx context.Context, page int, pageSize int) ([]*models.DeviceGroup, int64, error) {
	var deviceGroups []*models.DeviceGroup
	var total int64

	db := r.db.WithContext(ctx).Model(models.DeviceGroup{})

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Limit(pageSize).Offset(offset)
	}
	db = db.Order("created_at asc")

	err := db.Find(&deviceGroups).Error
	if err != nil {
		return nil, 0, err
	}

	return deviceGroups, total, nil
}
func (r *DeviceGroupRepository) GetDevcieGroupsByDeviceGroupIds(ctx context.Context, deviceGroupIds []int64, page int, pageSize int) ([]*models.DeviceGroup, int64, error) {
	var deviceGroups []*models.DeviceGroup
	var total int64

	db := r.db.WithContext(ctx).Model(models.DeviceGroup{}).Where("id IN (?)", deviceGroupIds)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Limit(pageSize).Offset(offset)
	}
	db = db.Order("created_at asc")

	err := db.Find(&deviceGroups).Error
	if err != nil {
		return nil, 0, err
	}

	return deviceGroups, total, nil
}

func (r *DeviceGroupRepository) GetDevcieGroupsByIds(ctx context.Context, deviceGroupIds []int64) ([]*models.DeviceGroup, error) {
	var deviceGroups []*models.DeviceGroup

	db := r.db.WithContext(ctx).Model(models.DeviceGroup{}).Where("id IN (?)", deviceGroupIds)

	db = db.Order("created_at asc")

	err := db.Find(&deviceGroups).Error
	if err != nil {
		return nil, err
	}

	return deviceGroups, nil
}
