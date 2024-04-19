package llm

import (
	"errors"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/module/llm/common"
	"github.com/apicat/apicat/v2/backend/module/llm/openai"
)

const (
	OPENAI      = "openai"
	AZUREOPENAI = "azure-openai"
)

type LLM struct {
	Driver      string
	OpenAI      openai.OpenAI
	AzureOpenAI openai.AzureOpenAI
}

func NewLLM(cfg LLM) (common.Provider, error) {
	slog.Debug("llm.NewLLM", "cfg", cfg)

	if cfg.Driver == OPENAI {
		if o := openai.NewOpenAI(cfg.OpenAI); o != nil {
			return o, nil
		} else {
			return nil, errors.New("openai.NewOpenAI failed")
		}
	} else if cfg.Driver == AZUREOPENAI {
		if o := openai.NewAzureOpenAI(cfg.AzureOpenAI); o != nil {
			return o, nil
		} else {
			return nil, errors.New("openai.NewAzureOpenAI failed")
		}
	}
	return nil, errors.New("llm driver not found")
}
