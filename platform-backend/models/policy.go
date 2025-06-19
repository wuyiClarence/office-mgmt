package models

import (
	"platform-backend/dto/enum"
	"time"
)

type Policy struct {
	Base
	PolicyName    string                 `gorm:"type:varchar(100);not null;default:'';column:policy_name" json:"policy_name"`
	ActionType    enum.ActionType        `gorm:"type:varchar(100);not null;default:'';column:action_type" json:"action_type" comment:"开机，关机"`
	AssociateType enum.AssociateType     `gorm:"type:tinyint;not null;default:0;column:associate_type" json:"associate_type" comment:"1:单个设备 2:服务器组"`
	Status        int8                   `gorm:"type:tinyint;not null;default:0;column:status" json:"status" comment:"0:未启用 1:启用中"`
	ExecuteType   enum.TimeConditionType `gorm:"type:int;not null;default:0;column:execute_type" json:"execute_type" comment:"仅一次， 每天， 星期几"`
	StartDate     *time.Time             `gorm:"type:DATETIME;default:null;column:start_date" json:"start_date"`
	EndDate       *time.Time             `gorm:"type:DATETIME;default:null;column:end_date" json:"end_date"`
	ExecuteTime   string                 `gorm:"type:varchar(50);not null;default:'';column:execute_time" json:"execute_time" comment:"策略执行时间"`
	DayOfWeek     int                    `gorm:"type:int;not null;default:0;column:day_of_week" json:"day_of_week" comment:"周一到周日，1 bit代表星期1， 7 bit代表星期天"`
	CreateUser    int64                  `gorm:"type:bigint;column:create_user;not null;default:0" json:"create_user"`
}

func (d *Policy) TableName() string {
	return "policy"
}
