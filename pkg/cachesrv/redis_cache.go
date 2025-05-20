package cachesrv

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/thanhphuocnguyen/go-eshop/config"
)

type RedisCache struct {
	Storage *cache.Cache
}

// Delete implements Cache.
func (r *RedisCache) Delete(c context.Context, key string) error {
	return r.Storage.Delete(c, key)
}

// Get implements Cache.
func (r *RedisCache) Get(c context.Context, key string, value interface{}) error {
	return r.Storage.Get(c, key, value)
}

// Set implements Cache.
func (r *RedisCache) Set(c context.Context, key string, value interface{}, expireIn *time.Duration) error {
	if expireIn == nil {
		expireIn = &DEFAULT_EXPIRATION
	}
	return r.Storage.Set(&cache.Item{
		Ctx:   c,
		Key:   key,
		Value: value,
		TTL:   *expireIn,
	})
}

func NewRedisCache(cfg config.Config) Cache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"api_cache_server": cfg.RedisUrl,
		},
	})

	storage := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &RedisCache{
		storage,
	}
}
