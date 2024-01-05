package spec

import (
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

func (s *Schema) DereferenceSchema(sub *Schema) {
	if sub == nil {
		return
	}

	id := strconv.Itoa(int(sub.ID))
	// If the type of root refers to sub
	if s.Schema.IsRefId(id) {
		s.Schema = jsonschema.Create("object")
	}

	refs := s.Schema.FindRefById(id)
	if len(refs) == 0 {
		// sub_schema with this id  not find in parent_schema
		return
	}

	// if sub's refers to itself, Dereference it self
	sub.dereferenceSelf()

	// replace all referenced sub.Schema with dereferenced sub.Schema
	for i := range refs {
		*refs[i] = *sub.Schema
	}

}

func (s *Schema) RemoveSchema(sub *Schema) {
	if sub == nil {
		return
	}
	id := strconv.Itoa(int(sub.ID))

	if s.Schema.IsRefId(id) {
		s.Schema = jsonschema.Create("object")
	}

	s.Schema.RemovePropertyByRefId(id)

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
