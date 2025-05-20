package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

const (
	authorization        = "Authorization"
	authorizationType    = "Bearer"
	authorizationPayload = "authorization_payload"
)

func authenticateMiddleware(tokenGenerator auth.TokenGenerator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(authorization)
		if len(authorization) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, createErrorResponse[bool](UnauthorizedCode, "authorize failed", fmt.Errorf("authorization header is not provided")))
			return
		}
		authGroup := strings.Fields(authorization)
		if len(authGroup) != 2 || authGroup[0] != authorizationType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, createErrorResponse[bool](UnauthorizedCode, "authorize failed", fmt.Errorf("authorization header is not valid format")))
			return
		}

		payload, err := tokenGenerator.VerifyToken(authGroup[1])
		if err != nil {
			log.Error().Err(err).Msg("verify token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, createErrorResponse[bool](UnauthorizedCode, "authorize failed", err))
			return
		}

		ctx.Set(authorizationPayload, payload)
		ctx.Next()
	}
}

func authorizeMiddleware(repo repository.Repository, roles ...repository.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, createErrorResponse[bool](PermissionDeniedCode, "authorize failed", fmt.Errorf("authorization payload is not provided")))
			return
		}
		user, err := repo.GetUserByID(c, authPayload.UserID)
		if err != nil {
			log.Error().Err(err).Msg("get user by ID")
			c.AbortWithStatusJSON(http.StatusForbidden, createErrorResponse[bool](PermissionDeniedCode, "authorize failed", fmt.Errorf("user not found")))
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
			c.AbortWithStatusJSON(http.StatusForbidden, createErrorResponse[bool](PermissionDeniedCode, "authorize failed", fmt.Errorf("user does not have permission")))
			return
		}
		c.Next()
	}
}

// Setup CORS configuration
func (sv *Server) corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "http://localhost:8080"},
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		AllowFiles:       true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		// AllowAllOrigins:  sv.config.Env == "development",
	})
}

// Setup environment mode based on configuration
func (sv *Server) setEnvModeMiddleware(router *gin.Engine) {
	router.Use(gin.Recovery())
	if sv.config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	if sv.config.Env == "development" {
		router.Use(gin.Logger())
	}
}
