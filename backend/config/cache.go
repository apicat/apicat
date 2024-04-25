package config

import (
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/module/cache/redis"
)

type Cache struct {
	Host     string `yaml:"Host"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

func SetCache(cacheConfig *Cache) {
	globalConf.Cache = cacheConfig
}

func (c *Cache) ToCfg() cache.Cache {
	return cache.Cache{
		Driver: cache.REDIS,
		Redis: redis.RedisOpt{
			Host:     c.Host,
			Password: c.Password,
			DB:       c.DB,
		},
	}
}
