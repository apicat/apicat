package model

import (
	"errors"
)

const (
	OPENAI       = "openai"
	AZURE_OPENAI = "azure-openai"
	BAICHUAN     = "baichuan"
	MOONSHOT     = "moonshot"
)

type Model struct {
	Driver      string
	OpenAI      OpenAI
	AzureOpenAI AzureOpenAI
	Baichuan    Baichuan
	Moonshot    Moonshot
}

func NewModel(cfg Model) (Provider, error) {
	if cfg.Driver == OPENAI {
		if o := newOpenAI(cfg.OpenAI); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewOpenAI failed")
		}
	} else if cfg.Driver == AZURE_OPENAI {
		if o := newAzureOpenAI(cfg.AzureOpenAI); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewAzureOpenAI failed")
		}
	} else if cfg.Driver == BAICHUAN {
		if o := newBaichuan(cfg.Baichuan); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewBaichuan failed")
		}
	} else if cfg.Driver == MOONSHOT {
		if o := newMoonshot(cfg.Moonshot); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewMoonshot failed")
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
	case BAICHUAN:
		switch modelType {
		case "llm":
			for _, v := range BAICHUAN_LLM_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		case "embedding":
			for _, v := range BAICHUAN_EMBEDDING_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		}
		return false
	case MOONSHOT:
		switch modelType {
		case "llm":
			for _, v := range MOONSHOT_LLM_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		case "embedding":
			for _, v := range MOONSHOT_EMBEDDING_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		}
		return false
	}
	return false
}
