package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	addr     string
	password string
	db       int
	client   *redis.Client
	ctx      context.Context
}

func NewRedis(cfg map[string]interface{}) (*redisCache, error) {
	for _, v := range []string{"Host", "Password", "DB"} {
		if _, ok := cfg[v]; !ok {
			return nil, fmt.Errorf("redis config %s is required", v)
		}
	}
	return &redisCache{
		addr:     cfg["Host"].(string),
		password: cfg["Password"].(string),
		db:       int(cfg["DB"].(float64)),
		ctx:      context.Background(),
	}, nil
}

func (r *redisCache) init() {
	r.ctx = context.Background()
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.addr,
		Password: r.password,
		DB:       r.db,
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
