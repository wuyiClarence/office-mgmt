package middleware

import (
	"fmt"

	"runtime/debug"

	"github.com/gin-gonic/gin"

	myerrors "platform-backend/errors"
	"platform-backend/utils/format"
	"platform-backend/utils/log"
)

func Recover() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.SystemLog(fmt.Sprintf("%s", err))
				if gin.IsDebugging() {
					debug.PrintStack()
				}
				format.NewResponseJson(context).Error(myerrors.ErrInternalError)
			}
		}()
		context.Next()
	}
}
