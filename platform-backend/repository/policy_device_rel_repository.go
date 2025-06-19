package repository

import (
	"context"
	"platform-backend/models"

	"gorm.io/gorm"
)

type PolicyDeviceRelRepository struct {
	*BaseRepository[models.PolicyDevice]
}

func NewPolicyDeviceRelRepository(db *gorm.DB) *PolicyDeviceRelRepository {
	return &PolicyDeviceRelRepository{BaseRepository: NewBaseRepository[models.PolicyDevice](db)}
}

func (d *PolicyDeviceRelRepository) GetDeviceIDs(ctx context.Context, policyId int64) ([]int64, error) {
	var allIDs []int64
	err := d.db.WithContext(ctx).Model(models.PolicyDevice{}).Where("policy_id = ?", policyId).Select("device_id").Find(&allIDs).Error

	return allIDs, err
}
