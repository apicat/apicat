package openai

import (
	"context"
	"errors"
	"strings"

	"github.com/apicat/apicat/backend/config"
	openAI "github.com/sashabaranov/go-openai"
	"golang.org/x/exp/slog"
)

type OpenAI struct {
	source                 string
	key                    string
	endpoint               string
	language               string
	maxTokens              int
	ChatCompletionResponse openAI.ChatCompletionResponse
}

func NewOpenAI(configs config.OpenAI, language string) *OpenAI {
	return &OpenAI{
		source:    configs.Source,
		key:       configs.Key,
		endpoint:  configs.Endpoint,
		language:  strings.ToLower(language),
		maxTokens: 1000,
	}
}

func (o *OpenAI) CreateApi(apiName string) (string, error) {
	message := o.genCreateApiMessage(apiName)
	err := o.createChatCompletion(message)
	if err != nil {
		return "", err
	}

	// The message content like: ```yaml \n xxx \n```
	// The ```yaml on the first line and the `` on the last line need to be removed
	replacer := strings.NewReplacer("```yaml\n", "", "```\n", "", "```", "")
	return replacer.Replace(o.ChatCompletionResponse.Choices[0].Message.Content), nil
}

func (o *OpenAI) CreateApiBySchema(apiName, apiPath, apiMethod, schemaContent string) (string, error) {
	message := o.genCreateApiBySchemaMessage(apiName, apiPath, apiMethod, schemaContent)
	err := o.createChatCompletion(message)
	if err != nil {
		return "", err
	}

	// The message content like: ```yaml \n xxx \n```
	// The ```yaml on the first line and the `` on the last line need to be removed
	replacer := strings.NewReplacer("```yaml\n", "", "```\n", "", "```", "")
	result := replacer.Replace(o.ChatCompletionResponse.Choices[0].Message.Content)
	return result, nil
}

func (o *OpenAI) CreateSchema(schemaName string) (string, error) {
	message := o.genCreateSchemaMessage(schemaName)
	err := o.createChatCompletion(message)
	if err != nil {
		return "", err
	}
	if strings.Contains(o.ChatCompletionResponse.Choices[0].Message.Content, "invalid content") {
		return "", errors.New("invalid content")
	}

	return o.ChatCompletionResponse.Choices[0].Message.Content, nil
}

func (o *OpenAI) ListApiBySchema(schemaName string) (string, error) {
	message := o.genListApiBySchemaMessage(schemaName)
	err := o.createChatCompletion(message)
	if err != nil {
		return "", err
	}

	return o.ChatCompletionResponse.Choices[0].Message.Content, nil
}

func (o *OpenAI) SetMaxTokens(maxTokens int) {
	o.maxTokens = maxTokens
}

func (o *OpenAI) createChatCompletion(messages []openAI.ChatCompletionMessage) error {
	var (
		err    error
		client *openAI.Client
	)

	if o.source == "openai" {
		client = openAI.NewClient(o.key)
	} else if o.source == "azure" {
		config := openAI.DefaultAzureConfig(o.key, o.endpoint)
		client = openAI.NewClientWithConfig(config)
	} else {
		slog.Debug("The OpenAI source is invalid")
		return errors.New("The OpenAI source is invalid")
	}

	o.ChatCompletionResponse, err = client.CreateChatCompletion(
		context.Background(),
		openAI.ChatCompletionRequest{
			Model:       openAI.GPT3Dot5Turbo,
			Messages:    messages,
			Temperature: 0,
		},
	)

	if err != nil {
		slog.Warn("ChatCompletion error: " + err.Error())
		return err
	}

	return nil
}
