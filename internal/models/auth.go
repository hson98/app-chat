package models

type RefreshTokenBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
