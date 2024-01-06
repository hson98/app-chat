package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"hson98/app-chat/pkg/httperrs"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func (mw *MiddlewareManager) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httperrs.NewUnauthorizedError(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httperrs.NewUnauthorizedError(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httperrs.NewUnauthorizedError(err))
			return
		}

		accessToken := fields[1]
		payload, err := mw.jwtMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httperrs.NewUnauthorizedError(err))
			return
		}
		userId, err := mw.clientRedis.Get(ctx.Request.Context(), payload.ID).Result()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httperrs.NewUnauthorizedError(err))
			return
		}

		if userId != payload.UserID.String() {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "login session has been blocked.")
			return
		}
		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}

}
