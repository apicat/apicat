package cache

import (
	"errors"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/module/cache/common"
	"github.com/apicat/apicat/v2/backend/module/cache/local"
	"github.com/apicat/apicat/v2/backend/module/cache/redis"
)

const (
	MEMORY = "memory"
	REDIS  = "redis"
)

type Cache struct {
	Driver string
	Redis  redis.RedisOpt
}

func NewCache(cfg Cache) (common.Cache, error) {
	slog.Debug("cache.NewCache", "cfg", cfg)

	switch cfg.Driver {
	case REDIS:
		return redis.NewRedis(cfg.Redis)
	case MEMORY:
		return local.NewLocal()
	default:
		return nil, errors.New("cache driver not found")
	}
}

func Init(cfg Cache) error {
	switch cfg.Driver {
	case REDIS:
		if c, err := redis.NewRedis(cfg.Redis); err != nil {
			return err
		} else {
			return c.Check()
		}
	case MEMORY:
		if c, err := local.NewLocal(); err != nil {
			return err
		} else {
			return c.Check()
		}
	default:
		return errors.New("cache driver not found")
	}
}
