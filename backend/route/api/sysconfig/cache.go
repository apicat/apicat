package sysconfig

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/sysconfig"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/module/cache/redis"
	protosysconfig "github.com/apicat/apicat/v2/backend/route/proto/sysconfig"
	sysconfigbase "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/base"
	sysconfigrequest "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/request"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type cacheApiImpl struct{}

func NewCacheApi() protosysconfig.CacheApi {
	return &cacheApiImpl{}
}

func (c *cacheApiImpl) Get(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigbase.ConfigList, error) {
	list, err := sysconfig.GetList(ctx, "cache")
	if err != nil {
		slog.ErrorContext(ctx, "sysconfig.GetList", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.FailedToGetCacheList"))
	}
	slist := make(sysconfigbase.ConfigList, 0, len(list))
	for _, v := range list {
		if v.Config == "" {
			slist = append(slist, &sysconfigbase.ConfigDetail{
				Driver: v.Driver,
				Use:    v.BeingUsed,
				Config: map[string]interface{}{},
			})
		} else {
			slist = append(slist, &sysconfigbase.ConfigDetail{
				Driver: v.Driver,
				Use:    v.BeingUsed,
				Config: cfgFormat(v),
			})
		}
	}
	return &slist, nil
}

func (c *cacheApiImpl) UpdateMemory(ctx *gin.Context, _ *ginrpc.Empty) (*ginrpc.Empty, error) {
	cacheConfig := &config.Cache{
		Driver: cache.LOCAL,
	}

	cache := &sysconfig.Sysconfig{
		Type:      "cache",
		Driver:    cache.LOCAL,
		BeingUsed: true,
		Config:    "{}",
	}
	if err := sysconfig.UpdateOrCreate(ctx, cache); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.CacheUpdateFailed"))
	}
	config.SetCache(cacheConfig)
	return nil, nil
}

func (c *cacheApiImpl) UpdateRedis(ctx *gin.Context, opt *sysconfigrequest.RedisOption) (*ginrpc.Empty, error) {
	if i := strings.Index(opt.Host, ":"); i < 7 {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.RedisConfigInvalid"))
	}

	cacheConfig := &config.Cache{
		Driver: cache.REDIS,
		Redis: &config.Redis{
			Host:     opt.Host,
			Password: opt.Password,
			DB:       opt.Database,
		},
	}

	cc := cacheConfig.ToMapInterface()
	if r, err := redis.NewRedis(cc["Redis"].(map[string]interface{})); err != nil {
		slog.ErrorContext(ctx, "redis.NewRedis", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.CacheUpdateFailed"))
	} else {
		if err := r.Check(); err != nil {
			slog.ErrorContext(ctx, "r.Check", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.RedisConfigInvalid"))
		}
	}

	jsonData, err := json.Marshal(opt)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.CacheUpdateFailed"))
	}

	cache := &sysconfig.Sysconfig{
		Type:      "cache",
		Driver:    cache.REDIS,
		BeingUsed: true,
		Config:    string(jsonData),
	}
	if err := sysconfig.UpdateOrCreate(ctx, cache); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.CacheUpdateFailed"))
	}

	config.SetCache(cacheConfig)
	return nil, nil
}
