package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/processors"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	cachesrv "github.com/thanhphuocnguyen/go-eshop/pkg/cache"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
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
	config            config.Config
	router            chi.Router
	repo              repository.Store
	tokenGenerator    auth.TokenGenerator
	uploadService     upload.CdnUploader
	paymentSrv        *payment.PaymentManager
	cacheSrv          cachesrv.CacheContainer
	taskDistributor   worker.TaskDistributor
	discountProcessor *processors.DiscountProcessor
}

func NewAPI(
	cfg config.Config,
	repo repository.Store,
	cachesrv cachesrv.CacheContainer,
	taskDistributor worker.TaskDistributor,
	uploadService upload.CdnUploader,
	paymentSrv *payment.PaymentManager,
	discountProcessor *processors.DiscountProcessor,
) (*Server, error) {
	tokenGenerator, err := auth.NewJwtGenerator(cfg.SymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		tokenGenerator:    tokenGenerator,
		repo:              repo,
		config:            cfg,
		taskDistributor:   taskDistributor,
		uploadService:     uploadService,
		cacheSrv:          cachesrv,
		paymentSrv:        paymentSrv,
		discountProcessor: discountProcessor,
	}
	server.initializeRouter()
	return server, nil
}

func (s *Server) Server(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
}

type DashboardData struct {
	Categories  []dto.CategoryDetail `json:"categories"`
	Collections []dto.CategoryDetail `json:"collections"`
}
