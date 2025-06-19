package repository

import (
	"context"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/models"
)

type RoleResourcePermissionRepository struct {
	*BaseRepository[models.RoleResourcePermission]
}

func NewRoleResourcePermissionRepository(db *gorm.DB) *RoleResourcePermissionRepository {
	return &RoleResourcePermissionRepository{BaseRepository: NewBaseRepository[models.RoleResourcePermission](db)}
}

func (r *RoleResourcePermissionRepository) GetPermissionKeys(ctx context.Context, roleIds []int64, resourceType enum.ResourceType, resourceId int64) ([]string, error) {
	var permissionKeys []string
	if err := r.db.Model(&models.RoleResourcePermission{}).
		Distinct("permission_key").
		Where("role_id IN ?", roleIds).
		Where("resource_type", resourceType).Where("resource_id", resourceId).
		Pluck("permission_key", &permissionKeys).Error; err != nil {
		return nil, err
	}
	return permissionKeys, nil
}
func (r *RoleResourcePermissionRepository) GetRolePermissionKeys(ctx context.Context, resourceType enum.ResourceType, resourceId int64, includeRoleIDs []int64, excludeRoleIDs []int64, page int, pageSize int) (int64, []dto.RolePermission, error) {
	var rolePermissionKeys []dto.RolePermissionKey
	var totalCount int64
	var selectRoleId []int64
	var roles []*models.Role

	db := r.db.WithContext(ctx).Model(models.Role{}).
		Where("id in ?", includeRoleIDs)

		// 动态添加排除条件
	if len(excludeRoleIDs) > 0 {
		db = db.Where("id not in ?", excludeRoleIDs)
	}

	db = db.Order("id desc")
	if err := db.Count(&totalCount).Error; err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return 0, nil, err
	}
	selectRoleId = lo.Map(roles, func(p *models.Role, _ int) int64 {
		return p.ID
	})

	if err := r.db.WithContext(ctx).Model(models.Role{}).
		Joins("JOIN role_resource_permission ON role.id = role_resource_permission.role_id").
		Where("resource_type", resourceType).
		Where("resource_id", resourceId).
		Where("role_id in ?", selectRoleId).
		Select("role.id as role_id,role.role_name,role_resource_permission.permission_key,role_resource_permission.permission_name").
		Find(&rolePermissionKeys).Error; err != nil {
		return 0, nil, err
	}
	grouped := make(map[int64][]dto.RolePermissionKey)
	for _, upk := range rolePermissionKeys {
		grouped[upk.RoleID] = append(grouped[upk.RoleID], upk)
	}
	var rolePermissions []dto.RolePermission

	for _, role := range roles {
		item := dto.RolePermission{}
		item.RoleId = role.ID
		item.RoleName = role.RoleName
		item.Keys = []enum.PermissionKey{}
		value, exists := grouped[role.ID]
		if exists {
			for _, p := range value {
				item.Keys = append(item.Keys, p.PermissionKey)
			}
		}
		rolePermissions = append(rolePermissions, item)
	}
	return totalCount, rolePermissions, nil
}

// func (r *RoleResourcePermissionRepository) GetRoleWithResourcePermissions(ctx context.Context, resourceType enum.ResourceType, permissionKey enum.PermissionKey, roleId int64) ([]*models.RoleResourcePermission, error) {
// 	results := make([]*models.RoleResourcePermission, 0)
// 	err := r.db.WithContext(ctx).Model(&models.RoleResourcePermission{}).Where("role_id", roleId).
// 		Where("resource_type", resourceType).
// 		Where("permission_key", permissionKey).Find(&results).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return results, nil
// }

func (r *RoleResourcePermissionRepository) GetRoleIdsWithResourcePermissions(ctx context.Context, resourceType enum.ResourceType, permissionKey enum.PermissionKey, roleIds []int64) ([]*models.RoleResourcePermission, error) {
	results := make([]*models.RoleResourcePermission, 0)
	err := r.db.WithContext(ctx).Model(&models.RoleResourcePermission{}).
		Where("role_id IN ?", roleIds).
		Where("resource_type", resourceType).
		Where("permission_key", permissionKey).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *RoleResourcePermissionRepository) GetRoleIdsWithResourceIdsPermissions(ctx context.Context, resourceType enum.ResourceType, resourceIDs []int64, roleIds []int64) ([]*models.RoleResourcePermission, error) {
	results := make([]*models.RoleResourcePermission, 0)
	err := r.db.WithContext(ctx).Model(&models.RoleResourcePermission{}).
		Where("role_id IN ?", roleIds).
		Where("resource_type", resourceType).
		Where("resource_id in ?", resourceIDs).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *RoleResourcePermissionRepository) GetRolesOwnPermissionSetResourceIds(ctx context.Context, resourceType enum.ResourceType, roleIds []int64) ([]int64, error) {
	results := make([]int64, 0)
	err := r.db.WithContext(ctx).Model(&models.RoleResourcePermission{}).
		Where("role_id IN ?", roleIds).
		Where("permission_key", enum.PermissionKeyPermissionMgmt).
		Where("resource_type", resourceType).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *RoleResourcePermissionRepository) GetRolesOwnResourcePermissionKeys(ctx context.Context, roleIds []int64, resourceType enum.ResourceType, resourceId int64) ([]enum.PermissionKey, error) {
	var permissionKeys []enum.PermissionKey
	if err := r.db.Model(&models.RoleResourcePermission{}).
		Distinct("permission_key").
		Where("role_id IN ?", roleIds).
		Where("resource_type", resourceType).Where("resource_id", resourceId).
		Pluck("permission_key", &permissionKeys).Error; err != nil {
		return nil, err
	}
	return permissionKeys, nil
}
func (r *RoleResourcePermissionRepository) GetRoleOwnResourcePermissionKeys(ctx context.Context, roleId int64, resourceType enum.ResourceType, resourceId int64) ([]enum.PermissionKey, error) {
	var permissionKeys []enum.PermissionKey
	if err := r.db.Model(&models.RoleResourcePermission{}).
		Distinct("permission_key").
		Where("role_id", roleId).
		Where("resource_type", resourceType).Where("resource_id", resourceId).
		Pluck("permission_key", &permissionKeys).Error; err != nil {
		return nil, err
	}
	return permissionKeys, nil
}

func (r *RoleResourcePermissionRepository) RoleOwnResourcePermissionKey(ctx context.Context, roleIds []int64, resourceType enum.ResourceType, resourceId int64, key enum.PermissionKey) (bool, error) {
	var count int64
	err := r.db.Model(&models.RoleResourcePermission{}).
		Where("role_id IN ?", roleIds).
		Where("resource_type = ?", resourceType).
		Where("resource_id = ?", resourceId).
		Where("permission_key = ?", key).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
