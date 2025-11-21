package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/cachesrv"
	"github.com/thanhphuocnguyen/go-eshop/pkg/pmgateway"
	"github.com/thanhphuocnguyen/go-eshop/pkg/upload"
)

// gin-swagger middleware
// swagger embed files

// @BasePath /api/v1
// @title           E-Commerce API
// @description     This is a sample server for a simple e-commerce API.
// @BasePath /api/v1
// @host      localhost:4000
type Server struct {
	config          config.Config
	router          *gin.Engine
	repo            repository.Repository
	tokenGenerator  auth.TokenGenerator
	uploadService   upload.UploadService
	paymentCtx      *pmgateway.PaymentContext
	cachesrv        cachesrv.Cache
	taskDistributor worker.TaskDistributor
}

func NewAPI(
	cfg config.Config,
	repo repository.Repository,
	cachesrv cachesrv.Cache,
	taskDistributor worker.TaskDistributor,
	uploadService upload.UploadService,
	paymentCtx *pmgateway.PaymentContext,
) (*Server, error) {
	tokenGenerator, err := auth.NewJwtGenerator(cfg.SymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		tokenGenerator:  tokenGenerator,
		repo:            repo,
		config:          cfg,
		taskDistributor: taskDistributor,
		uploadService:   uploadService,
		cachesrv:        cachesrv,
		paymentCtx:      paymentCtx,
	}
	server.initializeRouter()
	return server, nil
}

func (s *Server) Server(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: s.router.Handler(),
	}
}

type DashboardData struct {
	Categories  []CategoryResponse `json:"categories"`
	Collections []CategoryResponse `json:"collections"`
}
