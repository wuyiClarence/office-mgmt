package models

import "platform-backend/dto/enum"

type RoleMenuPermission struct {
	Base
	RoleID        int64              `gorm:"type:bigint;column:role_id;not null;default:0;index:uk_role_permission" json:"role_id"`
	PermissionKey enum.PermissionKey `gorm:"type:varchar(128);not null;column:permission_key;uniqueIndex:uk_role_permission" json:"permission_key"`
	CreateUser    int64              `gorm:"type:bigint;column:create_user;not null;default:0" json:"create_user"`
}

func (d *RoleMenuPermission) TableName() string {
	return "role_menu_permission"
}
