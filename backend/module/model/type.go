package model

type ChatCompletionMessage struct {
	Role    string
	Content string
}

type ChatCompletionMessages []ChatCompletionMessage

func NewChatCompletionMessages(role, content string) ChatCompletionMessages {
	m := make(ChatCompletionMessages, 0)
	m = append(m, ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
	return m
}

type ChatCompletionOption struct {
	Temperature float32
	MaxTokens   int
	Messages    ChatCompletionMessages
}

func NewChatCompletionOption(msg ChatCompletionMessages) *ChatCompletionOption {
	return &ChatCompletionOption{
		Temperature: 0.3,
		MaxTokens:   5000,
		Messages:    msg,
	}
}
