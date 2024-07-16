package config

import (
	"errors"
	"os"

	"github.com/apicat/apicat/v2/backend/module/llm"
	"github.com/apicat/apicat/v2/backend/module/llm/openai"
)

type Model struct {
	LLMDriver       string
	EmbeddingDriver string
	OpenAI          *OpenAI
	AzureOpenAI     *AzureOpenAI
}

type OpenAI struct {
	ApiKey         string `json:"apiKey"`
	OrganizationID string `json:"organizationID"`
	ApiBase        string `json:"apiBase"`
	LLM            string `json:"llm"`
	Embedding      string `json:"embedding"`
}

type AzureOpenAI struct {
	ApiKey    string `json:"apiKey"`
	Endpoint  string `json:"endpoint"`
	LLM       string `json:"llm"`
	Embedding string `json:"embedding"`
}

func LoadModelConfig() {
	globalConf.Model = &Model{}

	if v, exists := os.LookupEnv("LLM_DRIVER"); exists {
		switch v {
		case llm.OPENAI:
			globalConf.Model.LLMDriver = llm.OPENAI
			loadOpenAIConfig()
		case llm.AZUREOPENAI:
			globalConf.Model.LLMDriver = llm.AZUREOPENAI
			loadAzureOpenAIConfig()
		}
	}

	if v, exists := os.LookupEnv("EMBEDDING_DRIVER"); exists {
		switch v {
		case llm.OPENAI:
			globalConf.Model.EmbeddingDriver = llm.OPENAI
			loadOpenAIConfig()
		case llm.AZUREOPENAI:
			globalConf.Model.EmbeddingDriver = llm.AZUREOPENAI
			loadAzureOpenAIConfig()
		}
	}
}

func loadOpenAIConfig() {
	globalConf.Model.OpenAI = &OpenAI{}
	if v, exists := os.LookupEnv("OPENAI_API_KEY"); exists {
		globalConf.Model.OpenAI.ApiKey = v
	}
	if v, exists := os.LookupEnv("OPENAI_ORGANIZATION_ID"); exists {
		globalConf.Model.OpenAI.OrganizationID = v
	}
	if v, exists := os.LookupEnv("OPENAI_API_BASE"); exists {
		globalConf.Model.OpenAI.ApiBase = v
	}
	if v, exists := os.LookupEnv("OPENAI_LLM"); exists {
		globalConf.Model.OpenAI.LLM = v
	}
	if v, exists := os.LookupEnv("OPENAI_EMBEDDING"); exists {
		globalConf.Model.OpenAI.Embedding = v
	}
}

func loadAzureOpenAIConfig() {
	globalConf.Model.AzureOpenAI = &AzureOpenAI{}
	if v, exists := os.LookupEnv("AZURE_OPENAI_API_KEY"); exists {
		globalConf.Model.AzureOpenAI.ApiKey = v
	}
	if v, exists := os.LookupEnv("AZURE_OPENAI_ENDPOINT"); exists {
		globalConf.Model.AzureOpenAI.Endpoint = v
	}
	if v, exists := os.LookupEnv("AZURE_OPENAI_LLM"); exists {
		globalConf.Model.AzureOpenAI.LLM = v
	}
	if v, exists := os.LookupEnv("AZURE_OPENAI_EMBEDDING"); exists {
		globalConf.Model.AzureOpenAI.Embedding = v
	}
}

func CheckModelConfig() error {
	if globalConf.Model.LLMDriver != "" {
		switch globalConf.Model.LLMDriver {
		case llm.OPENAI:
			if globalConf.Model.OpenAI == nil {
				return errors.New("openai config is empty")
			}
			if globalConf.Model.OpenAI.ApiKey == "" {
				return errors.New("openai api key is empty")
			}
			if globalConf.Model.OpenAI.LLM == "" {
				return errors.New("openai llm is empty")
			}
		case llm.AZUREOPENAI:
			if globalConf.Model.AzureOpenAI == nil {
				return errors.New("azure openai config is empty")
			}
			if globalConf.Model.AzureOpenAI.ApiKey == "" {
				return errors.New("azure openai api key is empty")
			}
			if globalConf.Model.AzureOpenAI.Endpoint == "" {
				return errors.New("azure openai endpoint is empty")
			}
			if globalConf.Model.AzureOpenAI.LLM == "" {
				return errors.New("azure openai llm is empty")
			}
		}
	}
	if globalConf.Model.EmbeddingDriver != "" {
		switch globalConf.Model.EmbeddingDriver {
		case llm.OPENAI:
			if globalConf.Model.OpenAI == nil {
				return errors.New("openai config is empty")
			}
			if globalConf.Model.OpenAI.ApiKey == "" {
				return errors.New("openai api key is empty")
			}
			if globalConf.Model.OpenAI.Embedding == "" {
				return errors.New("openai embedding is empty")
			}
		case llm.AZUREOPENAI:
			if globalConf.Model.AzureOpenAI == nil {
				return errors.New("azure openai config is empty")
			}
			if globalConf.Model.AzureOpenAI.ApiKey == "" {
				return errors.New("azure openai api key is empty")
			}
			if globalConf.Model.AzureOpenAI.Endpoint == "" {
				return errors.New("azure openai endpoint is empty")
			}
			if globalConf.Model.AzureOpenAI.Embedding == "" {
				return errors.New("azure openai embedding is empty")
			}
		}
	}
	return nil
}

func SetLLMModel(m *Model) {
	globalConf.Model.LLMDriver = m.LLMDriver
	switch m.LLMDriver {
	case llm.OPENAI:
		globalConf.Model.OpenAI = m.OpenAI
	case llm.AZUREOPENAI:
		globalConf.Model.AzureOpenAI = m.AzureOpenAI
	}
}

func SetEmbeddingModel(m *Model) {
	globalConf.Model.EmbeddingDriver = m.EmbeddingDriver
	switch m.EmbeddingDriver {
	case llm.OPENAI:
		globalConf.Model.OpenAI = m.OpenAI
	case llm.AZUREOPENAI:
		globalConf.Model.AzureOpenAI = m.AzureOpenAI
	}
}

func (m *Model) ToCfg() llm.LLM {
	if m == nil {
		return llm.LLM{}
	}

	switch l.Driver {
	case llm.OPENAI:
		return llm.LLM{
			Driver: l.Driver,
			OpenAI: openai.OpenAI{
				ApiKey:         l.OpenAI.ApiKey,
				OrganizationID: l.OpenAI.OrganizationID,
				ApiBase:        l.OpenAI.ApiBase,
				LLMName:        l.OpenAI.LLMName,
				EmbeddingName:  l.OpenAI.EmbeddingName,
			},
		}
	case llm.AZUREOPENAI:
		return llm.LLM{
			Driver: l.Driver,
			AzureOpenAI: openai.AzureOpenAI{
				ApiKey:        l.AzureOpenAI.ApiKey,
				Endpoint:      l.AzureOpenAI.Endpoint,
				LLMName:       l.AzureOpenAI.LLMName,
				EmbeddingName: l.AzureOpenAI.EmbeddingName,
			},
		}
	default:
		return llm.LLM{}
	}
}
