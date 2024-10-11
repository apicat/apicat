package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	oai "github.com/sashabaranov/go-openai"
)

type DeepSeek struct {
	ApiKey    string
	LLM       string
	Embedding string
}

type deepseek struct {
	llm       string
	embedding string
	client    *oai.Client
	ctx       context.Context
}

var DEEPSEEK_LLM_SUPPORTS = []string{
	"deepseek-chat",
}

var DEEPSEEK_EMBEDDING_SUPPORTS = []string{}

func newDeepSeek(cfg DeepSeek) *deepseek {
	clientConfig := oai.DefaultConfig(cfg.ApiKey)
	clientConfig.BaseURL = "https://api.deepseek.com/v1"
	clientConfig.HTTPClient.Timeout = time.Second * 30

	return &deepseek{
		llm:       cfg.LLM,
		embedding: cfg.Embedding,
		client:    oai.NewClientWithConfig(clientConfig),
		ctx:       context.Background(),
	}
}

func (d *deepseek) Check(modelType string) error {
	switch modelType {
	case "llm":
		if !ModelAvailable(DEEPSEEK, modelType, d.llm) {
			return fmt.Errorf("llm model %s not supported", d.llm)
		}
		return d.checkLLM()
	case "embedding":
		if !ModelAvailable(DEEPSEEK, modelType, d.embedding) {
			return fmt.Errorf("embedding model %s not supported", d.embedding)
		}
		return d.checkEmbedding()
	default:
		return fmt.Errorf("unknown model type: %s", modelType)
	}
}

func (d *deepseek) checkLLM() error {
	if d.llm == "" {
		return errors.New("llm name not set")
	}

	msg := NewChatCompletionMessages(oai.ChatMessageRoleUser, "Hello")
	_, err := d.ChatCompletionRequest(NewChatCompletionOption(msg))
	return err
}

func (d *deepseek) checkEmbedding() error {
	if d.embedding == "" {
		return errors.New("embedding name not set")
	}
	_, err := d.CreateEmbeddings("Hello")
	return err
}

func (d *deepseek) CreateEmbeddings(input string) ([]float32, error) {
	resp, err := d.client.CreateEmbeddings(d.ctx, oai.EmbeddingRequest{
		Input: []string{input},
		Model: oai.EmbeddingModel(d.embedding),
	})
	if err != nil {
		return nil, err
	}
	return resp.Data[0].Embedding, nil
}

func (d *deepseek) ChatCompletionRequest(r *ChatCompletionOption) (string, error) {
	resp, err := d.client.CreateChatCompletion(
		d.ctx,
		oai.ChatCompletionRequest{
			Model:    d.llm,
			Messages: d.compileMessages(r.Messages),
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (d *deepseek) ChatMessageRoleSystem() string {
	return oai.ChatMessageRoleSystem
}

func (d *deepseek) ChatMessageRoleAssistant() string {
	return oai.ChatMessageRoleAssistant
}

func (d *deepseek) ChatMessageRoleUser() string {
	return oai.ChatMessageRoleUser
}

func (d *deepseek) compileMessages(ms ChatCompletionMessages) []oai.ChatCompletionMessage {
	messages := make([]oai.ChatCompletionMessage, len(ms))
	for k, v := range ms {
		messages[k] = oai.ChatCompletionMessage{
			Role:    v.Role,
			Content: v.Content,
		}
	}
	return messages
}
