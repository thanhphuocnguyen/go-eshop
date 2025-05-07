package cache

import (
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/metrics"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
	"github.com/thanhphuocnguyen/go-eshop/config"
)

type RedisCache struct {
	marshal *marshaler.Marshaler
}

// Delete implements Cache.
func (r *RedisCache) Delete(key string) error {
	panic("unimplemented")
}

// Get implements Cache.
func (r *RedisCache) Get(key string) (string, error) {
	panic("unimplemented")
}

// Set implements Cache.
func (r *RedisCache) Set(key string, value interface{}) error {
	panic("unimplemented")
}

func NewRedisCache(cfg config.Config) Cache {
	redisStore := redis_store.NewRedis(redis.NewClient(&redis.Options{
		Addr: cfg.RedisUrl,
	}))
	promMetrics := metrics.NewPrometheus("my-test-app")

	// Initialize metric cache
	cacheManager := cache.NewMetric[any](
		promMetrics,
		cache.New[any](redisStore),
	)

	marshal := marshaler.New(cacheManager)

	return &RedisCache{
		marshal,
	}
}
