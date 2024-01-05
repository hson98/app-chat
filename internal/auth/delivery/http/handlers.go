package auth_http

import (
	"github.com/gin-gonic/gin"
	auth_usecase "hson98/app-chat/internal/auth/usecase"
)

type Handler interface {
	Login() gin.HandlerFunc
	Register() gin.HandlerFunc
	Logout() gin.HandlerFunc
}

type authHandler struct {
	authUC auth_usecase.UseCase
}

func (a authHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (a authHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (a authHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func NewAuthHandler(authUC auth_usecase.UseCase) Handler {
	return &authHandler{authUC: authUC}
}
