package controller

import (
	"platform-backend/dto"
	"platform-backend/utils"
	"platform-backend/utils/format"

	"github.com/gin-gonic/gin"
)

// PolicyCreate godoc
// @Summary PolicyCreate
// @Description PolicyCreate Api
// @Tags Policy
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param PolicyCreate body dto.PolicyCreateReq true "Create Policy request body"
// @Success 200 {object} nil
// @Router /api/policy/create [post]
func (ctrl *V1Controller) PolicyCreate(c *gin.Context) {
	req := &dto.PolicyCreateReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.PolicyCreate(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// PolicyUpdate godoc
// @Summary PolicyUpdate
// @Description PolicyUpdate Api
// @Tags Policy
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param PolicyUpdate body dto.PolicyUpdateReq true "Update Policy request body"
// @Success 200 {object} nil
// @Router /api/policy/update [post]
func (ctrl *V1Controller) PolicyUpdate(c *gin.Context) {
	req := &dto.PolicyUpdateReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.PolicyUpdate(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}

// PolicyList godoc
// @Summary PolicyList
// @Description PolicyList Api
// @Tags Policy
// @Accept  json
// @Produce  json
// @Param pageIndex query int false "页索引"
// @Param pageSize query int false "页大小"
// @Success 200 {object} dto.PolicyListRes
// @Router /api/policy/list [get]
func (ctrl *V1Controller) PolicyList(c *gin.Context) {
	req := &dto.ListReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	utils.HandlePagination(req)

	res, err := ctrl.svc.PolicyList(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(res)
}

// PolicyDelete godoc
// @Summary PolicyDelete
// @Description PolicyDelete Api
// @Tags Policy
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param PolicyDelete body dto.PolicyDelReq true "Del Policy request body"
// @Success 200 {object} nil
// @Router /api/policy/delete [delete]
func (ctrl *V1Controller) PolicyDelete(c *gin.Context) {
	req := &dto.PolicyDelReq{}
	if err := utils.BindAndValidate(c, &req); err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	err := ctrl.svc.PolicyDelete(c, req)
	if err != nil {
		format.NewResponseJson(c).Error(err)
		return
	}

	format.NewResponseJson(c).Success(nil)
}
