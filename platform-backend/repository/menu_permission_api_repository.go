package repository

import (
	"context"
	"platform-backend/models"

	"gorm.io/gorm"
)

type MenuPermissionAPIApiRepository struct {
	*BaseRepository[models.MenuPermissionAPI]
}

func NewMenuPermissionApiRepository(db *gorm.DB) *MenuPermissionAPIApiRepository {
	return &MenuPermissionAPIApiRepository{BaseRepository: NewBaseRepository[models.MenuPermissionAPI](db)}
}

func (r *MenuPermissionAPIApiRepository) UpdateAllUsedStatus(ctx context.Context) error {

	err := r.db.Exec("UPDATE menu_permission_api SET used = ?", "0").Error
	return err
}

func (r *MenuPermissionAPIApiRepository) DeleteUnUsedPermissionApi(ctx context.Context) error {

	err := r.db.WithContext(ctx).Model(&models.MenuPermissionAPI{}).Where("used = ?", 0).Delete(&models.MenuPermissionAPI{}).Error

	return err
}
