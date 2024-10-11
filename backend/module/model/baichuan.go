package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	oai "github.com/sashabaranov/go-openai"
)

type Baichuan struct {
	ApiKey    string
	LLM       string
	Embedding string
}

type baichuan struct {
	llm       string
	embedding string
	client    *oai.Client
	ctx       context.Context
}

var BAICHUAN_LLM_SUPPORTS = []string{
	"Baichuan4",
	"Baichuan3-Turbo",
	"Baichuan2-Turbo",
}

var BAICHUAN_EMBEDDING_SUPPORTS = []string{
	"Baichuan-Text-Embedding",
}

func newBaichuan(cfg Baichuan) *baichuan {
	clientConfig := oai.DefaultConfig(cfg.ApiKey)
	clientConfig.BaseURL = "https://api.baichuan-ai.com/v1"
	clientConfig.HTTPClient.Timeout = time.Second * 30

	return &baichuan{
		llm:       cfg.LLM,
		embedding: cfg.Embedding,
		client:    oai.NewClientWithConfig(clientConfig),
		ctx:       context.Background(),
	}
}

func (b *baichuan) Check(modelType string) error {
	switch modelType {
	case "llm":
		if !ModelAvailable(BAICHUAN, modelType, b.llm) {
			return fmt.Errorf("llm model %s not supported", b.llm)
		}
		return b.checkLLM()
	case "embedding":
		if !ModelAvailable(BAICHUAN, modelType, b.embedding) {
			return fmt.Errorf("embedding model %s not supported", b.embedding)
		}
		return b.checkEmbedding()
	default:
		return fmt.Errorf("unknown model type: %s", modelType)
	}
}

func (b *baichuan) checkLLM() error {
	if b.llm == "" {
		return errors.New("llm name not set")
	}

	msg := NewChatCompletionMessages(oai.ChatMessageRoleUser, "Hello")
	_, err := b.ChatCompletionRequest(NewChatCompletionOption(msg))
	return err
}

func (b *baichuan) checkEmbedding() error {
	if b.embedding == "" {
		return errors.New("embedding name not set")
	}
	_, err := b.CreateEmbeddings("Hello")
	return err
}

func (b *baichuan) CreateEmbeddings(input string) ([]float32, error) {
	resp, err := b.client.CreateEmbeddings(b.ctx, oai.EmbeddingRequest{
		Input: []string{input},
		Model: oai.EmbeddingModel(b.embedding),
	})
	if err != nil {
		return nil, err
	}
	return resp.Data[0].Embedding, nil
}

func (b *baichuan) ChatCompletionRequest(r *ChatCompletionOption) (string, error) {
	resp, err := b.client.CreateChatCompletion(
		b.ctx,
		oai.ChatCompletionRequest{
			Model:    b.llm,
			Messages: b.compileMessages(r.Messages),
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (b *baichuan) ChatMessageRoleSystem() string {
	return oai.ChatMessageRoleSystem
}

func (b *baichuan) ChatMessageRoleAssistant() string {
	return oai.ChatMessageRoleAssistant
}

func (b *baichuan) ChatMessageRoleUser() string {
	return oai.ChatMessageRoleUser
}

func (b *baichuan) compileMessages(m ChatCompletionMessages) []oai.ChatCompletionMessage {
	messages := make([]oai.ChatCompletionMessage, len(m))
	for k, v := range m {
		messages[k] = oai.ChatCompletionMessage{
			Role:    v.Role,
			Content: v.Content,
		}
	}
	return messages
}
