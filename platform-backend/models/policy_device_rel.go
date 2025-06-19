package models

type PolicyDevice struct {
	Base
	PolicyId int64 `gorm:"type:bigint;not null;default:0;column:policy_id;uniqueIndex:uk_policy_device,priority:1" json:"policy_id"`
	DeviceId int64 `gorm:"type:bigint;not null;default:0;column:device_id;uniqueIndex:uk_policy_device,priority:2" json:"device_id"`
}

func (d *PolicyDevice) TableName() string {
	return "policy_device_rel"
}
