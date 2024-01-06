package server

import (
	"github.com/gin-gonic/gin"
	auth_http "hson98/app-chat/internal/auth/delivery/http"
	auth_repository "hson98/app-chat/internal/auth/repository"
	auth_usecase "hson98/app-chat/internal/auth/usecase"
	"hson98/app-chat/internal/middlewares"
	"net/http"
)

func (s *Server) MapHandlers(r *gin.Engine) {
	//auth
	authRedisRepo := auth_repository.NewAuthRedisRepo(s.redisClient)
	authRepo := auth_repository.NewAuthRepo(s.db)
	authUC := auth_usecase.NewAuthUC(authRepo, authRedisRepo, s.config, s.jwtMaker)
	authHandler := auth_http.NewAuthHandler(authUC)
	mw := middlewares.NewMiddlewareManager([]string{"*"}, s.jwtMaker, s.redisClient)
	r.Use(middlewares.CORSMiddleware())
	v1 := r.Group("/api/v1")
	{
		//auth
		authGroup := v1.Group("/auth")
		auth_http.MapAuthRoutes(authGroup, authHandler, mw)
	}
	health := v1.Group("/health")
	health.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
}
