package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	oai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	ApiKey         string
	OrganizationID string
	ApiBase        string
	LLM            string
	Embedding      string
}

type AzureOpenAI struct {
	ApiKey    string
	Endpoint  string
	LLM       string
	Embedding string
}

type openai struct {
	llm       string
	embedding string
	client    *oai.Client
}

func NewOpenAI(cfg OpenAI) *openai {
	clientConfig := oai.DefaultConfig(cfg.ApiKey)
	clientConfig.HTTPClient.Timeout = time.Second * 30

	return &openai{
		llm:       cfg.LLM,
		embedding: cfg.Embedding,
		client:    oai.NewClientWithConfig(clientConfig),
	}
}

func NewAzureOpenAI(cfg AzureOpenAI) *openai {
	clientConfig := oai.DefaultAzureConfig(cfg.ApiKey, cfg.Endpoint)
	clientConfig.HTTPClient.Timeout = time.Second * 30

	return &openai{
		llm:       cfg.LLM,
		embedding: cfg.Embedding,
		client:    oai.NewClientWithConfig(clientConfig),
	}
}

func (o *openai) Check(modelType string) error {
	switch modelType {
	case "llm":
		return o.CheckLLM()
	case "embedding":
		return o.CheckEmbedding()
	default:
		return fmt.Errorf("unknown model type: %s", modelType)
	}
}

func (o *openai) CheckLLM() error {
	if o.llm == "" {
		return errors.New("llm name not set")
	}

	msg := NewChatCompletionMessages(oai.ChatMessageRoleUser, "Hello")
	_, err := o.ChatCompletionRequest(NewChatCompletionOption(msg))
	return err
}

func (o *openai) CheckEmbedding() error {
	if o.embedding == "" {
		return errors.New("embedding name not set")
	}
	// TODO: check embedding
	return nil
}

func (o *openai) ChatCompletionRequest(r *ChatCompletionOption) (string, error) {
	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		oai.ChatCompletionRequest{
			Model:    o.llm,
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

func compileMessages(m ChatCompletionMessages) []oai.ChatCompletionMessage {
	messages := make([]oai.ChatCompletionMessage, len(m))
	for k, v := range m {
		messages[k] = oai.ChatCompletionMessage{
			Role:    v.Role,
			Content: v.Content,
		}
	}
	return messages
}
