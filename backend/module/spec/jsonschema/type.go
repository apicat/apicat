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

// Check if the type is one-dimensional
// Such as int array string
func (s *SchemaType) IsOneDimensional() bool {
	switch s.First() {
	case T_NULL, T_BOOL, T_NUM, T_INT, T_STR:
		return true
	}
	return false
}

func (s *SchemaType) Equal(a *SchemaType) bool {
	if a == nil {
		return false
	}
	if len(s.types) != len(a.types) {
		return false
	}
	stypes := make(map[string]bool, 0)
	for _, t := range s.types {
		stypes[t] = true
	}
	for _, t := range a.types {
		if _, ok := stypes[t]; !ok {
			return false
		}
	}
	return true
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
		return json.Marshal(s.First())
	}
}
