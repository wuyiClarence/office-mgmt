package controller

import (
	"github.com/gin-gonic/gin"

	"platform-backend/dto"
	"platform-backend/utils"
	"platform-backend/utils/format"
)

// UserLogIn godoc
// @Summary UserLogIn
// @Description UserLogIn Api 用户登陆
// @Tags User
// @Accept  json
// @Produce  json
// @Param UserLogin body dto.UserLoginApiReq true "UserLogin request body"
// @Success 200 {object} dto.UserLoginApiRes
// @Router /api/user/login [post]
func (ctrl *V1Controller) UserLogIn(c *gin.Context) {

	req := &dto.UserLoginApiReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	res, err := ctrl.svc.UserLogIn(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// UserUpdatePassword godoc
// @Summary UserUpdatePassword
// @Description UserUpdatePassword Api 更新密码
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param UserUpdatePassword body dto.UserUpdatePasswordApiReq true "UserUpdatePassword request body"
// @Success 200 {object} nil
// @Router /api/user/update_password [post]
func (ctrl *V1Controller) UserUpdatePassword(c *gin.Context) {

	req := &dto.UserUpdatePasswordApiReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.UserUpdatePassword(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// UserRefreshToken godoc
// @Summary UserRefreshToken
// @Description UserRefreshToken Api 刷新access-token
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param UserRefreshToken body dto.UserRefreshTokenReq true "UserRefreshToken request body"
// @Success 200 {object} dto.UserRefreshTokenRes
// @Router /api/user/refresh_token [post]
func (ctrl *V1Controller) UserRefreshToken(c *gin.Context) {

	req := &dto.UserRefreshTokenReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	res, err := ctrl.svc.UserRefreshToken(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// UserCreate godoc
// @Summary UserCreate
// @Description UserCreate Api 创建用户
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param UserCreate body dto.UserCreateReq true "Create User request body"
// @Success 200 {object} nil
// @Router /api/user/create [post]
func (ctrl *V1Controller) UserCreate(c *gin.Context) {

	req := &dto.UserCreateReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.UserCreate(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// UserDelete godoc
// @Summary UserDelete
// @Description UserDelete Api 删除用户
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param UserDelete body dto.UserDeleteReq true "Delete User request body"
// @Success 200 {object} nil
// @Router /api/user/delete [delete]
func (ctrl *V1Controller) UserDelete(c *gin.Context) {

	req := &dto.UserDeleteReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.UserDelete(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// UserUpdateInfo godoc
// @Summary UserUpdateInfo
// @Description UserUpdateInfo Api 更新用户信息
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param UserUpdateInfo body dto.UserUpdateInfoReq true "Update User Info request body"
// @Success 200 {object} nil
// @Router /api/user/update_info [put]
func (ctrl *V1Controller) UserUpdateInfo(c *gin.Context) {

	req := &dto.UserUpdateInfoReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.UserUpdateInfo(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// UserResetPassword godoc
// @Summary UserResetPassword
// @Description UserResetPassword Api 重置密码
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param UserResetPassword body dto.UserResetPasswordReq true "Reset User Password request body"
// @Success 200 {object} nil
// @Router /api/user/reset_password [post]
func (ctrl *V1Controller) UserResetPassword(c *gin.Context) {

	req := &dto.UserResetPasswordReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.UserResetPassword(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// UserInfo godoc
// @Summary UserInfo
// @Description UserInfo Api 查询用户的信息，包含用户权限
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.UserInfoRes
// @Router /api/user/info [get]
func (ctrl *V1Controller) UserInfo(c *gin.Context) {
	res, err := ctrl.svc.UserInfo(c)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// UserLogOut godoc
// @Summary UserLogOut
// @Description UserLogOut Api 用户退出登录
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} nil
// @Router /api/user/logout [post]
func (ctrl *V1Controller) UserLogOut(c *gin.Context) {

	err := ctrl.svc.UserLogOut(c)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// UserList godoc
// @Summary UserList
// @Description UserList Api 查询用户的信息，包含用户权限
// @Tags User
// @Accept  json
// @Produce  json
// @Param pageIndex query int false "页索引"
// @Param pageSize query int false "页大小"
// @Success 200 {object} dto.UserListRes
// @Router /api/user/list [get]
func (ctrl *V1Controller) UserList(c *gin.Context) {
	req := &dto.UserListReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}
	utils.HandlePagination(&req.ListReq)

	res, err := ctrl.svc.UserList(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}
