package openai

import (
	"context"
	"errors"
	"strings"

	openAI "github.com/sashabaranov/go-openai"
	"golang.org/x/exp/slog"
)

type OpenAI struct {
	token                  string
	language               string
	maxTokens              int
	CompletionResponse     openAI.CompletionResponse
	ChatCompletionResponse openAI.ChatCompletionResponse
}

func NewOpenAI(token, language string) *OpenAI {
	return &OpenAI{
		token:     token,
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
	prompt := o.generatePrompt("createApiBySchema", apiName, apiPath, apiMethod, schemaContent)
	err := o.createCompletion(prompt)
	if err != nil {
		return "", err
	}

	result := strings.Split(o.CompletionResponse.Choices[0].Text, "\n")
	if len(result) < 3 {
		slog.Debug("invalid result: " + o.CompletionResponse.Choices[0].Text)
		return "", errors.New("invalid result")
	}

	return result[1], nil
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
	var err error

	client := openAI.NewClient(o.token)
	o.ChatCompletionResponse, err = client.CreateChatCompletion(
		context.Background(),
		openAI.ChatCompletionRequest{
			Model:       openAI.GPT3Dot5Turbo,
			Messages:    messages,
			Temperature: 1,
		},
	)

	if err != nil {
		slog.Warn("ChatCompletion error: " + err.Error())
		return err
	}

	return nil
}

func (o *OpenAI) createCompletion(prompt string) error {
	var err error

	c := openAI.NewClient(o.token)
	ctx := context.Background()

	req := openAI.CompletionRequest{
		Model:           openAI.GPT3TextDavinci003,
		MaxTokens:       o.maxTokens,
		Prompt:          prompt,
		Temperature:     0.7,
		TopP:            1.0,
		PresencePenalty: 0.0,
		Stop:            []string{"\"\"\""},
	}
	o.CompletionResponse, err = c.CreateCompletion(ctx, req)
	if err != nil {
		return err
	}

	if o.CompletionResponse.Usage.TotalTokens > o.maxTokens {
		return errors.New("tokens used more than maxTokens")
	}

	return nil
}
