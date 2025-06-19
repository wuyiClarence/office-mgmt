package permission

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/samber/lo"
	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/dto/permission"
	"platform-backend/models"
	"platform-backend/repository"
)

// GetUserWithResourceIds 获取用户能访问的其中一类资源的Ids用于列表展示
func GetUserWithResourceIds(ctx context.Context, resourceType enum.ResourceType, permissionKey enum.PermissionKey, userID int64, isSuperAdmin bool) ([]int64, error) {
	var allIds []int64
	roleRep := repository.NewRoleRepository(db.MysqlDB.DB())
	userRoleRep := repository.NewUserRoleRepository(db.MysqlDB.DB())
	userRep := repository.NewUserRepository(db.MysqlDB.DB())
	userResourceRep := repository.NewUserResourcePermissionRepository(db.MysqlDB.DB())
	roleResourceRep := repository.NewRoleResourcePermissionRepository(db.MysqlDB.DB())

	if resourceType == enum.ResourceTypeRole {
		// 获取自己及其下级用户创建的RoleId
		subUserIds, err := userRep.GetChildUserIds(ctx, userID, isSuperAdmin)
		if err != nil {
			return nil, err
		}
		subUserIds = append(subUserIds, userID)

		subRoleIds, err := roleRep.GetUserIdsCreateRoleIds(ctx, subUserIds)
		if err != nil {
			return nil, err
		}
		allIds = append(allIds, subRoleIds...)
	}
	if resourceType == enum.ResourceTypeUser {
		// 获取自己及其下级用户创建的UserId
		subUserIds, err := userRep.GetChildUserIds(ctx, userID, isSuperAdmin)
		if err != nil {
			return nil, err
		}
		allIds = append(allIds, subUserIds...)
	}
	if resourceType == enum.ResourceTypeDeviceGroup {
		// 获取自己创建的device group
		deviceGroupIDs, err := repository.NewDeviceGroupRepository(db.MysqlDB.DB()).GetUserDeviceGroupIDs(ctx, userID, isSuperAdmin)
		if err != nil {
			return nil, err
		}

		allIds = append(allIds, deviceGroupIDs...)
	}
	if resourceType == enum.ResourceTypeDevice && isSuperAdmin {
		deviceIDs, err := repository.NewDeviceRepository(db.MysqlDB.DB()).GetAllDeviceIDs(ctx)
		if err != nil {
			return nil, err
		}

		allIds = append(allIds, deviceIDs...)
	}
	if resourceType == enum.ResourceTypePolicy {
		// 获取自己创建的device policy
		policyIDs, err := repository.NewPolicyRepository(db.MysqlDB.DB()).GetUserPolicyIDs(ctx, userID, isSuperAdmin)
		if err != nil {
			return nil, err
		}

		allIds = append(allIds, policyIDs...)
	}

	//获取所有人都可以操作的资源

	// 获取其他人授予我管理的资源
	userResources, err := userResourceRep.GetUserWithResourcePermissions(ctx, resourceType, permissionKey, userID)
	if err != nil {
		return nil, err
	}
	userWithRoleIds := lo.Map(userResources, func(p *models.UserResourcePermission, _ int) int64 {
		return p.ResourceID
	})
	allIds = lo.Union(allIds, userWithRoleIds)

	if !isSuperAdmin {

		userRoleIds, err := userRoleRep.GetUserRoleIds(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("用户角色获取失败")
		}

		//获取自己所属角色被授予的管理权限
		roleResources, err := roleResourceRep.GetRoleIdsWithResourcePermissions(ctx, resourceType, permissionKey, userRoleIds)
		if err != nil {
			return nil, err
		}
		roleWithRoleIds := lo.Map(roleResources, func(p *models.RoleResourcePermission, _ int) int64 {
			return p.ResourceID
		})
		allIds = lo.Union(allIds, roleWithRoleIds)
	} else {

	}

	return allIds, nil
}

