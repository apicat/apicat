package rag

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/core/prompt"
	"github.com/apicat/apicat/v2/backend/core/rag/utils"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/plugin/openapi"
	"github.com/apicat/apicat/v2/backend/module/vector"
)

type CollectionGenerator struct {
	collection        *collection.Collection
	similarCollection *collection.Collection
	openapiContent    string
	language          string
	ctx               context.Context
}

func NewCollectionGenerator(c *collection.Collection, lang string) (*CollectionGenerator, error) {
	if c == nil {
		return nil, errors.New("collection is nil")
	}

	return &CollectionGenerator{
		collection: c,
		language:   lang,
		ctx:        context.Background(),
	}, nil
}

func (cg *CollectionGenerator) Generate() (*collection.Collection, error) {
	if len(cg.collection.Title) < 5 {
		// Reference: login are 5 characters long
		return nil, errors.New("title is too short")
	}

	cid, err := cg.similaritySearch()
	if err != nil {
		return nil, err
	}
	if cid > 0 {
		if err := cg.getSimilarityContent(cid); err != nil {
			return nil, err
		}
	}

	return cg.gen()
}

func (cg *CollectionGenerator) similaritySearch() (int, error) {
	var text string
	if cg.collection.Path == "" {
		text = cg.collection.Title
	} else {
		text = fmt.Sprintf("%s\n%s", cg.collection.Title, cg.collection.Path)
	}

	search, err := utils.NewSimilaritySearch(cg.collection.ProjectID, text)
	if err != nil {
		return 0, err
	}

	wantProperties := apiContentProperty{
		CollectionID: 1,
	}
	fields := wantProperties.GetPropertyNames()
	search.WithFields(fields)
	search.WithAdditionalFields([]string{"distance"})
	search.WithLimit(1)
	// search.WithDistance(0.5)

	where := &vector.WhereCondition{
		PropertyName: fields[0],
		Operator:     ">",
		Value:        vector.T_INT(0),
	}
	search.WithWhere([]*vector.WhereCondition{where})
	result, err := search.Do()
	if err != nil {
		return 0, err
	}
	if result == "" {
		return 0, nil
	}

	type additionalField struct {
		Distance float32 `json:"distance"`
	}
	type searchResult struct {
		apiContentProperty
		additionalField `json:"_additional"`
	}
	var searchResults []searchResult
	if err := json.Unmarshal([]byte(result), &searchResults); err != nil {
		return 0, err
	}

	for _, v := range searchResults {
		if v.CollectionID > 0 {
			return int(v.CollectionID), nil
		}
	}
	return 0, nil
}

func (cg *CollectionGenerator) getSimilarityContent(cid int) error {
	cg.similarCollection = &collection.Collection{ID: uint(cid), ProjectID: cg.collection.ProjectID}
	if exist, err := cg.similarCollection.Get(cg.ctx); err != nil || !exist {
		if err != nil {
			return err
		}
		return errors.New("collection not exist")
	}

	specDefinitions := &spec.Definitions{}
	var err error
	specDefinitions.Schemas, err = definition.GetDefinitionSchemasWithSpec(cg.ctx, cg.collection.ProjectID)
	if err != nil {
		return err
	}
	specDefinitions.Responses, err = definition.GetDefinitionResponsesWithSpec(cg.ctx, cg.collection.ProjectID)
	if err != nil {
		return err
	}
	specGlobalParameters, err := global.GetGlobalParametersWithSpec(cg.ctx, cg.collection.ProjectID)
	if err != nil {
		return err
	}

	if specCollection, err := cg.similarCollection.ToSpec(); err != nil {
		slog.ErrorContext(cg.ctx, "c.ToSpec", "err", err)
	} else {
		if err := specCollection.Content.DeepDerefAll(specGlobalParameters, specDefinitions); err != nil {
			return err
		}
		specDoc := spec.NewEmptySpec()
		specDoc.Collections = append(specDoc.Collections, specCollection)
		if openapiContent, err := openapi.Generate(specDoc, "3.0.0", "yaml"); err != nil {
			return err
		} else {
			cg.openapiContent = string(openapiContent)
		}
	}
	return nil
}

func (cg *CollectionGenerator) gen() (*collection.Collection, error) {
	infomation := map[string]string{
		"title": cg.collection.Title,
		"path":  cg.collection.Path,
	}

	tmpl := "api_generate.tmpl"
	if cg.openapiContent != "" {
		infomation["demoTitle"] = cg.similarCollection.Title
		infomation["demoPath"] = cg.similarCollection.Path
		infomation["demoAPI"] = cg.openapiContent
		tmpl = "api_generate_similar.tmpl"
	}

	builder := prompt.NewPrompt(tmpl, cg.language, infomation)
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
		return nil, errors.New("API generate failed")
	}

	re := regexp.MustCompile("(?s)```yaml\\n(.*?)```")
	matches := re.FindAllStringSubmatch(result, 1)
	if len(matches) == 0 {
		return nil, errors.New("API generate failed")
	}

	apiSpec, err := openapi.Parse([]byte(matches[0][1]))
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
