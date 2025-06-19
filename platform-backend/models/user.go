package models

import (
	"database/sql"
	"platform-backend/dto/enum"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Base
	UserName          string          `gorm:"type:varchar(255);not null;default:'';column:user_name;uniqueIndex:uk_user_name" json:"user_name"`
	UserDisplayName   string          `gorm:"type:varchar(255);not null;default:'';column:user_display_name" json:"user_display_name"`
	Status            enum.UserStatus `gorm:"type:tinyint;not null;default:0;column:status;comment:'1:正常状态 2:已删除'" json:"status"`
	Email             string          `gorm:"type:varchar(255);not null;default:'';column:email" json:"email"`
	PhoneNumber       string          `gorm:"type:varchar(255);not null;default:'';column:phone_number" json:"phone_number"`
	Password          sql.NullString  `gorm:"type:varchar(255);not null;default:'';column:password" json:"password"`
	Permissions       int64           `gorm:"type:bigint;not null;default:0;column:permissions" json:"permissions" comment:"按位表示权限"`
	PasswordUpdatedAt *time.Time      `gorm:"column:password_updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"password_updated_at"`
	DeletedAt         *time.Time      `gorm:"column:deleted_at;type:datetime;default:NULL" json:"deleted_at"`
	LogOutAt          *time.Time      `gorm:"column:logout_at;type:datetime;default:NULL" json:"logout_at"`
	CreateUser        int64           `gorm:"type:bigint;column:create_user;not null;default:0" json:"create_user"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) VerifyPassword(password string) bool {
	return u.Password.Valid &&
		bcrypt.CompareHashAndPassword([]byte(u.Password.String), []byte(password)) == nil
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		u.Password = sql.NullString{}
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = sql.NullString{String: string(hash), Valid: true}
	return nil
}
