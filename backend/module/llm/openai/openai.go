package openai

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/apicat/apicat/v2/backend/module/llm/common"

	oai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	ApiKey         string
	OrganizationID string
	ApiBase        string
	LLMName        string
	EmbeddingName  string
	Timeout        int
}

type AzureOpenAI struct {
	ApiKey        string
	Endpoint      string
	LLMName       string
	EmbeddingName string
	Timeout       int
}

type openai struct {
	llmName       string
	embeddingName string
	client        *oai.Client
}

func NewOpenAI(cfg OpenAI) *openai {
	clientConfig := oai.DefaultConfig(cfg.ApiKey)

	if cfg.Timeout > 0 {
		clientConfig.HTTPClient.Timeout = time.Second * time.Duration(cfg.Timeout)
	} else {
		clientConfig.HTTPClient.Timeout = time.Second * 30
	}

	return &openai{
		llmName:       cfg.LLMName,
		embeddingName: cfg.EmbeddingName,
		client:        oai.NewClientWithConfig(clientConfig),
	}
}

func NewAzureOpenAI(cfg AzureOpenAI) *openai {
	clientConfig := oai.DefaultAzureConfig(cfg.ApiKey, cfg.Endpoint)

	if cfg.Timeout > 0 {
		clientConfig.HTTPClient.Timeout = time.Second * time.Duration(cfg.Timeout)
	} else {
		clientConfig.HTTPClient.Timeout = time.Second * 30
	}

	return &openai{
		llmName:       cfg.LLMName,
		embeddingName: cfg.EmbeddingName,
		client:        oai.NewClientWithConfig(clientConfig),
	}
}

func (o *openai) Check() error {
	if o.llmName == "" {
		return errors.New("model name not set")
	}
	if _, err := o.client.GetModel(context.Background(), o.llmName); err != nil {
		slog.Error("openai.Check", "err", err)
		return fmt.Errorf("%s model not found", o.llmName)
	}
	return nil
}

func (o *openai) ChatCompletionRequest(r *common.ChatCompletionRequest) (string, error) {
	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		oai.ChatCompletionRequest{
			Model:    o.llmName,
			Messages: compileMessages(r.Messages),
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (o *openai) ChatMessageRoleSystem() string {
	return oai.ChatMessageRoleSystem
}

func (o *openai) ChatMessageRoleAssistant() string {
	return oai.ChatMessageRoleAssistant
}

func (o *openai) ChatMessageRoleUser() string {
	return oai.ChatMessageRoleUser
}

func compileMessages(m []common.ChatCompletionMessage) []oai.ChatCompletionMessage {
	messages := make([]oai.ChatCompletionMessage, len(m))
	for k, v := range m {
		messages[k] = oai.ChatCompletionMessage{
			Role:    v.Role,
			Content: v.Content,
		}
	}
	return messages
}
