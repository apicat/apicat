package spec

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

const TYPE_MODEL = "schema"

type DefinitionModel struct {
	ID          int64              `json:"id,omitempty" yaml:"id,omitempty"`
	ParentId    uint64             `json:"parentid,omitempty" yaml:"parentid,omitempty"`
	Name        string             `json:"name,omitempty" yaml:"name,omitempty"`
	Type        string             `json:"type,omitempty" yaml:"type,omitempty"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	Schema      *jsonschema.Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	Items       DefinitionModels   `json:"items,omitempty" yaml:"items,omitempty"`
}

type DefinitionModels []*DefinitionModel

func NewModelFromJson(str string) (*DefinitionModel, error) {
	if str == "" {
		return nil, errors.New("empty json content")
	}
	s := &DefinitionModel{}
	if err := json.Unmarshal([]byte(str), s); err != nil {
		return nil, err
	}
	s.Schema.ID = s.ID
	return s, nil
}

// The id referenced by the model itself and its child elements
func (s *DefinitionModel) RefIDs() []int64 {
	return s.Schema.DeepGetRefID()
}

// Model(s) refers to Model(ref), removes the reference relationship between s and ref, and replaces the content of ref into s.
func (s *DefinitionModel) Deref(ref *DefinitionModel) error {
	if ref == nil {
		return errors.New("model is nil")
	}

	s.Schema.ID = s.ID
	ref.Schema.ID = ref.ID

	refSchemas := s.Schema.DeepFindRefById(strconv.FormatInt(ref.ID, 10))
	if len(refSchemas) > 0 {
		for _, schema := range refSchemas {
			if err := schema.ReplaceRef(ref.Schema); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *DefinitionModel) DeepDeref(refs DefinitionModels) error {
	if len(refs) == 0 {
		return nil
	}

	helper := jsonschema.NewDerefHelper(refs.ToJsonSchemaMap())
	new, err := helper.DeepDeref(s.Schema)
	if err != nil {
		return err
	}
	s.Schema = &new
	return nil
}

func (s *DefinitionModel) DelRef(ref *DefinitionModel) {
	if ref == nil {
		return
	}

	s.Schema.ID = s.ID
	ref.Schema.ID = ref.ID
	s.Schema.DelRef(ref.Schema)
}

func (s *DefinitionModel) SetXDiff(x string) {
	if s.Schema != nil {
		s.Schema.SetXDiff(x)
	}
}

func (s *DefinitionModel) ItemsTreeToList() DefinitionModels {
	list := make(DefinitionModels, 0)

	if s.Type == TYPE_MODEL {
		list = append(list, s)
	}

	if s.Items != nil && len(s.Items) > 0 {
		for _, v := range s.Items {
			list = append(list, v.ItemsTreeToList()...)
		}
	}
	return list
}

func (s *DefinitionModels) FindByName(name string) *DefinitionModel {
	for _, v := range *s {
		if v.Type == TYPE_CATEGORY {
			return v.Items.FindByName(name)
		}
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (s *DefinitionModels) FindByID(id int64) *DefinitionModel {
	for _, v := range *s {
		if v.Type == TYPE_CATEGORY {
			return v.Items.FindByID(id)
		}
		if id == v.ID {
			return v
		}
	}
	return nil
}

func (s *DefinitionModels) DelByID(id int64) {
	for i, v := range *s {
		if v.Type == TYPE_CATEGORY {
			v.Items.DelByID(id)
		} else {
			if v.ID == id {
				*s = append((*s)[:i], (*s)[i+1:]...)
				return
			}
		}
	}
}

func (s *DefinitionModels) RemoveDir() DefinitionModels {
	result := DefinitionModels{}
	for _, v := range *s {
		if v.Type == TYPE_CATEGORY {
			result = append(result, v.Items.RemoveDir()...)
		}
		result = append(result, v)
	}
	return result
}

func (s *DefinitionModels) ToJsonSchemaMap() map[int64]*jsonschema.Schema {
	result := make(map[int64]*jsonschema.Schema)
	for _, v := range *s {
		if v.Type == TYPE_MODEL {
			v.Schema.ID = v.ID
			result[v.ID] = v.Schema
		}
	}
	return result
}
