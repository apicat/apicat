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
