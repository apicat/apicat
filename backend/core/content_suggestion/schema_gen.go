package content_suggestion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/core/content_suggestion/utils"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
	array_operation "github.com/apicat/apicat/v2/backend/utils/array"
)

type SchemaGenerator struct {
	projectID         string
	rootSchema        *jsonschema.Schema
	exceptType        string
	exceptID          int64
	focusParentKey    string
	focusParentSchema *jsonschema.Schema
	similarSchema     *jsonschema.Schema
	ctx               context.Context
}

func NewSchemaGenerator(projectID string, js *jsonschema.Schema, exceptType string, exceptID int64) (*SchemaGenerator, error) {
	if js == nil {
		return nil, errors.New("jsonschema is nil")
	}

	return &SchemaGenerator{
		projectID:  projectID,
		rootSchema: js,
		exceptType: exceptType,
		exceptID:   exceptID,
		ctx:        context.Background(),
	}, nil
}

func (sg *SchemaGenerator) Generate() (*jsonschema.Schema, error) {
	sg.focusParentKey, sg.focusParentSchema = sg.findFocusParentSchema("root", sg.rootSchema)
	if sg.focusParentKey == "" || sg.focusParentSchema == nil {
		return nil, nil
	}

	ids, err := sg.similaritySearch()
	if err != nil {
		return nil, err
	}

	if !sg.compareResult(ids) {
		return nil, nil
	}

	return sg.augmentSchema()
}

func (sg *SchemaGenerator) findFocusParentSchema(k string, js *jsonschema.Schema) (string, *jsonschema.Schema) {
	if js.Reference != nil {
		// reference schema is not supported
		slog.ErrorContext(sg.ctx, "findFocusParentSchema reference schema is not supported")
		return "", nil
	}

	if ok, err := sg.hasFocusSchema(js); ok && err == nil {
		return k, js
	}

	if len(js.AllOf) > 0 {
		for _, v := range js.AllOf {
			if targetK, targetJS := sg.findFocusParentSchema(k, v); targetK != "" && targetJS != nil {
				return targetK, targetJS
			}
		}
	}

	if len(js.AnyOf) > 0 {
		for i, v := range js.AnyOf {
			if targetK, targetJS := sg.findFocusParentSchema(fmt.Sprintf("type%d", i+1), v); targetK != "" && targetJS != nil {
				return targetK, targetJS
			}
		}
	}

	if len(js.OneOf) > 0 {
		for i, v := range js.OneOf {
			if targetK, targetJS := sg.findFocusParentSchema(fmt.Sprintf("type%d", i), v); targetK != "" && targetJS != nil {
				return targetK, targetJS
			}
		}
	}

	if js.Properties != nil {
		for pk, v := range js.Properties {
			if targetK, targetJS := sg.findFocusParentSchema(pk, v); targetK != "" && targetJS != nil {
				return targetK, targetJS
			}
		}
	}

	if js.Items != nil && !js.Items.IsBool() {
		if targetK, targetJS := sg.findFocusParentSchema("items", js.Items.Value()); targetK != "" && targetJS != nil {
			return targetK, targetJS
		}
	}

	return "", nil
}

func (sg *SchemaGenerator) hasFocusSchema(js *jsonschema.Schema) (bool, error) {
	if len(js.AllOf) > 0 {
		for _, v := range js.AllOf {
			if v.Reference != nil {
				return false, errors.New("reference schema is not supported")
			}
			if v.XFocus {
				return true, nil
			}
		}
	}

	if len(js.AnyOf) > 0 {
		for _, v := range js.AnyOf {
			if v.Reference != nil {
				return false, errors.New("reference schema is not supported")
			}
			if v.XFocus {
				return true, nil
			}
		}
	}

	if len(js.OneOf) > 0 {
		for _, v := range js.OneOf {
			if v.Reference != nil {
				return false, errors.New("reference schema is not supported")
			}
			if v.XFocus {
				return true, nil
			}
		}
	}

	if js.Properties != nil {
		for _, v := range js.Properties {
			if v.Reference != nil {
				return false, errors.New("reference schema is not supported")
			}
			if v.XFocus {
				return true, nil
			}
		}
	}

	if js.Items != nil && !js.Items.IsBool() {
		if js.Items.Value().Reference != nil {
			return false, errors.New("reference schema is not supported")
		}
		if js.Items.Value().XFocus {
			return true, nil
		}
	}

	return false, nil
}

