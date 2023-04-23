package openai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	openAI "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	token              string
	language           string
	maxTokens          int
	CompletionResponse openAI.CompletionResponse
}

var createApiBySchemaEn = "\"\"\"\nPlease generate a \"%s\" API based on the json schema content I provided below, and return it in the data format of openapi3.0.0.\n%s\n\"\"\"\n"
var createApiBySchemaZh = "\"\"\"\n请根据我下面提供的json schema内容生成一个 \"%s\" 的api，并以openapi3.0.0的数据格式返回。\n%s\n\"\"\"\n"

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

func (o *OpenAI) CreateApiBySchema(apiName, schemaContent string) (string, error) {
	prompt := o.generatePrompt("createApiBySchema", apiName, schemaContent)
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

func (o *OpenAI) ListApiBySchema(schemaName, schemaContent string) (string, error) {
	prompt := o.generatePrompt("listApiBySchema", schemaName, schemaContent)
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

func (o *OpenAI) generatePrompt(action string, text ...string) string {
	switch action {
	case "createApi":
		if o.language == "zh" {
			return fmt.Sprintf(createApiPrompt, text[0])
		}
		return fmt.Sprintf(createApiPrompt, text[0])
	case "createSchema":
		if o.language == "zh" {
			return fmt.Sprintf(createSchemaPrompt, text[0]+"\nIf there is a description field in the returned content, please translate the content into Chinese.")
		}
		return fmt.Sprintf(createSchemaPrompt, text[0])
	case "createApiBySchema":
		if o.language == "zh" {
			return fmt.Sprintf(createApiBySchemaZh, text[0], text[1])
		}
		return fmt.Sprintf(createApiBySchemaEn, text[0], text[1])
	case "listApiBySchema":
		if o.language == "zh" {
			return fmt.Sprintf(listApiBySchemaPrompt, text[0], text[1]+"\nPlease translate the content in the description field into Chinese.")
		}
		return fmt.Sprintf(listApiBySchemaPrompt, text[0], text[1])
	default:
		return ""
	}
}
