package repository

import (
	"context"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/models"
)

type UserResourcePermissionRepository struct {
	*BaseRepository[models.UserResourcePermission]
}

func NewUserResourcePermissionRepository(db *gorm.DB) *UserResourcePermissionRepository {
	return &UserResourcePermissionRepository{BaseRepository: NewBaseRepository[models.UserResourcePermission](db)}
}

func (r *UserResourcePermissionRepository) GetPermissionKeys(ctx context.Context, userIDs []int64, resourceType enum.ResourceType, resourceId int64) ([]enum.PermissionKey, error) {
	var permissionKeys []enum.PermissionKey
	if err := r.db.Model(&models.UserResourcePermission{}).
		Distinct("permission_key").
		Where("user_id IN ?", userIDs).
		Where("resource_type", resourceType).Where("resource_id", resourceId).
		Pluck("permission_key", &permissionKeys).Error; err != nil {
		return nil, err
	}
	return permissionKeys, nil
}

func (r *UserResourcePermissionRepository) GetUserPermissionKeys(ctx context.Context, resourceType enum.ResourceType, resourceId int64, includeUserIDs []int64, excludeUserIDs []int64, page int, pageSize int) (int64, []dto.UserPermission, error) {
	var userPermissionKeys []dto.UserPermissionKey
	var totalCount int64
	var selectUserId []int64
	var users []*models.User

	db := r.db.WithContext(ctx).Model(models.User{}).
		Where("status = ?", 0).
		Where("id in ?", includeUserIDs)

	// 动态添加排除条件
	if len(excludeUserIDs) > 0 {
		db = db.Where("id not in ?", excludeUserIDs)
	}

	db = db.Order("id desc")

	if err := db.Count(&totalCount).Error; err != nil {
		return 0, nil, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Limit(pageSize).Offset(offset)
	}
	db = db.Order("user.created_at desc")
	err := db.Find(&users).Error
	if err != nil {
		return 0, nil, err
	}
	selectUserId = lo.Map(users, func(p *models.User, _ int) int64 {
		return p.ID
	})

	if err := r.db.WithContext(ctx).Model(&models.User{}).
		Joins("JOIN user_resource_permission ON user.id = user_resource_permission.user_id").
		Where("user.status = ?", 0).
		Where("user.id in ?", selectUserId).
		Where("resource_type", resourceType).
		Where("resource_id", resourceId).
		Select("user.id as user_id,user.user_display_name,user_resource_permission.permission_key,user_resource_permission.permission_name").
		Find(&userPermissionKeys).Error; err != nil {
		return 0, nil, err
	}
	grouped := make(map[int64][]dto.UserPermissionKey)
	for _, upk := range userPermissionKeys {
		grouped[upk.UserID] = append(grouped[upk.UserID], upk)
	}
	var userPermissions []dto.UserPermission
	for _, user := range users {
		item := dto.UserPermission{}
		item.UserID = user.ID
		item.UserDisplayName = user.UserDisplayName
		item.UserName = user.UserName
		item.Keys = []enum.PermissionKey{}
		value, exists := grouped[user.ID]
		if exists {
			for _, p := range value {
				item.Keys = append(item.Keys, p.PermissionKey)
			}
		}
		userPermissions = append(userPermissions, item)
	}
	return totalCount, userPermissions, nil
}

func (r *UserResourcePermissionRepository) GetUserWithResourcePermissions(ctx context.Context, resourceType enum.ResourceType, permissionKey enum.PermissionKey, userId int64) ([]*models.UserResourcePermission, error) {
	results := make([]*models.UserResourcePermission, 0)
	var userIDs []int64
	userIDs = append(userIDs, userId, 0)

	err := r.db.WithContext(ctx).Model(&models.UserResourcePermission{}).Where("user_id in ?", userIDs).
		Where("resource_type", resourceType).
		Where("permission_key", permissionKey).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *UserResourcePermissionRepository) GetUserWithResourceIdsPermissions(ctx context.Context, resourceType enum.ResourceType, resourceIDs []int64, userId int64) ([]*models.UserResourcePermission, error) {
	results := make([]*models.UserResourcePermission, 0)
	var userIDs []int64
	userIDs = append(userIDs, userId, 0)
	err := r.db.WithContext(ctx).Model(&models.UserResourcePermission{}).Where("user_id in ?", userIDs).
		Where("resource_type", resourceType).
		Where("resource_id in ?", resourceIDs).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *UserResourcePermissionRepository) GetUserOwnPermissionSetResourceIds(ctx context.Context, resourceType enum.ResourceType, userId int64) ([]int64, error) {
	results := make([]int64, 0)

	var userIDs []int64
	userIDs = append(userIDs, userId, 0)
	err := r.db.WithContext(ctx).Model(&models.UserResourcePermission{}).Where("user_id in ?", userIDs).
		Where("permission_key", enum.PermissionKeyPermissionMgmt).
		Where("resource_type", resourceType).
		Distinct("resource_id").
		Pluck("resource_id", &results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *UserResourcePermissionRepository) GetUserOwnResourcePermissionKeys(ctx context.Context, userId int64, resourceType enum.ResourceType, resourceId int64, containAllUser bool) ([]enum.PermissionKey, error) {
	var permissionKeys []enum.PermissionKey

	var userIDs []int64
	if containAllUser {
		userIDs = append(userIDs, userId, 0)
	} else {
		userIDs = append(userIDs, userId)
	}

	if err := r.db.Model(&models.UserResourcePermission{}).
		Distinct("permission_key").
		Where("user_id in ?", userIDs).
		Where("resource_type", resourceType).Where("resource_id", resourceId).
		Pluck("permission_key", &permissionKeys).Error; err != nil {
		return nil, err
	}
	return permissionKeys, nil
}

func (r *UserResourcePermissionRepository) UserOwnResourcePermissionKey(ctx context.Context, userId int64, resourceType enum.ResourceType, resourceId int64, key enum.PermissionKey) (bool, error) {
	var count int64
	var userIDs []int64
	userIDs = append(userIDs, userId, 0)

	err := r.db.Model(&models.UserResourcePermission{}).
		Where("user_id in ?", userIDs).
		Where("resource_type = ?", resourceType).
		Where("resource_id = ?", resourceId).
		Where("permission_key = ?", key).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
