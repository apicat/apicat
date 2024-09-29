package model

import (
	"context"
	"time"

	oai "github.com/sashabaranov/go-openai"
)

type AzureOpenAI struct {
	ApiKey              string
	Endpoint            string
	LLM                 string
	LLMDeployName       string
	Embedding           string
	EmbeddingDeployName string
}

func newAzureOpenAI(cfg AzureOpenAI) *openai {
	clientConfig := oai.DefaultAzureConfig(cfg.ApiKey, cfg.Endpoint)
	clientConfig.HTTPClient.Timeout = time.Second * 30
	clientConfig.AzureModelMapperFunc = func(model string) string {
		azureModelMapping := make(map[string]string)
		if cfg.LLM != "" && cfg.LLMDeployName != "" {
			azureModelMapping[cfg.LLM] = cfg.LLMDeployName
		}
		if cfg.Embedding != "" && cfg.EmbeddingDeployName != "" {
			azureModelMapping[cfg.Embedding] = cfg.EmbeddingDeployName
		}
		return azureModelMapping[model]
	}

	return &openai{
		llm:       cfg.LLM,
		embedding: cfg.Embedding,
		client:    oai.NewClientWithConfig(clientConfig),
		ctx:       context.Background(),
	}
}
