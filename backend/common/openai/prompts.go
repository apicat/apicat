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

func (o *OpenAI) genCreateApiBySchemaMessage(apiName, apiPath, apiMethod, schemaContent string) []openAI.ChatCompletionMessage {
	var prompt []string

	prompt = append(prompt, fmt.Sprintf("Generate an HTTP API for %s based on the JSON Schema enclosed in triple backticks.", apiName))
	prompt = append(prompt, fmt.Sprintf("The path of the API is <%s>, and the method of the API is <%s>.", apiPath, apiMethod))
	prompt = append(prompt, "Provide them in OpenAPI 3.0 YAML format, and the content of YAML must be complete, including basic information of an HTTP API.")
	prompt = append(prompt, "No explanation is needed in the generated content, only the YAML content itself should be returned.")
	if o.language == "zh" {
		prompt = append(prompt, "The content of the 'description' and 'title' fields in YAML must be translated into Chinese.")
	}
	prompt = append(prompt, fmt.Sprintf("JSON Schema: ```\n%s\n```", schemaContent))

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

func (o *OpenAI) genListApiBySchemaMessage(title string) []openAI.ChatCompletionMessage {
	var prompt = []string{
		"Generate a list of HTTP APIs based on the data model name enclosed in triple backticks.",
		"The generated API should be reasonable and have practical value in actual use.",
		"Including only API descriptions, request methods, and paths.",
		"Provide them in JSON format with the following keys: description, method, path.",
		"No explanation is needed in the generated content, only the JSON Schema itself should be returned.",
		"For example:",
		`[{"description": "create user", "method": "POST", "path": "/users"}]`,
	}

	if o.language == "zh" {
		prompt = append(prompt, "The 'description' field must be translated into Chinese.")
	}

	prompt = append(prompt, fmt.Sprintf("```%s```", title))
	prompt = append(prompt, "JSON format:")

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
