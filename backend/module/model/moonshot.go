package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	oai "github.com/sashabaranov/go-openai"
)

type Moonshot struct {
	ApiKey    string
	LLM       string
	Embedding string
}

type moonshot struct {
	llm       string
	embedding string
	client    *oai.Client
	ctx       context.Context
}

var MOONSHOT_LLM_SUPPORTS = []string{
	"moonshot-v1-8k",
	"moonshot-v1-32k",
}

var MOONSHOT_EMBEDDING_SUPPORTS = []string{}

func newMoonshot(cfg Moonshot) *moonshot {
	clientConfig := oai.DefaultConfig(cfg.ApiKey)
	clientConfig.BaseURL = "https://api.moonshot.cn/v1"
	clientConfig.HTTPClient.Timeout = time.Second * 30

	return &moonshot{
		llm:       cfg.LLM,
		embedding: cfg.Embedding,
		client:    oai.NewClientWithConfig(clientConfig),
		ctx:       context.Background(),
	}
}

func (m *moonshot) Check(modelType string) error {
	switch modelType {
	case "llm":
		if !ModelAvailable(MOONSHOT, modelType, m.llm) {
			return fmt.Errorf("llm model %s not supported", m.llm)
		}
		return m.checkLLM()
	case "embedding":
		if !ModelAvailable(MOONSHOT, modelType, m.embedding) {
			return fmt.Errorf("embedding model %s not supported", m.embedding)
		}
		return m.checkEmbedding()
	default:
		return fmt.Errorf("unknown model type: %s", modelType)
	}
}

func (m *moonshot) checkLLM() error {
	if m.llm == "" {
		return errors.New("llm name not set")
	}

	msg := NewChatCompletionMessages(oai.ChatMessageRoleUser, "Hello")
	_, err := m.ChatCompletionRequest(NewChatCompletionOption(msg))
	return err
}

func (m *moonshot) checkEmbedding() error {
	if m.embedding == "" {
		return errors.New("embedding name not set")
	}
	_, err := m.CreateEmbeddings("Hello")
	return err
}

func (m *moonshot) CreateEmbeddings(input string) ([]float32, error) {
	resp, err := m.client.CreateEmbeddings(m.ctx, oai.EmbeddingRequest{
		Input: []string{input},
		Model: oai.EmbeddingModel(m.embedding),
	})
	if err != nil {
		return nil, err
	}
	return resp.Data[0].Embedding, nil
}

func (m *moonshot) ChatCompletionRequest(r *ChatCompletionOption) (string, error) {
	resp, err := m.client.CreateChatCompletion(
		m.ctx,
		oai.ChatCompletionRequest{
			Model:    m.llm,
			Messages: m.compileMessages(r.Messages),
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (m *moonshot) ChatMessageRoleSystem() string {
	return oai.ChatMessageRoleSystem
}

func (m *moonshot) ChatMessageRoleAssistant() string {
	return oai.ChatMessageRoleAssistant
}

func (m *moonshot) ChatMessageRoleUser() string {
	return oai.ChatMessageRoleUser
}

func (m *moonshot) compileMessages(ms ChatCompletionMessages) []oai.ChatCompletionMessage {
	messages := make([]oai.ChatCompletionMessage, len(ms))
	for k, v := range ms {
		messages[k] = oai.ChatCompletionMessage{
			Role:    v.Role,
			Content: v.Content,
		}
	}
	return messages
}
