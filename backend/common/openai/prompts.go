package openai

import (
	"fmt"
	"strings"

	openAI "github.com/sashabaranov/go-openai"
)

func (o *OpenAI) genCreateApiMessage(title string) []openAI.ChatCompletionMessage {
	var prompt = []string{
		"Generate a complete content represented in OpenAPI 3.0 YAML format based on the text within triple backticks.",
		"The content of YAML must be complete, including basic information of an HTTP API.",
		"No explanation is needed in the generated content, only the YAML content itself should be returned.",
	}

	if o.language == "zh" {
		prompt = append(prompt, "The content of the 'description' and 'title' fields in YAML must be translated into Chinese.")
	}

	prompt = append(prompt, fmt.Sprintf("```an HTTP API for %s```", title))

	message := []openAI.ChatCompletionMessage{
		{
			Role:    openAI.ChatMessageRoleSystem,
			Content: "You are a programming assistant.",
		},
		{
			Role:    openAI.ChatMessageRoleUser,
			Content: strings.Join(prompt, "\n"),
		},
	}

	return message
}

var createApiBySchemaPrompt = `
Generate an HTTP API for %s based on the following JSON Schema information.
The path of the API is "%s", and the method of the API is "%s".
Provide them in OpenAPI3.0 format%s
JSON Schema: """%s
"""
`

func (o *OpenAI) genCreateSchemaMessage(title string) []openAI.ChatCompletionMessage {
	var prompt = []string{
		"Generate a data model for use in HTTP API requests based on the text enclosed in triple backticks.",
		"If the content enclosed by triple quotes is not a noun or noun phrase, return <invaild content> and end the task.",
		"This model should include its commonly used attributes.",
		"Provide the result in JSON Schema format, including the “title” and “description” fields, and the content must be complete.",
		"No explanation is needed in the generated content, only the JSON Schema content itself should be returned.",
	}

	if o.language == "zh" {
		prompt = append(prompt, "The 'description' field must be translated into Chinese, and the 'title' field must be in pure English.")
	}

	prompt = append(prompt, fmt.Sprintf("```%s```", title))

	message := []openAI.ChatCompletionMessage{
		{
			Role:    openAI.ChatMessageRoleSystem,
			Content: "You are a programming assistant.",
		},
		{
			Role:    openAI.ChatMessageRoleUser,
			Content: strings.Join(prompt, "\n"),
		},
	}

	return message
}

var listApiBySchemaPrompt = `
Generate a list of HTTP APIs used to handle %s, including only API descriptions, request methods, and paths.
Provide them in JSON format with the following keys: description, method, path.%s
For example:
[
	{"description": "create user", "method": "POST", "path": "/users"}
]
JSON format:
`

func (o *OpenAI) generatePrompt(action string, text ...string) string {
	switch action {
	case "createApiBySchema":
		if o.language == "zh" {
			return fmt.Sprintf(createApiBySchemaPrompt, text[0], text[1], text[2], ", the description and title field must be translated into Chinese.", text[3])
		}
		return fmt.Sprintf(createApiBySchemaPrompt, text[0], text[1], text[2], ".", text[3])
	case "listApiBySchema":
		if o.language == "zh" {
			return fmt.Sprintf(listApiBySchemaPrompt, text[0], "\nThe description field must be translated into Chinese.")
		}
		return fmt.Sprintf(listApiBySchemaPrompt, text[0], "")
	default:
		return ""
	}
}
