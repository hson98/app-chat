package auth_usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"hson98/app-chat/config"
	auth_repository "hson98/app-chat/internal/auth/repository"
	"hson98/app-chat/internal/models"
	"hson98/app-chat/pkg/httperrs"
	"hson98/app-chat/pkg/myjwt"
	"hson98/app-chat/pkg/utils"
	"time"
)

type UseCase interface {
	Login(ctx context.Context, userLogin *models.UserLogin) (*models.UserWithToken, error)
	Register(ctx context.Context, user *models.UserCreate) (*models.UserWithToken, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	Logout(ctx context.Context, rtPayload *myjwt.Payload, atPayload *myjwt.Payload) error
}

type authUC struct {
	authRepo      auth_repository.Repository
	authRedisRepo auth_repository.RedisRepository
	config        *config.Config
	jwtMaker      myjwt.JwtMaker
}

func NewAuthUC(authRepo auth_repository.Repository, authRedisRepo auth_repository.RedisRepository, config *config.Config, jwtMaker myjwt.JwtMaker) UseCase {
	return &authUC{
		authRepo:      authRepo,
		authRedisRepo: authRedisRepo,
		config:        config,
		jwtMaker:      jwtMaker,
	}
}

func (a *authUC) Login(ctx context.Context, userLogin *models.UserLogin) (*models.UserWithToken, error) {
	findUser, err := a.authRepo.FindUser(ctx, map[string]interface{}{"email": userLogin.Email})
	if err != nil {
		return nil, errors.New(httperrs.ErrUsernameOrPasswordInvalid)
	}
	errCompare := utils.CheckPassword(findUser.Password, userLogin.Password)
	if errCompare != nil {
		return nil, errors.New(httperrs.ErrUsernameOrPasswordInvalid)
	}
	//create access token
	accessToken, accessPayload, err := CreateTokenAndSaveToRedis(ctx, a, findUser.ID, a.config.AccessTokenDuration)
	if err != nil {
		return nil, err
	}
	//create refresh token
	refreshToken, refreshPayload, err := CreateTokenAndSaveToRedis(ctx, a, findUser.ID, a.config.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}
	return &models.UserWithToken{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  findUser,
	}, nil
}
func CreateTokenAndSaveToRedis(c context.Context, a *authUC, idUser uuid.UUID, duration time.Duration) (string, *myjwt.Payload, error) {
	//Create token
	token, payload, err := a.jwtMaker.CreateToken(idUser, duration)
	if err != nil {
		return "", nil, err
	}
	//Save to redis
	errRedis := a.authRedisRepo.SaveIDSession(c, payload.ID, idUser, duration)
	if errRedis != nil {
		return "", nil, errRedis
	}
	return token, payload, nil
}
func (a *authUC) Register(ctx context.Context, data *models.UserCreate) (*models.UserWithToken, error) {
	user, _ := a.authRepo.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if user != nil {
		return nil, errors.New(httperrs.ErrEmailExisted)
	}
	if data.Password != data.ConfirmPassword {
		return nil, errors.New(httperrs.PassAndConfirmPassNotMatch)
	}
	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}
	userCreated, err := a.authRepo.CreateUser(ctx, &models.User{
		Email:    data.Email,
		Password: hashedPassword,
		Name:     data.FullName})
	if err != nil {
		return nil, err
	}
	//create access token
	accessToken, accessPayload, err := CreateTokenAndSaveToRedis(ctx, a, userCreated.ID, a.config.AccessTokenDuration)
	if err != nil {
		return nil, err
	}
	//create refresh token
	refreshToken, refreshPayload, err := CreateTokenAndSaveToRedis(ctx, a, userCreated.ID, a.config.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}
	return &models.UserWithToken{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  userCreated,
	}, nil
}

func (a *authUC) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return a.authRepo.FindUser(ctx, map[string]interface{}{"id": userID.String()})
}

func (a *authUC) Logout(ctx context.Context, rtPayload *myjwt.Payload, atPayload *myjwt.Payload) error {
	//check id session refresh token exits
	userID, err := a.authRedisRepo.GetUserID(ctx, rtPayload.ID)
	if err != nil {
		return err
	}
	//check user
	if userID.String() != atPayload.UserID.String() {
		return errors.New("userID different")
	}
	err = a.authRedisRepo.DeleteIDSession(ctx, atPayload.ID)
	if err != nil {
		return err
	}
	err = a.authRedisRepo.DeleteIDSession(ctx, rtPayload.ID)
	if err != nil {
		return err
	}

	return nil
}
