package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisOpt struct {
	Host     string
	Password string
	DB       int
}

type redisCache struct {
	cfg    RedisOpt
	client *redis.Client
	ctx    context.Context
}

func NewRedis(cfg RedisOpt) (*redisCache, error) {
	return &redisCache{
		cfg: cfg,
		ctx: context.Background(),
	}, nil
}

func (r *redisCache) init() {
	r.ctx = context.Background()
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.cfg.Host,
		Password: r.cfg.Password,
		DB:       r.cfg.DB,
	})
}

func (r *redisCache) Check() error {
	r.init()
	if _, err := r.client.Ping(r.ctx).Result(); err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Set(k string, data string, du time.Duration) error {
	r.init()
	err := r.client.Set(r.ctx, k, data, du).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Get(key string) (string, bool, error) {
	r.init()
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", false, fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", false, err
	} else {
		return val, true, nil
	}
}

func (r *redisCache) Del(key string) error {
	r.init()
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
