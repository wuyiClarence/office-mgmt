package models

type DeviceGroup struct {
	Base
	GroupName  string `gorm:"type:varchar(100);not null;default:'';column:group_name;uniqueIndex:uk_group_name" json:"group_name"`
	CreateUser int64  `gorm:"type:bigint;column:create_user;not null;default:0" json:"create_user"`
}

func (s *DeviceGroup) TableName() string {
	return "device_group"
}