func (sg *SchemaGenerator) similaritySearch() ([]map[string]int, error) {
	schemaText := utils.SchemaToText(sg.focusParentKey, sg.focusParentSchema)
	search, err := utils.NewSimilaritySearch(sg.projectID, schemaText)
	if err != nil {
		return nil, err
	}

	wantProperties := apiContentProperty{
		CollectionID:      1,
		DefinitionModelID: 1,
	}
	search.WithFields(wantProperties.GetPropertyNames())
	search.WithAdditionalFields([]string{"distance"})
	search.WithLimit(3)
	// search.WithDistance(0.5)
	result, err := search.Do()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, nil
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
		return nil, err
	}

	ids := make([]map[string]int, 0)
	for _, v := range searchResults {
		if v.CollectionID > 0 {
			if sg.exceptType == "collection" && int64(v.CollectionID) == sg.exceptID {
				continue
			}
			ids = append(ids, map[string]int{
				"collection_id": int(v.CollectionID),
			})
		} else if v.DefinitionModelID > 0 {
			if sg.exceptType == "model" && int64(v.DefinitionModelID) == sg.exceptID {
				continue
			}
			ids = append(ids, map[string]int{
				"definition_model_id": int(v.DefinitionModelID),
			})
		}
	}
	return ids, nil
}

func (sg *SchemaGenerator) compareResult(ids []map[string]int) bool {
	specDefinitions := &spec.Definitions{}
	var err error
	specDefinitions.Schemas, err = definition.GetDefinitionSchemasWithSpec(sg.projectID)
	if err != nil {
		slog.ErrorContext(sg.ctx, "definition.GetDefinitionSchemasWithSpec", "err", err)
		return false
	}
	specDefinitions.Responses, err = definition.GetDefinitionResponsesWithSpec(sg.projectID)
	if err != nil {
		slog.ErrorContext(sg.ctx, "definition.GetDefinitionResponsesWithSpec", "err", err)
		return false
	}
	specGlobalParameters, err := global.GetGlobalParametersWithSpec(sg.projectID)
	if err != nil {
		slog.ErrorContext(sg.ctx, "global.GetGlobalParametersWithSpec", "err", err)
		return false
	}

	skip := map[string]bool{
		"none":                     true,
		"raw":                      true,
		"application/octet-stream": true,
		"text/plain":               true,
		"text/html":                true,
	}

	for _, v := range ids {
		if collectionID, ok := v["collection_id"]; ok {
			c := &collection.Collection{ID: uint(collectionID), ProjectID: sg.projectID}
			if exist, err := c.Get(sg.ctx); err != nil || !exist {
				slog.ErrorContext(sg.ctx, "c.Get", "err", err)
				continue
			}
			if specContent, err := c.ContentToSpec(); err != nil {
				slog.ErrorContext(sg.ctx, "c.ContentToSpec", "err", err)
			} else {
				if err := specContent.DeepDerefAll(specGlobalParameters, specDefinitions); err != nil {
					slog.ErrorContext(sg.ctx, "specContent.DeepDerefAll", "err", err)
				}
				for _, node := range specContent {
					switch node.NodeType() {
					case spec.NODE_HTTP_REQUEST:
						req := node.ToHttpRequest()
						for contentType, v := range req.Attrs.Content {
							if _, ok := skip[contentType]; ok {
								continue
							}
							if v.Schema != nil {
								if sg.findSimilarSchema(v.Schema) {
									return true
								}
							}
						}
					case spec.NODE_HTTP_RESPONSE:
						res := node.ToHttpResponse()
						for _, r := range res.Attrs.List {
							for contentType, v := range r.Content {
								if _, ok := skip[contentType]; ok {
									continue
								}
								if v.Schema != nil {
									if sg.findSimilarSchema(v.Schema) {
										return true
									}
								}
							}
						}
					}
				}
			}
		} else if definitionModelID, ok := v["definition_model_id"]; ok {
			dm := &definition.DefinitionSchema{ID: uint(definitionModelID), ProjectID: sg.projectID}
			exist, err := dm.Get(sg.ctx)
			if err != nil || !exist {
				slog.ErrorContext(sg.ctx, "dm.Get", "err", err)
				continue
			}
			if specModel, err := dm.ToSpec(); err != nil {
				slog.ErrorContext(sg.ctx, "dm.ToSpec", "err", err)
			} else {
				if specModel.DeepDeref(specDefinitions.Schemas) != nil {
					slog.ErrorContext(sg.ctx, "specModel.DeepDeref", "err", err)
				}
				if sg.findSimilarSchema(specModel.Schema) {
					return true
				}
			}
		}
	}
	return false
}

