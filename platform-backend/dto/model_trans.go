package dto

import (
	"time"
)

type UserWithRoleTrans struct {
	ID                int64      `json:"id"`
	UserName          string     `json:"user_name"`
	UserDisplayName   string     `json:"user_display_name"`
	Email             string     `json:"email"`
	PhoneNumber       string     `json:"phone_number"`
	Password          string     `json:"password"`
	PasswordUpdatedAt *time.Time `json:"password_updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}
