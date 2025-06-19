package controller

import (
	"github.com/gin-gonic/gin"

	"platform-backend/dto"
	"platform-backend/utils"
	"platform-backend/utils/format"
)

// DeviceGroupCreate godoc
// @Summary DeviceGroupCreate
// @Description DeviceGroupCreate Api
// @Tags Device-Group
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param DeviceGroupCreate body dto.DeviceGroupCreateReq true "Create DeviceGroup request body"
// @Success 200 {object} nil
// @Router /api/device_group/create [put]
func (ctrl *V1Controller) DeviceGroupCreate(c *gin.Context) {
	req := &dto.DeviceGroupCreateReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.DeviceGroupCreate(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// DeviceGroupUpdate godoc
// @Summary DeviceGroupUpdate
// @Description DeviceGroupUpdate Api
// @Tags Device-Group
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param DeviceGroupUpdate body dto.DeviceGroupUpdateReq true "Update DeviceGroup request body"
// @Success 200 {object} nil
// @Router /api/device_group/update [post]
func (ctrl *V1Controller) DeviceGroupUpdate(c *gin.Context) {
	req := &dto.DeviceGroupUpdateReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.DeviceGroupUpdate(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// DeviceGroupDelete godoc
// @Summary DeviceGroupDelete
// @Description DeviceGroupDelete Api
// @Tags Device-Group
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param DeviceGroupDelete body dto.DeviceGroupDelReq true "Del DeviceGroup request body"
// @Success 200 {object} nil
// @Router /api/device_group/delete [delete]
func (ctrl *V1Controller) DeviceGroupDelete(c *gin.Context) {
	req := &dto.DeviceGroupDelReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.DeviceGroupDelete(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// DeviceGroupList godoc
// @Summary DeviceGroupList
// @Description DeviceGroupList Api
// @Tags Device-Group
// @Accept  json
// @Produce  json
// @Param pageIndex query int false "页索引"
// @Param pageSize query int false "页大小"
// @Success 200 {object} dto.DeviceGroupListRes
// @Router /api/device_group/list [get]
func (ctrl *V1Controller) DeviceGroupList(c *gin.Context) {
	req := &dto.ListReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	utils.HandlePagination(req)

	res, err := ctrl.svc.DeviceGroupList(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// DeviceGroupDetail godoc
// @Summary DeviceGroupDetail
// @Description DeviceGroupDetail Api
// @Tags Device-Group
// @Accept  json
// @Produce  json
// @Param device_group_id query int false "id"
// @Success 200 {object} dto.DeviceGroupDetail
// @Router /api/device_group/detail [get]
func (ctrl *V1Controller) DeviceGroupDetail(c *gin.Context) {
	req := &dto.DeviceGroupDetailReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	res, err := ctrl.svc.DeviceGroupDetail(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// DeviceGroupPowerOn godoc
// @Summary DeviceGroupPowerOn
// @Description DeviceGroupPowerOn Api
// @Tags Device-Group
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param DeviceGroupPowerOn body dto.DeviceGroupPowerOnReq true "DeviceGroup PowerOn request body"
// @Success 200 {object} nil
// @Router /api/device_group/poweron [post]
func (ctrl *V1Controller) DeviceGroupPowerOn(c *gin.Context) {
	req := &dto.DeviceGroupPowerOnReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.DeviceGroupPowerOn(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// DeviceGroupPowerOff godoc
// @Summary DeviceGroupPowerOff
// @Description DeviceGroupPowerOff Api
// @Tags Device-Group
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param DeviceGroupPowerOff body dto.DeviceGroupPowerOffReq true "DeviceGroup PowerOff request body"
// @Success 200 {object} nil
// @Router /api/device_group/poweron [post]
func (ctrl *V1Controller) DeviceGroupPowerOff(c *gin.Context) {
	req := &dto.DeviceGroupPowerOffReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.DeviceGroupPowerOff(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}
