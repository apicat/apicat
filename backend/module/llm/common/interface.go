package common

type ChatCompletionMessage struct {
	Role    string
	Content string
}

type ChatCompletionRequest struct {
	Model       string
	Temperature float32
	MaxTokens   int
	Messages    []ChatCompletionMessage
}

type Provider interface {
	ChatMessageRoleSystem() string
	ChatMessageRoleAssistant() string
	ChatMessageRoleUser() string
	ChatCompletionRequest(r *ChatCompletionRequest) (string, error)
	Check() error
}
