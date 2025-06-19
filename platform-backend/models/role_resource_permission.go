package models

import "platform-backend/dto/enum"

type RoleResourcePermission struct {
	Base
	RoleID         int64              `gorm:"type:bigint;column:role_id;not null;default:0;uniqueIndex:uk_role_key_resource,priority:1" json:"role_id"`
	PermissionKey  enum.PermissionKey `gorm:"type:varchar(128);not null;column:permission_key;uniqueIndex:uk_role_key_resource,priority:3" json:"permission_key"`
	PermissionName string             `gorm:"type:varchar(128);not null;column:permission_name;" json:"permission_name"`
	ResourceType   enum.ResourceType  `gorm:"type:varchar(128);not null;column:resource_type;uniqueIndex:uk_role_key_resource,priority:4" json:"resource_type"`
	ResourceID     int64              `gorm:"type:bigint;column:resource_id;not null;uniqueIndex:uk_role_key_resource,priority:2" json:"resource_id"`
	CreateUser     int64              `gorm:"type:bigint;column:create_user;not null;default:0" json:"create_user"`
}

func (d *RoleResourcePermission) TableName() string {
	return "role_resource_permission"
}
