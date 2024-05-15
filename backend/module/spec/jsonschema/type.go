package jsonschema

import (
	"encoding/json"
)

const (
	T_NULL = "null"
	T_BOOL = "boolean"
	T_OBJ  = "object"
	T_ARR  = "array"
	T_NUM  = "number"
	T_INT  = "integer"
	T_STR  = "string"
)

type SchemaType struct {
	types   []string
	isSlice bool
}

func NewSchemaType(v ...string) *SchemaType {
	return &SchemaType{
		types:   v,
		isSlice: len(v) > 1,
	}
}

func (s *SchemaType) List() []string {
	if s == nil {
		return []string{}
	}
	return s.types
}

func (s *SchemaType) First() string {
	l := s.List()
	if len(l) == 0 {
		return T_NULL
	} else {
		return l[0]
	}
}

func (s *SchemaType) Set(v ...string) {
	if len(v) > 1 {
		s.isSlice = true
	}
	s.types = v
}

func (s *SchemaType) UnmarshalJSON(raw []byte) error {
	if len(raw) == 0 {
		return nil
	}
	if p := raw[0]; p == '[' {
		s.isSlice = true
		return json.Unmarshal(raw, &s.types)
	}
	var o string
	if err := json.Unmarshal(raw, &o); err != nil {
		return err
	}
	s.types = []string{o}
	return nil
}

func (s SchemaType) MarshalJSON() ([]byte, error) {
	if len(s.types) == 0 {
		return []byte(T_NULL), nil
	}
	if s.isSlice {
		return json.Marshal(s.types)
	} else {
		return json.Marshal(s.types[0])
	}
}