// GetUserWithResourcePermission 获取用户对这些资源拥有的权限
func GetUserWithResourcePermission(ctx context.Context, resourceType enum.ResourceType, resourceIDs []int64, userID int64, isSuperAdmin bool) (map[int64]map[enum.PermissionKey]dto.Permission, error) {
	var withAllIds []int64

	allKeys, ok := permission.ResourceTypeKeyMap[resourceType]
	if !ok {
		return nil, errors.New("请求的资源类型错误")
	}

	result := make(map[int64]map[enum.PermissionKey]dto.Permission)
	for _, resourceID := range resourceIDs {
		result[resourceID] = make(map[enum.PermissionKey]dto.Permission)
	}

	roleRep := repository.NewRoleRepository(db.MysqlDB.DB())
	userRoleRep := repository.NewUserRoleRepository(db.MysqlDB.DB())
	userRep := repository.NewUserRepository(db.MysqlDB.DB())
	userResourceRep := repository.NewUserResourcePermissionRepository(db.MysqlDB.DB())
	roleResourceRep := repository.NewRoleResourcePermissionRepository(db.MysqlDB.DB())

	if resourceType == enum.ResourceTypeRole {
		// 获取自己及其下级用户创建的RoleId
		subUserIds, err := userRep.GetChildUserIds(ctx, userID, isSuperAdmin)
		if err != nil {
			return nil, err
		}
		subUserIds = append(subUserIds, userID)

		subRoleIds, err := roleRep.GetUserIdsCreateRoleIds(ctx, subUserIds)
		if err != nil {
			return nil, err
		}
		withAllIds = lo.Intersect(resourceIDs, subRoleIds)
	}
	if resourceType == enum.ResourceTypeUser {
		// 获取自己及其下级用户创建的UserId
		subUserIds, err := userRep.GetChildUserIds(ctx, userID, isSuperAdmin)
		if err != nil {
			return nil, err
		}
		withAllIds = lo.Intersect(resourceIDs, subUserIds)
	}

	if isSuperAdmin {
		withAllIds = resourceIDs
	}

	if resourceType == enum.ResourceTypeDeviceGroup {
		// 获取自己创建的device group
		deviceGroupIDs, err := repository.NewDeviceGroupRepository(db.MysqlDB.DB()).GetUserDeviceGroupIDs(ctx, userID, isSuperAdmin)
		if err != nil {
			return nil, err
		}
		withAllIds = lo.Intersect(resourceIDs, deviceGroupIDs)
	}

	if resourceType == enum.ResourceTypePolicy {
		// 获取自己创建的device policy
		policyIDs, err := repository.NewPolicyRepository(db.MysqlDB.DB()).GetUserPolicyIDs(ctx, userID, isSuperAdmin)
		if err != nil {
			return nil, err
		}
		withAllIds = lo.Intersect(resourceIDs, policyIDs)
	}

	for _, resourceID := range withAllIds {
		for _, key := range allKeys {
			result[resourceID][key.Key] = dto.Permission{Key: key.Key, Name: key.Name, Order: key.Order}
		}
	}

	// 获取其他人授予我管理的资源
	userResources, err := userResourceRep.GetUserWithResourceIdsPermissions(ctx, resourceType, resourceIDs, userID)
	if err != nil {
		return nil, err
	}
	for _, userResource := range userResources {
		order, err := GetResourcePermissionOrder(resourceType, userResource.PermissionKey)
		if err != nil {
			continue
		}
		result[userResource.ResourceID][userResource.PermissionKey] = dto.Permission{
			Key:   userResource.PermissionKey,
			Name:  userResource.PermissionName,
			Order: order,
		}
	}

	if !isSuperAdmin {
		userRoleIds, err := userRoleRep.GetUserRoleIds(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("用户角色获取失败")
		}

		//获取自己所属角色被授予的管理权限
		roleResources, err := roleResourceRep.GetRoleIdsWithResourceIdsPermissions(ctx, resourceType, resourceIDs, userRoleIds)
		if err != nil {
			return nil, err
		}
		for _, roleResource := range roleResources {
			order, err := GetResourcePermissionOrder(resourceType, roleResource.PermissionKey)
			if err != nil {
				continue
			}
			result[roleResource.ResourceID][roleResource.PermissionKey] = dto.Permission{
				Key:   roleResource.PermissionKey,
				Name:  roleResource.PermissionName,
				Order: order}
		}
	}

	return result, nil
}

// CheckPermission 检查用户是否对一类资源拥有操作权限，如批量删除用户等
func CheckPermission(ctx context.Context, resourceType enum.ResourceType, resourceId interface{}, permissionKey enum.PermissionKey, userID int64, isSuperAdmin bool) error {
	var resourceIds []int64
	switch v := resourceId.(type) {
	case []int64:
		resourceIds = v
	case int64:
		resourceIds = append(resourceIds, v)
	default:
		return fmt.Errorf("参数传递错误: %v", v)
	}

	allIds, err := GetUserWithResourceIds(ctx, resourceType, permissionKey, userID, isSuperAdmin)
	if err != nil {
		return err
	}
	unPermission := lo.Without(resourceIds, allIds...)
	if len(unPermission) > 0 {
		return fmt.Errorf("没有权限进行操作: %v", unPermission)
	}
	return nil
}

// GetAllUserPermission 获取一个类型的一个资源授予给所有人可操作的权限
func GetAllUserPermission(ctx context.Context, resourceType enum.ResourceType, resourceId int64, userID int64, isSuperAdmin bool) ([]enum.PermissionKey, bool, error) {
	userResourcePermissionRep := repository.NewUserResourcePermissionRepository(db.MysqlDB.DB())
	var userIds []int64
	canSetAllUser := false

	if isSuperAdmin {
		canSetAllUser = true
	}
	createUser, err := GetResourceCreateUser(ctx, resourceType, resourceId)
	if err != nil {
		return nil, false, err
	}
	if userID == createUser {
		canSetAllUser = true
	}

	userIds = append(userIds, 0)
	userPermission, err := userResourcePermissionRep.GetPermissionKeys(ctx, userIds, resourceType, resourceId)
	if err != nil {
		return nil, false, err
	}
	return userPermission, canSetAllUser, nil
}

