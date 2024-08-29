package rag

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/core/rag/utils"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
	"github.com/apicat/apicat/v2/backend/module/vector"
	array_operation "github.com/apicat/apicat/v2/backend/utils/array"
)

type ReferenceMatcher struct {
	projectID      string
	matchs         map[uint]*jsonschema.Schema
	embeddingModel model.Provider
	vectorDB       vector.VectorApi
	ctx            context.Context
}

func NewReferenceMatcher(ctx context.Context, projectID string) (*ReferenceMatcher, error) {
	embeddingModel, err := model.NewModel(config.GetModel().ToModuleStruct("embedding"))
	if err != nil {
		slog.ErrorContext(ctx, "model.NewModel", "err", err)
		return nil, err
	}

	vectorDB, err := vector.NewVector(config.GetVector().ToModuleStruct())
	if err != nil {
		slog.ErrorContext(ctx, "vector.NewVector", "err", err)
		return nil, err
	}

	if ok, _ := vectorDB.CheckCollectionExist(projectID); !ok {
		if err := vectorDB.CreateCollection(projectID, getAPIContentProperties()); err != nil {
			slog.ErrorContext(ctx, "vectorDB.CreateCollection", "err", err)
			return nil, err
		}
	}

	return &ReferenceMatcher{
		projectID:      projectID,
		matchs:         make(map[uint]*jsonschema.Schema),
		embeddingModel: embeddingModel,
		vectorDB:       vectorDB,
		ctx:            ctx,
	}, nil
}

func (rm *ReferenceMatcher) Match(title string, schema *jsonschema.Schema) (*jsonschema.Schema, error) {
	if schema == nil {
		return nil, errors.New("schema is nil")
	}

	if matchIDs, err := rm.similaritySearch(title, schema); err != nil {
		slog.ErrorContext(rm.ctx, "rm.similaritySearch", "err", err)
		return nil, err
	} else {
		if len(matchIDs) == 0 {
			return nil, nil
		}

		if refIDs := schema.DeepGetRefID(); len(refIDs) > 0 {
			for _, refID := range refIDs {
				if array_operation.InArray(uint(refID), matchIDs) {
					array_operation.Remove(uint(refID), matchIDs)
				}
			}
			if len(matchIDs) == 0 {
				return nil, nil
			}
		}

		if err := rm.getJsonSchema(matchIDs); err != nil {
			slog.ErrorContext(rm.ctx, "rm.getJsonSchema", "err", err)
			return nil, err
		}
	}

	matched := false
	for id := range rm.matchs {
		if rm.match(schema, id) {
			matched = true
		}
	}

	if matched {
		return schema, nil
	}
	return nil, nil
}

func (rm *ReferenceMatcher) similaritySearch(title string, s *jsonschema.Schema) ([]uint, error) {
	textList := make([]string, 0)

	if title != "" {
		textList = append(textList, title)
	}

	textList = append(textList, utils.SchemaToTextList("root", s)...)
	schemaText := strings.Join(textList, "\n")

	search, err := utils.NewSimilaritySearch(rm.projectID, schemaText)
	if err != nil {
		return nil, err
	}

	wantProperties := apiContentProperty{
		DefinitionModelID: 1,
	}
	fields := wantProperties.GetPropertyNames()
	search.WithFields(fields)
	search.WithAdditionalFields([]string{"distance"})
	search.WithLimit(3)
	// search.WithDistance(0.5)

	where := &vector.WhereCondition{
		PropertyName: fields[0],
		Operator:     ">",
		Value:        vector.T_INT(0),
	}
	search.WithWhere([]*vector.WhereCondition{where})

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

	ids := make([]uint, 0)
	for _, v := range searchResults {
		if v.DefinitionModelID > 0 {
			ids = append(ids, uint(v.DefinitionModelID))
		}
	}
	return ids, nil
}

func (rm *ReferenceMatcher) getJsonSchema(ids []uint) error {
	for _, id := range ids {
		dm := &definition.DefinitionSchema{ID: id, ProjectID: rm.projectID}
		exist, err := dm.Get(rm.ctx)
		if err != nil || !exist {
			slog.ErrorContext(rm.ctx, "dm.Get", "err", err)
			continue
		}
		if specModel, err := dm.ToSpec(); err != nil {
			slog.ErrorContext(rm.ctx, "dm.ToSpec", "err", err)
		} else {
			rm.matchs[id] = specModel.Schema
		}
	}

	return nil
}

