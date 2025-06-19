package repository

import (
	"context"
	"platform-backend/dto/enum"
	"platform-backend/models"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	*BaseRepository[models.Device]
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{BaseRepository: NewBaseRepository[models.Device](db)}
}

func (d *DeviceRepository) GetDevicesByGroupID(ctx context.Context, groupID int64) ([]*models.Device, error) {
	var devices []*models.Device
	err := d.db.WithContext(ctx).Model(models.Device{}).
		Select("device.*").
		Joins("JOIN device_device_group_rel ON device_device_group_rel.device_id = device.id").
		Where("device_device_group_rel.group_id = ?", groupID).Order("created_at desc").
		Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (d *DeviceRepository) GetDevicesByIds(ctx context.Context, deviceIds []int64) ([]*models.Device, error) {
	var devices []*models.Device

	db := d.db.WithContext(ctx).Model(models.Device{}).Where("id IN (?)", deviceIds)

	db = db.Order("created_at asc")

	err := db.Find(&devices).Error
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (d *DeviceRepository) GetAllDeviceIDs(ctx context.Context) ([]int64, error) {
	var allIDs []int64
	err := d.db.WithContext(ctx).Model(models.Device{}).Select("id").Find(&allIDs).Error

	return allIDs, err
}

func (r *DeviceRepository) GetDevicesByDeviceIds(ctx context.Context, deviceIds []int64, page int, pageSize int) ([]*models.Device, int64, error) {
	var devices []*models.Device
	var total int64

	db := r.db.WithContext(ctx).Model(models.Device{}).Where("id IN (?)", deviceIds)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Limit(pageSize).Offset(offset)
	}
	db = db.Order("created_at asc")

	err := db.Find(&devices).Error
	if err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

func (r *DeviceRepository) GetDevcies(ctx context.Context, page int, pageSize int) ([]*models.Device, int64, error) {
	var devices []*models.Device
	var total int64

	db := r.db.WithContext(ctx).Model(models.Device{})

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Limit(pageSize).Offset(offset)
	}
	db = db.Order("created_at asc")

	err := db.Find(&devices).Error
	if err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

func (d *DeviceRepository) GetPhysicalDevices(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device
	err := d.db.WithContext(ctx).Model(models.Device{}).
		Where("device_type = ? and status = ?", enum.DeviceTypePhysical, enum.DeviceStatusOnLine).Order("created_at desc").
		Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}
