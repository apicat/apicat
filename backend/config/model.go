package config

import (
	"errors"
	"os"

	"github.com/apicat/apicat/v2/backend/module/model"
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
		case model.OPENAI:
			globalConf.Model.LLMDriver = model.OPENAI
			loadOpenAIConfig()
		case model.AZUREOPENAI:
			globalConf.Model.LLMDriver = model.AZUREOPENAI
			loadAzureOpenAIConfig()
		}
	}

	if v, exists := os.LookupEnv("EMBEDDING_DRIVER"); exists {
		switch v {
		case model.OPENAI:
			globalConf.Model.EmbeddingDriver = model.OPENAI
			loadOpenAIConfig()
		case model.AZUREOPENAI:
			globalConf.Model.EmbeddingDriver = model.AZUREOPENAI
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
		case model.OPENAI:
			if err := checkOpenAI("llm"); err != nil {
				return err
			}
		case model.AZUREOPENAI:
			if err := checkAzureOpenAI("llm"); err != nil {
				return err
			}
		}
	}
	if globalConf.Model.EmbeddingDriver != "" {
		switch globalConf.Model.EmbeddingDriver {
		case model.OPENAI:
			if err := checkOpenAI("embedding"); err != nil {
				return err
			}
		case model.AZUREOPENAI:
			if err := checkAzureOpenAI("embedding"); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkOpenAI(modelType string) error {
	if globalConf.Model.OpenAI == nil {
		return errors.New("openai config is empty")
	}
	if globalConf.Model.OpenAI.ApiKey == "" {
		return errors.New("openai api key is empty")
	}
	if modelType == "llm" && globalConf.Model.OpenAI.LLM == "" {
		return errors.New("openai llm is empty")
	}
	if modelType == "embedding" && globalConf.Model.OpenAI.Embedding == "" {
		return errors.New("openai embedding is empty")
	}
	return nil
}

func checkAzureOpenAI(modelType string) error {
	if globalConf.Model.AzureOpenAI == nil {
		return errors.New("azure openai config is empty")
	}
	if globalConf.Model.AzureOpenAI.ApiKey == "" {
		return errors.New("azure openai api key is empty")
	}
	if globalConf.Model.AzureOpenAI.Endpoint == "" {
		return errors.New("azure openai endpoint is empty")
	}
	if modelType == "llm" && globalConf.Model.AzureOpenAI.LLM == "" {
		return errors.New("azure openai llm is empty")
	}
	if modelType == "embedding" && globalConf.Model.AzureOpenAI.Embedding == "" {
		return errors.New("azure openai embedding is empty")
	}
	return nil
}

func SetLLMModel(m *Model) {
	globalConf.Model.LLMDriver = m.LLMDriver
	switch m.LLMDriver {
	case model.OPENAI:
		globalConf.Model.OpenAI = m.OpenAI
	case model.AZUREOPENAI:
		globalConf.Model.AzureOpenAI = m.AzureOpenAI
	}
}

func SetEmbeddingModel(m *Model) {
	globalConf.Model.EmbeddingDriver = m.EmbeddingDriver
	switch m.EmbeddingDriver {
	case model.OPENAI:
		globalConf.Model.OpenAI = m.OpenAI
	case model.AZUREOPENAI:
		globalConf.Model.AzureOpenAI = m.AzureOpenAI
	}
}

func (m *Model) ToCfg(modelType string) model.Model {
	if m == nil {
		return model.Model{}
	}

	var driver string
	if modelType == "llm" {
		driver = m.LLMDriver
	} else if modelType == "embedding" {
		driver = m.EmbeddingDriver
	} else {
		return model.Model{}
	}

	switch driver {
	case model.OPENAI:
		return m.toOpenAICfg()
	case model.AZUREOPENAI:
		return m.toAzureOpenAICfg()
	default:
		return model.Model{}
	}
}

func (m *Model) toOpenAICfg() model.Model {
	return model.Model{
		Driver: model.OPENAI,
		OpenAI: model.OpenAI{
			ApiKey:         m.OpenAI.ApiKey,
			OrganizationID: m.OpenAI.OrganizationID,
			ApiBase:        m.OpenAI.ApiBase,
			LLM:            m.OpenAI.LLM,
			Embedding:      m.OpenAI.Embedding,
		},
	}
}

func (m *Model) toAzureOpenAICfg() model.Model {
	return model.Model{
		Driver: model.AZUREOPENAI,
		AzureOpenAI: model.AzureOpenAI{
			ApiKey:    m.AzureOpenAI.ApiKey,
			Endpoint:  m.AzureOpenAI.Endpoint,
			LLM:       m.AzureOpenAI.LLM,
			Embedding: m.AzureOpenAI.Embedding,
		},
	}
}
