package auth_usecase

import (
	"context"
	"github.com/gin-gonic/gin"
	auth_repository "hson98/app-chat/internal/auth/repository"
	"hson98/app-chat/internal/models"
	"hson98/app-chat/pkg/jwt"
)

type authUC struct {
	authRepo auth_repository.Repository
}

type UseCase interface {
	Login(cxt context.Context, userLogin *models.UserLogin) (*models.UserWithToken, error)
	Logout(ctx *gin.Context, rtPayload *jwt.Payload)
}

func NewAuthUC(authRepo auth_repository.Repository) {

}
