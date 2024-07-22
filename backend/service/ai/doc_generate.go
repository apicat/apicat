package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"

	"github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/spec/plugin/openapi"
	"github.com/gin-gonic/gin"
)

func DocGenerate(ctx *gin.Context, prompt string) (*collection.Collection, error) {
	tpl := NewTpl("api_generate.tmpl", jwt.GetUser(ctx).Language, prompt)
	messages, err := tpl.Prompt()
	if err != nil {
		return nil, err
	}

	m, err := model.NewModel(config.GetModel().ToModuleStruct("llm"))
	if err != nil {
		return nil, err
	}

	result, err := m.ChatCompletionRequest(model.NewChatCompletionOption(messages))
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, errors.New("API generate failed")
	}

	result = strings.TrimSuffix(result, "```")
	apiSpec, err := openapi.Parse([]byte(result))
	if err != nil {
		return nil, fmt.Errorf("openapi.Parse failed: %s content:\n%s", err.Error(), result)
	}

	if len(apiSpec.Collections) == 0 {
		return nil, fmt.Errorf("collection not found, original openapi content:\n%s", result)
	}

	if err := apiSpec.Collections[0].Content.DeepDerefAll(apiSpec.Globals.Parameters, apiSpec.Definitions); err != nil {
		return nil, fmt.Errorf("collection content DeepDerefAll failed: %s", err.Error())
	}
	apiSpec.Collections[0].Content.SortResponses()

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
