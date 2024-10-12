package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	arkruntimemodel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	vengine "github.com/volcengine/volcengine-go-sdk/volcengine"
)

type VolcanoEngine struct {
	ApiKey              string
	Region              string
	BaseUrl             string
	LLM                 string
	LLMEndpointID       string
	Embedding           string
	EmbeddingEndpointID string
}

type volcengine struct {
	llm                 string
	llmEndpointID       string
	embedding           string
	embeddingEndpointID string
	client              *arkruntime.Client
	ctx                 context.Context
}

var VOLCANOENGINE_LLM_SUPPORTS = []string{
	"Doubao-lite-4k",
	"Doubao-lite-32k",
	"Doubao-lite-128k",
	"Doubao-pro-4k",
	"Doubao-pro-32k",
	"Doubao-pro-128k",
}

var VOLCANOENGINE_EMBEDDING_SUPPORTS = []string{
	"Doubao-embedding",
	"Doubao-embedding-large",
}

func newVolcanoEngine(cfg VolcanoEngine) *volcengine {
	return &volcengine{
		llm:                 cfg.LLM,
		llmEndpointID:       cfg.LLMEndpointID,
		embedding:           cfg.Embedding,
		embeddingEndpointID: cfg.EmbeddingEndpointID,
		client: arkruntime.NewClientWithApiKey(
			cfg.ApiKey,
			arkruntime.WithBaseUrl(cfg.BaseUrl),
		),
		ctx: context.Background(),
	}
}

func (v *volcengine) Check(modelType string) error {
	switch modelType {
	case "llm":
		if !ModelAvailable(VOLCANOENGINE, modelType, v.llm) {
			return fmt.Errorf("llm model %s not supported", v.llm)
		}
		return v.checkLLM()
	case "embedding":
		if !ModelAvailable(VOLCANOENGINE, modelType, v.embedding) {
			return fmt.Errorf("embedding model %s not supported", v.embedding)
		}
		return v.checkEmbedding()
	default:
		return fmt.Errorf("unknown model type: %s", modelType)
	}
}

func (v *volcengine) checkLLM() error {
	if v.llm == "" {
		return errors.New("llm name not set")
	}

	msg := NewChatCompletionMessages(arkruntimemodel.ChatMessageRoleUser, "Hello")
	_, err := v.ChatCompletionRequest(NewChatCompletionOption(msg))
	return err
}

func (v *volcengine) checkEmbedding() error {
	if v.embedding == "" {
		return errors.New("embedding name not set")
	}
	_, err := v.CreateEmbeddings("Hello")
	return err
}

func (v *volcengine) CreateEmbeddings(input string) ([]float32, error) {
	resp, err := v.client.CreateEmbeddings(v.ctx, arkruntimemodel.EmbeddingRequestStrings{
		Input: []string{input},
		Model: v.embeddingEndpointID,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data[0].Embedding, nil
}

func (v *volcengine) ChatCompletionRequest(r *ChatCompletionOption) (string, error) {
	resp, err := v.client.CreateChatCompletion(
		v.ctx,
		arkruntimemodel.ChatCompletionRequest{
			Model:    v.llmEndpointID,
			Messages: v.compileMessages(r.Messages),
		},
	)

	if err != nil {
		return "", err
	}

	return *resp.Choices[0].Message.Content.StringValue, nil
}

func (v *volcengine) ChatMessageRoleSystem() string {
	return arkruntimemodel.ChatMessageRoleSystem
}

func (v *volcengine) ChatMessageRoleAssistant() string {
	return arkruntimemodel.ChatMessageRoleAssistant
}

func (v *volcengine) ChatMessageRoleUser() string {
	return arkruntimemodel.ChatMessageRoleUser
}

func (v *volcengine) compileMessages(ms ChatCompletionMessages) []*arkruntimemodel.ChatCompletionMessage {
	messages := make([]*arkruntimemodel.ChatCompletionMessage, len(ms))
	for k, v := range ms {
		messages[k] = &arkruntimemodel.ChatCompletionMessage{
			Role: v.Role,
			Content: &arkruntimemodel.ChatCompletionMessageContent{
				StringValue: vengine.String(v.Content),
			},
		}
	}
	return messages
}
