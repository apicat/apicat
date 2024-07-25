package model

import (
	"errors"
)

const (
	OPENAI       = "openai"
	AZURE_OPENAI = "azure-openai"
)

type Model struct {
	Driver      string
	OpenAI      OpenAI
	AzureOpenAI AzureOpenAI
}

func NewModel(cfg Model) (Provider, error) {
	if cfg.Driver == OPENAI {
		if o := NewOpenAI(cfg.OpenAI); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewOpenAI failed")
		}
	} else if cfg.Driver == AZURE_OPENAI {
		if o := NewAzureOpenAI(cfg.AzureOpenAI); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewAzureOpenAI failed")
		}
	}

	return nil, errors.New("model driver not found")
}

func ModelAvailable(driver, modelType, modelName string) bool {
	switch driver {
	case OPENAI, AZURE_OPENAI:
		switch modelType {
		case "llm":
			for _, v := range OPENAI_LLM_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		case "embedding":
			for _, v := range OPENAI_EMBEDDING_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		}
		return false
	}
	return false
}
