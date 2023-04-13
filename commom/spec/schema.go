package spec

import "github.com/apicat/apicat/commom/spec/jsonschema"

type Schema struct {
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Required    bool               `json:"required,omitempty"`
	Schema      *jsonschema.Schema `json:"schema,omitempty"`
}

type Schemas []*Schema

func (s *Schemas) Lookup(name string) *Schema {
	for _, v := range *s {
		if name == v.Name {
			return v
		}
	}
	return nil
}

func (s *Schemas) Length() int {
	return len(*s)
}
