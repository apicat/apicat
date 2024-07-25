package model

type Provider interface {
	ChatMessageRoleSystem() string
	ChatMessageRoleAssistant() string
	ChatMessageRoleUser() string
	ChatCompletionRequest(*ChatCompletionOption) (string, error)
	CreateEmbeddings(string) ([]float32, error)
	Check(string) error
}
