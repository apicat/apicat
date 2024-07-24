package model

type Provider interface {
	ChatMessageRoleSystem() string
	ChatMessageRoleAssistant() string
	ChatMessageRoleUser() string
	ChatCompletionRequest(r *ChatCompletionOption) (string, error)
	Check(modelType string) error
}
