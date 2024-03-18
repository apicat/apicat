package request

type ModelNameOption struct {
	LLMName       string `json:"llmName" binding:"required,gt=1"`
	EmbeddingName string `json:"embeddingName" binding:"omitempty"`
}

type AzureOpenAIOption struct {
	ApiKey   string `json:"apiKey" binding:"required,gt=1"`
	Endpoint string `json:"endpoint" binding:"required,url"`
	ModelNameOption
}

type OpenAIOption struct {
	ApiKey         string `json:"apiKey" binding:"required,gt=1"`
	OrganizationID string `json:"organizationID" binding:"omitempty"`
	ApiBase        string `json:"apiBase" binding:"omitempty"`
	ModelNameOption
}
