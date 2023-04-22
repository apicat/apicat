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

var createApiPromptEn = "\"\"\"\nDesign a http api for %s and return content in OpenAPI3.0.0 format.\n\"\"\"\n"
var createApiPromptZh = "\"\"\"\n为%s设计一个 http api，并以 OpenAPI3.0.0 的格式返回内容。 \n\"\"\"\n"

var createSchemaPromptEn = "\"\"\"\nDesign a json schema format for the %v and return.\n\"\"\"\n"
var createSchemaPromptZh = "\"\"\"\n为%s设计一个 json schema 格式，并返回内容。 \n\"\"\"\n"

var createApiBySchemaEn = "\"\"\"\nPlease generate a \"%s\" API based on the json schema content I provided below, and return it in the data format of openapi3.0.0.\n%s\n\"\"\"\n"
var createApiBySchemaZh = "\"\"\"\n请根据我下面提供的json schema内容生成一个 \"%s\" 的api，并以openapi3.0.0的数据格式返回。\n%s\n\"\"\"\n"

var listApiBySchemaEn = "\"\"\"\nBelow I will provide a json schema content named %s, which APIs can be generated according to this schema? Please directly provide the API names that can be generated in the form of an array, Inappropriate parameters can be ignored.\n%s\n\"\"\"\n"
var listApiBySchemaZh = "\"\"\"\n下面我会提供一个名为%s的json schema内容，根据这个schema可以生成哪些API？ 请直接以数组形式给出可生成的API名称。\n%s\n\"\"\"\n"

func NewOpenAI(token, language string) *OpenAI {
	return &OpenAI{
		token:     token,
		language:  strings.ToLower(language),
		maxTokens: 500,
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
			return fmt.Sprintf(createApiPromptZh, text[0])
		}
		return fmt.Sprintf(createApiPromptEn, text[0])
	case "createSchema":
		if o.language == "zh" {
			return fmt.Sprintf(createSchemaPromptZh, text[0])
		}
		return fmt.Sprintf(createSchemaPromptEn, text[0])
	case "createApiBySchema":
		if o.language == "zh" {
			return fmt.Sprintf(createApiBySchemaZh, text[0], text[1])
		}
		return fmt.Sprintf(createApiBySchemaEn, text[0], text[1])
	case "listApiBySchema":
		if o.language == "zh" {
			return fmt.Sprintf(listApiBySchemaZh, text[0], text[1])
		}
		return fmt.Sprintf(listApiBySchemaEn, text[0], text[1])
	default:
		return ""
	}
}
