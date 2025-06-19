package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/repository"
	sp "platform-backend/service/permission"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) RoleDeleteByRoleId(c *gin.Context, roleID int64) error {
	return repository.DoInTx(c.Request.Context(), db.MysqlDB.DB(), func(ctx context.Context, tx *gorm.DB) error {
		// 1.role
		err := repository.NewRoleRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"id": roleID,
		})
		if err != nil {
			return err
		}

		err = repository.NewUserResourcePermissionRepository(tx).DeleteByCondition(ctx, map[string]interface{}{
			"resource_id":   roleID,
			"resource_type": enum.ResourceTypeRole,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (svc *V1ServiceImpl) RoleDelete(c *gin.Context, req *dto.RoleDelReq) error {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}

	err = sp.CheckPermission(c, enum.ResourceTypeRole, req.RoleIDs, enum.PermissionKeyDelete, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	roleRep := repository.NewRoleRepository(db.MysqlDB.DB())
	for _, roleID := range req.RoleIDs {
		role, err := roleRep.FindOne(c.Request.Context(), map[string]interface{}{
			"id": roleID,
		})
		if err != nil {
			return fmt.Errorf("不存在的角色ID:%d,不能进行删除操作", roleID)
		}
		if role.SystemCreate == 1 {
			return fmt.Errorf("角色ID:%d是系统角色,不能进行删除操作", roleID)
		}
		total, userWithRoleTrans, err := repository.NewUserRoleRepository(db.MysqlDB.DB()).GetWithRoleUserList(c.Request.Context(), roleID, 1, 5)
		if err != nil {
			return err
		}
		if total != 0 {
			var userinfo string
			for _, user := range userWithRoleTrans {
				userinfo = fmt.Sprintf("%s%s ", userinfo, user.UserName)
			}
			return fmt.Errorf("不能对角色ID:%d进行删除操作,存在使用的用户%s", roleID, userinfo)
		}
		err = svc.RoleDeleteByRoleId(c, roleID)
		if err != nil {
			return err
		}
	}
	return nil
}
