package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	db "platform-backend/db"
	"platform-backend/repository"
	"platform-backend/utils"
)

func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.IsSuperAdmin(c) {
			c.Next()
		} else {
			urlPath := c.Request.URL.Path
			urlMethod := c.Request.Method
			curUserID, err := utils.GetUserIDFromCtx(c)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err})
				c.Abort()
				return
			}

			ok, err := repository.NewRoleMenuPermissionRepository(db.MysqlDB.DB()).FindUserRolePermission(c.Request.Context(), curUserID, urlPath, urlMethod)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err})
				c.Abort()
				return
			}
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("无权访问")})
				c.Abort()
				return
			}

			c.Next()
		}
	}
}
