package spec2

import (
	"encoding/json"

	"github.com/apicat/apicat/v2/backend/module/spec2/jsonschema"
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

func NewSchemaFromJson(str string) (*Model, error) {
	s := &Model{}
	return s, json.Unmarshal([]byte(str), s)
}

func (s *Model) Deref(wantToDeref ...*jsonschema.Schema) {
}
