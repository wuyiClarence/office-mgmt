package dto

import "platform-backend/dto/enum"

type DeviceListRes struct {
	Total int64           `json:"total"`
	List  []*DeviceDetail `json:"list"`
}
type DeviceGroupInfo struct {
	ID              int64  `json:"id"`
	DeviceGroupName string `json:"device_group_name"`
}

type DeviceDetail struct {
	ID              int64             `json:"id"`
	DeviceName      string            `json:"device_name"`
	AliasName       string            `json:"alias_name"`
	Mac             string            `json:"mac"`
	Ip              string            `json:"ip"`
	OsType          string            `json:"os_type"`
	DeviceType      enum.DeviceType   `json:"device_type"`
	Status          enum.DeviceStatus `json:"status" enums:"0,1"`
	DeviceGroupInfo []DeviceGroupInfo `json:"device_group_info"`
	Permissions     []Permission      `json:"permissions" gorm:"-"`
}
type DeviceDelReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}

type DevicePowerOnReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}
type DevicePowerOffReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}

type DeviceUpdateReq struct {
	ID        int64  `json:"id" binding:"required"`
	AliasName string `json:"alias_name"`
	Mac       string `json:"mac"`
	Ip        string `json:"ip"`
}
