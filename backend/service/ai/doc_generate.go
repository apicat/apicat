package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"

	"github.com/apicat/apicat/v2/backend/module/llm"
	llmcommon "github.com/apicat/apicat/v2/backend/module/llm/common"
	"github.com/apicat/apicat/v2/backend/module/spec/plugin/openapi"
	"github.com/gin-gonic/gin"
)

func DocGenerate(ctx *gin.Context, prompt string) (*collection.Collection, error) {
	tpl := NewTpl("api_generate.tmpl", jwt.GetUser(ctx).Language, prompt)
	messages, err := tpl.Prompt()
	if err != nil {
		return nil, err
	}

	a, err := llm.NewLLM(config.Get().LLM.ToCfg())
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
		return nil, errors.New("API generate failed")
	}

	result = strings.TrimSuffix(result, "```")
	apiSpec, err := openapi.Decode([]byte(result))
	if err != nil {
		return nil, fmt.Errorf("openapi.Decode failed: %s content:\n%s", err.Error(), result)
	}

	if len(apiSpec.Collections) == 0 {
		return nil, fmt.Errorf("collection not found, original openapi content:\n%s", result)
	}

	if err := apiSpec.Collections[0].DerefSchema(apiSpec.Definitions.Schemas...); err != nil {
		return nil, fmt.Errorf("DerefSchema failed: %s", err.Error())
	}
	if err := apiSpec.Collections[0].DerefResponse(apiSpec.Definitions.Responses...); err != nil {
		return nil, fmt.Errorf("DerefResponse failed: %s", err.Error())
	}

	r, err := json.Marshal(apiSpec.Collections[0].Content)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %s", err.Error())
	}

	return &collection.Collection{
		Title:   apiSpec.Collections[0].Title,
		Type:    string(apiSpec.Collections[0].Type),
		Content: string(r),
	}, nil
}
