package routers

import (
	"platform-backend/controller"
	"platform-backend/middleware"
	"platform-backend/service"

	"github.com/gin-gonic/gin"
)

// v1版本接口 /api
func routerV1(group *gin.RouterGroup) {
	v1Controller := controller.NewV1Controller(service.NewV2ServiceImpl())

	userGroup := group.Group("/user")
	deviceGroup := group.Group("/device")
	deviceGroupGroup := group.Group("/device_group")
	policyGroup := group.Group("/policy")
	roleGroup := group.Group("/role")
	permissionGroup := group.Group("/permission")

	/*
		No Auth Api
	*/
	{
		userGroup.POST("/login", v1Controller.UserLogIn)
		userGroup.POST("/refresh_token", v1Controller.UserRefreshToken)
	}

	/*
		User Api /user
	*/
	userGroup.Use(middleware.JWTAuthMiddleware())
	userGroup.Use(middleware.PermissionMiddleware())
	{
		userGroup.POST("/create", v1Controller.UserCreate)
		userGroup.DELETE("/delete", v1Controller.UserDelete)
		userGroup.PUT("/update_info", v1Controller.UserUpdateInfo)
		userGroup.POST("/reset_password", v1Controller.UserResetPassword)
		userGroup.POST("/update_password", v1Controller.UserUpdatePassword)
		userGroup.GET("/list", v1Controller.UserList)
		userGroup.GET("/info", v1Controller.UserInfo)
		userGroup.POST("/logout", v1Controller.UserLogOut)
	}

	/*
		Device Api /device
	*/
	deviceGroup.Use(middleware.JWTAuthMiddleware())
	deviceGroup.Use(middleware.PermissionMiddleware())
	{
		deviceGroup.GET("/list", v1Controller.DeviceList)
		deviceGroup.PUT("/update", v1Controller.DeviceUpdate)
		deviceGroup.DELETE("/delete", v1Controller.DeviceDelete)
		deviceGroup.POST("/poweron", v1Controller.DevicePowerOn)
		deviceGroup.POST("/poweroff", v1Controller.DevicePowerOff)
		deviceGroup.POST("/restart")
	}

	/*
		Server Group Api /device_group
	*/
	deviceGroupGroup.Use(middleware.JWTAuthMiddleware())
	deviceGroupGroup.Use(middleware.PermissionMiddleware())
	{
		deviceGroupGroup.POST("/create", v1Controller.DeviceGroupCreate)
		deviceGroupGroup.PUT("/update", v1Controller.DeviceGroupUpdate)
		deviceGroupGroup.DELETE("/delete", v1Controller.DeviceGroupDelete)
		deviceGroupGroup.GET("/list", v1Controller.DeviceGroupList)
		deviceGroupGroup.GET("/detail", v1Controller.DeviceGroupDetail)
		deviceGroupGroup.POST("/poweron", v1Controller.DeviceGroupPowerOn)
		deviceGroupGroup.POST("/poweroff", v1Controller.DeviceGroupPowerOff)
	}

	/*
		Policy Api /policy
	*/
	policyGroup.Use(middleware.JWTAuthMiddleware())
	policyGroup.Use(middleware.PermissionMiddleware())
	{
		policyGroup.POST("/create", v1Controller.PolicyCreate)
		policyGroup.PUT("/update", v1Controller.PolicyUpdate)
		policyGroup.GET("/list", v1Controller.PolicyList)
		policyGroup.DELETE("/delete", v1Controller.PolicyDelete)
	}

	/*
		Role Api /role
	*/
	roleGroup.Use(middleware.JWTAuthMiddleware())
	roleGroup.Use(middleware.PermissionMiddleware())
	{
		roleGroup.POST("/create", v1Controller.RoleCreate)
		roleGroup.PUT("/update", v1Controller.RoleUpdate)
		roleGroup.GET("/list", v1Controller.RoleList)
		roleGroup.DELETE("/delete", v1Controller.RoleDelete)

		roleGroup.GET("/menu_permission", v1Controller.RoleMenuPermission)
		roleGroup.PUT("/menu_permission", v1Controller.UpdateRoleMenuPermission)
	}

	/*
		Permission Api /permission
	*/
	permissionGroup.Use(middleware.JWTAuthMiddleware())
	{
		permissionGroup.GET("/list", v1Controller.PermissionList)
		permissionGroup.GET("/alluser", v1Controller.PermissionAllUser)
		permissionGroup.GET("/user", v1Controller.PermissionUser)
		permissionGroup.GET("/role", v1Controller.PermissionRole)
		permissionGroup.PUT("/alluser", v1Controller.PutPermissionAllUser)
		permissionGroup.PUT("/user", v1Controller.PutPermissionUser)
		permissionGroup.PUT("/role", v1Controller.PutPermissionRole)
		permissionGroup.GET("/role/own", v1Controller.GetRoleOwnResourcePermission)
		permissionGroup.GET("/user/own", v1Controller.GetUserOwnResourcePermission)
		permissionGroup.POST("/role/own", v1Controller.PostRoleOwnResourcePermission)
		permissionGroup.POST("/user/own", v1Controller.PostUserOwnResourcePermission)
	}
}
