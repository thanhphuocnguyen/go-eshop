package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

func GetImageName(originFileName string) string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), strings.ReplaceAll(originFileName, " ", "-"))
}

func StringPtr(s string) *string {
	return &s
}

func GetLimitOffset(c *gin.Context) (limit, offset int) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}
	return
}

func GetPaymentGatewayInstanceFromName(name repository.PaymentGateway, cfg config.Config) payment.PaymentStrategy {
	switch name {
	case repository.PaymentGatewayStripe:
		instance, err := payment.NewStripePayment(cfg.StripeSecretKey)
		if err != nil {
			return nil
		}
		return instance
	default:
		return nil
	}
}
