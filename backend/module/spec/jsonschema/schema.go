package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Schema struct {
	// Meta Data
	Title       string `json:"title,omitempty" yaml:"title,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Default     any    `json:"default,omitempty" yaml:"default,omitempty"`
	WriteOnly   *bool  `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	ReadOnly    *bool  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	Examples    any    `json:"examples,omitempty" yaml:"examples,omitempty"`
	Deprecated  *bool  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	// Core
	Reference *string `json:"$ref,omitempty" yaml:"$ref,omitempty"`

	// Applicator
	AllOf                AllOf                    `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf                AnyOf                    `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	OneOf                OneOf                    `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	Not                  *Schema                  `json:"not,omitempty" yaml:"not,omitempty"`
	Properties           map[string]*Schema       `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties *ValueOrBoolean[*Schema] `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Items                *ValueOrBoolean[*Schema] `json:"items,omitempty" yaml:"items,omitempty"` // 3.1 schema or bool

	// Validation
	Type             *SchemaType              `json:"type,omitempty" yaml:"type,omitempty"` // 3.1 []string 2,3.0 string
	Enum             []any                    `json:"enum,omitempty" yaml:"enum,omitempty"`
	Pattern          string                   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MinLength        *int64                   `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength        *int64                   `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	ExclusiveMaximum *ValueOrBoolean[float64] `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"` // 3.0 bool 3.1 int
	MultipleOf       *float64                 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	ExclusiveMinimum *ValueOrBoolean[float64] `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"` // 3.0 bool 3.1 int
	Maximum          *float64                 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	Minimum          *float64                 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	MaxProperties    *int64                   `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties    *int64                   `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required         []string                 `json:"required,omitempty" yaml:"required,omitempty"`
	MaxItems         *int64                   `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems         *int64                   `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems      *bool                    `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`

	// Format Annotation
	Format string `json:"format,omitempty" yaml:"format,omitempty"`

	// Extension
	ID       int64    `json:"id,omitempty" yaml:"id,omitempty"`
	XOrder   []string `json:"x-apicat-orders,omitempty" yaml:"x-apicat-orders,omitempty"`
	XMock    string   `json:"x-apicat-mock,omitempty" yaml:"x-apicat-mock,omitempty"`
	XDiff    string   `json:"x-apicat-diff,omitempty" yaml:"x-apicat-diff,omitempty"`
	Nullable *bool    `json:"nullable,omitempty" yaml:"nullable,omitempty"`
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
	if typ == "" {
		return &Schema{
			Type: NewSchemaType(T_OBJ),
		}
	} else {
		return &Schema{
			Type: NewSchemaType(typ),
		}
	}
}

func NewSchemaFromJson(str string) (*Schema, error) {
	s := &Schema{}
	if err := json.Unmarshal([]byte(str), s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Schema) Ref() bool { return s != nil && s.Reference != nil }

func (s *Schema) DeepRef() bool {
	if s != nil && s.Ref() {
		return true
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			if v.DeepRef() {
				return true
			}
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		return s.Items.Value().DeepRef()
	}
	return false
}

// Check if the schema refers to this id
func (s *Schema) IsRefID(id string) bool {
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

func (s *Schema) DeepFindRefById(id string) (refs []*Schema) {
	if s == nil {
		return
	}

	if s.IsRefID(id) {
		refs = append(refs, s)
		return
	}

	if s.Properties != nil {
		for k := range s.Properties {
			refs = append(refs, s.Properties[k].DeepFindRefById(id)...)
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		refs = append(refs, s.Items.Value().DeepFindRefById(id)...)
	}
	return
}

func (s *Schema) GetRefID() int64 {
	if !s.Ref() {
		return 0
	}

	i := strings.LastIndex(*s.Reference, "/")
	if i != -1 {
		id, _ := strconv.ParseInt((*s.Reference)[i+1:], 10, 64)
		return id
	}
	return 0
}

func (s *Schema) DeepGetRefID() (ids []int64) {
	if s == nil {
		return
	}

	if id := s.GetRefID(); id > 0 {
		ids = append(ids, id)
		return
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			ids = append(ids, v.DeepGetRefID()...)
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		ids = append(ids, s.Items.Value().DeepGetRefID()...)
	}
	return
}

func (s *Schema) ReplaceRef(ref *Schema) error {
	if !s.Ref() || ref == nil {
		return errors.New("schema is not a reference or ref is nil")
	}

	refID := s.GetRefID()
	if refID != ref.ID {
		return errors.New("ref id does not match")
	}

	*s = *ref
	return nil
}

func (s *Schema) DelRootRef(ref *Schema) {
	if !s.Ref() || ref == nil {
		return
	}

	refID := s.GetRefID()
	if refID != ref.ID {
		return
	}

	s.Reference = nil
	s.Type = NewSchemaType(T_OBJ)
}

func (s *Schema) DelChildrenRef(ref *Schema) {
	if s == nil || ref == nil {
		return
	}

	propertyKeys := []string{}
	if s.Properties != nil {
		for k, v := range s.Properties {
			if v.IsRefID(strconv.FormatInt(ref.ID, 10)) {
				propertyKeys = append(propertyKeys, k)
			}
			v.DelChildrenRef(ref)
		}

		for _, k := range propertyKeys {
			delete(s.Properties, k)
			s.DelXOrderByName(k)
		}
	}

	if s.Items != nil && !s.Items.IsBool() && s.Items.Value() != nil {
		if s.Items.Value().IsRefID(strconv.FormatInt(ref.ID, 10)) {
			s.Items.SetValue(NewSchema(ref.Type.First()))
		}
	}
}

func (s *Schema) DelXOrderByName(name string) {
	if s == nil || name == "" {
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

func (s *Schema) CheckAllOf() bool { return s != nil && len(s.AllOf) > 0 }

func (s *Schema) DeepCheckAllOf() bool {
	if s == nil {
		return false
	}

	if s.CheckAllOf() {
		return true
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			if v.DeepCheckAllOf() {
				return true
			}
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		return s.Items.Value().DeepCheckAllOf()
	}
	return false
}

func (s *Schema) MergeAllOf() {
	if s == nil {
		return
	}

	if s.CheckAllOf() {
		s.AllOf = s.AllOf.Merge()
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			v.MergeAllOf()
		}
	}
	if s.Items != nil && !s.Items.IsBool() {
		s.Items.Value().MergeAllOf()
	}
}

func (s *Schema) Validation(raw []byte) error {
	return nil
}

func (s *Schema) Valid() error {
	if s == nil {
		return errors.New("schema is nil")
	}

	if s.Ref() {
		return nil
	}

	for _, v := range s.Type.List() {
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
	if s == nil {
		return errors.New("schema is nil")
	}

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
	if s == nil || s.Items == nil {
		return nil
	}
	if s.Items.IsBool() {
		return nil
	}
	return s.Items.Value().Valid()
}

func (s *Schema) SetXDiff(x string) {
	if s == nil || x == "" {
		return
	}

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

func (s *Schema) ToJson() string {
	if s == nil {
		return ""
	}

	b, _ := json.Marshal(s)
	return string(b)
}
