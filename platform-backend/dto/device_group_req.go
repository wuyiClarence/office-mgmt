package dto

type DeviceGroupCreateReq struct {
	DeviceGroupName string  `json:"device_group_name" binding:"required"`
	DeviceIDS       []int64 `json:"device_ids"`
}

type DeviceGroupUpdateReq struct {
	ID        int64   `json:"id" binding:"required"`
	DeviceIDS []int64 `json:"device_ids"`
}

type DeviceGroupDelReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}

type DeviceGroupListRes struct {
	List  []*DeviceGroupDetail `json:"list"`
	Total int64                `json:"total"`
}

type DeviceGroupDetail struct {
	ID              int64           `json:"id"`
	DeviceGroupName string          `json:"device_group_name"`
	Devices         []*DeviceDetail `json:"devices"`
	Permissions     []Permission    `json:"permissions" gorm:"-"`
}

type DeviceGroupDetailReq struct {
	ID string `json:"id" binding:"required"`
}

type DeviceGroupPowerOffReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}

type DeviceGroupPowerOnReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}
