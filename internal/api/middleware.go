package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

const (
	authorization        = "Authorization"
	authorizationType    = "Bearer"
	authorizationPayload = "authorization_payload"
)

func authMiddleware(tokenGenerator auth.TokenGenerator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(authorization)
		if len(authorization) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, mapErrResp(fmt.Errorf("authorization header is not provided")))
			return
		}
		authGroup := strings.Fields(authorization)
		if len(authGroup) != 2 || authGroup[0] != authorizationType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, mapErrResp(fmt.Errorf("authorization header is not provided")))
			return
		}

		payload, err := tokenGenerator.VerifyToken(authGroup[1])
		fmt.Println(err)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, mapErrResp(err))
			return
		}

		ctx.Set(authorizationPayload, payload)
		ctx.Next()
	}
}

func roleMiddleware(roles ...repository.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, ok := ctx.MustGet(authorizationPayload).(*auth.Payload)
		if !ok {
			log.Error().Msg("Role middleware: cannot get authorization payload")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, mapErrResp(fmt.Errorf("authorization payload is not provided")))
			return
		}
		log.Info().Str("role", string(payload.Role)).Msg("Role middleware")

		for _, r := range roles {
			if r == payload.Role {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, mapErrResp(fmt.Errorf("user does not have permission")))
	}
}
