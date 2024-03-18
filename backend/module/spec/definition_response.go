package spec

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type HTTPBody map[string]*Schema

type HTTPResponseDefine struct {
	ID          int64               `json:"id,omitempty" yaml:"id,omitempty"`
	ParentId    uint64              `json:"parentid,omitempty" yaml:"parentid,omitempty"`
	Name        string              `json:"name,omitempty" yaml:"name,omitempty"`
	Type        string              `json:"type,omitempty" yaml:"type,omitempty"`
	Description string              `json:"description,omitempty" yaml:"description,omitempty"`
	Header      ParameterList       `json:"header,omitempty" yaml:"header,omitempty"`
	Content     HTTPBody            `json:"content" yaml:"content"`
	Items       HTTPResponseDefines `json:"items,omitempty" yaml:"items,omitempty"`
	Reference   *string             `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	XDiff       *string             `json:"x-apicat-diff,omitempty" yaml:"x-apicat-diff,omitempty"`
	Category    string              `json:"x-apicat-category,omitempty" yaml:"x-apicat-category,omitempty"`
}

func (h *HTTPResponseDefine) Ref() bool { return h.Reference != nil }

func (h *HTTPResponseDefine) IsRefID(id string) bool {
	if h == nil {
		return false
	}

	if h.Reference != nil {
		i := strings.LastIndex(*h.Reference, "/")
		if i != -1 {
			if id == (*h.Reference)[i+1:] {
				return true
			}
		}
	}
	return false
}

// @title Deref
// @description 删除 API 中的响应，如果引用了此 ID 对应的响应，将响应从 API 中删除
// @param sub *HTTPResponseDefine 不再引用的 Definition Response
// @return error
func (h *HTTPResponseDefine) Deref(refResponse *HTTPResponseDefine) error {
	if h == nil {
		return errors.New("sub response is nil")
	}
	id := strconv.Itoa(int(refResponse.ID))

	if !h.IsRefID(id) {
		return nil
	}

	*h = *refResponse
	return nil
}

// @title DerefSchema
// @description 解开响应中引用的 Definition Schema
// @param wantToDeref ...*Schema 不再引用的 Definition Response
// @return nil
func (h *HTTPResponseDefine) DerefSchema(wantToDeref ...*Schema) {
	for _, body := range h.Content {
		body.Deref(wantToDeref...)
	}
}

// @title DelRefSchema
// @description 删除响应中引用的 Definition Schema
// @param schema *Schema 要删除的 Definition Schema
// @return error
func (h *HTTPResponseDefine) DelRefSchema(schema *Schema) (err error) {
	for _, body := range h.Content {
		err = body.DelRef(schema)
		if err != nil {
			return err
		}
	}
	return nil
}

// 将 Response 树形结构转为一维的列表结构
// 此方法主要用于将 apicat 的 Response 结构转为 openapi 支持的 Response 列表结构
func (h *HTTPResponseDefine) ItemsTreeToList() (res HTTPResponseDefines) {
	if h.Type != string(CollectionItemTypeDir) {
		return append(res, h)
	}
	return h.itemsTreeToList(h.Name)
}

func (h *HTTPResponseDefine) itemsTreeToList(path string) (res HTTPResponseDefines) {
	if h.Items == nil || len(h.Items) == 0 {
		return res
	}

	for _, item := range h.Items {
		if item.Type == string(CollectionItemTypeDir) {
			res = append(res, item.itemsTreeToList(fmt.Sprintf("%s/%s", path, item.Name))...)
		} else {
			item.Category = path
			res = append(res, item)
		}
	}
	return res
}

func (h *HTTPResponseDefine) makeSelfTree(path string, category map[string]*HTTPResponseDefine) *HTTPResponseDefine {
	if path == "" {
		return h
	}
	i := strings.Index(path, "/")
	if i == -1 {
		parent := &HTTPResponseDefine{
			Name:  path,
			Items: HTTPResponseDefines{h},
			Type:  string(CollectionItemTypeDir),
		}
		category[path] = parent
		return parent
	}
	parent := &HTTPResponseDefine{
		Name:  path[:i],
		Items: HTTPResponseDefines{h.makeSelfTree(path[i+1:], category)},
		Type:  string(CollectionItemTypeDir),
	}
	category[path] = parent
	return parent
}

func (h *HTTPResponseDefine) SetXDiff(x *string) {
	h.Header.SetXDiff(x)
	h.Content.SetXDiff(x)
	h.XDiff = x
}

type HTTPResponseDefines []*HTTPResponseDefine

func (h HTTPResponseDefines) LookupByName(name string) *HTTPResponseDefine {
	for _, v := range h {
		if v.Type == string(CollectionItemTypeDir) {
			if res := v.Items.LookupByName(name); res != nil {
				return res
			}
		} else {
			if v.Name == name {
				return v
			}
		}
	}
	return nil
}

func (h HTTPResponseDefines) LookupByID(id int64) *HTTPResponseDefine {
	for _, v := range h {
		if v.Type == string(CollectionItemTypeDir) {
			if res := v.Items.LookupByID(id); res != nil {
				return res
			}
		} else {
			if v.ID == id {
				return v
			}
		}
	}
	return nil
}

func (h *HTTPResponseDefines) SetXDiff(x *string) {
	for _, v := range *h {
		v.SetXDiff(x)
	}
}

func (h *HTTPResponseDefines) ItemsListToTree() HTTPResponseDefines {
	root := &HTTPResponseDefine{
		Items: HTTPResponseDefines{},
	}

	if h == nil || len(*h) == 0 {
		return root.Items
	}

	category := map[string]*HTTPResponseDefine{
		"": root,
	}

	for _, v := range *h {
		if parent, ok := category[v.Category]; ok {
			parent.Items = append(parent.Items, v)
		} else {
			root.Items = append(root.Items, v.makeSelfTree(v.Category, category))
		}
	}

	return root.Items
}
