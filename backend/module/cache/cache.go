package cache

import (
	"errors"
	"log/slog"

	"github.com/apicat/apicat/backend/module/cache/common"
	"github.com/apicat/apicat/backend/module/cache/local"
	"github.com/apicat/apicat/backend/module/cache/redis"
)

const (
	LOCAL = "memory"
	REDIS = "redis"
)

func NewCache(cfg map[string]interface{}) (common.Cache, error) {
	slog.Debug("cache.NewCache", "cfg", cfg)
	if cfg == nil {
		return nil, errors.New("cache config is nil")
	}

	switch cfg["Driver"].(string) {
	case REDIS:
		return redis.NewRedis(cfg["Redis"].(map[string]interface{}))
	case LOCAL:
		return local.NewLocal()
	default:
		return nil, errors.New("cache driver not found")
	}
}

func Init(cfg map[string]interface{}) error {
	if cfg == nil {
		return errors.New("cache config is nil")
	}
	switch cfg["Driver"].(string) {
	case REDIS:
		if c, err := redis.NewRedis(cfg["Redis"].(map[string]interface{})); err != nil {
			return err
		} else {
			return c.Check()
		}
	case LOCAL:
		if c, err := local.NewLocal(); err != nil {
			return err
		} else {
			return c.Check()
		}
	default:
		return errors.New("cache driver not found")
	}
}
