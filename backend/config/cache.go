package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/module/cache/redis"
)

type Cache struct {
	Host     string
	Username string
	Password string
	DB       int
}

func LoadCacheConfig() {
	globalConf.Cache = &Cache{}
	if v, exists := os.LookupEnv("REDIS_HOST"); exists {
		globalConf.Cache.Host = v
	}
	if v, exists := os.LookupEnv("REDIS_USERNAME"); exists {
		globalConf.Cache.Username = v
	}
	if v, exists := os.LookupEnv("REDIS_PASSWORD"); exists {
		globalConf.Cache.Password = v
	}
	if v, exists := os.LookupEnv("REDIS_DB"); exists {
		if i, err := strconv.Atoi(v); err == nil {
			globalConf.Cache.DB = i
		} else {
			globalConf.Cache.DB = 0
		}
	}
}

func CheckCacheConfig() error {
	if globalConf.Cache.Host == "" {
		return errors.New("cache host is empty")
	}
	if globalConf.Cache.DB < 0 || globalConf.Cache.DB > 15 {
		return errors.New("cache db is invalid")
	}
	return nil
}

func (c *Cache) ToModuleStruct() cache.Cache {
	return cache.Cache{
		Driver: cache.REDIS,
		Redis: redis.RedisOpt{
			Host:     c.Host,
			Password: c.Password,
			DB:       c.DB,
		},
	}
}
