package openai

import (
	"context"
	"errors"
	"strings"

	openAI "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	token              string
	language           string
	maxTokens          int
	CompletionResponse openAI.CompletionResponse
}

func NewOpenAI(token, language string) *OpenAI {
	return &OpenAI{
		token:     token,
		language:  strings.ToLower(language),
		maxTokens: 1000,
	}
}

func (o *OpenAI) CreateApi(apiName string) (string, error) {
	prompt := o.generatePrompt("createApi", apiName)
	err := o.createCompletion(prompt)
	if err != nil {
		return "", err
	}

	return o.CompletionResponse.Choices[0].Text, nil
}

func (o *OpenAI) CreateApiBySchema(apiName, apiPath, apiMethod, schemaContent string) (string, error) {
	prompt := o.generatePrompt("createApiBySchema", apiName, apiPath, apiMethod, schemaContent)
	err := o.createCompletion(prompt)
	if err != nil {
		return "", err
	}

	return o.CompletionResponse.Choices[0].Text, nil
}

func (o *OpenAI) CreateSchema(schemaName string) (string, error) {
	prompt := o.generatePrompt("createSchema", schemaName)
	err := o.createCompletion(prompt)
	if err != nil {
		return "", err
	}

	return o.CompletionResponse.Choices[0].Text, nil
}

func (o *OpenAI) ListApiBySchema(schemaName string) (string, error) {
	prompt := o.generatePrompt("listApiBySchema", schemaName)
	err := o.createCompletion(prompt)
	if err != nil {
		return "", err
	}

	return o.CompletionResponse.Choices[0].Text, nil
}

func (o *OpenAI) SetMaxTokens(maxTokens int) {
	o.maxTokens = maxTokens
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
