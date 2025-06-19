package controller

import (
	"platform-backend/dto"
	"platform-backend/utils"
	"platform-backend/utils/format"

	"github.com/gin-gonic/gin"
)

// DeviceList godoc
// @Summary DeviceList
// @Description DeviceList Api 查询用户拥有的设备
// @Tags Device
// @Accept  json
// @Produce  json
// @Param pageIndex query int false "页索引"
// @Param pageSize query int false "页大小"
// @Success 200 {object} dto.DeviceListRes
// @Router /api/device/list [get]
func (ctrl *V1Controller) DeviceList(c *gin.Context) {

	req := &dto.ListReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	utils.HandlePagination(req)

	res, err := ctrl.svc.DeviceList(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// DeviceDelete godoc
// @Summary DeviceDelete
// @Description DeviceDelete Api
// @Tags Device
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param DeviceDelete body dto.PolicyDelReq true "Del Policy request body"
// @Success 200 {object} nil
// @Router /api/device/delete [delete]
func (ctrl *V1Controller) DeviceDelete(c *gin.Context) {
	req := &dto.DeviceDelReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.DeviceDelete(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// DevicePowerOn godoc
// @Summary DevicePowerOn
// @Description DevicePowerOn Api
// @Tags Device
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param DevicePowerOn body dto.DevicePowerOnReq true "Del Policy request body"
// @Success 200 {object} nil
// @Router /api/device/poweron [post]
func (ctrl *V1Controller) DevicePowerOn(c *gin.Context) {
	req := &dto.DevicePowerOnReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.DevicePowerOn(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// DevicePowerOff godoc
// @Summary DevicePowerOff
// @Description DevicePowerOff Api
// @Tags Device
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param DevicePowerOff body dto.DevicePowerOffReq true "Del Policy request body"
// @Success 200 {object} nil
// @Router /api/device/poweron [post]
func (ctrl *V1Controller) DevicePowerOff(c *gin.Context) {
	req := &dto.DevicePowerOffReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.DevicePowerOff(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// DeviceUpdate godoc
// @Summary DeviceUpdate
// @Description DeviceUpdate Api
// @Tags Device
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param DeviceUpdate body dto.DeviceUpdateReq true "Update Device request body"
// @Success 200 {object} nil
// @Router /api/device/update [put]
func (ctrl *V1Controller) DeviceUpdate(c *gin.Context) {
	req := &dto.DeviceUpdateReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.DeviceUpdate(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}
