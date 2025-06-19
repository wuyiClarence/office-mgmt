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

func (svc *V1ServiceImpl) RoleList(c *gin.Context, req *dto.RoleListReq) (*dto.RoleListRes, error) {
	userID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, err
	}

	var (
		res = &dto.RoleListRes{
			Total: 0,
			List:  make([]*dto.RoleListItem, 0),
		}

		condition = make(map[string]interface{})
	)

	if len(req.RoleName) > 0 {
		condition["role_name"] = req.RoleName
	}

	roleRep := repository.NewRoleRepository(db.MysqlDB.DB())

	allIds, err := sp.GetUserWithResourceIds(c, enum.ResourceTypeRole, enum.PermissionKeyView, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}

	allItems, total, err := roleRep.GetRoleByRoleIds(c.Request.Context(), allIds,
		req.PageIndex, req.PageSize)
	if err != nil {
		return nil, err
	}

	resourceIds := lo.Map(allItems, func(p *dto.RoleListItem, _ int) int64 {
		return p.ID
	})
	permissions, err := sp.GetUserWithResourcePermission(c, enum.ResourceTypeRole, resourceIds, userID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}

	for _, item := range allItems {
		for _, value := range permissions[item.ID] {
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
