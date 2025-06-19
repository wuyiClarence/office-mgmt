package models

import "platform-backend/dto/enum"

// MenuPermissionAPI 菜单API权限表
type MenuPermissionAPI struct {
	Base
	PermissionKey enum.PermissionKey `gorm:"type:varchar(128);not null;column:permission_key;uniqueIndex:uk_permission_api" json:"permission_key"`
	ApiMethod     string             `gorm:"type:varchar(16);not null;uniqueIndex:uk_permission_api" json:"api_method"`
	ApiPath       string             `gorm:"type:varchar(256);not null;uniqueIndex:uk_permission_api" json:"api_path"`
	Used          int8               `gorm:"type:tinyint;not null;default:0;column:used;" json:"used"`
}

func (d *MenuPermissionAPI) TableName() string {
	return "menu_permission_api"
}
