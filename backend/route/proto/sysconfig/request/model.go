package request

type ModelOption struct {
	LLM       string `json:"llm" binding:"omitempty,gt=1"`
	Embedding string `json:"embedding" binding:"omitempty,gt=1"`
}

type AzureOpenAIOption struct {
	ModelOption
	LLMDeployName       string `json:"llmDeployName" binding:"omitempty,gt=1"`
	EmbeddingDeployName string `json:"embeddingDeployName" binding:"omitempty,gt=1"`
	ApiKey              string `json:"apiKey" binding:"required,gt=1"`
	Endpoint            string `json:"endpoint" binding:"required,url"`
}

type OpenAIOption struct {
	ModelOption
	ApiKey         string `json:"apiKey" binding:"required,gt=1"`
	OrganizationID string `json:"organizationID" binding:"omitempty"`
	ApiBase        string `json:"apiBase" binding:"omitempty,url"`
}

type DefaultModelOption struct {
	Driver string `json:"driver" binding:"required,gt=1"`
	Model  string `json:"model" binding:"required,gt=1"`
}

type DefaultModelMapOption struct {
	LLM       DefaultModelOption `json:"llm" binding:"required"`
	Embedding DefaultModelOption `json:"embedding" binding:"required"`
}
