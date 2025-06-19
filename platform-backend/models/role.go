package models

type Role struct {
	Base
	RoleName     string `gorm:"type:varchar(128);not null;default:'';column:role_name;uniqueIndex:uk_role_name" json:"role_name"`
	Description  string `gorm:"type:varchar(255);not null;default:'';column:description" json:"description"`
	CreateUser   int64  `gorm:"type:bigint;column:create_user;not null;default:0" json:"create_user"`
	SystemCreate int8   `gorm:"type:tinyint;not null;default:0;column:system_create;" json:"system_create"`
}

func (r *Role) TableName() string {
	return "role"
}