func (rm *ReferenceMatcher) match(s *jsonschema.Schema, id uint) bool {
	if rm.contains(s, rm.matchs[id]) {
		if err := rm.replace(s, id); err != nil {
			slog.ErrorContext(rm.ctx, "rm.replace", "err", err)
			return false
		}
		return true
	}

	if len(s.AllOf) > 0 {
		for _, a := range s.AllOf {
			if rm.match(a, id) {
				return true
			}
		}
	}

	if len(s.AnyOf) > 0 {
		for _, a := range s.AnyOf {
			if rm.match(a, id) {
				return true
			}
		}
	}

	if len(s.OneOf) > 0 {
		for _, a := range s.OneOf {
			if rm.match(a, id) {
				return true
			}
		}
	}

	if s.Properties != nil {
		for _, p := range s.Properties {
			if rm.match(p, id) {
				return true
			}
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		if rm.match(s.Items.Value(), id) {
			return true
		}
	}

	return false
}

// if a contains b returns true else false
func (rm *ReferenceMatcher) contains(a *jsonschema.Schema, b *jsonschema.Schema) bool {
	if a == nil || (a.Properties == nil && len(a.AllOf) == 0) {
		return false
	}
	if b == nil || (b.Properties == nil && len(b.AllOf) == 0) {
		return false
	}

	aProperties := make(map[string]*jsonschema.Schema)
	if a.Properties != nil {
		for k, v := range a.Properties {
			aProperties[k] = v
		}
	}
	if len(a.AllOf) > 0 {
		for _, v := range a.AllOf {
			if v.Properties != nil {
				for k := range v.Properties {
					aProperties[k] = v.Properties[k]
				}
			}
			if refID, err := v.GetRefID(); err == nil && refID > 0 {
				aProperties[strconv.FormatInt(refID, 10)] = v
			}
		}
	}

	bProperties := make(map[string]*jsonschema.Schema)
	if b.Properties != nil {
		for k, v := range b.Properties {
			bProperties[k] = v
		}
	}
	if len(b.AllOf) > 0 {
		for _, v := range b.AllOf {
			if v.Properties != nil {
				for k := range v.Properties {
					bProperties[k] = v.Properties[k]
				}
			}
			if refID, err := v.GetRefID(); err == nil && refID > 0 {
				bProperties[strconv.FormatInt(refID, 10)] = v
			}
		}
	}

	count := 0
	for k, v := range bProperties {
		if av, ok := aProperties[k]; !ok {
			return false
		} else {
			if !av.Equal(v) {
				return false
			}
			count++
		}
	}
	return count == len(bProperties)
}

func (rm *ReferenceMatcher) replace(a *jsonschema.Schema, refID uint) error {
	if _, ok := rm.matchs[refID]; !ok {
		return errors.New("refID not exist")
	}

	refKeys := make([]string, 0)
	if rm.matchs[refID].Properties != nil {
		for k := range rm.matchs[refID].Properties {
			refKeys = append(refKeys, k)
		}
	}
	if len(rm.matchs[refID].AllOf) > 0 {
		for _, v := range rm.matchs[refID].AllOf {
			if v.Properties != nil {
				for k := range v.Properties {
					refKeys = append(refKeys, k)
				}
			}
			if refID, err := v.GetRefID(); err == nil && refID > 0 {
				refKeys = append(refKeys, strconv.FormatInt(refID, 10))
			}
		}
	}

	new := &jsonschema.Schema{}
	new.SetDefinitionModelRef(strconv.FormatUint(uint64(refID), 10))
	new.XSuggestion = true

	if len(a.AllOf) > 0 {
		for _, v := range a.AllOf {
			if v.Properties != nil {
				for _, k := range v.XOrder {
					if _, ok := v.Properties[k]; ok {
						if !array_operation.InArray(k, refKeys) {
							if new.AllOf == nil {
								new.AllOf = make([]*jsonschema.Schema, 0)
							}

							tmp := jsonschema.NewSchema(jsonschema.T_OBJ)
							tmp.Properties = map[string]*jsonschema.Schema{
								k: v.Properties[k],
							}
							tmp.XOrder = []string{k}
							if array_operation.InArray(k, v.Required) {
								tmp.Required = []string{k}
							}
							new.AllOf = append(new.AllOf, tmp)
						} else {
							// new.Reference 初始化的时候不为 nil，如果 new.Reference 为 nil，说明 Reference 被移动到了 allOf 中
							// new.AllOf 默认是 nil，如果 new.AllOf 不为 nil，说明要被替换的 Schema a 中的参数多于对应 ref model 中的参数
							if new.Reference == nil || new.AllOf == nil {
								continue
							}
							tmp := &jsonschema.Schema{}
							tmp.SetDefinitionModelRef(strconv.FormatUint(uint64(refID), 10))
							tmp.XSuggestion = true
							new.AllOf = append(new.AllOf, tmp)
							new.Reference = nil
							new.XSuggestion = false
						}
					}
				}
			}
			if refID, err := v.GetRefID(); err == nil && refID > 0 {
				if !array_operation.InArray(strconv.FormatInt(refID, 10), refKeys) {
					if new.AllOf == nil {
						new.AllOf = make([]*jsonschema.Schema, 0)
					}
					new.AllOf = append(new.AllOf, v)
				} else {
					if new.Reference == nil || new.AllOf == nil {
						continue
					}
					tmp := &jsonschema.Schema{}
					tmp.SetDefinitionModelRef(strconv.FormatUint(uint64(refID), 10))
					tmp.XSuggestion = true
					new.AllOf = append(new.AllOf, tmp)
					new.Reference = nil
					new.XSuggestion = false
				}
			}
		}
	}

	if a.Properties != nil {
		for _, k := range a.XOrder {
			if _, ok := a.Properties[k]; ok {
				if !array_operation.InArray(k, refKeys) {
					if new.Properties == nil {
						new.Properties = make(map[string]*jsonschema.Schema)
					}

					tmp := jsonschema.NewSchema(jsonschema.T_OBJ)
					tmp.Properties = map[string]*jsonschema.Schema{
						k: a.Properties[k],
					}
					tmp.XOrder = []string{k}
					if array_operation.InArray(k, a.Required) {
						tmp.Required = []string{k}
					}
					new.AllOf = append(new.AllOf, tmp)
				} else {
					if new.Reference == nil || new.AllOf == nil {
						continue
					}
					tmp := &jsonschema.Schema{}
					tmp.SetDefinitionModelRef(strconv.FormatUint(uint64(refID), 10))
					tmp.XSuggestion = true
					new.AllOf = append(new.AllOf, tmp)
					new.Reference = nil
					new.XSuggestion = false
				}
			}
		}
	}

	*a = *new
	return nil
}
