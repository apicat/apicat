package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

type Referencer interface {
	Ref() bool
}

type Schema struct {
	ID          int64               `json:"id,omitempty" yaml:"id,omitempty"`
	ParentId    uint64              `json:"parentid,omitempty" yaml:"parentid,omitempty"`
	Name        string              `json:"name,omitempty" yaml:"name,omitempty"`
	Type        string              `json:"type,omitempty" yaml:"type,omitempty"`
	Description string              `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool                `json:"required,omitempty" yaml:"required,omitempty"`
	Example     any                 `json:"example,omitempty" yaml:"example,omitempty"`
	Examples    map[string]*Example `json:"examples,omitempty" yaml:"examples,omitempty"`
	Schema      *jsonschema.Schema  `json:"schema,omitempty" yaml:"schema,omitempty"`
	Items       Schemas             `json:"items,omitempty" yaml:"items,omitempty"`
	Reference   *string             `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

type Example struct {
	Summary string  `json:"summary" yaml:"summary"`
	Value   any     `json:"value" yaml:"value"`
	XDiff   *string `json:"x-apicat-diff,omitempty" yaml:"x-apicat-diff,omitempty"`
}

func NewSchemaWithJson(str string) (*Schema, error) {
	s := &Schema{}
	return s, json.Unmarshal([]byte(str), s)
}

func (s *Schema) Ref() bool { return s.Reference != nil }

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

// @title Deref
// @description 给一个引用列表 wantToDeref，将目标 Schema 中所有引用 wantToDeref 提到的 Schema 全部解开
// @param wantToDeref ...*Schema 不再引用的 Definition Schema
// @return nil
func (s *Schema) Deref(wantToDeref ...*Schema) {
	// 自己直接引自己
	if s.IsRefID(strconv.Itoa(int(s.ID))) {
		s.Schema = jsonschema.Create("object")
		return
	}

	if len(wantToDeref) == 0 {
		return
	}

	jsonSchemaDeref(s.Schema, wantToDeref, make(map[int64]*Schema))
}

// @title jsonSchemaDeref
// @description 根据 wantToDeref 中的 Schema，将 jsonschema 中的 $ref 以及嵌套的 $ref 替换为非 $ref 的方式
// @param s *jsonschema.Schema Definition Schema 中的 Schema 参数内容
// @param wantToDeref Schemas 不再引用的 Definition Response
// @param repeat map[int64]*Schema 已经解过引用的 Definition Response ID 标记
// @return nil
func jsonSchemaDeref(s *jsonschema.Schema, wantToDeref Schemas, repeat map[int64]*Schema) {
	if s == nil {
		return
	}

	if s.Reference != nil {
		refSchemaID := mustGetRefID(*s.Reference)
		if v, ok := repeat[refSchemaID]; ok {
			// 解过就不再解了，避免出现死循环，针对 Scheme 自己引用自己的情况
			refSchemaType := "object"
			if len(v.Schema.Type.Value()) > 0 {
				refSchemaType = v.Schema.Type.Value()[0]
			}
			*s = *jsonschema.Create(refSchemaType)
			return
		}

		refSchema := wantToDeref.LookupByID(refSchemaID)
		if refSchema == nil {
			// 没找到，就不解了
			return
		}
		// 找到了，标记已经解过引用了
		repeat[refSchemaID] = refSchema

		result := jsonschema.Schema{}
		refJsonSchema, _ := json.Marshal(refSchema.Schema)
		_ = json.Unmarshal(refJsonSchema, &result)
		*s = result
		jsonSchemaDeref(s, wantToDeref, repeat)
		delete(repeat, refSchemaID)
		return
	}

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
			jsonSchemaDeref(s.Items.Value(), wantToDeref, repeat)
		}
	case "object":
		for _, v := range s.Properties {
			jsonSchemaDeref(v, wantToDeref, repeat)
		}
	}
}

// @title DelRef
// @description 删除引用的 Definition Schema
// @param refSchema *Schema 要删除的 Definition Schema
// @return error
func (s *Schema) DelRef(refSchema *Schema) error {
	if s == nil {
		return errors.New("schema is nil")
	}
	if s.Type == string(CollectionItemTypeDir) {
		return errors.New("schema type is dir")
	}
	id := strconv.Itoa(int(refSchema.ID))

	stype := "object"
	if len(refSchema.Schema.Type.Value()) > 0 {
		stype = refSchema.Schema.Type.Value()[0]
	}
	if s.Schema.IsRefId(id) {
		s.Schema = jsonschema.Create(stype)
		return nil
	}

	s.Schema.DelPropertyByRefId(id, stype)
	return nil
}

// this schema's type must be dir
func (s *Schema) ItemsTreeToList() (res Schemas) {
	if s.Type != string(CollectionItemTypeDir) {
		return append(res, s)
	}
	return s.itemsTreeToList(s.Name)
}

func (s *Schema) itemsTreeToList(path string) (res Schemas) {
	if s.Items == nil || len(s.Items) == 0 {
		return res
	}

	for _, item := range s.Items {
		if item.Type == string(CollectionItemTypeDir) {
			res = append(res, item.itemsTreeToList(fmt.Sprintf("%s/%s", path, item.Name))...)
		} else {
			item.Schema.XCategory = path
			res = append(res, item)
		}
	}

	return res
}

type Schemas []*Schema

func (s *Schemas) LookupByName(name string) *Schema {
	if s == nil {
		return nil
	}
	for _, v := range *s {
		if v.Type == string(CollectionItemTypeDir) {
			if res := v.Items.LookupByName(name); res != nil {
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

func (s *Schemas) LookupByID(id int64) *Schema {
	if s == nil {
		return nil
	}
	for _, v := range *s {
		if v.Type == string(CollectionItemTypeDir) {
			if res := v.Items.LookupByID(id); res != nil {
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

func (s *Schemas) DelById(id int64) {
	if s == nil {
		return
	}
	for i, v := range *s {
		if v.Type == string(CollectionItemTypeDir) {
			v.Items.DelById(id)
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
			Type:  string(CollectionItemTypeDir),
		}
		category[path] = parent
		return parent
	}
	parent := &Schema{
		Name:  path[:i],
		Items: Schemas{s.makeSelfTree(path[i+1:], category)},
		Type:  string(CollectionItemTypeDir),
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
	if s.Description != o.Description || s.Required != o.Required || fmt.Sprintf("%v", s.Example) != fmt.Sprintf("%v", o.Example) {
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
		} else if fmt.Sprintf("%v", se.Value) != fmt.Sprintf("%v", oe.Value) {
			s := "!"
			oe.XDiff = &s
			b = false
		}
	}
	return b
}
