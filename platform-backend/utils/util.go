package utils

import (
	"errors"
	"platform-backend/dto"
	"platform-backend/dto/enum"
	"platform-backend/models"

	"github.com/gin-gonic/gin"
)

func BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		if err = c.ShouldBindQuery(obj); err != nil {

			return err
		}
	}

	return nil
}

// HandlePagination 统一处理分页参数
func HandlePagination(req *dto.ListReq) {
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}
}
func GetUserIDFromCtx(c *gin.Context) (int64, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("invalid user")
	}

	if v, ok := userID.(int64); ok {
		return int64(v), nil
	}

	return 0, errors.New("invalid user")
}

func GetUserNameFromCtx(c *gin.Context) (string, error) {
	userName, exists := c.Get("userName")
	if !exists {
		return "", errors.New("invalid user")
	}

	if v, ok := userName.(string); ok {
		return v, nil
	}

	return "", errors.New("invalid user")
}

func IsSuperAdmin(c *gin.Context) bool {
	userName, err := GetUserNameFromCtx(c)
	if err != nil {
		return false
	}

	return userName == enum.SuperAdminUserAccName
}

// HasXPermission 判断指定位置的权限位是否为1
// permission 是表示权限的 int64 值，index 是要检查的位（从右往左，1 表示最右边一位）
func HasXPermission(permission int64, index int) bool {
	index -= 1
	if index < 0 || index > 63 {
		return false
	}

	mask := int64(1 << index)
	return (permission & mask) != 0
}

// SetXPermission 将指定位置权限设置为1
// permission 是表示权限的 int64 值，index 是要设置的位（从右往左，1 表示最右边一位）
func SetXPermission(permission int64, index int) int64 {
	index -= 1
	permission |= 1 << index
	return permission
}

// SetBatchPermission 将指定多个位置权限设置为1
// permission 是表示权限的 int64 值，indexes 是要设置的位（从右往左，1 表示最右边一位）
func SetBatchPermission(permission int64, indexes []int) int64 {
	for _, index := range indexes {
		index -= 1
		permission |= 1 << index
	}

	return permission
}

// ClearXPermission 将指定位置权限设置为0
// permission 是表示权限的 int64 值，index 是要设置的位（从右往左，1 表示最右边一位）
func ClearXPermission(permission int64, index int) int64 {
	index -= 1
	mask := ^(1 << index)
	return permission & int64(mask)
}

func MakeRoleResource(original *[]models.RoleMenuPermission, roleID int64, resourceType enum.ResourceType, resource string, condition bool) {
	// if !condition {
	// 	return
	// }

	// *original = append(*original, models.RoleResource{
	// 	RoleID:       roleID,
	// 	ResourceType: resourceType,
	// 	ResourceKey:  permission.ResourceMap[resource],
	// })
}
