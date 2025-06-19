package models

type UserRole struct {
	Base
	UserID int64 `gorm:"type:bigint;not null;default:0;column:user_id;uniqueIndex:uk_user_role,priority:1" json:"user_id"`
	RoleID int64 `gorm:"type:bigint;not null;default:0;column:role_id;uniqueIndex:uk_user_role,priority:2" json:"role_id" comment:"用户身上绑定的角色id"`
}

func (u *UserRole) TableName() string {
	return "user_role"
}
