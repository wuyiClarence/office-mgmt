package repository

import (
	"context"
	"platform-backend/models"
	"strings"

	"gorm.io/gorm"
)

type RoleMenuPermissionRepository struct {
	*BaseRepository[models.RoleMenuPermission]
}

func NewRoleMenuPermissionRepository(db *gorm.DB) *RoleMenuPermissionRepository {
	return &RoleMenuPermissionRepository{BaseRepository: NewBaseRepository[models.RoleMenuPermission](db)}
}

func (r *RoleMenuPermissionRepository) FindUserRolePermission(ctx context.Context, userID int64, apiPath string, apiMethod string) (bool, error) {
	var (
		roleids           []int64
		permissionKeys    []string
		ownPermissionKeys []string
	)
	apiMethod = strings.ToUpper(apiMethod)

	err := r.db.WithContext(ctx).Model(&models.MenuPermissionAPI{}).Select("permission_key").Where("api_path = ? and api_method = ?", apiPath, apiMethod).Find(&permissionKeys).Error
	if err != nil {
		return false, err
	}
	if len(permissionKeys) == 0 {
		// 不需要验证的API接口
		return true, nil
	}
	err = r.db.WithContext(ctx).Model(&models.UserRole{}).Select("role_id").Where("user_id = ?", userID).Find(&roleids).Error
	if err != nil {
		return false, err
	}
	err = r.db.WithContext(ctx).Model(&models.RoleMenuPermission{}).Select("permission_key").Where("role_id IN (?)", roleids).Where("permission_key IN (?)", permissionKeys).Find(&ownPermissionKeys).Error
	if err != nil {
		return false, err
	}

	if len(ownPermissionKeys) == 0 {
		// 验证失败，没有权限
		return false, nil
	}

	return true, nil
}

func (r *RoleMenuPermissionRepository) FindUserAllRolePermission(ctx context.Context, userID int64, isSuperAdmin bool) ([]models.RoleMenuPermission, error) {
	var (
		resources []models.RoleMenuPermission
	)

	subQuery := r.db.WithContext(ctx).Model(&models.UserRole{}).Select("role_id").Where("user_id = ?", userID)
	if err := r.db.WithContext(ctx).Where("role_id IN (?)", subQuery).Find(&resources).Error; err != nil {
		return nil, err
	}

	return resources, nil
}
