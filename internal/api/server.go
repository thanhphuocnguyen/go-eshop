package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"

	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/processors"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
	cache "github.com/thanhphuocnguyen/go-eshop/pkg/cache"
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
	router            *chi.Mux
	repo              repository.Store
	uploadService     upload.CdnUploader
	paymentSrv        *payment.PaymentManager
	tokenAuth         *jwtauth.JWTAuth
	cacheSrv          cache.CacheContainer
	taskDistributor   worker.TaskDistributor
	discountProcessor *processors.DiscountProcessor
	validator         *validator.Validate
}

func NewAPI(
	cfg config.Config,
	repo repository.Store,
	taskDistributor worker.TaskDistributor,
	uploadService upload.CdnUploader,
	paymentSrv *payment.PaymentManager,
) (*Server, error) {
	// Add nil checks for critical dependencies
	if repo == nil {
		return nil, fmt.Errorf("repository cannot be nil")
	}
	if taskDistributor == nil {
		return nil, fmt.Errorf("task distributor cannot be nil")
	}
	if uploadService == nil {
		return nil, fmt.Errorf("upload service cannot be nil")
	}
	if paymentSrv == nil {
		return nil, fmt.Errorf("payment service cannot be nil")
	}
	if cfg.SymmetricKey == "" {
		return nil, fmt.Errorf("symmetric key cannot be empty")
	}

	discountProcessor := processors.NewDiscountProcessor(repo)
	if discountProcessor == nil {
		return nil, fmt.Errorf("failed to create discount processor")
	}

	cacheService := cache.NewRedisCache(cfg)
	if cacheService == nil {
		return nil, fmt.Errorf("failed to create cache service")
	}

	tokenAuth := jwtauth.New("HS256", []byte(cfg.SymmetricKey), nil)
	if tokenAuth == nil {
		return nil, fmt.Errorf("failed to create token auth")
	}

	server := &Server{
		repo:              repo,
		config:            cfg,
		taskDistributor:   taskDistributor,
		uploadService:     uploadService,
		cacheSrv:          cacheService,
		tokenAuth:         tokenAuth,
		paymentSrv:        paymentSrv,
		discountProcessor: discountProcessor,
		validator:         validator.New(),
	}

	if server.validator == nil {
		return nil, fmt.Errorf("failed to create validator")
	}

	server.initializeRouter()

	if server.router == nil {
		return nil, fmt.Errorf("failed to initialize router")
	}

	return server, nil
}

func (s *Server) Server(addr string) *http.Server {
	if s.router == nil {
		panic("router is nil - server not properly initialized")
	}
	return &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
}

type DashboardData struct {
	Categories  []dto.CategoryDetail `json:"categories"`
	Collections []dto.CategoryDetail `json:"collections"`
}
