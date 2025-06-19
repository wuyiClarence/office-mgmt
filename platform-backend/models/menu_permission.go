package models

import "platform-backend/dto/enum"

// MenuPermission 菜单权限表
type MenuPermission struct {
	Base
	PermissionName string             `gorm:"type:varchar(128);not null;column:permission_name;" json:"permission_name"`
	PermissionKey  enum.PermissionKey `gorm:"type:varchar(128);not null;column:permission_key;uniqueIndex:uk_permission_key" json:"permission_key"`
	Used           int8               `gorm:"type:tinyint;not null;default:0;column:used;" json:"used"`
	ParentID       int64              `gorm:"type:bigint;not null;default:0;column:parent_id;uniqueIndex:uk_permission_type_key" json:"parent_id"`
	Order          int                `gorm:"type:int unsigned;not null;default:0;column:order;" json:"order"`
}

func (d *MenuPermission) TableName() string {
	return "menu_permission"
}
