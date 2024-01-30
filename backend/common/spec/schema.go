package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
}

type Example struct {
	Summary string  `json:"summary"`
	Value   any     `json:"value"`
	XDiff   *string `json:"x-apicat-diff,omitempty"`
}

func (s *Schema) Ref() bool { return s.Reference != nil }

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

func (s *Schema) DereferenceSchema(sub *Schema) error {
	if sub == nil {
		return errors.New("sub schema is nil")
	}

	if sub.Type == string(ContentItemTypeDir) {
		return errors.New("sub schema type is dir")
	}

	id := strconv.Itoa(int(sub.ID))

	// If the type of root refers to sub
	if s.Schema.IsRefId(id) || s.IsRefId(id) {
		*s.Schema = *sub.Schema
		return nil
	}

	refs := s.Schema.FindRefById(id)

	// if not find, skip
	// if len(refs) == 0 {
	// 	// sub_schema with this id  not find in parent_schema
	// 	return errors.New("sub schema not find in parent schema")
	// }

	// if sub's refers to itself, Dereference it self
	sub.dereferenceSelf()

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
			stype := "object"
			if len(ps.Schema.Type.Value()) > 0 {
				stype = ps.Schema.Type.Value()[0]
			}
			*s = *jsonschema.Create(stype)
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

func (s *Schema) RemoveSchema(sub *Schema) error {
	if s == nil {
		return errors.New("schema is nil")
	}
	if s.Type == string(ContentItemTypeDir) {
		return errors.New("schema type is dir")
	}
	id := strconv.Itoa(int(sub.ID))
	stype := "object"
	if len(sub.Schema.Type.Value()) > 0 {
		stype = sub.Schema.Type.Value()[0]
	}

	if s.Schema.IsRefId(id) {
		// need change to sub.type
		s.Schema = jsonschema.Create(stype)
		return nil
	}

	s.Schema.RemovePropertyByRefId(id, stype)
	return nil
}

func (s *Schema) dereferenceSelf() {
	if s.Schema == nil || s == nil {
		return
	}

	stype := "object"
	if len(s.Schema.Type.Value()) > 0 {
		stype = s.Schema.Type.Value()[0]
	}

	id := strconv.Itoa(int(s.ID))
	if s.IsRefId(id) {
		s.Schema = jsonschema.Create(stype)
		return
	}

	s.Schema.RemovePropertyByRefId(id, stype)
}

// this schema's type must be dir
func (s *Schema) ItemsTreeToList() (res Schemas) {
	if s.Type != string(ContentItemTypeDir) {
		return append(res, s)
	}
	return s.itemsTreeToList(s.Name)
}

func (s *Schema) itemsTreeToList(path string) (res Schemas) {
	if s.Items == nil || len(s.Items) == 0 {
		return res
	}

	for _, item := range s.Items {
		if item.Type == string(ContentItemTypeDir) {
			res = append(res, item.itemsTreeToList(fmt.Sprintf("%s/%s", path, item.Name))...)
		} else {
			item.Schema.XCategory = path
			res = append(res, item)
		}
	}

	return res
}

// basic info just have name and type
func IsChangedBasic(a, b *Schema) bool {
	return a.Name != b.Name || a.Type != b.Type || jsonschema.IsChangedBasic(a.Schema, b.Schema)
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

// this schemas must not have category
func (s *Schemas) ItemsListToTree() Schemas {
	root := &Schema{
		Items: Schemas{},
	}
	if s == nil || len(*s) == 0 {
		return root.Items
	}

	category := map[string]*Schema{
		"": root,
	}

	for _, v := range *s {
		if parent, ok := category[v.Schema.XCategory]; ok {
			parent.Items = append(parent.Items, v)
		} else {
			root.Items = append(root.Items, v.makeSelfTree(v.Schema.XCategory, category))
		}
	}

	return root.Items
}

func (s *Schema) makeSelfTree(path string, category map[string]*Schema) *Schema {
	if path == "" {
		return s
	}
	i := strings.Index(path, "/")
	if i == -1 {
		parent := &Schema{
			Name:  path,
			Items: Schemas{s},
			Type:  string(ContentItemTypeDir),
		}
		category[path] = parent
		return parent
	}
	parent := &Schema{
		Name:  path[:i],
		Items: Schemas{s.makeSelfTree(path[i+1:], category)},
		Type:  string(ContentItemTypeDir),
	}
	category[path] = parent
	return parent
}

func (s *Schema) SetXDiff(x *string) {
	if s.Schema != nil {
		s.Schema.SetXDiff(x)
	}
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
