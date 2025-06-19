package repository

import (
	"context"
	"platform-backend/models"

	"gorm.io/gorm"
)

type MenuPermissionRepository struct {
	*BaseRepository[models.MenuPermission]
}

func NewMenuPermissionRepository(db *gorm.DB) *MenuPermissionRepository {
	return &MenuPermissionRepository{BaseRepository: NewBaseRepository[models.MenuPermission](db)}
}

func (r *MenuPermissionRepository) UpdateAllUsedStatus(ctx context.Context) error {
	err := r.db.Exec("UPDATE menu_permission SET used = ?", "0").Error

	return err
}

func (r *MenuPermissionRepository) FindUserMenuPermission(ctx context.Context, userID int64, isSuperAdmin bool) ([]*models.MenuPermission, error) {
	var permissions []*models.MenuPermission
	var err error
	if isSuperAdmin {
		err = r.db.WithContext(ctx).
			Table("menu_permission").
			Select("menu_permission.*").
			Order("menu_permission.created_at DESC").
			Find(&permissions).Error
	} else {
		err = r.db.WithContext(ctx).
			Table("menu_permission").
			Select("menu_permission.*").
			Joins("JOIN role_menu_permission ON menu_permission.permission_key = role_menu_permission.permission_key").
			Joins("JOIN user_role ON role_menu_permission.role_id = user_role.role_id").
			Where("user_role.user_id = ?", userID).
			Order("role_menu_permission.created_at DESC").
			Find(&permissions).Error
	}

	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *MenuPermissionRepository) FindRoleMenuPermission(ctx context.Context, roleId int64) ([]*models.MenuPermission, error) {
	var permissions []*models.MenuPermission

	err := r.db.WithContext(ctx).
		Table("menu_permission").
		Select("menu_permission.*").
		Joins("JOIN role_menu_permission ON menu_permission.permission_key = role_menu_permission.permission_key").
		Where("role_menu_permission.role_id = ?", roleId).
		Order("menu_permission.created_at DESC").
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *MenuPermissionRepository) DeleteUnUsedPermission(ctx context.Context) error {
	subQuery := r.db.WithContext(ctx).Model(&models.MenuPermission{}).Select("permission_key").Where("used = ?", 0)
	err := r.db.WithContext(ctx).Model(&models.RoleMenuPermission{}).Where("permission_key IN (?)", subQuery).Delete(&models.RoleMenuPermission{}).Error
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Model(&models.MenuPermission{}).Where("used = ?", 0).Delete(&models.MenuPermission{}).Error
	if err != nil {
		return err
	}
	return err
}
