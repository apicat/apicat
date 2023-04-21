package openai

import (
	"context"
	"errors"
	"fmt"

	openAI "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	token              string
	language           string
	maxTokens          int
	CompletionResponse openAI.CompletionResponse
}

var createApiPromptEn = "\"\"\"\nDesign a http api for %s and return content in OpenAPI3.0.0 format.\n\"\"\"\n"
var createApiPromptZh = "\"\"\"\n为%s设计一个 http api，并以 OpenAPI3.0.0 的格式返回内容。 \n\"\"\"\n"

var createSchemaPromptEn = "\"\"\"\nDesign a json schema format for the %v and return.\n\"\"\"\n"
var createSchemaPromptZh = "\"\"\"\n为%s设计一个 json schema 格式，并返回内容。 \n\"\"\"\n"

func NewOpenAI(token, language string) *OpenAI {
	return &OpenAI{
		token:     token,
		language:  language,
		maxTokens: 500,
	}
}

func (o *OpenAI) CreateApi(apiName string) (string, error) {
	prompt := o.generatePrompt(apiName, "createApi")
	err := o.createCompletion(prompt)
	if err != nil {
		return "", err
	}

	return o.CompletionResponse.Choices[0].Text, nil
}

func (o *OpenAI) CreateSchema(schemaName string) (string, error) {
	prompt := o.generatePrompt(schemaName, "createSchema")
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
		Temperature:     0,
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

func (o *OpenAI) generatePrompt(text, action string) string {
	switch action {
	case "createApi":
		if o.language == "zh" {
			return fmt.Sprintf(createApiPromptZh, text)
		}
		return fmt.Sprintf(createApiPromptEn, text)
	case "createSchema":
		if o.language == "zh" {
			return fmt.Sprintf(createSchemaPromptZh, text)
		}
		return fmt.Sprintf(createSchemaPromptEn, text)
	case "createApiBySchema":
		if o.language == "zh" {
			return fmt.Sprintf(createApiPromptZh, text)
		}
		return fmt.Sprintf(createApiPromptEn, text)
	default:
		return ""
	}
}
