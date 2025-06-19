package models

import "platform-backend/dto/enum"

type UserResourcePermission struct {
	Base
	UserID         int64              `gorm:"type:bigint;column:user_id;not null;default:0;uniqueIndex:uk_user_key_resource,priority:1" json:"user_id"`
	PermissionKey  enum.PermissionKey `gorm:"type:varchar(128);not null;column:permission_key;uniqueIndex:uk_user_key_resource,priority:3" json:"permission_key"`
	PermissionName string             `gorm:"type:varchar(128);not null;column:permission_name;" json:"permission_name"`
	ResourceType   enum.ResourceType  `gorm:"type:varchar(128);not null;column:resource_type;uniqueIndex:uk_user_key_resource,priority:4" json:"resource_type"`
	ResourceID     int64              `gorm:"type:bigint;column:resource_id;not null;uniqueIndex:uk_user_key_resource,priority:2" json:"resource_id"`
	CreateUser     int64              `gorm:"type:bigint;column:create_user;not null;default:0" json:"create_user"`
}

func (d *UserResourcePermission) TableName() string {
	return "user_resource_permission"
}
