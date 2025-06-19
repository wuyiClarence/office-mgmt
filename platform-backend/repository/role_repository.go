package repository

import (
	"context"
	"platform-backend/dto"
	"platform-backend/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	*BaseRepository[models.Role]
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{BaseRepository: NewBaseRepository[models.Role](db)}
}

func (r *RoleRepository) GetUserIdsCreateRoleIds(ctx context.Context, userIDs []int64) ([]int64, error) {
	var allids []int64
	if err := r.db.WithContext(ctx).Model(models.Role{}).Where("create_user IN (?)", userIDs).Select("id").Find(&allids).Error; err != nil {
		return nil, err
	}
	return allids, nil
}
func (r *RoleRepository) GetRoleByRoleIds(ctx context.Context, roleIds []int64, page int, pageSize int) ([]*dto.RoleListItem, int64, error) {
	var results []*dto.RoleListItem
	var total int64

	db := r.db.WithContext(ctx).Model(models.Role{})
	db = db.Joins("left JOIN user ON role.create_user = user.id").
		Where("role.id IN (?)", roleIds)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Limit(pageSize).Offset(offset)
	}
	db = db.Order("role.created_at asc")

	err := db.Select("role.*,user.user_name as create_user_name,user.user_display_name as create_user_display_name").Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
func (r *RoleRepository) GetRoleList(ctx context.Context, isSuperAdmin bool, userID int64, page, pageSize int) (int64, []*dto.RoleListItem, error) {
	results := make([]*dto.RoleListItem, 0)
	var totalCount int64
	db := r.db.WithContext(ctx).Model(models.Role{})
	if isSuperAdmin {
		db = db.Joins("JOIN user ON user.id = role.create_user")
		if err := db.Count(&totalCount).Error; err != nil {
			return 0, nil, err
		}
	}

	offset := (page - 1) * pageSize
	if err := db.Select("role.*,user.user_name as create_user_name,user.user_display_name as create_user_display_name").
		Order("user.created_at desc").Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		return 0, nil, err
	}

	return totalCount, results, nil
}
