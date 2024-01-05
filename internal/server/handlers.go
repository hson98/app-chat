package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"hson98/app-chat/config"
	"hson98/app-chat/pkg/myjwt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	gin         *gin.Engine
	db          *gorm.DB
	redisClient *redis.Client
	jwtMaker    myjwt.Maker
	config      *config.Config
}

func NewServer(gin *gin.Engine, db *gorm.DB, redisClient *redis.Client, jwtMaker myjwt.Maker, config *config.Config) *Server {
	return &Server{
		gin:         gin,
		db:          db,
		redisClient: redisClient,
		jwtMaker:    jwtMaker,
		config:      config,
	}
}
func (s *Server) Run() error {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", s.config.Port),
		Handler: s.gin,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
	return err
}
