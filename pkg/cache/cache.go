package cachesrv

import (
	"context"
	"errors"
	"time"
)

// ErrCacheMiss is returned when a requested cache key is not found
var ErrCacheMiss = errors.New("cache miss")

type Cache interface {
	Set(c context.Context, key string, value interface{}, expireIn *time.Duration) error
	Get(c context.Context, key string, value interface{}) error
	Delete(c context.Context, key string) error
}

var DEFAULT_EXPIRATION time.Duration = 5 * time.Minute

const (
	USER_KEY_PREFIX             = "user:"
	PRODUCT_KEY_PREFIX          = "product:"
	ORDER_KEY_PREFIX            = "order:"
	ORDER_ITEM_KEY_PREFIX       = "order_item:"
	PRODUCT_CATEGORY_KEY_PREFIX = "product_category:"
)
