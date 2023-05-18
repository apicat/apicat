package jsonschema

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type Schema struct {
	// 3.1 []string 2,3.0 string
	Type *SliceOrOneValue[string] `json:"type,omitempty"`

	XOrder []string `json:"x-apicat-orders,omitempty"`
	// 3.1 schema or bool
	Items *ValueOrBoolean[*Schema] `json:"items,omitempty"`

	// 先不支持 后期再说
	// AllOf []*Schema `json:"AllOf,omitempty"`
	// AnyOf []*Schema `json:"anyOf,omitempty"`
	// OneOf []*Schema `json:"oneOf,omitempty"`

	// 3.0 bool 3.1 int
	ExclusiveMaximum *ValueOrBoolean[int64] `json:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum *ValueOrBoolean[int64] `json:"exclusiveMinimum,omitempty"`

	// 所有版本
	Not                  *Schema                  `json:"not,omitempty"`
	Properties           map[string]*Schema       `json:"properties,omitempty"`
	Title                string                   `json:"title,omitempty"`
	MultipleOf           *int64                   `json:"multipleOf,omitempty"`
	Maximum              *int64                   `json:"maximum,omitempty"`
	Minimum              *int64                   `json:"minimum,omitempty"`
	MaxLength            *int64                   `json:"maxLength,omitempty"`
	MinLength            *int64                   `json:"minLength,omitempty"`
	Format               string                   `json:"format,omitempty"`
	Pattern              string                   `json:"pattern,omitempty"`
	MaxItems             *int64                   `json:"maxItems,omitempty"`
	MinItems             *int64                   `json:"minItems,omitempty"`
	UniqueItems          *int64                   `json:"uniqueItems,omitempty"`
	MaxProperties        *int64                   `json:"maxProperties,omitempty"`
	MinProperties        *int64                   `json:"minProperties,omitempty"`
	Required             []string                 `json:"required,omitempty"`
	Enum                 []any                    `json:"enum,omitempty"`
	AdditionalProperties *ValueOrBoolean[*Schema] `json:"additionalProperties,omitempty"`
	Description          string                   `json:"description,omitempty"`
	Default              any                      `json:"default,omitempty"`
	Nullable             *bool                    `json:"nullable,omitempty"`
	ReadOnly             bool                     `json:"readOnly,omitempty"`
	WriteOnly            bool                     `json:"writeOnly,omitempty"`
	Example              any                      `json:"example,omitempty"`
	Deprecated           bool                     `json:"deprecated,omitempty"`

	Reference *string `json:"$ref,omitempty"`
}

func (s *Schema) Ref() bool { return s.Reference != nil }

var coreTypes = []string{
	"string",
	"integer",
	"number",
	"boolean",
	"object",
	"array",
	"null",
}

func (s *Schema) Validation(raw []byte) error {
	return nil
}

func (s *Schema) Valid() error {
	if s.Reference != nil {
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
	if s.Properties == nil || s.AdditionalProperties == nil {
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
