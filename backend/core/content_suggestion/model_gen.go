package content_suggestion

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/core/prompt"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

type ModelGenerator struct {
	model    *definition.DefinitionSchema
	language string
	ctx      context.Context
}

func NewModelGenerator(m *definition.DefinitionSchema, lang string) (*ModelGenerator, error) {
	if m == nil {
		return nil, errors.New("model is nil")
	}

	return &ModelGenerator{
		model:    m,
		language: lang,
		ctx:      context.Background(),
	}, nil
}

func (mg *ModelGenerator) Generate() (*definition.DefinitionSchema, error) {
	if len(mg.model.Name) < 5 {
		return nil, errors.New("name is too short")
	}

	builder := prompt.NewPrompt("model_generate.tmpl", mg.language, mg.model.Name)
	messages, err := builder.Prompt()
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
		return nil, errors.New("model generate failed")
	}

	re := regexp.MustCompile("(?s)```json\\n(.*?)```")
	matches := re.FindAllStringSubmatch(result, 1)
	if len(matches) == 0 {
		return nil, errors.New("model generate failed")
	}

	jsonschema := &jsonschema.Schema{}
	if err := json.Unmarshal([]byte(matches[0][1]), jsonschema); err != nil {
		return nil, err
	}

	return &definition.DefinitionSchema{
		Name:        jsonschema.Title,
		Description: jsonschema.Description,
		Type:        definition.SchemaSchema,
		Schema:      strings.TrimSpace(matches[0][1]),
	}, nil
}
