package cacheDriver

import (
	"example/core/helper/enviromentHelper"

	"github.com/go-redis/redis/v8"

	"example/core/config"
)

type CacheDriver interface {
	GetCache() *redis.Client
}

type cacheDriver struct {
	cache *redis.Client
}

func New(config config.Config) CacheDriver {
	var cache *redis.Client
	cacheConfig := config.GetCache()

	if enviromentHelper.IsDevelopment() {
		cache = redis.NewClient(&redis.Options{
			Addr: cacheConfig.Addr,
			DB:   0,
		})
	} else {
		cache = redis.NewClient(&redis.Options{
			Addr:     cacheConfig.Addr,
			Password: cacheConfig.Password,
			DB:       0,
		})
	}

	return cacheDriver{
		cache,
	}
}

func (c cacheDriver) GetCache() *redis.Client {
	return c.cache
}
