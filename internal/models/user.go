package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Name     string    `json:"name"`
	Base
}

type UserLogin struct {
	Email    string `json:"email" binding:"required""`   //Email
	Password string `json:"password" binding:"required"` //Mật khẩu
}

type UserWithToken struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  *User     `json:"user"`
}

func (User) TableName() string {
	return "users"
}
