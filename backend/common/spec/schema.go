package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/apicat/apicat/backend/common/spec/jsonschema"
)

type Referencer interface {
	Ref() bool
}

type Schema struct {
	ID          int64               `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Type        string              `json:"type,omitempty"`
	ParentId    uint64              `json:"parentid,omitempty"`
	Description string              `json:"description,omitempty"`
	Required    bool                `json:"required,omitempty"`
	Example     any                 `json:"example,omitempty"`
	Examples    map[string]*Example `json:"examples,omitempty"`
	Schema      *jsonschema.Schema  `json:"schema,omitempty"`
	Items       Schemas             `json:"items,omitempty"`
	Reference   *string             `json:"$ref,omitempty"`
	XDiff       *string             `json:"x-apicat-diff,omitempty"`
}

type Example struct {
	Summary string  `json:"summary"`
	Value   any     `json:"value"`
	XDiff   *string `json:"x-apicat-diff,omitempty"`
}

func (s *Schema) Ref() bool { return s.Reference != nil }

func (s *Schema) DereferenceSchema(sub *Schema) error {
	if sub == nil {
		return errors.New("sub schema is nil")
	}

	if sub.Type == string(ContentItemTypeDir) {
		return errors.New("sub schema type is dir")
	}

	id := strconv.Itoa(int(sub.ID))
	// If the type of root refers to sub
	if s.Schema.IsRefId(id) {
		s.Schema = jsonschema.Create("object")
	}

	refs := s.Schema.FindRefById(id)
	if len(refs) == 0 {
		// sub_schema with this id  not find in parent_schema
		return errors.New("sub schema not find in parent schema")
	}

	// if sub's refers to itself, Dereference it self
	// sub.dereferenceSelf()

	// replace all referenced sub.Schema with dereferenced sub.Schema
	for i := range refs {
		*refs[i] = *sub.Schema
	}

	return nil
}

func (s *Schema) UnpackDereferenceSchema(sub Schemas) error {
	if s == nil {
		return errors.New("schema is nil")
	}
	if s.Type == string(ContentItemTypeDir) {
		return errors.New("schema type is dir")
	}
	if len(sub) == 0 {
		return errors.New("sub schema is nil")
	}
	dfsDereferenceSchemas(s.Schema, s.ID, sub, make(map[int64]bool))
	return nil
}

func dfsDereferenceSchemas(s *jsonschema.Schema, id int64, sub Schemas, refs map[int64]bool) {
	if s == nil {
		return
	}

	if _, ok := refs[id]; ok && id != -1 {
		return
	}

	refs[id] = true
	if s.Reference != nil {
		ps := sub.LookupID(mustGetRefID(*s.Reference))
		if ps == nil {
			return
		}
		if _, ok := refs[ps.ID]; ok {
			*s = *jsonschema.Create("object")
			return
		}

		// Prevent the data in sub from being modified by pointers
		// jsonmarshal just like deep copy
		ss := jsonschema.Schema{}
		b, _ := json.Marshal(ps.Schema)
		_ = json.Unmarshal(b, &ss)
		*s = ss
		dfsDereferenceSchemas(s, ps.ID, sub, refs)
		delete(refs, ps.ID)
		return
	}
	// no type
	if s.Type == nil {
		return
	}
	t := s.Type.Value()[0]
	switch t {
	case "array":
		if s.Items != nil {
			if s.Items.IsBool() {
				return
			}
			dfsDereferenceSchemas(s.Items.Value(), -1, sub, refs)
		}
	case "object":
		for _, v := range s.Properties {
			dfsDereferenceSchemas(v, -1, sub, refs)
		}
	}
}

func (s *Schema) RemoveSchema(s_id int64) error {
	if s == nil {
		return errors.New("schema is nil")
	}
	if s.Type == string(ContentItemTypeDir) {
		return errors.New("schema type is dir")
	}
	id := strconv.Itoa(int(s_id))

	if s.Schema.IsRefId(id) {
		s.Schema = jsonschema.Create("object")
	}

	s.Schema.RemovePropertyByRefId(id)
	return nil
}

func (s *Schema) dereferenceSelf() {
	if s.Schema == nil {
		return
	}

	id := strconv.Itoa(int(s.ID))

	s.Schema.RemovePropertyByRefId(id)
}

// this schema's type is dir
func (s *Schema) ItemsTreeToList() (res Schemas) {
	if s.Items == nil || len(s.Items) == 0 {
		return res
	}

	for _, item := range s.Items {
		if item.Type == string(ContentItemTypeDir) {
			res = append(res, item.ItemsTreeToList()...)
		} else {
			res = append(res, item)
		}
	}

	return res
}

func (s *Schema) JoinNameId() string {
	return fmt.Sprintf("%s-%d", s.Name, s.ID)
}

type Schemas []*Schema

func (s *Schemas) Lookup(name string) *Schema {
	if s == nil {
		return nil
	}
	for _, v := range *s {
		if v.Type == string(ContentItemTypeDir) {
			if res := v.Items.Lookup(name); res != nil {
				return res
			}
		} else {
			if name == v.Name {
				return v
			}
		}
	}
	return nil
}

func (s *Schemas) LookupID(id int64) *Schema {
	if s == nil {
		return nil
	}
	for _, v := range *s {
		if v.Type == string(ContentItemTypeDir) {
			if res := v.Items.LookupID(id); res != nil {
				return res
			}
		} else {
			if id == v.ID {
				return v
			}
		}
	}
	return nil
}

func (s *Schemas) RemoveId(id int64) {
	if s == nil {
		return
	}
	for i, v := range *s {
		if v.Type == string(ContentItemTypeDir) {
			v.Items.RemoveId(id)
		} else {
			if v.ID == id {
				*s = append((*s)[:i], (*s)[i+1:]...)
			}
		}
	}
}

func (s *Schemas) Length() int {
	return len(*s)
}

func (s *Schemas) SetXDiff(x *string) {
	for _, v := range *s {
		v.SetXDiff(x)
	}
}

func (s *Schema) SetXDiff(x *string) {
	if s.Schema != nil {
		s.Schema.SetXDiff(x)
	}
	s.XDiff = x
}

func (s *Schemas) UnpackDereferenceSchema(sub Schemas) (err error) {
	for _, v := range *s {
		err = v.UnpackDereferenceSchema(sub)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *Schema) FindExample(summary string) (*Example, bool) {
	for _, v := range s.Examples {
		if v.Summary == summary {
			return v, true
		}
	}
	return nil, false
}

func (s *Schema) EqualNomal(o *Schema) (b bool) {
	b = true
	if s.Description != o.Description || s.Required != o.Required || s.Example != o.Example {
		b = false
	}
	return b && s.EqualExamples(o)
}

func (s *Schema) EqualExamples(o *Schema) (b bool) {
	b = true
	names := map[string]struct{}{}
	for _, v := range s.Examples {
		names[v.Summary] = struct{}{}
	}
	for _, v := range o.Examples {
		names[v.Summary] = struct{}{}
	}

	for k := range names {
		se, s_has := s.FindExample(k)
		oe, o_has := o.FindExample(k)
		if !s_has && o_has {
			s := "+"
			oe.XDiff = &s
			b = false
		} else if s_has && !o_has {
			s := "-"
			se.XDiff = &s
			if o.Examples == nil {
				o.Examples = make(map[string]*Example)
			}
			o.Examples[strconv.Itoa(len(o.Examples))] = se
			b = false
		} else if se.Value != oe.Value {
			s := "!"
			oe.XDiff = &s
			b = false
		}
	}
	return b
}
