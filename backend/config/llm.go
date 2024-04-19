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
	ApiKey         string `yaml:"ApiKey"`
	OrganizationID string `yaml:"OrganizationID"`
	ApiBase        string `yaml:"ApiBase"`
	LLMName        string `yaml:"LLMName"`
	EmbeddingName  string `yaml:"EmbeddingName"`
	Timeout        int    `yaml:"Timeout"`
}

type AzureOpenAI struct {
	ApiKey        string `yaml:"ApiKey"`
	Endpoint      string `yaml:"Endpoint"`
	LLMName       string `yaml:"LLMName"`
	EmbeddingName string `yaml:"EmbeddingName"`
	Timeout       int    `yaml:"Timeout"`
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
