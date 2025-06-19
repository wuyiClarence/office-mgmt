package service

import (
	"platform-backend/dto"

	"github.com/gin-gonic/gin"
)

type V1Service interface {
	UserLogIn(*gin.Context, *dto.UserLoginApiReq) (*dto.UserLoginApiRes, error)
	UserUpdatePassword(*gin.Context, *dto.UserUpdatePasswordApiReq) error
	UserRefreshToken(*gin.Context, *dto.UserRefreshTokenReq) (*dto.UserRefreshTokenRes, error)
	UserCreate(*gin.Context, *dto.UserCreateReq) error
	UserDelete(*gin.Context, *dto.UserDeleteReq) error
	UserResetPassword(*gin.Context, *dto.UserResetPasswordReq) error
	UserUpdateInfo(*gin.Context, *dto.UserUpdateInfoReq) error
	UserInfo(*gin.Context) (*dto.UserInfoRes, error)
	UserLogOut(*gin.Context) error
	UserList(*gin.Context, *dto.UserListReq) (*dto.UserListRes, error)

	RoleCreate(*gin.Context, *dto.RoleCreateReq) error
	RoleUpdate(*gin.Context, *dto.RoleUpdateReq) error
	RoleList(*gin.Context, *dto.RoleListReq) (*dto.RoleListRes, error)

	RoleDelete(*gin.Context, *dto.RoleDelReq) error
	RoleMenuPermissionList(*gin.Context, *dto.RolePermissionReq) (*dto.PermissionListRes, error)
	RoleMenuPermissionUpdate(*gin.Context, *dto.RoleUpdateMenuPermissionReq) error

	DeviceList(*gin.Context, *dto.ListReq) (*dto.DeviceListRes, error)
	DeviceDelete(*gin.Context, *dto.DeviceDelReq) error
	DeviceUpdate(*gin.Context, *dto.DeviceUpdateReq) error
	DevicePowerOn(*gin.Context, *dto.DevicePowerOnReq) error
	DevicePowerOff(*gin.Context, *dto.DevicePowerOffReq) error

	DeviceGroupCreate(*gin.Context, *dto.DeviceGroupCreateReq) error
	DeviceGroupUpdate(*gin.Context, *dto.DeviceGroupUpdateReq) error
	DeviceGroupDelete(*gin.Context, *dto.DeviceGroupDelReq) error
	DeviceGroupList(*gin.Context, *dto.ListReq) (*dto.DeviceGroupListRes, error)
	DeviceGroupDetail(*gin.Context, *dto.DeviceGroupDetailReq) (*dto.DeviceGroupDetail, error)
	DeviceGroupPowerOn(*gin.Context, *dto.DeviceGroupPowerOnReq) error
	DeviceGroupPowerOff(*gin.Context, *dto.DeviceGroupPowerOffReq) error

	PolicyCreate(*gin.Context, *dto.PolicyCreateReq) error
	PolicyUpdate(*gin.Context, *dto.PolicyUpdateReq) error
	PolicyList(*gin.Context, *dto.ListReq) (*dto.PolicyListRes, error)

	PolicyDelete(*gin.Context, *dto.PolicyDelReq) error

	PermissionList(*gin.Context) (*dto.PermissionListRes, error)
	GetResourceAllUserPermissionInfo(*gin.Context, *dto.ResourcePermissionReq) (*dto.ResourcePermissionKeys, error)
	GetResourceUserPermissionInfo(*gin.Context, *dto.ResourcePermissionReq) (*dto.UserResourcePermissionRes, error)
	GetResourceRolePermissionInfo(*gin.Context, *dto.ResourcePermissionReq) (*dto.RoleResourcePermissionRes, error)
	PutResourceAllUserPermission(*gin.Context, *dto.ResourcePermissionPutReq) error
	PutResourceUserPermission(*gin.Context, *dto.ResourcePermissionPutUserReq) error
	PutResourceRolePermission(*gin.Context, *dto.ResourcePermissionPutRoleReq) error

	GetRoleOwnResourcePermission(*gin.Context, *dto.RoleOwnResourcePermissionReq) (*dto.AssignResourcePermissionRes, error)
	GetUserOwnResourcePermission(*gin.Context, *dto.UserOwnResourcePermissionReq) (*dto.AssignResourcePermissionRes, error)
	PostUserOwnResourcePermission(*gin.Context, *dto.PostUserOwnResourcePermissionReq) error
	PostRoleOwnResourcePermission(*gin.Context, *dto.PostRoleOwnResourcePermissionReq) error
}
