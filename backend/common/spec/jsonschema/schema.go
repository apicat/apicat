package jsonschema

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Schema struct {
	// 3.1 []string 2,3.0 string
	Type *SliceOrOneValue[string] `json:"type,omitempty"`

	XOrder []string `json:"x-apicat-orders,omitempty"`
	XMock  string   `json:"x-apicat-mock,omitempty"`
	// diff 如果有值就代表有变化
	XDiff *string `json:"x-apicat-diff,omitempty"`
	// category path
	XCategory string `json:"x-apicat-category,omitempty"`
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

	// this s is object
	if s.Properties != nil {
		for k := range s.Properties {
			refs = append(refs, s.Properties[k].FindRefById(id)...)
		}
	}

	// this s is array
	if s.Items != nil {
		if s.Items.IsBool() {
			return
		}
		refs = append(refs, s.Items.value.FindRefById(id)...)
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

func IsChangedBasic(a, b *Schema) bool {
	if a == nil && b == nil {
		return false
	}

	if a == nil && b != nil {
		return true
	}

	if a != nil && b == nil {
		return true
	}

	if !slices.Equal(a.Type.Value(), b.Type.Value()) {
		return true
	}

	if len(a.Type.Value()) == 0 {
		return false
	}
	at := a.Type.Value()[0]
	bt := b.Type.Value()[0]
	if at != bt {
		return true
	}

	switch bt {
	case "object":
		names := map[string]struct{}{}
		for v := range a.Properties {
			names[v] = struct{}{}
		}
		for v := range b.Properties {
			names[v] = struct{}{}
		}

		for v := range names {
			as, a_has := a.Properties[v]
			bs, b_has := b.Properties[v]
			if !a_has && b_has {
				return true
			}
			if a_has && !b_has {
				return true
			}
			if IsChangedBasic(as, bs) {
				return true
			}
		}
	case "array":
		if a.Items != nil && b.Items != nil {
			if IsChangedBasic(a.Items.Value(), b.Items.Value()) {
				return true
			}
		}
	}
	return false
}

func (s *Schema) RemovePropertyByRefId(id string) {
	if s == nil {
		return
	}
	ks := []string{}
	if s.Properties != nil {
		for k, v := range s.Properties {
			if v.IsRefId(id) {
				ks = append(ks, k)
			}
			v.RemovePropertyByRefId(id)
		}
		for _, v := range ks {
			delete(s.Properties, v)
			s.RemoveXOrderByName(v)
		}
	}

	if s.Items != nil {
		if !s.Items.IsBool() && s.Items.value.IsRefId(id) {
			s.Items.value = Create("object")
		}
	}
}

func (s *Schema) RemoveXOrderByName(name string) {
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
	s.XDiff = x
}
