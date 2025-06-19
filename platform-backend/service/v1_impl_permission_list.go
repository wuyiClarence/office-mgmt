package service

import (
	"context"
	"errors"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/models"
	"platform-backend/repository"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) PermissionList(c *gin.Context) (*dto.PermissionListRes, error) {
	isSuperAdmin := utils.IsSuperAdmin(c)

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, err
	}

	userPermission, err := repository.NewMenuPermissionRepository(db.MysqlDB.DB()).FindUserMenuPermission(c.Request.Context(), curUserID, isSuperAdmin)
	if err != nil {
		return nil, err
	}

	res := dto.PermissionListRes{}

	if userPermission == nil {
		return &res, nil
	}

	res.Menu = lo.Map(userPermission, func(p *models.MenuPermission, _ int) enum.PermissionKey {
		return p.PermissionKey
	})

	topmenu := lo.Filter(userPermission, func(p *models.MenuPermission, _ int) bool {
		return p.ParentID == 0
	})

	sort.Slice(topmenu, func(i, j int) bool {
		return topmenu[i].Order < topmenu[j].Order
	})

	for _, menu := range topmenu {
		res.PermissionTree = append(res.PermissionTree, dto.Permission{
			Key:      menu.PermissionKey,
			Name:     menu.PermissionName,
			Children: buildPermissionTree(menu.ID, userPermission),
		})
	}

	return &res, nil
}

func buildPermissionTree(parentId int64, permissions []*models.MenuPermission) []dto.Permission {
	var roots []dto.Permission
	menu := lo.Filter(permissions, func(p *models.MenuPermission, _ int) bool {
		return p.ParentID == parentId
	})

	sort.Slice(menu, func(i, j int) bool {
		return menu[i].Order < menu[j].Order
	})

	for _, m := range menu {
		roots = append(roots, dto.Permission{
			Name:     m.PermissionName,
			Key:      m.PermissionKey,
			Children: buildPermissionTree(m.ID, permissions),
		})
	}
	return roots
}

func (svc *V1ServiceImpl) RoleMenuPermissionList(c *gin.Context, req *dto.RolePermissionReq) (*dto.PermissionListRes, error) {

	rolePermission, err := repository.NewMenuPermissionRepository(db.MysqlDB.DB()).FindRoleMenuPermission(c.Request.Context(), req.RoleID)
	if err != nil {
		return nil, err
	}

	res := dto.PermissionListRes{}

	if rolePermission == nil {
		return &res, nil
	}

	res.Menu = lo.Map(rolePermission, func(p *models.MenuPermission, _ int) enum.PermissionKey {
		return p.PermissionKey
	})
	userPermission, err := svc.PermissionList(c)
	if err != nil {
		return nil, err
	}
	res.PermissionTree = userPermission.PermissionTree

	return &res, nil
}

func (svc *V1ServiceImpl) RoleMenuPermissionUpdate(c *gin.Context, req *dto.RoleUpdateMenuPermissionReq) error {
	permissionRepo := repository.NewMenuPermissionRepository(db.MysqlDB.DB())
	roleMenuPermissionRepo := repository.NewRoleMenuPermissionRepository(db.MysqlDB.DB())
	ctx := c.Request.Context()

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	res, err := svc.PermissionList(c)
	if err != nil {
		return err
	}
	haveRight := lo.Every(res.Menu, req.Menu)
	if !haveRight {
		return errors.New("没有操作权限")
	}

	err = repository.DoInTx(ctx, db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		rolePermission, err := permissionRepo.FindRoleMenuPermission(ctx, req.RoleID)
		if err != nil {
			return err
		}
		// 提取 PermissionKey 的集合
		permissionKeys := lo.Map(rolePermission, func(p *models.MenuPermission, _ int) enum.PermissionKey {
			return p.PermissionKey
		})
		addMenu := lo.Without(req.Menu, permissionKeys...)
		if len(addMenu) > 0 {

			roleMenuPermissions := make([]models.RoleMenuPermission, 0, len(addMenu))
			for _, item := range addMenu {
				roleMenuPermissions = append(roleMenuPermissions, models.RoleMenuPermission{
					RoleID:        req.RoleID,
					PermissionKey: item,
					CreateUser:    curUserID,
				})
			}
			err = roleMenuPermissionRepo.BatchInsert(ctx, roleMenuPermissions)
			if err != nil {
				return err
			}
		}

		delMenu := lo.Without(permissionKeys, req.Menu...)
		// fmt.Println("删除的:", delMenu)
		if len(delMenu) > 0 {

			err = roleMenuPermissionRepo.DeleteByCondition(ctx, map[string]interface{}{
				"role_id":        req.RoleID,
				"permission_key": delMenu,
			})
			if err != nil {
				return err
			}
		}
		return err
	})

	return err
}
