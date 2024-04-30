package jsonschema

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Schema struct {
	// Meta Data
	Title       string `json:"title,omitempty" yaml:"title,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Default     any    `json:"default,omitempty" yaml:"default,omitempty"`
	WriteOnly   bool   `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	ReadOnly    bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	Examples    any    `json:"examples,omitempty" yaml:"examples,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	// Core
	Reference string `json:"$ref,omitempty" yaml:"$ref,omitempty"`

	// Applicator
	AllOf                AllOf                    `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf                AnyOf                    `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	OneOf                OneOf                    `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	Not                  *Schema                  `json:"not,omitempty" yaml:"not,omitempty"`
	Properties           map[string]*Schema       `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties *ValueOrBoolean[*Schema] `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Items                *ValueOrBoolean[*Schema] `json:"items,omitempty" yaml:"items,omitempty"` // 3.1 schema or bool

	// Validation
	Type             *SliceOrOneValue[string] `json:"type,omitempty" yaml:"type,omitempty"` // 3.1 []string 2,3.0 string
	Enum             []any                    `json:"enum,omitempty" yaml:"enum,omitempty"`
	Pattern          string                   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MinLength        int64                    `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength        int64                    `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	ExclusiveMaximum *ValueOrBoolean[int64]   `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"` // 3.0 bool 3.1 int
	MultipleOf       int64                    `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	ExclusiveMinimum *ValueOrBoolean[int64]   `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"` // 3.0 bool 3.1 int
	Maximum          int64                    `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	Minimum          int64                    `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	MaxProperties    int64                    `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties    int64                    `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required         []string                 `json:"required,omitempty" yaml:"required,omitempty"`
	MaxItems         int64                    `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems         int64                    `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems      int64                    `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`

	// Format Annotation
	Format string `json:"format,omitempty" yaml:"format,omitempty"`

	// Extension
	XOrder   []string `json:"x-apicat-orders,omitempty" yaml:"x-apicat-orders,omitempty"`
	XMock    string   `json:"x-apicat-mock,omitempty" yaml:"x-apicat-mock,omitempty"`
	XDiff    string   `json:"x-apicat-diff,omitempty" yaml:"x-apicat-diff,omitempty"`
	Nullable bool     `json:"nullable,omitempty" yaml:"nullable,omitempty"`
}

var coreTypes = []string{
	"string",
	"integer",
	"number",
	"boolean",
	"object",
	"array",
	"null",
}

func NewSchema(typ string) *Schema {
	return &Schema{
		Type: NewSliceOrOne(typ),
	}
}

func (s *Schema) Ref() bool { return s != nil && s.Reference != "" }

func (s *Schema) FindRefById(id string) (refs []*Schema) {
	if s == nil {
		return nil
	}

	if s.IsRefId(id) {
		refs = append(refs, s)
		return refs
	}

	if s.Properties != nil {
		for k := range s.Properties {
			refs = append(refs, s.Properties[k].FindRefById(id)...)
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		refs = append(refs, s.Items.Value().FindRefById(id)...)
	}
	return refs
}

// check this schema reference this id
func (s *Schema) IsRefId(id string) bool {
	if s == nil || s.Reference == "" {
		return false
	}

	i := strings.LastIndex(s.Reference, "/")
	if i != -1 {
		if id == (s.Reference)[i+1:] {
			return true
		}
	}
	return false
}

func (s *Schema) DelPropertyByRefId(id string, stype string) {
	if s == nil {
		return
	}

	ks := []string{}
	if s.Properties != nil {
		for k, v := range s.Properties {
			if v.IsRefId(id) {
				ks = append(ks, k)
			}
			v.DelPropertyByRefId(id, stype)
		}

		for _, v := range ks {
			delete(s.Properties, v)
			s.DelXOrderByName(v)
		}
	}

	if s.Items != nil && s.Items.Value() != nil {
		if s.Items.value.IsRefId(id) {
			s.Items.value = NewSchema(stype)
		}
	}
}

func (s *Schema) DelXOrderByName(name string) {
	if s == nil {
		return
	}
	if s.XOrder != nil {
		i := 0
		for i < len(s.XOrder) {
			if s.XOrder[i] == name {
				s.XOrder = append(s.XOrder[:i], s.XOrder[i+1:]...)
				continue
			}
			i++
		}
	}
}

func (s *Schema) Validation(raw []byte) error {
	return nil
}

func (s *Schema) Valid() error {
	if s.Ref() {
		return nil
	}

	for _, v := range s.Type.Value() {
		if !slices.Contains(coreTypes, v) {
			return fmt.Errorf("unkowan type %s", v)
		}
		switch v {
		case "array":
			return s.checkArray()
		case "object":
			return s.checkObject()
		}
	}
	return nil
}

func (s *Schema) checkObject() error {
	if s.Ref() || s.AdditionalProperties == nil {
		return nil
	}
	proplen := 0
	if s.Properties != nil {
		proplen = len(s.Properties)
	}
	if orderlen := len(s.XOrder); proplen > 0 {
		for name, prop := range s.Properties {
			if err := prop.Valid(); err != nil {
				return err
			}
			if orderlen > 0 {
				if !slices.Contains(s.XOrder, name) {
					return fmt.Errorf("x-apicat-order does not match the properties")
				}
			}
		}
		// check required?
	}
	if s.AdditionalProperties != nil &&
		!s.AdditionalProperties.IsBool() {
		return s.AdditionalProperties.Value().Valid()
	}
	return nil
}

func (s *Schema) checkArray() error {
	if s.Items == nil {
		return nil
	}
	if s.Items.IsBool() {
		return nil
	}
	return s.Items.Value().Valid()
}

func (s *Schema) SetXDiff(x string) {
	if s.Properties != nil {
		for _, v := range s.Properties {
			v.SetXDiff(x)
		}
	}
	if s.Items != nil && !s.Items.IsBool() {
		s.Items.value.SetXDiff(x)
	}
	s.XDiff = x
}
