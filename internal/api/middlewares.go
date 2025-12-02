package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

func authenticateMiddleware(tokenGenerator auth.TokenGenerator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(constants.Authorization)
		if len(authorization) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.CreateErr(UnauthorizedCode, fmt.Errorf("authorization header is not provided")))
			return
		}
		authGroup := strings.Fields(authorization)
		if len(authGroup) != 2 || authGroup[0] != constants.AuthorizationType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.CreateErr(UnauthorizedCode, fmt.Errorf("authorization header is not valid format")))
			return
		}

		payload, err := tokenGenerator.VerifyToken(authGroup[1])
		if err != nil {
			log.Error().Err(err).Msg("verify token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.CreateErr(UnauthorizedCode, err))
			return
		}

		ctx.Set(constants.AuthPayLoad, payload)
		ctx.Set(constants.UserRole, payload.RoleCode)
		ctx.Next()
	}
}

func authorizeMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.CreateErr(PermissionDeniedCode, fmt.Errorf("authorization payload is not provided")))
			return
		}

		hasRole := false
		for _, role := range roles {
			if authPayload.RoleCode == role {
				hasRole = true
				break
			}
		}
		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.CreateErr(PermissionDeniedCode, fmt.Errorf("user does not have permission")))
			return
		}
		c.Next()
	}
}

// Setup CORS configuration
func corsMiddleware() gin.HandlerFunc {
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
