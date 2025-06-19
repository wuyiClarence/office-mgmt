package service

import (
	"errors"
	"fmt"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/dto/permission"
	sp "platform-backend/service/permission"
	"platform-backend/utils"

	"github.com/gin-gonic/gin"
)

func (svc *V1ServiceImpl) GetResourceAllUserPermissionInfo(c *gin.Context, req *dto.ResourcePermissionReq) (*dto.ResourcePermissionKeys, error) {

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, errors.New("无效用户")
	}

	res := &dto.ResourcePermissionKeys{}

	_, ok := permission.ResourceTypeKeyMap[req.ResourceType]
	if !ok {
		return nil, fmt.Errorf("请求的资源类型错误%s", req.ResourceType)
	}

	permissions, canSetAllUser, err := sp.GetAllUserPermission(c.Request.Context(), req.ResourceType, req.ResourceID, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}
	res.CanSetAllUser = canSetAllUser

	res.Keys = permissions

	return res, nil
}

func (svc *V1ServiceImpl) GetResourceUserPermissionInfo(c *gin.Context, req *dto.ResourcePermissionReq) (*dto.UserResourcePermissionRes, error) {

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, errors.New("无效用户")
	}

	res, err := sp.GetResourceUserPermissionInfo(c.Request.Context(), req.ResourceType, req.ResourceID, req.PageIndex, req.PageSize, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (svc *V1ServiceImpl) GetResourceRolePermissionInfo(c *gin.Context, req *dto.ResourcePermissionReq) (*dto.RoleResourcePermissionRes, error) {

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, errors.New("无效用户")
	}

	res, err := sp.GetResourceRolePermissionInfo(c.Request.Context(), req.ResourceType, req.ResourceID, req.PageIndex, req.PageSize, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *V1ServiceImpl) PutResourceAllUserPermission(c *gin.Context, req *dto.ResourcePermissionPutReq) error {

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return errors.New("无效用户")
	}

	err = sp.CheckPermission(c, req.ResourceType, req.ResourceID, enum.PermissionKeyPermissionMgmt, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	err = sp.CheckResourceTypeAndKey(req.ResourceType, req.PermissionKey)
	if err != nil {
		return err
	}
	err = sp.PutUserPermission(c.Request.Context(), req.ResourceType, req.ResourceID, req.PermissionKey, req.Enable, 0, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	return nil
}

func (svc *V1ServiceImpl) PutResourceUserPermission(c *gin.Context, req *dto.ResourcePermissionPutUserReq) error {

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return errors.New("无效用户")
	}

	err = sp.CheckPermission(c, req.ResourceType, req.ResourceID, enum.PermissionKeyPermissionMgmt, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	err = sp.PutUserPermission(c.Request.Context(), req.ResourceType, req.ResourceID, req.PermissionKey, req.Enable, req.UserID, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	return nil
}

func (svc *V1ServiceImpl) PutResourceRolePermission(c *gin.Context, req *dto.ResourcePermissionPutRoleReq) error {

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return errors.New("无效用户")
	}

	err = sp.CheckPermission(c, req.ResourceType, req.ResourceID, enum.PermissionKeyPermissionMgmt, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	err = sp.PutRolePermission(c.Request.Context(), req.ResourceType, req.ResourceID, req.PermissionKey, req.Enable, req.RoleID, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	return nil
}

func (svc *V1ServiceImpl) GetRoleOwnResourcePermission(c *gin.Context, req *dto.RoleOwnResourcePermissionReq) (*dto.AssignResourcePermissionRes, error) {
	_, ok := permission.ResourceTypeKeyMap[req.ResourceType]
	if !ok {
		return nil, fmt.Errorf("请求的资源类型错误%s", req.ResourceType)
	}

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, errors.New("无效用户")
	}
	result, err := sp.GetAssignRoleResourcePermissionInfo(c.Request.Context(), req.ResourceType, req.RoleID, req.PageIndex, req.PageSize, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (svc *V1ServiceImpl) GetUserOwnResourcePermission(c *gin.Context, req *dto.UserOwnResourcePermissionReq) (*dto.AssignResourcePermissionRes, error) {
	_, ok := permission.ResourceTypeKeyMap[req.ResourceType]
	if !ok {
		return nil, fmt.Errorf("请求的资源类型错误%s", req.ResourceType)
	}

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return nil, errors.New("无效用户")
	}
	result, err := sp.GetAssignUserResourcePermissionInfo(c.Request.Context(), req.ResourceType, req.UserID, req.PageIndex, req.PageSize, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (svc *V1ServiceImpl) PostRoleOwnResourcePermission(c *gin.Context, req *dto.PostRoleOwnResourcePermissionReq) error {
	_, ok := permission.ResourceTypeKeyMap[req.ResourceType]
	if !ok {
		return fmt.Errorf("请求的资源类型错误%s", req.ResourceType)
	}

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return errors.New("无效用户")
	}

	err = sp.CheckPermission(c, req.ResourceType, req.ResourceID, enum.PermissionKeyPermissionMgmt, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	err = sp.AssignRoleResourcePermission(c.Request.Context(), req.ResourceType, req.ResourceID, req.AssignKey, req.RoleID, curUserID, utils.IsSuperAdmin(c), req.Enable)
	if err != nil {
		return err
	}
	return nil
}
func (svc *V1ServiceImpl) PostUserOwnResourcePermission(c *gin.Context, req *dto.PostUserOwnResourcePermissionReq) error {
	_, ok := permission.ResourceTypeKeyMap[req.ResourceType]
	if !ok {
		return fmt.Errorf("请求的资源类型错误%s", req.ResourceType)
	}

	curUserID, err := utils.GetUserIDFromCtx(c)
	if err != nil {
		return errors.New("无效用户")
	}

	err = sp.CheckPermission(c, req.ResourceType, req.ResourceID, enum.PermissionKeyPermissionMgmt, curUserID, utils.IsSuperAdmin(c))
	if err != nil {
		return err
	}

	err = sp.AssignUserResourcePermission(c.Request.Context(), req.ResourceType, req.ResourceID, req.AssignKey, req.UserID, curUserID, utils.IsSuperAdmin(c), req.Enable)
	if err != nil {
		return err
	}
	return nil
}
