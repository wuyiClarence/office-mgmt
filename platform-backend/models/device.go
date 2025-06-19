package models

import (
	"platform-backend/dto/enum"
	"time"
)

type Device struct {
	Base
	UniqueId     string              `gorm:"type:varchar(128);default:'';uniqueIndex:uk_unique_id;column:unique_id;NOT NULL" json:"unique_id"`
	DeviceName   string              `gorm:"type:varchar(512);default:'';column:device_name;NOT NULL" json:"device_name"`
	AliasName    string              `gorm:"type:varchar(512);default:'';column:alias_name;NOT NULL" json:"alias_name"`
	Mac          string              `gorm:"type:varchar(128);default:'';column:mac;NOT NULL" json:"mac"`
	Ip           string              `gorm:"type:varchar(50);default:'';column:ip;NOT NULL" json:"ip"`
	OsType       string              `gorm:"type:varchar(50);default:'';column:os_type;" json:"os_type"`
	DeviceType   enum.DeviceType     `gorm:"type:tinyint;default:0;column:device_type;default:0;NOT NULL" json:"device_type"` // 1:实体服务器 2:虚拟机
	OpStatus     enum.DeviceOpStatus `gorm:"type:tinyint;default:0;column:op_status;default:0;NOT NULL" json:"op_status"`     // 1:运行中
	Status       enum.DeviceStatus   `gorm:"type:tinyint;default:0;column:status;default:0;NOT NULL" json:"status"`           // 1:运行中
	HostDeviceId int64               `gorm:"type:bigint;not null;default:0;column:host_device_id;" json:"host_device_id"`
	WakeDeviceId int64               `gorm:"type:bigint;not null;default:0;column:wake_device_id;" json:"wake_device_id"`
	OnLineTime   *time.Time          `gorm:"column:online_time;type:datetime;default:NULL"`
}

func (d *Device) TableName() string {
	return "device"
}
