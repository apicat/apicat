package openai

import (
	"apicat-cloud/backend/module/llm/common"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	oai "github.com/sashabaranov/go-openai"
)

type openAI struct {
	client        *oai.Client
	llmName       string
	embeddingName string
}

func NewOpenAI(cfg map[string]interface{}) (*openAI, error) {
	var clientConfig oai.ClientConfig

	if _, ok := cfg["ApiKey"]; !ok {
		return nil, errors.New("openai config ApiKey is required")
	}

	if _, ok := cfg["Endpoint"]; ok {
		clientConfig = oai.DefaultAzureConfig(cfg["ApiKey"].(string), cfg["Endpoint"].(string))
	} else {
		clientConfig = oai.DefaultConfig(cfg["ApiKey"].(string))
	}

	if _, ok := cfg["Timeout"]; ok {
		clientConfig.HTTPClient.Timeout = time.Second * time.Duration(cfg["Timeout"].(int))
	} else {
		clientConfig.HTTPClient.Timeout = time.Second * 30
	}

	o := &openAI{
		client: oai.NewClientWithConfig(clientConfig),
	}
	if _, ok := cfg["LLMName"]; ok {
		o.llmName = cfg["LLMName"].(string)
	}
	if _, ok := cfg["EmbeddingName"]; ok {
		o.embeddingName = cfg["EmbeddingName"].(string)
	}

	return o, nil
}

func (o *openAI) Check() error {
	if o.llmName == "" {
		return errors.New("model name not set")
	}
	if _, err := o.client.GetModel(context.Background(), o.llmName); err != nil {
		slog.Error("openai.Check", "err", err)
		return fmt.Errorf("%s model not found", o.llmName)
	}
	return nil
}

func (o *openAI) ChatCompletionRequest(r *common.ChatCompletionRequest) (string, error) {
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

func (o *openAI) ChatMessageRoleSystem() string {
	return oai.ChatMessageRoleSystem
}

func (o *openAI) ChatMessageRoleAssistant() string {
	return oai.ChatMessageRoleAssistant
}

func (o *openAI) ChatMessageRoleUser() string {
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
