package auth_http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	auth_usecase "hson98/app-chat/internal/auth/usecase"
	"hson98/app-chat/internal/middlewares"
	"hson98/app-chat/internal/models"
	"hson98/app-chat/pkg/httperrs"
	"hson98/app-chat/pkg/myjwt"
	"hson98/app-chat/pkg/response"
	"hson98/app-chat/pkg/utils"
	"net/http"
)

type Handlers interface {
	Login() gin.HandlerFunc
	Register() gin.HandlerFunc
	Logout() gin.HandlerFunc
	GetProfile() gin.HandlerFunc
}

type authHandler struct {
	authUC   auth_usecase.UseCase
	jwtMaker myjwt.Maker
}

func (a authHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin models.UserLogin
		if err := c.ShouldBind(&userLogin); err != nil {
			out := utils.ErrorsBindParamOrBody(err)
			response := response.BuildErrorResponse(httperrs.ErrBody, out, models.EmptyObj{})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		loginUser, err := a.authUC.Login(c, &userLogin)
		if err != nil {
			response := response.BuildErrorResponse(httperrs.HasErrTryAgain, err.Error(), models.EmptyObj{})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		response := response.BuildResponse(true, "Login successful", loginUser)
		c.JSON(http.StatusOK, response)
	}
}

func (a authHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userBody models.UserCreate
		if err := c.ShouldBind(&userBody); err != nil {
			out := utils.ErrorsBindParamOrBody(err)
			response := response.BuildErrorResponse(httperrs.ErrBody, out, models.EmptyObj{})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		userCreated, err := a.authUC.Register(c, &userBody)
		if err != nil {
			response := response.BuildErrorResponse(err.Error(), "", models.EmptyObj{})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		response := response.BuildResponse(true, "Account created successfully.", userCreated)
		c.JSON(http.StatusCreated, response)
	}
}

func (a authHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var refreshTokenBody models.RefreshTokenBody
		if err := c.ShouldBind(&refreshTokenBody); err != nil {
			out := utils.ErrorsBindParamOrBody(err)
			response := response.BuildErrorResponse(httperrs.ErrBody, out, models.EmptyObj{})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		refreshPayload, err := a.jwtMaker.VerifyToken(refreshTokenBody.RefreshToken)
		if err != nil {
			log.Printf("Refresh token expired.")
		}
		atPayload := c.MustGet(middlewares.AuthorizationPayloadKey).(*myjwt.Payload)
		errLogout := a.authUC.Logout(c, refreshPayload, atPayload)
		if errLogout != nil {
			response := response.BuildErrorResponse(httperrs.HasErrTryAgain, errLogout.Error(), models.EmptyObj{})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		response := response.BuildResponse(true, "Logout from the system successful.", models.EmptyObj{})
		c.JSON(http.StatusOK, response)
	}
}
func (a authHandler) GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userPayload := c.MustGet(middlewares.AuthorizationPayloadKey).(*myjwt.Payload)
		user, err := a.authUC.GetUserByID(c, userPayload.UserID)
		if err != nil {
			response := response.BuildErrorResponse(httperrs.HasErrTryAgain, err.Error(), models.EmptyObj{})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		response := response.BuildResponse(true, "Successfully retrieved user information.", user)
		c.JSON(http.StatusOK, response)
	}
}

func NewAuthHandler(authUC auth_usecase.UseCase, jwtMaker myjwt.Maker) Handlers {
	return &authHandler{authUC: authUC, jwtMaker: jwtMaker}
}
