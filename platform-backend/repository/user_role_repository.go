package repository

import (
	"context"

	"gorm.io/gorm"

	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/models"
)

type UserRoleRepository struct {
	*BaseRepository[models.UserRole]
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{BaseRepository: NewBaseRepository[models.UserRole](db)}
}

func (r *UserRoleRepository) GetWithRoleUserList(ctx context.Context, roleID int64, page, pageSize int) (int64, []*dto.UserWithRoleTrans, error) {
	results := make([]*dto.UserWithRoleTrans, 0)
	var totalCount int64

	db := r.db.WithContext(ctx).Model(models.UserRole{})
	db = db.Joins("JOIN user ON user.id = user_role.user_id").
		Where("user_role.role_id = ? and user.status = ?", roleID, enum.UserStatusOk)

	if err := db.Count(&totalCount).Error; err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize
	if err := db.Select("user.*").
		Order("user.created_at desc").Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		return 0, nil, err
	}

	return totalCount, results, nil
}

func (r *UserRoleRepository) GetUsersRoleIds(ctx context.Context, userIDs []int64) ([]int64, error) {
	var roleIDs []int64
	if err := r.db.Model(&models.UserRole{}).
		Distinct("role_id").
		Where("user_id IN ?", userIDs).
		Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, err
	}
	return roleIDs, nil
}

func (r *UserRoleRepository) GetUserRoleInfos(ctx context.Context, userID int64) ([]*dto.UserRoleInfo, error) {
	results := make([]*dto.UserRoleInfo, 0)

	db := r.db.WithContext(ctx).Model(models.UserRole{})
	db = db.Joins("JOIN role ON role.id = user_role.role_id").
		Where("user_role.user_id = ?", userID)
	if err := db.Select("user_role.role_id,role.role_name as role_name").
		Order("user_role.role_id asc").Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *UserRoleRepository) GetUserRoleIds(ctx context.Context, userID int64) ([]int64, error) {
	var roleIDs []int64
	if err := r.db.Model(&models.UserRole{}).
		Distinct("role_id").
		Where("user_id", userID).
		Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, err
	}
	return roleIDs, nil
}
