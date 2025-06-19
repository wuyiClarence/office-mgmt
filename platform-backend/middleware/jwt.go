package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	db "platform-backend/db"
	"platform-backend/dto/enum"
	myerrors "platform-backend/errors"
	"platform-backend/repository"
	"platform-backend/utils/format"
	myjwt "platform-backend/utils/jwt"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization Header 中获取 Token
		tokenString := c.Query("token")
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")

			if authHeader == "" {
				format.NewResponseJson(c).Error(myerrors.ErrUnauthorizedToken)
				c.Abort()
				return
			}

			// 解析 Bearer Token
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				format.NewResponseJson(c).Error(myerrors.ErrUnauthorizedToken)
				c.Abort()
				return
			}
		}

		// 验证 Token
		claims, err := myjwt.VerifyToken(tokenString)
		if err != nil {
			format.NewResponseJson(c).Error(myerrors.ErrUnauthorizedToken)
			c.Abort()
			return
		}

		// TODO 做缓存 重新登陆删除缓存
		userRepo := repository.NewUserRepository(db.MysqlDB.DB())
		user, err := userRepo.FindOne(c.Request.Context(), map[string]interface{}{
			"id":        claims.UserID,
			"user_name": claims.UserName,
		})
		if err != nil {
			format.NewResponseJson(c).Error(myerrors.ErrReLogin)
			c.Abort()
			return
		}

		if claims.IssuedAt.Before(*user.PasswordUpdatedAt) {
			format.NewResponseJson(c).Error(myerrors.ErrReLogin)
			c.Abort()
			return
		}

		if user.LogOutAt != nil && claims.IssuedAt.Before(*user.LogOutAt) {
			format.NewResponseJson(c).Error(myerrors.ErrReLogin)
			c.Abort()
			return
		}

		if user.Status == enum.UserStatusDeleted {
			format.NewResponseJson(c).Error(errors.New("账号已注销"))
			c.Abort()
			return
		}

		// 将用户信息 存储在上下文中，以供后续的处理使用
		c.Set("userID", claims.UserID)
		c.Set("userName", claims.UserName)
		c.Next()
	}
}
