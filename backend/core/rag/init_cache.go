package rag

import (
	"fmt"
	"strconv"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/module/cache/common"
)

const (
	INIT_FINISHED = "finished"
	INIT_RUNNING  = "running"
)

type initCache struct {
	connect                 common.Cache
	statusKey               string
	collectionEmbedLaterKey string
	modelEmbedLaterKey      string
}

func newInitCache(projectID string) (*initCache, error) {
	c, err := cache.NewCache(config.Get().Cache.ToModuleStruct())
	if err != nil {
		return nil, err
	}
	if err := c.Check(); err != nil {
		return nil, err
	}

	return &initCache{
		connect:                 c,
		statusKey:               fmt.Sprintf("vector_init_%s", projectID),
		collectionEmbedLaterKey: fmt.Sprintf("vector_init_%s_collection_later", projectID),
		modelEmbedLaterKey:      fmt.Sprintf("vector_init_%s_model_later", projectID),
	}, nil
}

func (ic *initCache) Finished() error {
	return ic.connect.Del(ic.statusKey)
}

func (ic *initCache) SetStatus(status string) error {
	return ic.connect.Set(ic.statusKey, status, 30*time.Second)
}

func (ic *initCache) GetStatus() (string, error) {
	content, ok, err := ic.connect.Get(ic.statusKey)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	return content, nil
}

func (ic *initCache) SetCollectionLater(id uint) error {
	return ic.setLater(ic.collectionEmbedLaterKey, id)
}

func (ic *initCache) GetCollectionLater() (uint, error) {
	return ic.getLater(ic.collectionEmbedLaterKey)
}

func (ic *initCache) SetModelLater(id uint) error {
	return ic.setLater(ic.modelEmbedLaterKey, id)
}

func (ic *initCache) GetModelLater() (uint, error) {
	return ic.getLater(ic.modelEmbedLaterKey)
}

func (ic *initCache) setLater(k string, id uint) error {
	if err := ic.connect.LPush(k, id); err != nil {
		return err
	}
	return ic.connect.Expire(k, 10*time.Minute)
}

func (ic *initCache) getLater(k string) (uint, error) {
	content, ok, err := ic.connect.RPop(k)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	if id, err := strconv.ParseUint(content, 10, 64); err != nil {
		return 0, err
	} else {
		return uint(id), nil
	}
}