// GetResourceUserPermissionInfo 获取用户对一个资源的操作权限,返回所有用户拥有该资源的情况
func GetResourceUserPermissionInfo(ctx context.Context, resourceType enum.ResourceType, resourceId int64, page int, pageSize int, userID int64, isSuperAdmin bool) (*dto.UserResourcePermissionRes, error) {
	res := &dto.UserResourcePermissionRes{}
	var includeUserIds []int64
	var excludeUserIds []int64

	containSubUer := true
	rep := repository.NewUserRepository(db.MysqlDB.DB())
	if resourceType == enum.ResourceTypeRole {
		containSubUer = true
		roleRep := repository.NewRoleRepository(db.MysqlDB.DB())
		role, err := roleRep.FindOne(ctx, map[string]interface{}{
			"id": resourceId,
		})
		if err != nil {
			return nil, err
		}
		excludeUserIds = append(excludeUserIds, role.CreateUser)
	}
	if resourceType == enum.ResourceTypeUser {
		containSubUer = true
		userRep := repository.NewUserRepository(db.MysqlDB.DB())
		user, err := userRep.FindOne(ctx, map[string]interface{}{
			"id": resourceId,
		})
		if err != nil {
			return nil, err
		}
		excludeUserIds = append(excludeUserIds, user.CreateUser)
		excludeUserIds = append(excludeUserIds, userID)
	}
	if resourceType == enum.ResourceTypeDeviceGroup {
		// containSubUer = true
		deviceGroupRep := repository.NewDeviceGroupRepository(db.MysqlDB.DB())
		deviceGroup, err := deviceGroupRep.FindOne(ctx, map[string]interface{}{
			"id": resourceId,
		})
		if err != nil {
			return nil, err
		}
		excludeUserIds = append(excludeUserIds, deviceGroup.CreateUser)
	}

	if resourceType == enum.ResourceTypePolicy {
		policyRep := repository.NewPolicyRepository(db.MysqlDB.DB())
		policy, err := policyRep.FindOne(ctx, map[string]interface{}{
			"id": resourceId,
		})
		if err != nil {
			return nil, err
		}
		excludeUserIds = append(excludeUserIds, policy.CreateUser)
	}

	if containSubUer {
		subUserIDs, err := rep.GetChildUserIds(ctx, userID, isSuperAdmin)
		if err != nil {
			return nil, err
		}
		includeUserIds = append(includeUserIds, subUserIDs...)
	}

	userResourcePermissionRep := repository.NewUserResourcePermissionRepository(db.MysqlDB.DB())
	total, userPermission, err := userResourcePermissionRep.GetUserPermissionKeys(ctx, resourceType, resourceId, includeUserIds, excludeUserIds, page, pageSize)
	if err != nil {
		return nil, err
	}
	res.UserPermission = userPermission
	res.Total = total

	return res, nil
}

// GetResourceRolePermissionInfo 获取角色对一个资源的操作权限,返回所有角色拥有该资源的情况
func GetResourceRolePermissionInfo(ctx context.Context, resourceType enum.ResourceType, resourceId int64, page int, pageSize int, userID int64, isSuperAdmin bool) (*dto.RoleResourcePermissionRes, error) {
	var includeRoleIds []int64
	var excludeRoleIds []int64

	containSubRole := true
	res := &dto.RoleResourcePermissionRes{}

	if resourceType == enum.ResourceTypeRole {
		containSubRole = true
		excludeRoleIds = append(excludeRoleIds, resourceId)
	}

	if resourceType == enum.ResourceTypeUser {
		containSubRole = true
		userRoleRep := repository.NewUserRoleRepository(db.MysqlDB.DB())
		userRole, err := userRoleRep.FindOne(ctx, map[string]interface{}{
			"user_id": resourceId,
		})
		if err != nil {
			return nil, err
		}

		excludeRoleIds = append(excludeRoleIds, userRole.RoleID)
	}

	if containSubRole {
		// 获取我创建的所有下级用户
		rep := repository.NewUserRepository(db.MysqlDB.DB())
		subUserIds, err := rep.GetChildUserIds(ctx, userID, isSuperAdmin)
		if err != nil {
			return nil, err
		}
		// 获取自己及其下级用户创建的RoleId
		subUserIds = append(subUserIds, userID)
		roleRep := repository.NewRoleRepository(db.MysqlDB.DB())
		roleIds, err := roleRep.GetUserIdsCreateRoleIds(ctx, subUserIds)
		if err != nil {
			return nil, err
		}
		includeRoleIds = append(includeRoleIds, roleIds...)
	}

	roleResourcePermissionRep := repository.NewRoleResourcePermissionRepository(db.MysqlDB.DB())
	total, rolePermission, err := roleResourcePermissionRep.GetRolePermissionKeys(ctx, resourceType, resourceId, includeRoleIds, excludeRoleIds, page, pageSize)
	if err != nil {
		return nil, err
	}
	res.RolePermission = rolePermission
	res.Total = total

	return res, nil
}

