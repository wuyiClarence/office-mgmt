package permission

import "platform-backend/dto/enum"

/*
	这里是api resource的枚举
*/

// API 定义每个 API 的结构
type API struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// Permission 定义资源的结构
type Permission struct {
	Name      string             `json:"name"`
	Key       enum.PermissionKey `json:"key"`
	RoleAdmin bool               `json:"roleAdmin,omitempty"` // 使用 `omitempty` 忽略空值字段
	RoleUser  bool               `json:"roleUser,omitempty"`  // 使用 `omitempty` 忽略空值字段
	RoleAll   bool               `json:"roleAll,omitempty"`   // 使用 `omitempty` 忽略空值字段
	Apis      []API              `json:"apis"`
	Children  []Permission       `json:"children,omitempty"` // 子资源，可能为空
}

// PermissionData 整体数据结构
type PermissionData struct {
	PermissionList []Permission `json:"permissionList"`
}

type PermissionKeyName struct {
	Name  string             `json:"name"`
	Key   enum.PermissionKey `json:"key"`
	Order int                `json:"order"`
}

var ResourceTypeKeyMap = map[enum.ResourceType]map[enum.PermissionKey]PermissionKeyName{
	enum.ResourceTypeRole: {
		enum.PermissionKeyView:           {Name: "查看", Key: enum.PermissionKeyView, Order: 1},
		enum.PermissionKeyEdit:           {Name: "编辑", Key: enum.PermissionKeyEdit, Order: 2},
		enum.PermissionKeyDelete:         {Name: "删除", Key: enum.PermissionKeyDelete, Order: 3},
		enum.PermissionKeyPermissionMgmt: {Name: "权限管理", Key: enum.PermissionKeyPermissionMgmt, Order: 4},
	},
	enum.ResourceTypeDevice: {
		enum.PermissionKeyView:           {Name: "查看", Key: enum.PermissionKeyView, Order: 1},
		enum.PermissionKeyEdit:           {Name: "编辑", Key: enum.PermissionKeyEdit, Order: 2},
		enum.PermissionKeyDelete:         {Name: "删除", Key: enum.PermissionKeyDelete, Order: 3},
		enum.PermissionKeyPowerOn:        {Name: "开机", Key: enum.PermissionKeyPowerOn, Order: 4},
		enum.PermissionKeyPowerOff:       {Name: "关机", Key: enum.PermissionKeyPowerOff, Order: 5},
		enum.PermissionKeyPermissionMgmt: {Name: "权限管理", Key: enum.PermissionKeyPermissionMgmt, Order: 6},
	},
	enum.ResourceTypeDeviceGroup: {
		enum.PermissionKeyView:           {Name: "查看", Key: enum.PermissionKeyView, Order: 1},
		enum.PermissionKeyEdit:           {Name: "编辑", Key: enum.PermissionKeyEdit, Order: 2},
		enum.PermissionKeyDelete:         {Name: "删除", Key: enum.PermissionKeyDelete, Order: 3},
		enum.PermissionKeyPowerOn:        {Name: "开机", Key: enum.PermissionKeyPowerOn, Order: 4},
		enum.PermissionKeyPowerOff:       {Name: "关机", Key: enum.PermissionKeyPowerOff, Order: 5},
		enum.PermissionKeyPermissionMgmt: {Name: "权限管理", Key: enum.PermissionKeyPermissionMgmt, Order: 6},
	},
	enum.ResourceTypePolicy: {
		enum.PermissionKeyView:           {Name: "查看", Key: enum.PermissionKeyView, Order: 1},
		enum.PermissionKeyEdit:           {Name: "编辑", Key: enum.PermissionKeyEdit, Order: 2},
		enum.PermissionKeyDelete:         {Name: "删除", Key: enum.PermissionKeyDelete, Order: 3},
		enum.PermissionKeyPermissionMgmt: {Name: "权限管理", Key: enum.PermissionKeyPermissionMgmt, Order: 4},
	},
	enum.ResourceTypeUser: {
		enum.PermissionKeyView:           {Name: "查看", Key: enum.PermissionKeyView, Order: 1},
		enum.PermissionKeyEdit:           {Name: "编辑", Key: enum.PermissionKeyEdit, Order: 2},
		enum.PermissionKeyDelete:         {Name: "删除", Key: enum.PermissionKeyDelete, Order: 3},
		enum.PermissionKeyPermissionMgmt: {Name: "权限管理", Key: enum.PermissionKeyPermissionMgmt, Order: 4},
	},
}
