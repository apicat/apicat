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

type openai struct {
	llm       string
	embedding string
	client    *oai.Client
	ctx       context.Context
}

var OPENAI_LLM_SUPPORTS = []string{
	oai.GPT4,
	oai.GPT4o,
	oai.GPT4oMini,
	oai.GPT4Turbo,
	oai.GPT3Dot5Turbo,
}

var OPENAI_EMBEDDING_SUPPORTS = []string{
	string(oai.SmallEmbedding3),
	string(oai.LargeEmbedding3),
	string(oai.AdaEmbeddingV2),
}

func newOpenAI(cfg OpenAI) *openai {
	clientConfig := oai.DefaultConfig(cfg.ApiKey)
	if cfg.OrganizationID != "" {
		clientConfig.OrgID = cfg.OrganizationID
	}
	if cfg.ApiBase != "" {
		clientConfig.BaseURL = cfg.ApiBase
	}
	clientConfig.HTTPClient.Timeout = time.Second * 30

	return &openai{
		llm:       cfg.LLM,
		embedding: cfg.Embedding,
		client:    oai.NewClientWithConfig(clientConfig),
		ctx:       context.Background(),
	}
}

func (o *openai) Check(modelType string) error {
	switch modelType {
	case "llm":
		if !ModelAvailable(OPENAI, modelType, o.llm) {
			return fmt.Errorf("llm model %s not supported", o.llm)
		}
		return o.checkLLM()
	case "embedding":
		if !ModelAvailable(OPENAI, modelType, o.embedding) {
			return fmt.Errorf("embedding model %s not supported", o.embedding)
		}
		return o.checkEmbedding()
	default:
		return fmt.Errorf("unknown model type: %s", modelType)
	}
}

func (o *openai) checkLLM() error {
	if o.llm == "" {
		return errors.New("llm name not set")
	}

	msg := NewChatCompletionMessages(oai.ChatMessageRoleUser, "Hello")
	_, err := o.ChatCompletionRequest(NewChatCompletionOption(msg))
	return err
}

func (o *openai) checkEmbedding() error {
	if o.embedding == "" {
		return errors.New("embedding name not set")
	}
	_, err := o.CreateEmbeddings("Hello")
	return err
}

func (o *openai) CreateEmbeddings(input string) ([]float32, error) {
	resp, err := o.client.CreateEmbeddings(o.ctx, oai.EmbeddingRequest{
		Input: []string{input},
		Model: oai.EmbeddingModel(o.embedding),
	})
	if err != nil {
		return nil, err
	}
	return resp.Data[0].Embedding, nil
}

func (o *openai) ChatCompletionRequest(r *ChatCompletionOption) (string, error) {
	resp, err := o.client.CreateChatCompletion(
		o.ctx,
		oai.ChatCompletionRequest{
			Model:    o.llm,
			Messages: o.compileMessages(r.Messages),
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

func (o *openai) compileMessages(m ChatCompletionMessages) []oai.ChatCompletionMessage {
	messages := make([]oai.ChatCompletionMessage, len(m))
	for k, v := range m {
		messages[k] = oai.ChatCompletionMessage{
			Role:    v.Role,
			Content: v.Content,
		}
	}
	return messages
}
