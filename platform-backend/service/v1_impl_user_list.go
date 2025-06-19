package service

import (
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	db "platform-backend/db"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/repository"
	sp "platform-backend/service/permission"
	"platform-backend/utils"
)

func (svc *V1ServiceImpl) UserList(c *gin.Context, req *dto.UserListReq) (*dto.UserListRes, error) {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, err
	}

	userRep := repository.NewUserRepository(db.MysqlDB.DB())

	var (
		res = &dto.UserListRes{
			Total: 0,
			List:  make([]*dto.UserListItem, 0),
		}
	)

	allIds, err := sp.GetUserWithResourceIds(c, enum.ResourceTypeUser, enum.PermissionKeyView, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}
	allItems, total, err := userRep.GetUserByUserIds(c.Request.Context(), allIds,
		req.PageIndex, req.PageSize)
	if err != nil {
		return nil, err
	}

	resourceIds := lo.Map(allItems, func(p *dto.UserListItem, _ int) int64 {
		return p.UserID
	})
	permissions, err := sp.GetUserWithResourcePermission(c, enum.ResourceTypeUser, resourceIds, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}
	for _, item := range allItems {
		//获取User的角色
		item.RoleInfos, err = repository.NewUserRoleRepository(db.MysqlDB.DB()).GetUserRoleInfos(c.Request.Context(), item.UserID)
		if err != nil {
			return nil, err
		}
		for _, value := range permissions[item.UserID] {
			item.Permissions = append(item.Permissions, value)
		}
		sort.Slice(item.Permissions, func(i, j int) bool {
			return item.Permissions[i].Order < item.Permissions[j].Order
		})
	}
	res.List = allItems
	res.Total = total

	return res, nil
}
