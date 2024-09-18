package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisOpt struct {
	Host     string
	Password string
	DB       int
}

type redisCache struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedis(cfg RedisOpt) (*redisCache, error) {
	return &redisCache{
		ctx: context.Background(),
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Host,
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
	}, nil
}

func (r *redisCache) Check() error {
	if _, err := r.client.Ping(r.ctx).Result(); err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Set(k string, data string, du time.Duration) error {
	err := r.client.Set(r.ctx, k, data, du).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Get(k string) (string, bool, error) {
	val, err := r.client.Get(r.ctx, k).Result()
	if err == redis.Nil {
		// key does not exist
		return "", false, nil
	} else if err != nil {
		return "", false, err
	} else {
		return val, true, nil
	}
}

func (r *redisCache) LPush(k string, values ...interface{}) error {
	err := r.client.LPush(r.ctx, k, values).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) RPop(k string) (string, bool, error) {
	result, err := r.client.RPop(r.ctx, k).Result()
	if err == redis.Nil {
		// key does not exist
		return "", false, nil
	} else if err != nil {
		return "", false, err
	} else {
		return result, true, nil
	}
}

func (r *redisCache) Expire(k string, du time.Duration) error {
	err := r.client.Expire(r.ctx, k, du).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Del(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
