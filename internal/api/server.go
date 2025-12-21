package api

import (
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
	router            chi.Router
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
	discountProcessor := processors.NewDiscountProcessor(repo)
	cacheService := cache.NewRedisCache(cfg)
	tokenAuth := jwtauth.New("HS256", []byte(cfg.SymmetricKey), nil)

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
