package repository

import (
	"context"
	"platform-backend/models"

	"gorm.io/gorm"
)

type PolicyDeviceGroupRelRepository struct {
	*BaseRepository[models.PolicyDeviceGroup]
}

func NewPolicyDeviceGroupRelRepository(db *gorm.DB) *PolicyDeviceGroupRelRepository {
	return &PolicyDeviceGroupRelRepository{BaseRepository: NewBaseRepository[models.PolicyDeviceGroup](db)}
}
func (d *PolicyDeviceGroupRelRepository) GetDeviceGroupIDs(ctx context.Context, policyId int64) ([]int64, error) {
	var allIDs []int64
	err := d.db.WithContext(ctx).Model(models.PolicyDeviceGroup{}).Where("policy_id = ?", policyId).Select("device_group_id").Find(&allIDs).Error

	return allIDs, err
}
