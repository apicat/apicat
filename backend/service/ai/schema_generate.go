package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"

	"github.com/apicat/apicat/v2/backend/module/llm"
	llmcommon "github.com/apicat/apicat/v2/backend/module/llm/common"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
	"github.com/gin-gonic/gin"
)

func SchemaGenerate(ctx *gin.Context, prompt string) (*definition.DefinitionSchema, error) {
	tpl := NewTpl("schema_generate.tmpl", jwt.GetUser(ctx).Language, prompt)
	messages, err := tpl.Prompt()
	if err != nil {
		return nil, err
	}

	a, err := llm.NewLLM(config.Get().LLM.ToMapInterface())
	if err != nil {
		return nil, err
	}

	result, err := a.ChatCompletionRequest(&llmcommon.ChatCompletionRequest{
		Temperature: 0.3,
		MaxTokens:   3000,
		Messages:    messages,
	})
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, errors.New("schema generate failed")
	}

	result = strings.TrimSuffix(result, "```")
	jsonschema := &jsonschema.Schema{}
	if err := json.Unmarshal([]byte(result), jsonschema); err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed: %s", err.Error())
	}

	return &definition.DefinitionSchema{
		Name:        jsonschema.Title,
		Description: jsonschema.Description,
		Type:        "schema",
		Schema:      strings.TrimSpace(result),
	}, nil
}
