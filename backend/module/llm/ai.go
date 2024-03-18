package llm

import (
	"apicat-cloud/backend/module/llm/common"
	"apicat-cloud/backend/module/llm/openai"
	"errors"
	"log/slog"
)

const (
	OPENAI      = "openai"
	AZUREOPENAI = "azure-openai"
)

func NewLLM(cfg map[string]interface{}) (common.Provider, error) {
	slog.Debug("llm.NewLLM", "cfg", cfg)
	if cfg == nil {
		return nil, errors.New("llm config is nil")
	}

	if cfg["Driver"] == OPENAI {
		return openai.NewOpenAI(cfg["OpenAI"].(map[string]interface{}))
	} else if cfg["Driver"] == AZUREOPENAI {
		return openai.NewOpenAI(cfg["AzureOpenAI"].(map[string]interface{}))
	}
	return nil, errors.New("llm driver not found")
}
