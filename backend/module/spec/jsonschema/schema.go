package jsonschema

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Schema struct {
	// 3.1 []string 2,3.0 string
	Type *SliceOrOneValue[string] `json:"type,omitempty" yaml:"type,omitempty"`

	XOrder []string `json:"x-apicat-orders,omitempty" yaml:"x-apicat-orders,omitempty"`
	XMock  string   `json:"x-apicat-mock,omitempty" yaml:"x-apicat-mock,omitempty"`
	// diff 如果有值就代表有变化
	XDiff *string `json:"x-apicat-diff,omitempty" yaml:"x-apicat-diff,omitempty"`
	// category path
	XCategory string `json:"x-apicat-category,omitempty" yaml:"x-apicat-category,omitempty"`
	// 3.1 schema or bool
	Items *ValueOrBoolean[*Schema] `json:"items,omitempty" yaml:"items,omitempty"`

	// 先不支持 后期再说
	// AllOf []*Schema `json:"AllOf,omitempty" yaml:"AllOf,omitempty"`
	// AnyOf []*Schema `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	// OneOf []*Schema `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`

	// 3.0 bool 3.1 int
	ExclusiveMaximum *ValueOrBoolean[int64] `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum *ValueOrBoolean[int64] `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`

	// 所有版本
	Not                  *Schema                  `json:"not,omitempty" yaml:"not,omitempty"`
	Properties           map[string]*Schema       `json:"properties,omitempty" yaml:"properties,omitempty"`
	Title                string                   `json:"title,omitempty" yaml:"title,omitempty"`
	MultipleOf           *int64                   `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Maximum              *int64                   `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	Minimum              *int64                   `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	MaxLength            *int64                   `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength            *int64                   `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Format               string                   `json:"format,omitempty" yaml:"format,omitempty"`
	Pattern              string                   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MaxItems             *int64                   `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems             *int64                   `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems          *int64                   `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	MaxProperties        *int64                   `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        *int64                   `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required             []string                 `json:"required,omitempty" yaml:"required,omitempty"`
	Enum                 []any                    `json:"enum,omitempty" yaml:"enum,omitempty"`
	AdditionalProperties *ValueOrBoolean[*Schema] `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Description          string                   `json:"description,omitempty" yaml:"description,omitempty"`
	Default              any                      `json:"default,omitempty" yaml:"default,omitempty"`
	Nullable             *bool                    `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	ReadOnly             bool                     `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly            bool                     `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	Example              any                      `json:"example,omitempty" yaml:"example,omitempty"`
	Deprecated           bool                     `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	Reference *string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
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

func (s *Schema) Ref() bool { return s != nil && s.Reference != nil }

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
	if s == nil || s.Reference == nil {
		return false
	}

	i := strings.LastIndex(*s.Reference, "/")
	if i != -1 {
		if id == (*s.Reference)[i+1:] {
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
			s.Items.value = Create(stype)
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

func Create(typ string) *Schema {
	return &Schema{
		Type: CreateSliceOrOne(typ),
	}
}

func (s *Schema) SetXDiff(x *string) {
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
