package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
)

type Server struct {
	config          config.Config
	postgres        postgres.Store
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
	tokenGenerator  auth.TokenGenerator
}

func NewAPI(cfg config.Config, postgres postgres.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenGenerator := auth.NewPasetoTokenGenerator()
	server := &Server{
		tokenGenerator:  tokenGenerator,
		postgres:        postgres,
		config:          cfg,
		taskDistributor: taskDistributor,
	}
	server.initializeRouter()
	return server, nil
}

func (sv *Server) initializeRouter() {
	router := gin.Default()
	router.Use(gin.Recovery())

	cors := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})

	router.Use(cors)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	sv.router = router
}

func (s *Server) Server(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: s.router.Handler(),
	}
}