func CheckResourceTypeAndKey(resourceType enum.ResourceType, key enum.PermissionKey) error {
	_, ok := permission.ResourceTypeKeyMap[resourceType]
	if !ok {
		return fmt.Errorf("请求的资源类型错误%s", resourceType)
	}

	for _, permissionItem := range permission.ResourceTypeKeyMap[resourceType] {
		if permissionItem.Key == key {
			return nil
		}
	}

	return fmt.Errorf("请求的资源操作类型错误%s", key)
}

func GetResourcePermissionName(resourceType enum.ResourceType, key enum.PermissionKey) (string, error) {
	_, ok := permission.ResourceTypeKeyMap[resourceType]
	if !ok {
		return "", fmt.Errorf("请求的资源类型错误%s", resourceType)
	}

	for _, permissionItem := range permission.ResourceTypeKeyMap[resourceType] {
		if permissionItem.Key == key {
			return permissionItem.Name, nil
		}
	}

	return "", fmt.Errorf("请求的资源操作类型错误%s", key)
}

func GetResourcePermissionOrder(resourceType enum.ResourceType, key enum.PermissionKey) (int, error) {
	_, ok := permission.ResourceTypeKeyMap[resourceType]
	if !ok {
		return 0, fmt.Errorf("请求的资源类型错误%s", resourceType)
	}

	for _, permissionItem := range permission.ResourceTypeKeyMap[resourceType] {
		if permissionItem.Key == key {
			return permissionItem.Order, nil
		}
	}

	return 0, fmt.Errorf("请求的资源操作类型错误%s", key)
}

