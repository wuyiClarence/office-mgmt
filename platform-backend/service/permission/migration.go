package permission

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto/enum"
	"platform-backend/dto/permission"
	"platform-backend/models"
	"platform-backend/repository"
)

func Migration() error {
	var myViper *viper.Viper
	var data permission.PermissionData
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	myViper = viper.New()
	myViper.AddConfigPath(configPath)
	myViper.SetConfigName("menu_permission")
	myViper.SetConfigType("json")

	err := myViper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return err
	}

	if err = myViper.Unmarshal(&data); err != nil {
		return err
	}

	err = TraverseData(data)
	return err
}

func UpdatePermissionItem(rep *repository.MenuPermissionRepository, ctx context.Context,
	permissionName string, permissionKey enum.PermissionKey, parentId int64, order int) (bool, int64, error) {
	newRecord := false
	res, err := rep.FindOne(ctx, map[string]interface{}{
		"permission_key": permissionKey,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, 0, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res = &models.MenuPermission{
			PermissionName: permissionName,
			PermissionKey:  permissionKey,
			ParentID:       parentId,
			Order:          order,
			Used:           1,
		}
		err = rep.Create(ctx, res)
		if err != nil {
			return false, 0, err
		}
		newRecord = true
	} else {
		res.PermissionName = permissionName
		res.ParentID = parentId
		res.Used = 1
		res.Order = order
		err = rep.Update(ctx, res)
		if err != nil {
			return newRecord, 0, errors.New("更新失败")
		}
	}

	return newRecord, res.ID, err
}

func UpdatePermissionApiItem(rep *repository.MenuPermissionAPIApiRepository, ctx context.Context,
	permissionKey enum.PermissionKey, apiMethod string, apiPath string) error {

	res, err := rep.FindOne(ctx, map[string]interface{}{
		"permission_key": permissionKey,
		"api_method":     apiMethod,
		"api_path":       apiPath,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res = &models.MenuPermissionAPI{
			PermissionKey: permissionKey,
			ApiMethod:     apiMethod,
			ApiPath:       apiPath,
			Used:          1,
		}
		err = rep.Create(ctx, res)
		if err != nil {
			return err
		}
	} else {
		res.Used = 1
		err = rep.Update(ctx, res)
		if err != nil {
			return errors.New("更新失败")
		}
	}

	return err
}

func UpdateSystemBuiltInRolePermission(newRecord bool, roleAdmin bool, roleUser bool, roleAll bool, permissionKey enum.PermissionKey) error {
	roleRepo := repository.NewRoleRepository(db.MysqlDB.DB())
	roleMenuPermissionRepo := repository.NewRoleMenuPermissionRepository(db.MysqlDB.DB())
	ctx := context.Background()
	adminRole, err := roleRepo.FindOne(ctx, map[string]interface{}{
		"role_name":     "管理员",
		"system_create": 1,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	userRole, err := roleRepo.FindOne(ctx, map[string]interface{}{
		"role_name":     "普通用户",
		"system_create": 1,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	allRole, err := roleRepo.FindAll(ctx, map[string]interface{}{
		"system_create": "0",
	}, false)

	if newRecord {
		if roleAdmin {
			roleRes := &models.RoleMenuPermission{
				PermissionKey: permissionKey,
				RoleID:        adminRole.ID,
				CreateUser:    1,
			}
			err = roleMenuPermissionRepo.Create(ctx, roleRes)
			if err != nil {
				return err
			}
		}

		if roleUser {
			roleRes := &models.RoleMenuPermission{
				PermissionKey: permissionKey,
				RoleID:        userRole.ID,
				CreateUser:    1,
			}
			err = roleMenuPermissionRepo.Create(ctx, roleRes)
			if err != nil {
				return err
			}
		}

		if roleAll {
			rolePer := make([]models.RoleMenuPermission, 0, len(allRole))
			for _, role := range allRole {
				rolePer = append(rolePer, models.RoleMenuPermission{
					PermissionKey: permissionKey,
					RoleID:        role.ID,
					CreateUser:    0,
				})
			}
			err = roleMenuPermissionRepo.BatchInsert(ctx, rolePer)
			if err != nil {
				return err
			}
		}
	} else {
		if !roleAdmin {
			err = roleMenuPermissionRepo.DeleteByCondition(ctx, map[string]interface{}{
				"permission_key": permissionKey,
				"role_id":        adminRole.ID,
				"create_user":    1,
			})
			if err != nil {
				return err
			}
		}
		if !roleUser {
			err = roleMenuPermissionRepo.DeleteByCondition(ctx, map[string]interface{}{
				"permission_key": permissionKey,
				"role_id":        userRole.ID,
				"create_user":    1,
			})
			if err != nil {
				return err
			}
		}

		if !roleAll {
			for _, role := range allRole {
				err = roleMenuPermissionRepo.DeleteByCondition(ctx, map[string]interface{}{
					"permission_key": permissionKey,
					"role_id":        role.ID,
					"create_user":    1,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}

func IterationResouceList(rep *repository.MenuPermissionRepository, repapi *repository.MenuPermissionAPIApiRepository, ctx context.Context,
	permissionList []permission.Permission, parentId int64) error {
	order := 0
	for _, permissionItem := range permissionList {
		// fmt.Printf("Permission Name: %s, Permission Key: %s\n", permissionItem.Name, permissionItem.Key)
		order++
		newRecord, id, err := UpdatePermissionItem(rep, ctx, permissionItem.Name, permissionItem.Key, parentId, order)
		if err != nil {
			return err
		}

		err = UpdateSystemBuiltInRolePermission(newRecord, permissionItem.RoleAdmin, permissionItem.RoleUser, permissionItem.RoleAll, permissionItem.Key)
		if err != nil {
			return err
		}

		// 遍历权限的 API 列表
		for _, api := range permissionItem.Apis {
			err = UpdatePermissionApiItem(repapi, ctx, permissionItem.Key, strings.ToUpper(api.Method), api.Path)
			if err != nil {
				return err
			}
			// fmt.Printf("API Method: %s, API Path: %s\n", api.Method, api.Path)
		}

		// 遍历子权限
		if len(permissionItem.Children) > 0 {
			err = IterationResouceList(rep, repapi, ctx, permissionItem.Children, id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DeleteUnUsedPermission() error {
	repo := repository.NewMenuPermissionRepository(db.MysqlDB.DB())

	ctx := context.Background()
	err := repo.DeleteUnUsedPermission(ctx)
	if err != nil {
		return err
	}

	apiRepo := repository.NewMenuPermissionApiRepository(db.MysqlDB.DB())
	err = apiRepo.DeleteUnUsedPermissionApi(ctx)
	return err
}

func TraverseData(data permission.PermissionData) error {
	menuPermissionRepo := repository.NewMenuPermissionRepository(db.MysqlDB.DB())
	menuPermissionApiRepo := repository.NewMenuPermissionApiRepository(db.MysqlDB.DB())
	ctx := context.Background()

	err := menuPermissionRepo.UpdateAllUsedStatus(ctx)
	if err != nil {
		return err
	}

	err = menuPermissionApiRepo.UpdateAllUsedStatus(ctx)
	if err != nil {
		return err
	}

	err = IterationResouceList(menuPermissionRepo, menuPermissionApiRepo, ctx, data.PermissionList, 0)
	if err != nil {
		return err
	}

	err = DeleteUnUsedPermission()
	if err != nil {
		return err
	}

	return nil
}
