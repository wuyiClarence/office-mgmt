package dto

type ListReq struct {
	PageIndex int `form:"pageIndex" json:"pageIndex" binding:"omitempty"`
	PageSize  int `form:"pageSize" json:"pageSize" binding:"omitempty"`
}
