package vector

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
)

type WeaviateOpt struct {
	Host   string
	ApiKey string
}

type weaviatedb struct {
	ctx    context.Context
	client *weaviate.Client
}

func NewWeaviate(cfg WeaviateOpt) (*weaviatedb, error) {
	c := weaviate.Config{
		Host:   cfg.Host,
		Scheme: "http",
	}

	if cfg.ApiKey != "" {
		c.AuthConfig = auth.ApiKey{Value: cfg.ApiKey}
	}

	if client, err := weaviate.NewClient(c); err == nil {
		return &weaviatedb{
			ctx:    context.Background(),
			client: client,
		}, nil
	} else {
		return nil, err
	}
}

func (w *weaviatedb) Check() error {
	if _, err := w.client.Misc().MetaGetter().Do(w.ctx); err != nil {
		return err
	}
	return nil
}
