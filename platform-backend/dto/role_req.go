package dto

import "platform-backend/dto/enum"

type RoleCreateReq struct {
	RoleName    string               `json:"role_name" binding:"required"`
	Description string               `json:"description"`
	Menu        []enum.PermissionKey `json:"menu"`
}

// RoleAssignReq 包括首次分配及后续更新
type RoleAssignReq struct {
	AssignedUserID int64                `json:"assigned_user_id" binding:"required"`
	RoleID         int64                `json:"role_id" binding:"required"`
	Permissions    []enum.PermissionKey `json:"permissions" binding:"required"`
}

type RoleDelReq struct {
	RoleIDs []int64 `json:"role_ids" binding:"required"`
}

type RoleUpdateReq struct {
	RoleID      int64                `json:"role_id" binding:"required"`
	RoleName    string               `json:"role_name"`
	Description string               `json:"description"`
	Menu        []enum.PermissionKey `json:"menu"`
}

type RoleListReq struct {
	RoleName string `json:"role_name"`
	ListReq
}

type RoleListRes struct {
	Total int64           `json:"total"`
	List  []*RoleListItem `json:"list"`
}

type RoleListItem struct {
	ID                    int64        `json:"id"`
	RoleName              string       `json:"role_name"`
	Description           string       `json:"description"`
	CreateUser            int64        `json:"create_user"`
	CreateUserName        string       `json:"create_user_name"`
	CreateUserDisplayName string       `json:"create_user_display_name"`
	CreatedAt             string       `json:"created_at"`
	SystemCreate          int64        `json:"system_create"`
	Permissions           []Permission `json:"permissions" gorm:"-"`
}

type RoleInfoReq struct {
	RoleID int64 `json:"role_id" binding:"required"`
	ListReq
}

type RoleDeviceGroupItem struct {
	DeviceGroupID   int64  `json:"device_group_id"`
	DeviceGroupName string `json:"device_group_name"`
}

type RoleUserItem struct {
	UserID          int64  `json:"user_id"`
	UserDisplayName string `json:"user_display_name"`
	PhoneNumber     string `json:"phone_number"`
	Email           string `json:"email"`
}

type RoleDeviceItem struct {
	DeviceID   int64           `json:"device_id"`
	DeviceName string          `json:"device_name"`
	Mac        string          `json:"mac"`
	IP         string          `json:"ip"`
	OsType     string          `json:"os_type"`
	DeviceType enum.DeviceType `json:"device_type"`
}

type UserRoleItem struct {
	UserID          int64  `json:"user_id"`
	UserDisplayName string `json:"user_display_name"`
	PhoneNumber     string `json:"phone_number"`
	Email           string `json:"email"`
}

type RolePermissionReq struct {
	RoleID int64 `json:"role_id" form:"role_id" binding:"required"`
}

type RoleUpdateMenuPermissionReq struct {
	RoleID int64                `json:"role_id" binding:"required"`
	Menu   []enum.PermissionKey `json:"menu"`
}
