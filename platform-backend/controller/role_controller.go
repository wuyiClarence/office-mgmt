package controller

import (
	"platform-backend/dto"
	"platform-backend/utils"
	"platform-backend/utils/format"

	"github.com/gin-gonic/gin"
)

// RoleCreate godoc
// @Summary RoleCreate
// @Description RoleCreate Api
// @Tags Role
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param RoleCreate body dto.RoleCreateReq true "创建用户"
// @Success 200 {object} nil
// @Router /api/role/create [post]
func (ctrl *V1Controller) RoleCreate(c *gin.Context) {
	req := &dto.RoleCreateReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.RoleCreate(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// RoleUpdate godoc
// @Summary RoleUpdate
// @Description RoleUpdate Api
// @Tags Role
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param RoleUpdate body dto.RoleUpdateReq true "更新角色"
// @Success 200 {object} nil
// @Router /api/role/update [put]
func (ctrl *V1Controller) RoleUpdate(c *gin.Context) {
	req := &dto.RoleUpdateReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.RoleUpdate(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// RoleList godoc
// @Summary RoleList
// @Description RoleList Api
// @Tags Role
// @Accept  json
// @Produce  json
// @Param pageIndex query int false "页索引"
// @Param pageSize query int false "页大小"
// @Param role_name query string false "角色名字"
// @Success 200 {object} dto.RoleListRes
// @Router /api/role/list [get]
func (ctrl *V1Controller) RoleList(c *gin.Context) {
	req := &dto.RoleListReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	utils.HandlePagination(&req.ListReq)

	res, err := ctrl.svc.RoleList(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// RoleDelete godoc
// @Summary RoleDelete
// @Description RoleDelete Api
// @Tags Role
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param RoleDelete body dto.RoleDelReq true "删除角色"
// @Success 200 {object} nil
// @Router /api/role/delete [delete]
func (ctrl *V1Controller) RoleDelete(c *gin.Context) {
	req := &dto.RoleDelReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.RoleDelete(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// RoleMenuPermission godoc
// @Summary RoleMenuPermission
// @Description RoleMenuPermission Api 查询角色拥有的菜单权限
// @Tags Role
// @Accept  json
// @Produce  json
// @Param role_id query int true "角色id"
// @Success 200 {object} dto.PermissionListRes
// @Router /api/role/menu_permission [get]
func (ctrl *V1Controller) RoleMenuPermission(c *gin.Context) {
	req := &dto.RolePermissionReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	res, err := ctrl.svc.RoleMenuPermissionList(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// UpdateRoleMenuPermission godoc
// @Summary UpdateRoleMenuPermission
// @Description UpdateRoleMenuPermission Api 更新角色拥有的菜单权限
// @Tags Role
// @Accept  json
// @Produce  json
// @Param RoleMenuPermissionUpdate body dto.RoleUpdateMenuPermissionReq true "Update Role Menu Permission request body"
// @Success 200 {object} nil
// @Router /api/role/menu_permission [put]
func (ctrl *V1Controller) UpdateRoleMenuPermission(c *gin.Context) {
	req := &dto.RoleUpdateMenuPermissionReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.RoleMenuPermissionUpdate(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}
