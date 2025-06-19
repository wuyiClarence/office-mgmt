package models

type RoleUser struct {
	Base
	RoleID int64 `gorm:"type:bigint;not null;default:0;column:role_id;uniqueIndex:uk_role_user,priority:1" json:"role_id"`
	UserID int64 `gorm:"type:bigint;not null;default:0;column:user_id;uniqueIndex:uk_role_user,priority:2" json:"user_id"`
}

func (r *RoleUser) TableName() string {
	return "role_user"
}
