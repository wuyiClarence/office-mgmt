package repository

import (
	"context"
	"platform-backend/models"

	"gorm.io/gorm"
)

type PolicyRepository struct {
	*BaseRepository[models.Policy]
}

func NewPolicyRepository(db *gorm.DB) *PolicyRepository {
	return &PolicyRepository{BaseRepository: NewBaseRepository[models.Policy](db)}
}

func (s *PolicyRepository) GetUserPolicyIDs(ctx context.Context, userID int64, isSuperAdmin bool) ([]int64, error) {
	var allIDs []int64

	if isSuperAdmin {
		err := s.db.WithContext(ctx).Model(models.Policy{}).Select("id").Find(&allIDs).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := s.db.WithContext(ctx).Model(models.Policy{}).Where("create_user = ?", userID).Select("id").Find(&allIDs).Error
		if err != nil {
			return nil, err
		}
	}

	return allIDs, nil
}
func (r *PolicyRepository) GetPolicys(ctx context.Context, page int, pageSize int) ([]*models.Policy, int64, error) {
	var policys []*models.Policy
	var total int64

	db := r.db.WithContext(ctx).Model(models.Policy{})

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Limit(pageSize).Offset(offset)
	}
	db = db.Order("created_at asc")

	err := db.Find(&policys).Error
	if err != nil {
		return nil, 0, err
	}

	return policys, total, nil
}

func (r *PolicyRepository) GetPolicysByPolicyIds(ctx context.Context, policyIds []int64, page int, pageSize int) ([]*models.Policy, int64, error) {
	var policys []*models.Policy
	var total int64

	db := r.db.WithContext(ctx).Model(models.Policy{}).Where("id IN (?)", policyIds)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Limit(pageSize).Offset(offset)
	}
	db = db.Order("created_at asc")

	err := db.Find(&policys).Error
	if err != nil {
		return nil, 0, err
	}

	return policys, total, nil
}
