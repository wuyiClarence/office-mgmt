package models

type DeviceDeviceGroupRel struct {
	Base
	DeviceId int64 `gorm:"type:bigint;not null;default:0;column:device_id;uniqueIndex:uk_group_device,priority:2" json:"device_id"`
	GroupId  int64 `gorm:"type:bigint;not null;default:0;column:group_id;uniqueIndex:uk_group_device,priority:1" json:"group_id"`
}

func (d *DeviceDeviceGroupRel) TableName() string {
	return "device_device_group_rel"
}
