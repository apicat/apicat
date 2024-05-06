package spec2

import (
	"encoding/json"
	"strconv"

	"github.com/apicat/apicat/v2/backend/module/spec2/jsonschema"
)

const (
	ModelTypeDir    = "category"
	ModelTypeSchema = "schema"
)

type Model struct {
	ID          int64              `json:"id,omitempty" yaml:"id,omitempty"`
	ParentId    uint64             `json:"parentid,omitempty" yaml:"parentid,omitempty"`
	Name        string             `json:"name,omitempty" yaml:"name,omitempty"`
	Type        string             `json:"type,omitempty" yaml:"type,omitempty"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	Schema      *jsonschema.Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	Items       Models             `json:"items,omitempty" yaml:"items,omitempty"`
}

type Models []*Model

func NewModelFromJson(str string) (*Model, error) {
	s := &Model{}
	if err := json.Unmarshal([]byte(str), s); err != nil {
		return nil, err
	}
	s.Schema.ID = s.ID
	return s, nil
}

// The id referenced by this model itself
func (s *Model) RefID() int64 {
	return s.Schema.GetRefID()
}

// The id referenced by the model itself and its child elements
func (s *Model) RefIDs() []int64 {
	return s.Schema.DeepGetRefID()
}

// Model(s) refers to Model(ref), removes the reference relationship between s and ref, and replaces the content of ref into s.
func (s *Model) Deref(ref *Model) {
	if s == nil || ref == nil {
		return
	}

	s.Schema.ID = s.ID
	ref.Schema.ID = ref.ID

	refSchemas := s.Schema.DeepFindRefById(strconv.FormatInt(ref.ID, 10))
	if len(refSchemas) > 0 {
		for _, schema := range refSchemas {
			schema.ReplaceRef(ref.Schema)
		}
	}
}

func (s *Model) DelRef(ref *Model) {
	if s == nil || ref == nil {
		return
	}

	s.Schema.ID = s.ID
	ref.Schema.ID = ref.ID

	if s.Schema.Ref() {
		s.Schema.DelRootRef(ref.Schema)
	}
	s.Schema.DelChildrenRef(ref.Schema)
}

func (s *Models) FindByName(name string) *Model {
	if s == nil {
		return nil
	}

	for _, v := range *s {
		if v.Type == ModelTypeDir {
			return v.Items.FindByName(name)
		} else {
			if v.Name == name {
				return v
			}
		}
	}
	return nil
}

func (s *Models) FindByID(id int64) *Model {
	if s == nil {
		return nil
	}

	for _, v := range *s {
		if v.Type == ModelTypeDir {
			return v.Items.FindByID(id)
		} else {
			if id == v.ID {
				return v
			}
		}
	}
	return nil
}

func (s *Models) DelByID(id int64) {
	if s == nil {
		return
	}
	for i, v := range *s {
		if v.Type == ModelTypeDir {
			v.Items.DelByID(id)
		} else {
			if v.ID == id {
				*s = append((*s)[:i], (*s)[i+1:]...)
			}
		}
	}
}
