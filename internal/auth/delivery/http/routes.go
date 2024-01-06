package auth_http

import (
	"github.com/gin-gonic/gin"
	"hson98/app-chat/internal/middlewares"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h Handlers, mw *middlewares.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.Use(mw.AuthMiddleware())
	authGroup.POST("/logout", h.Logout())
	authGroup.GET("/getProfile", h.GetProfile())
}
