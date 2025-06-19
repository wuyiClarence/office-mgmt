package models

type PolicyDeviceGroup struct {
	Base
	PolicyId      int64 `gorm:"type:bigint;not null;default:0;column:policy_id;uniqueIndex:uk_policy_device,priority:1" json:"policy_id"`
	DeviceGroupId int64 `gorm:"type:bigint;not null;default:0;column:device_group_id;uniqueIndex:uk_policy_device,priority:2" json:"device_group_id"`
}

func (d *PolicyDeviceGroup) TableName() string {
	return "policy_device_group_rel"
}
