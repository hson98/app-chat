package httperrs

import (
	"errors"
)

const (
	ErrUsernameOrPasswordInvalid = "email or password is incorrect"
	PasswordInvalid              = "invalid password"
	ErrEmailExisted              = "email already exists"
	CanNotSaveToStorage          = "an error occurred while storing data"
	PassAndConfirmPassNotMatch   = "password does not match the confirmed password"
)

var (
	InternalServerError = errors.New("Internal Server Error")
	Unauthorized        = errors.New("Unauthorized")
	InvalidJWTToken     = errors.New("Invalid JWT token")
	InvalidJWTClaims    = errors.New("Invalid JWT claims")
)
