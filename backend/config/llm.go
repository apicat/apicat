package config

import (
	"github.com/apicat/apicat/v2/backend/module/llm"
	"github.com/apicat/apicat/v2/backend/module/llm/openai"
)

type LLM struct {
	Driver      string       `yaml:"Driver"`
	OpenAI      *OpenAI      `yaml:"OpenAI"`
	AzureOpenAI *AzureOpenAI `yaml:"AzureOpenAI"`
}

type OpenAI struct {
	ApiKey         string `yaml:"ApiKey" json:"apiKey"`
	OrganizationID string `yaml:"OrganizationID" json:"organizationID"`
	ApiBase        string `yaml:"ApiBase" json:"apiBase"`
	LLMName        string `yaml:"LLMName" json:"llmName"`
	EmbeddingName  string `yaml:"EmbeddingName" json:"embeddingName"`
	Timeout        int    `yaml:"Timeout" json:"timeout"`
}

type AzureOpenAI struct {
	ApiKey        string `yaml:"ApiKey" json:"apiKey"`
	Endpoint      string `yaml:"Endpoint" json:"endpoint"`
	LLMName       string `yaml:"LLMName" json:"llmName"`
	EmbeddingName string `yaml:"EmbeddingName" json:"embeddingName"`
	Timeout       int    `yaml:"Timeout" json:"timeout"`
}

func SetLLM(c *LLM) {
	globalConf.LLM = c
}

func (l *LLM) ToCfg() llm.LLM {
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
	}
	return llm.LLM{}
}