func (sg *SchemaGenerator) findSimilarSchema(js *jsonschema.Schema) bool {
	var targetJS *jsonschema.Schema
	if js.Properties != nil {
		targetJS = js
	} else if js.AllOf != nil {
		if len(js.AllOf) > 1 {
			slog.ErrorContext(sg.ctx, "allOf is not merge", "schema", js)
			return false
		}
		if js.AllOf[0].Properties == nil {
			slog.ErrorContext(sg.ctx, "allOf children is invalid", "schema", js)
			return false
		}
		targetJS = js.AllOf[0]
	} else if js.Items != nil && !js.Items.IsBool() && js.Items.Value().Properties != nil {
		targetJS = js.Items.Value()
	} else {
		return false
	}

	if sg.findSimilar(targetJS) {
		return true
	}

	if targetJS.Properties != nil {
		for _, v := range targetJS.Properties {
			if sg.findSimilarSchema(v) {
				return true
			}
		}
	}
	return false
}

func (sg *SchemaGenerator) findSimilar(js *jsonschema.Schema) bool {
	if js.Properties == nil {
		return false
	}
	for k := range sg.focusParentSchema.Properties {
		if _, ok := js.Properties[k]; !ok {
			return false
		}
	}
	sg.similarSchema = js
	return true
}

func (sg *SchemaGenerator) augmentSchema() (*jsonschema.Schema, error) {
	if sg.similarSchema == nil || sg.focusParentSchema == nil {
		return nil, errors.New("similar schema or focus parent schema is nil")
	}

	result := jsonschema.NewSchema(jsonschema.T_OBJ)
	result.Properties = make(map[string]*jsonschema.Schema)
	result.XOrder = make([]string, 0)
	result.Required = make([]string, 0)

	for _, k := range sg.focusParentSchema.XOrder {
		if _, ok := sg.focusParentSchema.Properties[k]; !ok {
			continue
		}

		result.Properties[k] = sg.focusParentSchema.Properties[k]
		result.XOrder = append(result.XOrder, k)
		if array_operation.InArray(k, sg.focusParentSchema.Required) {
			result.Required = append(result.Required, k)
		}

		if sg.focusParentSchema.Properties[k].XFocus {
			for _, newK := range sg.similarSchema.XOrder {
				if _, ok := sg.similarSchema.Properties[newK]; !ok {
					continue
				}

				if _, ok := sg.focusParentSchema.Properties[newK]; ok {
					continue
				}

				result.Properties[newK] = sg.similarSchema.Properties[newK]
				result.Properties[newK].XSuggestion = true
				result.XOrder = append(result.XOrder, newK)
				if array_operation.InArray(newK, sg.similarSchema.Required) {
					result.Required = append(result.Required, newK)
				}
			}
		}
	}

	if len(result.Properties) < len(sg.focusParentSchema.Properties) {
		return nil, errors.New("properties is less than focus parent schema properties")
	}
	return result, nil
}
