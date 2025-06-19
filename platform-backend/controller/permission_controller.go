package controller

import (
	"platform-backend/dto"
	"platform-backend/utils"
	"platform-backend/utils/format"

	"github.com/gin-gonic/gin"
)

// PermissionList godoc
// @Summary PermissionList
// @Description PermissionList Api 查询用户可调用的api列表
// @Tags Permission
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.PermissionListRes
// @Router /api/permission/list [get]
func (ctrl *V1Controller) PermissionList(c *gin.Context) {
	res, err := ctrl.svc.PermissionList(c)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// PermissionAllUser godoc
// @Summary PermissionAllUser
// @Description PermissionAllUser Api 查询资源分配给所有人的权限情况
// @Tags Permission
// @Accept  json
// @Produce  json
// @Param resource_id query int true "资源ID"
// @Param resource_type query int true "资源类型"
// @Success 200 {object} dto.PermissionListRes
// @Router /api/permission/alluser [get]
func (ctrl *V1Controller) PermissionAllUser(c *gin.Context) {
	req := &dto.ResourcePermissionReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}
	res, err := ctrl.svc.GetResourceAllUserPermissionInfo(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// PutPermissionAllUser godoc
// @Summary PutPermissionAllUser
// @Description PutPermissionAllUser Api 修改资源分配给所有人的权限
// @Tags Permission
// @Accept  json
// @Produce  json
// @Param PutPermissionAllUser body dto.ResourcePermissionPutReq true "Put Permission AllUser request body"
// @Success 200 {object} nil
// @Router /api/permission/alluser [put]
func (ctrl *V1Controller) PutPermissionAllUser(c *gin.Context) {
	req := &dto.ResourcePermissionPutReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}
	err := ctrl.svc.PutResourceAllUserPermission(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// PermissionUser godoc
// @Summary PermissionUser
// @Description PermissionUser Api 查询资源分配给指定人的权限情况
// @Tags Permission
// @Accept  json
// @Produce  json
// @Param resource_id query int true "资源ID"
// @Param resource_type query int true "资源类型"
// @Success 200 {object} dto.UserResourcePermissionRes
// @Router /api/permission/user [get]
func (ctrl *V1Controller) PermissionUser(c *gin.Context) {
	req := &dto.ResourcePermissionReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}
	res, err := ctrl.svc.GetResourceUserPermissionInfo(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// PutPermissionUser godoc
// @Summary PutPermissionUser
// @Description PutPermissionUser Api 修改资源分配给指定用户的权限
// @Tags Permission
// @Accept  json
// @Produce  json
// @Param ResourcePermissionPutUserReq body dto.ResourcePermissionPutUserReq true "Put Permission User request body"
// @Success 200 {object} nil
// @Router /api/permission/user [put]
func (ctrl *V1Controller) PutPermissionUser(c *gin.Context) {
	req := &dto.ResourcePermissionPutUserReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}
	err := ctrl.svc.PutResourceUserPermission(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// PermissionRole godoc
// @Summary PermissionRole
// @Description PermissionRole Api 查询资源分配给指定角色的情况
// @Tags Permission
// @Accept  json
// @Produce  json
// @Param resource_id query int true "资源ID"
// @Param resource_type query int true "资源类型"
// @Success 200 {object} dto.RoleResourcePermissionRes
// @Router /api/permission/role [get]
func (ctrl *V1Controller) PermissionRole(c *gin.Context) {
	req := &dto.ResourcePermissionReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}
	res, err := ctrl.svc.GetResourceRolePermissionInfo(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// PutPermissionRole godoc
// @Summary PutPermissionRole
// @Description PutPermissionRole Api 修改资源分配给指定用户的权限
// @Tags Permission
// @Accept  json
// @Produce  json
// @Param ResourcePermissionPutUserReq body dto.ResourcePermissionPutUserReq true "Put Permission User request body"
// @Success 200 {object} nil
// @Router /api/permission/role [put]
func (ctrl *V1Controller) PutPermissionRole(c *gin.Context) {
	req := &dto.ResourcePermissionPutRoleReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}
	err := ctrl.svc.PutResourceRolePermission(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// GetRoleOwnResourcePermission godoc
// @Summary GetRoleOwnResourcePermission
// @Description GetRoleOwnResourcePermission Api 获取可以对指定角色分配的一类资源的权限集合
// @Tags Permission
// @Param role_id query int true "角色ID"
// @Param resource_type query int true "资源类型"
// @Success 200 {object} dto.ResourcePermissionRes
// @Router /api/permission/role/own [get]
func (ctrl *V1Controller) GetRoleOwnResourcePermission(c *gin.Context) {
	req := &dto.RoleOwnResourcePermissionReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}
	res, err := ctrl.svc.GetRoleOwnResourcePermission(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// GetUserOwnResourcePermission godoc
// @Summary GetUserOwnResourcePermission
// @Description GetUserOwnResourcePermission Api 获取可以对指定用户分配的一类资源的权限集合
// @Tags Permission
// @Param user_id query int true "用户ID"
// @Param resource_type query int true "资源类型"
// @Success 200 {object} dto.ResourcePermissionRes
// @Router /api/permission/user/own [get]
func (ctrl *V1Controller) GetUserOwnResourcePermission(c *gin.Context) {
	req := &dto.UserOwnResourcePermissionReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}
	res, err := ctrl.svc.GetUserOwnResourcePermission(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// PostUserOwnResourcePermission godoc
// @Summary PostUserOwnResourcePermission
// @Description PostUserOwnResourcePermission Api 设置用户分配的一类资源的权限集合
// @Tags Permission
// @Param UserOwnResourcePermission body dto.PostUserOwnResourcePermissionReq true "创建用户"
// @Success 200 {object} nil
// @Router /api/permission/user/own [post]
func (ctrl *V1Controller) PostUserOwnResourcePermission(c *gin.Context) {
	req := &dto.PostUserOwnResourcePermissionReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.PostUserOwnResourcePermission(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// PostRoleOwnResourcePermission godoc
// @Summary PostRoleOwnResourcePermission
// @Description PostRoleOwnResourcePermission Api 设置用户分配的一类资源的权限集合
// @Tags Permission
// @Param RoleOwnResourcePermission body dto.PostRoleOwnResourcePermissionReq true "创建用户"
// @Success 200 {object} nil
// @Router /api/permission/role/own [post]
func (ctrl *V1Controller) PostRoleOwnResourcePermission(c *gin.Context) {
	req := &dto.PostRoleOwnResourcePermissionReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.PostRoleOwnResourcePermission(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}
