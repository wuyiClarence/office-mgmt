package dto

import "platform-backend/dto/enum"

type API struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type Permission struct {
	Name     string             `json:"name"`
	Key      enum.PermissionKey `json:"key"`
	Order    int                `json:"order"`
	Children []Permission       `json:"children,omitempty"` // 子资源，可能为空
}

type PermissionListRes struct {
	Menu           []enum.PermissionKey `json:"menu"`
	PermissionTree []Permission         `json:"permission_tree"`
}

type ResourcePermissionReq struct {
	ResourceID   int64             `json:"resource_id" form:"resource_id" binding:"required"`
	ResourceType enum.ResourceType `json:"resource_type" form:"resource_type" binding:"required"`
	ListReq
}

type UserPermission struct {
	UserID          int64                `json:"user_id"`
	UserName        string               `json:"user_name"`
	UserDisplayName string               `json:"user_display_name"`
	Keys            []enum.PermissionKey `json:"keys"`
}

type RolePermission struct {
	RoleId   int64                `json:"role_id"`
	RoleName string               `json:"role_name"`
	Keys     []enum.PermissionKey `json:"keys"`
}

type UserResourcePermissionRes struct {
	Total          int64            `json:"total"`
	UserPermission []UserPermission `json:"user_permission"`
}

type RoleResourcePermissionRes struct {
	Total          int64            `json:"total"`
	RolePermission []RolePermission `json:"role_permission"`
}

type ResourcePermissionKeys struct {
	CanSetAllUser bool                 `json:"can_set_alluser"`
	Keys          []enum.PermissionKey `json:"keys"`
}

type UserPermissionKey struct {
	UserID          int64              `json:"user_id"`
	UserDisplayName string             `json:"user_display_name"`
	PermissionName  string             `json:"permission_name"`
	PermissionKey   enum.PermissionKey `json:"permission_key"`
}

type RolePermissionKey struct {
	RoleID         int64              `json:"role_id"`
	RoleName       string             `json:"role_name"`
	PermissionName string             `json:"permission_name"`
	PermissionKey  enum.PermissionKey `json:"permission_key"`
}

// ResourcePermissionPutReq 请求接口参数定义
type ResourcePermissionPutReq struct {
	ResourceID    int64              `json:"resource_id" form:"resource_id" binding:"required"`
	ResourceType  enum.ResourceType  `json:"resource_type" form:"resource_type" binding:"required"`
	PermissionKey enum.PermissionKey `json:"permission_key"`
	Enable        bool               `json:"enable"`
}

// ResourcePermissionPutUserReq 请求接口参数定义
type ResourcePermissionPutUserReq struct {
	UserID        int64              `json:"user_id"`
	ResourceID    int64              `json:"resource_id" form:"resource_id" binding:"required"`
	ResourceType  enum.ResourceType  `json:"resource_type" form:"resource_type" binding:"required"`
	PermissionKey enum.PermissionKey `json:"permission_key"`
	Enable        bool               `json:"enable"`
}

// ResourcePermissionPutRoleReq 请求接口参数定义
type ResourcePermissionPutRoleReq struct {
	RoleID        int64              `json:"role_id"`
	ResourceID    int64              `json:"resource_id" form:"resource_id" binding:"required"`
	ResourceType  enum.ResourceType  `json:"resource_type" form:"resource_type" binding:"required"`
	PermissionKey enum.PermissionKey `json:"permission_key"`
	Enable        bool               `json:"enable"`
}

type UserOwnResourcePermissionReq struct {
	UserID       int64             `json:"user_id" form:"user_id" binding:"required"`
	ResourceType enum.ResourceType `json:"resource_type" form:"resource_type" binding:"required"`
	ListReq
}

type RoleOwnResourcePermissionReq struct {
	RoleID       int64             `json:"role_id" form:"role_id" binding:"required"`
	ResourceType enum.ResourceType `json:"resource_type" form:"resource_type" binding:"required"`
	ListReq
}

type PermissionKeyName struct {
	Key   enum.PermissionKey `json:"key"`
	Name  string             `json:"name"`
	Order int                `json:"order"`
}

type ResourcePermissionKey struct {
	ResourceID   int64               `json:"resource_id"`
	ResourceName string              `json:"resource_name"`
	UserOwnKeys  []PermissionKeyName `json:"user_own_keys"`
	AssignKeys   []PermissionKeyName `json:"assigned_own_keys"`
}

type AssignResourcePermissionRes struct {
	Total              int64                   `json:"total"`
	ResourcePermission []ResourcePermissionKey `json:"resource_permission"`
}

type PostUserOwnResourcePermissionReq struct {
	UserID       int64              `json:"user_id" form:"user_id" binding:"required"`
	ResourceType enum.ResourceType  `json:"resource_type" form:"resource_type" binding:"required"`
	ResourceID   int64              `json:"resource_id"`
	AssignKey    enum.PermissionKey `json:"assigned_key"`
	Enable       bool               `json:"enable"`
}

type PostRoleOwnResourcePermissionReq struct {
	RoleID       int64              `json:"role_id" form:"role_id" binding:"required"`
	ResourceType enum.ResourceType  `json:"resource_type" form:"resource_type" binding:"required"`
	ResourceID   int64              `json:"resource_id"`
	AssignKey    enum.PermissionKey `json:"assigned_key"`
	Enable       bool               `json:"enable"`
}
