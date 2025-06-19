package dto

import (
	"platform-backend/dto/enum"
	"time"
)

type PolicyBaseInfo struct {
	PolicyName    string                 `json:"policy_name" binding:"required"`
	Status        int8                   `json:"status"`
	ActionType    enum.ActionType        `json:"action_type" binding:"required"`
	AssociateType enum.AssociateType     `json:"associate_type" binding:"required"`
	ExecuteTime   string                 `json:"execute_time" binding:"required,checkTime"`
	ExecuteType   enum.TimeConditionType `json:"execute_type" binding:"required"`
	DayOfWeek     int                    `json:"day_of_week"` // 周一到周日，1 bit代表星期1， 7 bit代表星期天
	StartDate     *time.Time             `json:"start_date" binding:"omitempty"`
	EndDate       *time.Time             `json:"end_date" binding:"omitempty"`
}

type PolicyCreateReq struct {
	PolicyBaseInfo
	DeviceIDS      []int64 `json:"device_ids" binding:"required"`
	DeviceGroupIDS []int64 `json:"device_group_ids" binding:"required"`
}

type PolicyUpdateReq struct {
	ID int64 `json:"id" binding:"required"`
	PolicyBaseInfo
	DeviceIDS      []int64 `json:"device_ids" binding:"required"`
	DeviceGroupIDS []int64 `json:"device_group_ids" binding:"required"`
}

type PolicyList struct {
	ID int64 `json:"id" binding:"required"`
	PolicyBaseInfo
	Devices      []*DeviceDetail      `json:"devices"`
	DeviceGroups []*DeviceGroupDetail `json:"device_groups"`
	Permissions  []Permission         `json:"permissions" gorm:"-"`
}

type PolicyListRes struct {
	List  []*PolicyList `json:"list"`
	Total int64         `json:"total"`
}

type PolicyDelReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}
