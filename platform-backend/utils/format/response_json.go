package format

import (
	"errors"
	"fmt"
	"net/http"
	myerrors "platform-backend/errors"

	"github.com/gin-gonic/gin"
)

func NewResponseJson(ctx *gin.Context) *responseJson {
	return &responseJson{
		context: ctx,
	}
}

type responseJson struct {
	context *gin.Context
}

func (r *responseJson) SetHeader(key, value string) *responseJson {
	r.context.Writer.Header().Set(key, value)
	return r
}

func (r *responseJson) Success(data interface{}) {
	r.context.JSON(http.StatusOK, ResultData{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func (r *responseJson) Error(err error, params ...interface{}) {
	var appErr *myerrors.AppError
	var ae *myerrors.AppError
	if errors.As(err, &ae) {
		appErr = ae
	} else {
		appErr = myerrors.New(500, err.Error())
	}

	msg := appErr.Message
	if len(params) > 0 {
		msg = fmt.Sprintf(appErr.Message, params...)
	}

	result := ResultData{
		Code: appErr.Code,
		Msg:  msg,
		Data: "",
	}
	r.context.Abort()
	r.context.JSON(http.StatusOK, result)
}

func (r *responseJson) ErrorWithHttpCode(statusCode, errorCode int, params ...interface{}) {
	msg := ""
	for _, param := range params {
		msg += fmt.Sprintf("%v", param)
	}

	result := ResultData{
		Code: errorCode,
		Msg:  msg,
		Data: "",
	}
	r.context.Abort()
	r.context.JSON(statusCode, result)
}

func (r *responseJson) Download(filename, path string) {
	r.SetHeader("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	r.SetHeader("Content-Type", "application/octet-stream")
	r.context.File(path)
}

type ResultData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
