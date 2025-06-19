package errors

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

// 错误定义
var (
	ErrUnauthorizedToken = New(401, "Invalid Token")
	ErrReLogin           = New(401, "请重新登陆")
	ErrBadRequest        = New(400000, "Bad Request")
	ErrInternalError     = New(500000, "Internal Server Error")
)
