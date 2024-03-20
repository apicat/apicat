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
