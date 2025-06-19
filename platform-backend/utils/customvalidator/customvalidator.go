package customvalidator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func CheckTime(fl validator.FieldLevel) bool {
	if timeStr, ok := fl.Field().Interface().(string); ok {
		_, err := time.Parse("15:04:05", timeStr)
		if err != nil {
			return false
		}
	}
	return true
}
func CheckDate(fl validator.FieldLevel) bool {
	if timeStr, ok := fl.Field().Interface().(string); ok {
		_, err := time.Parse("2006-01-02", timeStr)
		if err != nil {
			return false
		}
	}
	return true
}