// PutUserPermission 分配资源用户, 对此资源有权限管理的用户，可以进行操作 allowUser为0 代表所有用户
func PutUserPermission(ctx context.Context, resourceType enum.ResourceType, resourceId int64, permissionKey enum.PermissionKey, enable bool, allowUserId int64, userID int64, isSuperAdmin bool) error {
	err := CheckResourceTypeAndKey(resourceType, permissionKey)
	if err != nil {
		return err
	}

	err = CheckPermission(ctx, resourceType, resourceId, enum.PermissionKeyPermissionMgmt, userID, isSuperAdmin)
	if err != nil {
		return err
	}
	if allowUserId != 0 {
		userRep := repository.NewUserRepository(db.MysqlDB.DB())
		_, err = userRep.FindOne(ctx, map[string]interface{}{
			"id": allowUserId,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("授权的用户不存在")
			}
			return err
		}
	} else {
		// 所有用户的访问权限设置，只能是superadmin或者资源的创建者
		if !isSuperAdmin {
			createUser, err := GetResourceCreateUser(ctx, resourceType, resourceId)
			if err != nil {
				return err
			}
			if userID != createUser {
				return fmt.Errorf("没有操作权限")
			}
		}
	}

	permissionName, err := GetResourcePermissionName(resourceType, permissionKey)
	if err != nil {
		return err
	}
	userResourceRep := repository.NewUserResourcePermissionRepository(db.MysqlDB.DB())
	_, err = userResourceRep.FindOne(ctx, map[string]interface{}{
		"user_id":        allowUserId,
		"permission_key": permissionKey,
		"resource_type":  resourceType,
		"resource_id":    resourceId,
	})

	if enable {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userResource := models.UserResourcePermission{
				UserID:         allowUserId,
				PermissionKey:  permissionKey,
				PermissionName: permissionName,
				ResourceType:   resourceType,
				ResourceID:     resourceId,
				CreateUser:     userID,
			}
			err = userResourceRep.Create(ctx, &userResource)
			if err != nil {
				return err
			}
		}
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			//delete
			err = userResourceRep.DeleteByCondition(ctx, map[string]interface{}{
				"user_id":        allowUserId,
				"permission_key": permissionKey,
				"resource_type":  resourceType,
				"resource_id":    resourceId})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// PutRolePermission 分配资源角色, 对此资源有权限管理的用户，可以进行操作
func PutRolePermission(ctx context.Context, resourceType enum.ResourceType, resourceId int64, permissionKey enum.PermissionKey, enable bool, allowRoleId int64, userID int64, isSuperAdmin bool) error {
	err := CheckResourceTypeAndKey(resourceType, permissionKey)
	if err != nil {
		return err
	}

	err = CheckPermission(ctx, resourceType, resourceId, enum.PermissionKeyPermissionMgmt, userID, isSuperAdmin)
	if err != nil {
		return err
	}

	if allowRoleId != 0 {
		roleRep := repository.NewRoleRepository(db.MysqlDB.DB())
		_, err = roleRep.FindOne(ctx, map[string]interface{}{
			"id": allowRoleId,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("授权的角色不存在")
			}
			return err
		}
	}

	permissionName, err := GetResourcePermissionName(resourceType, permissionKey)
	if err != nil {
		return err
	}
	roleResourceRep := repository.NewRoleResourcePermissionRepository(db.MysqlDB.DB())
	_, err = roleResourceRep.FindOne(ctx, map[string]interface{}{
		"role_id":        allowRoleId,
		"permission_key": permissionKey,
		"resource_type":  resourceType,
		"resource_id":    resourceId,
	})

	if enable {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			roleResource := models.RoleResourcePermission{
				RoleID:         allowRoleId,
				PermissionKey:  permissionKey,
				PermissionName: permissionName,
				ResourceType:   resourceType,
				ResourceID:     resourceId,
				CreateUser:     userID,
			}
			err = roleResourceRep.Create(ctx, &roleResource)
			if err != nil {
				return err
			}
		}
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			//delete
			err = roleResourceRep.DeleteByCondition(ctx, map[string]interface{}{
				"role_id":        allowRoleId,
				"permission_key": permissionKey,
				"resource_type":  resourceType,
				"resource_id":    resourceId})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetAssignRoleResourcePermissionInfo
// 获取当前用户能分配给此角色的资源权限情况
func GetAssignRoleResourcePermissionInfo(ctx context.Context, resourceType enum.ResourceType, roleId int64, page int, pageSize int, userID int64, isSuperAdmin bool) (*dto.AssignResourcePermissionRes, error) {

	userRoleRep := repository.NewUserRoleRepository(db.MysqlDB.DB())
	userResourceRep := repository.NewUserResourcePermissionRepository(db.MysqlDB.DB())
	roleResourceRep := repository.NewRoleResourcePermissionRepository(db.MysqlDB.DB())

	res := &dto.AssignResourcePermissionRes{}

	userRoleIds, err := userRoleRep.GetUserRoleIds(ctx, userID)
	if err != nil {
		return nil, err
	}

	res.Total = 0
	res.ResourcePermission = make([]dto.ResourcePermissionKey, 0)

	deviceRep := repository.NewDeviceRepository(db.MysqlDB.DB())
	deviceGroupRep := repository.NewDeviceGroupRepository(db.MysqlDB.DB())
	policyRep := repository.NewPolicyRepository(db.MysqlDB.DB())

	if resourceType == enum.ResourceTypeDevice || resourceType == enum.ResourceTypeDeviceGroup ||
		resourceType == enum.ResourceTypePolicy {

		var allResouceIds []int64
		var resources []dto.ResourcePermissionKey

		var devices []*models.Device
		var deviceGroups []*models.DeviceGroup
		var policys []*models.Policy

		var total int64
		if isSuperAdmin {
			if resourceType == enum.ResourceTypeDevice {
				devices, total, err = deviceRep.GetDevcies(ctx, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
			if resourceType == enum.ResourceTypeDeviceGroup {
				deviceGroups, total, err = deviceGroupRep.GetDevcieGroups(ctx, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
			if resourceType == enum.ResourceTypePolicy {
				policys, total, err = policyRep.GetPolicys(ctx, page, pageSize)
				if err != nil {
					return nil, err
				}
			}

		} else {
			// 用户拥有的可设置权限的DeviceId
			userOwnResourceIds, err := userResourceRep.GetUserOwnPermissionSetResourceIds(ctx, resourceType, userID)
			if err != nil {
				return nil, err
			}
			// 用户所属角色拥有的可设置权限的DeviceId
			roleOwnResourceIds, err := roleResourceRep.GetRolesOwnPermissionSetResourceIds(ctx, resourceType, userRoleIds)
			if err != nil {
				return nil, err
			}

			allResouceIds = lo.Union(userOwnResourceIds, roleOwnResourceIds)

			if resourceType == enum.ResourceTypeDevice {
				devices, total, err = deviceRep.GetDevicesByDeviceIds(ctx, allResouceIds, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
			if resourceType == enum.ResourceTypeDeviceGroup {
				deviceGroups, total, err = deviceGroupRep.GetDevcieGroupsByDeviceGroupIds(ctx, allResouceIds, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
			if resourceType == enum.ResourceTypePolicy {
				policys, total, err = policyRep.GetPolicysByPolicyIds(ctx, allResouceIds, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
		}

		for _, device := range devices {
			item := dto.ResourcePermissionKey{}
			item.ResourceID = device.ID
			item.ResourceName = device.DeviceName

			resources = append(resources, item)
		}
		for _, deviceGroup := range deviceGroups {
			item := dto.ResourcePermissionKey{}
			item.ResourceID = deviceGroup.ID
			item.ResourceName = deviceGroup.GroupName

			resources = append(resources, item)
		}
		for _, devicePolicy := range policys {
			item := dto.ResourcePermissionKey{}
			item.ResourceID = devicePolicy.ID
			item.ResourceName = devicePolicy.PolicyName

			resources = append(resources, item)
		}

		res.Total = total
		for _, resource := range resources {
			item := dto.ResourcePermissionKey{}
			item.ResourceID = resource.ResourceID
			item.ResourceName = resource.ResourceName
			item.UserOwnKeys = []dto.PermissionKeyName{}
			item.AssignKeys = []dto.PermissionKeyName{}
			if isSuperAdmin {
				allKeys, ok := permission.ResourceTypeKeyMap[resourceType]
				if !ok {
					return nil, errors.New("请求的资源类型错误")
				}
				for _, key := range allKeys {
					item.UserOwnKeys = append(item.UserOwnKeys, dto.PermissionKeyName{Key: key.Key, Name: key.Name, Order: key.Order})
				}
			} else {
				// 获取其他人授予我管理的资源
				userOwnKeys, err := userResourceRep.GetUserOwnResourcePermissionKeys(ctx, userID, resourceType, item.ResourceID, true)
				if err != nil {
					return nil, err
				}
				userRoleOwnKeys, err := roleResourceRep.GetRolesOwnResourcePermissionKeys(ctx, userRoleIds, resourceType, item.ResourceID)
				if err != nil {
					return nil, err
				}
				allKeys := lo.Union(userOwnKeys, userRoleOwnKeys)
				for _, key := range allKeys {
					permissionInfo, exists := permission.ResourceTypeKeyMap[resourceType][key]
					if exists {
						item.UserOwnKeys = append(item.UserOwnKeys, dto.PermissionKeyName{Key: permissionInfo.Key, Name: permissionInfo.Name, Order: permissionInfo.Order})
					} else {
						continue
					}
				}
			}
			roleOwnKeys, err := roleResourceRep.GetRoleOwnResourcePermissionKeys(ctx, roleId, resourceType, item.ResourceID)
			if err != nil {
				return nil, err
			}
			for _, key := range roleOwnKeys {
				permissionInfo, exists := permission.ResourceTypeKeyMap[resourceType][key]
				if exists {
					item.AssignKeys = append(item.AssignKeys, dto.PermissionKeyName{Key: permissionInfo.Key, Name: permissionInfo.Name, Order: permissionInfo.Order})
				} else {
					continue
				}
			}
			// 使用 sort.Slice 排序
			sort.Slice(item.UserOwnKeys, func(i, j int) bool {
				return item.UserOwnKeys[i].Order < item.UserOwnKeys[j].Order
			})

			res.ResourcePermission = append(res.ResourcePermission, item)
		}
	}

	return res, nil
}

// GetAssignUserResourcePermissionInfo
// 获取当前用户能分配给此用户的资源权限情况
func GetAssignUserResourcePermissionInfo(ctx context.Context, resourceType enum.ResourceType, assignUser int64, page int, pageSize int, userID int64, isSuperAdmin bool) (*dto.AssignResourcePermissionRes, error) {

	userRoleRep := repository.NewUserRoleRepository(db.MysqlDB.DB())
	userResourceRep := repository.NewUserResourcePermissionRepository(db.MysqlDB.DB())
	roleResourceRep := repository.NewRoleResourcePermissionRepository(db.MysqlDB.DB())

	res := &dto.AssignResourcePermissionRes{}

	userRoleIds, err := userRoleRep.GetUserRoleIds(ctx, userID)
	if err != nil {
		return nil, err
	}

	res.Total = 0
	res.ResourcePermission = make([]dto.ResourcePermissionKey, 0)

	deviceRep := repository.NewDeviceRepository(db.MysqlDB.DB())
	deviceGroupRep := repository.NewDeviceGroupRepository(db.MysqlDB.DB())
	policyRep := repository.NewPolicyRepository(db.MysqlDB.DB())

	if resourceType == enum.ResourceTypeDevice || resourceType == enum.ResourceTypeDeviceGroup || resourceType == enum.ResourceTypePolicy {
		var allResouceIds []int64
		var resources []dto.ResourcePermissionKey

		var devices []*models.Device
		var deviceGroups []*models.DeviceGroup
		var policys []*models.Policy

		var total int64
		if isSuperAdmin {
			if resourceType == enum.ResourceTypeDevice {
				devices, total, err = deviceRep.GetDevcies(ctx, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
			if resourceType == enum.ResourceTypeDeviceGroup {
				deviceGroups, total, err = deviceGroupRep.GetDevcieGroups(ctx, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
			if resourceType == enum.ResourceTypePolicy {
				policys, total, err = policyRep.GetPolicys(ctx, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
		} else {
			// 用户拥有的可设置权限的DeviceId
			userOwnResourceIds, err := userResourceRep.GetUserOwnPermissionSetResourceIds(ctx, resourceType, userID)
			if err != nil {
				return nil, err
			}
			// 用户所属角色拥有的可设置权限的DeviceId
			roleOwnResourceIds, err := roleResourceRep.GetRolesOwnPermissionSetResourceIds(ctx, resourceType, userRoleIds)
			if err != nil {
				return nil, err
			}

			allResouceIds = lo.Union(userOwnResourceIds, roleOwnResourceIds)

			if resourceType == enum.ResourceTypeDevice {
				devices, total, err = deviceRep.GetDevicesByDeviceIds(ctx, allResouceIds, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
			if resourceType == enum.ResourceTypeDeviceGroup {
				deviceGroups, total, err = deviceGroupRep.GetDevcieGroupsByDeviceGroupIds(ctx, allResouceIds, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
			if resourceType == enum.ResourceTypePolicy {
				policys, total, err = policyRep.GetPolicysByPolicyIds(ctx, allResouceIds, page, pageSize)
				if err != nil {
					return nil, err
				}
			}
		}

		for _, device := range devices {
			item := dto.ResourcePermissionKey{}
			item.ResourceID = device.ID
			item.ResourceName = device.DeviceName

			resources = append(resources, item)
		}
		for _, deviceGroup := range deviceGroups {
			item := dto.ResourcePermissionKey{}
			item.ResourceID = deviceGroup.ID
			item.ResourceName = deviceGroup.GroupName

			resources = append(resources, item)
		}
		for _, policy := range policys {
			item := dto.ResourcePermissionKey{}
			item.ResourceID = policy.ID
			item.ResourceName = policy.PolicyName

			resources = append(resources, item)
		}

		res.Total = total
		for _, resource := range resources {
			item := dto.ResourcePermissionKey{}
			item.ResourceID = resource.ResourceID
			item.ResourceName = resource.ResourceName
			item.UserOwnKeys = []dto.PermissionKeyName{}
			item.AssignKeys = []dto.PermissionKeyName{}
			if isSuperAdmin {
				allKeys, ok := permission.ResourceTypeKeyMap[resourceType]
				if !ok {
					return nil, errors.New("请求的资源类型错误")
				}
				for _, key := range allKeys {
					item.UserOwnKeys = append(item.UserOwnKeys, dto.PermissionKeyName{Key: key.Key, Name: key.Name, Order: key.Order})
				}
			} else {
				// 获取其他人授予我管理的资源
				userOwnKeys, err := userResourceRep.GetUserOwnResourcePermissionKeys(ctx, userID, resourceType, item.ResourceID, true)
				if err != nil {
					return nil, err
				}
				userRoleOwnKeys, err := roleResourceRep.GetRolesOwnResourcePermissionKeys(ctx, userRoleIds, resourceType, item.ResourceID)
				if err != nil {
					return nil, err
				}
				allKeys := lo.Union(userOwnKeys, userRoleOwnKeys)
				for _, key := range allKeys {
					permissionInfo, exists := permission.ResourceTypeKeyMap[resourceType][key]
					if exists {
						item.UserOwnKeys = append(item.UserOwnKeys, dto.PermissionKeyName{Key: permissionInfo.Key, Name: permissionInfo.Name, Order: permissionInfo.Order})
					} else {
						continue
					}
				}
			}
			userOwnKeys, err := userResourceRep.GetUserOwnResourcePermissionKeys(ctx, assignUser, resourceType, item.ResourceID, false)
			if err != nil {
				return nil, err
			}
			for _, key := range userOwnKeys {
				permissionInfo, exists := permission.ResourceTypeKeyMap[resourceType][key]
				if exists {
					item.AssignKeys = append(item.AssignKeys, dto.PermissionKeyName{Key: permissionInfo.Key, Name: permissionInfo.Name, Order: permissionInfo.Order})
				} else {
					continue
				}
			}

			// 使用 sort.Slice 排序
			sort.Slice(item.UserOwnKeys, func(i, j int) bool {
				return item.UserOwnKeys[i].Order < item.UserOwnKeys[j].Order
			})

			res.ResourcePermission = append(res.ResourcePermission, item)
		}
	}

	return res, nil
}

// AssignUserResourcePermission
// 分配资源给用户
func AssignUserResourcePermission(ctx context.Context, resourceType enum.ResourceType, resourceId int64, assignKey enum.PermissionKey, assignUser int64, userID int64, isSuperAdmin bool, isAdd bool) error {

	userRoleRep := repository.NewUserRoleRepository(db.MysqlDB.DB())
	userResourceRep := repository.NewUserResourcePermissionRepository(db.MysqlDB.DB())
	roleResourceRep := repository.NewRoleResourcePermissionRepository(db.MysqlDB.DB())

	userRoleIds, err := userRoleRep.GetUserRoleIds(ctx, userID)
	if err != nil {
		return err
	}

	if resourceType == enum.ResourceTypeDevice {
		if !isSuperAdmin {
			// 用户拥有的权限
			userHavePermission, err := userResourceRep.UserOwnResourcePermissionKey(ctx, userID, resourceType, resourceId, assignKey)
			if err != nil {
				return err
			}

			// 用户所属角色拥有的权限
			roleHavePermission, err := roleResourceRep.RoleOwnResourcePermissionKey(ctx, userRoleIds, resourceType, resourceId, assignKey)
			if err != nil {
				return err
			}
			if !userHavePermission && !roleHavePermission {
				return errors.New("没有操作权限")
			}
		}
	}

	if isAdd {
		userResourcePermission := &models.UserResourcePermission{
			UserID:         assignUser,
			PermissionKey:  assignKey,
			PermissionName: permission.ResourceTypeKeyMap[resourceType][assignKey].Name,
			ResourceType:   resourceType,
			ResourceID:     resourceId,
			CreateUser:     userID,
		}
		if err = userResourceRep.Create(ctx, userResourcePermission); err != nil {
			return err
		}
	} else {
		err = userResourceRep.DeleteByCondition(ctx, map[string]interface{}{
			"user_id":        assignUser,
			"permission_key": assignKey,
			"resource_type":  resourceType,
			"resource_id":    resourceId,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// AssignRoleResourcePermission
// 分配资源给角色
func AssignRoleResourcePermission(ctx context.Context, resourceType enum.ResourceType, resourceId int64, assignKey enum.PermissionKey, assignRole int64, userID int64, isSuperAdmin bool, isAdd bool) error {

	userRoleRep := repository.NewUserRoleRepository(db.MysqlDB.DB())
	userResourceRep := repository.NewUserResourcePermissionRepository(db.MysqlDB.DB())
	roleResourceRep := repository.NewRoleResourcePermissionRepository(db.MysqlDB.DB())

	userRoleIds, err := userRoleRep.GetUserRoleIds(ctx, userID)
	if err != nil {
		return err
	}

	if resourceType == enum.ResourceTypeDevice {
		if !isSuperAdmin {
			// 用户拥有的权限
			userHavePermission, err := userResourceRep.UserOwnResourcePermissionKey(ctx, userID, resourceType, resourceId, assignKey)
			if err != nil {
				return err
			}

			// 用户所属角色拥有的权限
			roleHavePermission, err := roleResourceRep.RoleOwnResourcePermissionKey(ctx, userRoleIds, resourceType, resourceId, assignKey)
			if err != nil {
				return err
			}
			if !userHavePermission && !roleHavePermission {
				return errors.New("没有操作权限")
			}
		}
	}

	if isAdd {
		roleResourcePermission := &models.RoleResourcePermission{
			RoleID:         assignRole,
			PermissionKey:  assignKey,
			PermissionName: permission.ResourceTypeKeyMap[resourceType][assignKey].Name,
			ResourceType:   resourceType,
			ResourceID:     resourceId,
			CreateUser:     userID,
		}

		if err = roleResourceRep.Create(ctx, roleResourcePermission); err != nil {
			return err
		}
	} else {
		err = roleResourceRep.DeleteByCondition(ctx, map[string]interface{}{
			"role_id":        assignRole,
			"permission_key": assignKey,
			"resource_type":  resourceType,
			"resource_id":    resourceId,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func GetResourceCreateUser(ctx context.Context, resourceType enum.ResourceType, resourceId int64) (int64, error) {
	var CreateUser int64
	CreateUser = -1

	if resourceType == enum.ResourceTypeRole {
		roleRep := repository.NewRoleRepository(db.MysqlDB.DB())
		role, err := roleRep.FindOne(ctx, map[string]interface{}{
			"id": resourceId,
		})
		if err != nil {
			return 0, fmt.Errorf("不存在的角色ID:%d,不能进行操作", resourceId)
		}
		CreateUser = role.CreateUser

	} else if resourceType == enum.ResourceTypeUser {
		userRep := repository.NewUserRepository(db.MysqlDB.DB())
		user, err := userRep.FindOne(ctx, map[string]interface{}{
			"id": resourceId,
		})
		if err != nil {
			return 0, fmt.Errorf("不存在的用户ID:%d,不能进行操作", resourceId)
		}
		CreateUser = user.CreateUser

	} else if resourceType == enum.ResourceTypeDevice {
		deviceRep := repository.NewDeviceRepository(db.MysqlDB.DB())
		_, err := deviceRep.FindOne(ctx, map[string]interface{}{
			"id": resourceId,
		})
		if err != nil {
			return 0, fmt.Errorf("不存在的设备ID:%d,不能进行操作", resourceId)
		}
		CreateUser = -1
	} else if resourceType == enum.ResourceTypeDeviceGroup {
		deviceGroupRep := repository.NewDeviceGroupRepository(db.MysqlDB.DB())
		deviceGroup, err := deviceGroupRep.FindOne(ctx, map[string]interface{}{
			"id": resourceId,
		})
		if err != nil {
			return 0, fmt.Errorf("不存在的服务器组ID:%d,不能进行操作", resourceId)
		}
		CreateUser = deviceGroup.CreateUser
	} else if resourceType == enum.ResourceTypePolicy {
		policyRep := repository.NewPolicyRepository(db.MysqlDB.DB())
		policy, err := policyRep.FindOne(ctx, map[string]interface{}{
			"id": resourceId,
		})
		if err != nil {
			return 0, fmt.Errorf("不存在的设备策略ID:%d,不能进行操作", resourceId)
		}
		CreateUser = policy.CreateUser
	} else {
		return 0, fmt.Errorf("未知的资源类型:%s", resourceType)
	}
	return CreateUser, nil
}
