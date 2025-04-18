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
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", fmt.Errorf("authorization header is not provided")))
			return
		}
		authGroup := strings.Fields(authorization)
		if len(authGroup) != 2 || authGroup[0] != authorizationType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", fmt.Errorf("authorization header is not provided")))
			return
		}

		payload, err := tokenGenerator.VerifyToken(authGroup[1])
		fmt.Println(err)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", err))
			return
		}

		ctx.Set(authorizationPayload, payload)
		ctx.Next()
	}
}

func roleMiddleware(repo repository.Repository, roles ...repository.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, createErrorResponse(http.StatusForbidden, "", fmt.Errorf("authorization payload is not provided")))
			return
		}
		user, err := repo.GetUserByID(c, authPayload.UserID)
		if err != nil {
			log.Error().Err(err).Msg("get user by ID")
			c.AbortWithStatusJSON(http.StatusForbidden, createErrorResponse(http.StatusForbidden, "", fmt.Errorf("user not found")))
			return
		}

		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}
		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, createErrorResponse(http.StatusForbidden, "", fmt.Errorf("user does not have permission")))
			return
		}
		c.Next()
	}
}
