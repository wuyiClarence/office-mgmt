package repository

import (
	"context"
	"gorm.io/gorm"
	"platform-backend/models"
)

type RoleUserRepository struct {
	*BaseRepository[models.RoleUser]
}

func NewRoleUserRepository(db *gorm.DB) *RoleUserRepository {
	return &RoleUserRepository{BaseRepository: NewBaseRepository[models.RoleUser](db)}
}

func (r *RoleUserRepository) GetRoleUsersByRoleID(ctx context.Context, roleID int64, page, pageSize int) ([]models.User, int64, error) {
	var (
		users      []models.User
		totalCount int64
	)

	subQuery := r.db.WithContext(ctx).Model(&models.RoleUser{}).Select("user_id").Where("role_id = ?", roleID)

	// 计算总记录数
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id IN (?)", subQuery).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("id IN (?)", subQuery).Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}
