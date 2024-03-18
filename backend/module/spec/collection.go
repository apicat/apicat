package spec

import (
	"errors"
	"fmt"
)

// CollectItem 集合中的每一项结构定义
type Collection struct {
	ID       uint           `json:"id,omitempty" yaml:"id,omitempty"`
	ParentID uint           `json:"parentid,omitempty" yaml:"parentid,omitempty"`
	Title    string         `json:"title" yaml:"title"`
	Type     CollectionType `json:"type" yaml:"type"`
	Tags     []string       `json:"tag,omitempty" yaml:"tag,omitempty"`
	Content  []*NodeProxy   `json:"content,omitempty" yaml:"content,omitempty"`
	Items    []*Collection  `json:"items,omitempty" yaml:"items,omitempty"`
	XDiff    *string        `json:"x-apicat-diff,omitempty" yaml:"x-apicat-diff,omitempty"`
}

func (v *Collection) HasTag(tag string) bool {
	for _, t := range v.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// @title DerefResponses
// @description 将 API 中，所有引用了入参 Definition Response List 中的 Definition Response 全部解引用
// @param id int64 Definition Response ID
// @return error
func (c *Collection) DerefResponse(wantToDeref ...*HTTPResponseDefine) error {
	if c == nil {
		return errors.New("collect item is nil")
	}
	if wantToDeref == nil {
		return errors.New("wantToDeref is nil")
	}
	// if it type is "category", just return nil
	if c.Type == CollectionItemTypeDir {
		return errors.New("collect item type is dir")
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		// just this type to reference response
		case NAME_HTTP_RESPONSES:
			resps, err := node.ToHTTPResponsesNode()
			if err != nil {
				return err
			}
			err = resps.Deref(wantToDeref)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// @title DelResponseByRefId
// @description 删除 API 中的响应，如果引用了此 ID 对应的响应，将响应从 API 中删除
// @param id int64 Definition Response ID
// @return error
func (c *Collection) DelResponseByRefId(id int64) error {
	if c == nil {
		return errors.New("collect item is nil")
	}

	if c.Type == CollectionItemTypeDir {
		return errors.New("collect item type is dir")
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NAME_HTTP_RESPONSES:
			resps, err := node.ToHTTPResponsesNode()
			if err != nil {
				return err
			}
			err = resps.DelResponse(id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Collection) DerefSchema(wantToDeref ...*Schema) error {
	if c == nil {
		return errors.New("collect item is nil")
	}

	if c.Type == CollectionItemTypeDir {
		return errors.New("collect item type is dir")
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NAME_HTTP_REQUEST:
			req, err := node.ToHTTPRequestNode()
			if err != nil {
				return err
			}
			req.Deref(wantToDeref...)
		case NAME_HTTP_RESPONSES:
			resps, err := node.ToHTTPResponsesNode()
			if err != nil {
				return err
			}
			resps.DerefSchema(wantToDeref...)
		}
	}
	return nil
}

func (c *Collection) DelRefSchema(schema *Schema) error {
	if c == nil {
		return errors.New("collect item is nil")
	}

	if c.Type == CollectionItemTypeDir {
		return errors.New("collect item type is dir")
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NAME_HTTP_REQUEST:
			req, err := node.ToHTTPRequestNode()
			if err != nil {
				return err
			}
			req.DelRef(schema)
		case NAME_HTTP_RESPONSES:
			resps, err := node.ToHTTPResponsesNode()
			if err != nil {
				return err
			}
			resps.DelRefSchema(schema)
		}

	}
	return nil
}

// @title DelGlobalExceptID
// @description 删除原本排除的全局参数 ID
// @param in string 全局参数所在的问题 header cookie path
// @param id int64 全局参数 ID
// @return error
func (c *Collection) DelGlobalExceptID(in string, id int64) error {
	if c == nil {
		return errors.New("collect item is nil")
	}

	if c.Type == CollectionItemTypeDir {
		return errors.New("collect item type is dir")
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NAME_HTTP_REQUEST:
			req, err := node.ToHTTPRequestNode()
			if err != nil {
				return err
			}

			req.DelGlobalExceptID(in, id)
		}
	}
	return nil
}

func (c *Collection) GetGlobalExcept(in string) ([]int64, error) {
	if c == nil {
		return nil, errors.New("collection is nil")
	}

	if c.Type == CollectionItemTypeDir {
		return nil, errors.New("collection type is dir")
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NAME_HTTP_REQUEST:
			req, err := node.ToHTTPRequestNode()
			if err != nil {
				return nil, err
			}
			if req.GlobalExcepts == nil {
				return nil, errors.New("GlobalExcepts is nil")
			}
			if v, exist := req.GlobalExcepts[in]; exist {
				return v, nil
			}
			return nil, fmt.Errorf("%s not found", in)
		}
	}
	return nil, errors.New("requested part not found")
}

func (c *Collection) AddParameter(in string, parameter *Parameter) error {
	if c == nil {
		return errors.New("collect item is nil")
	}

	if c.Type == CollectionItemTypeDir {
		return errors.New("collect item type is dir")
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NAME_HTTP_REQUEST:
			req, err := node.ToHTTPRequestNode()
			if err != nil {
				return err
			}
			req.Parameters.Add(in, parameter)
		}
	}
	return nil
}

// @title DelParameterById
// @description 取消原本排除的全局参数
// @param in string 参数所在的问题 header cookie path query
// @param id int64 参数的 ID
// @return error
func (c *Collection) DelParameterById(in string, id int64) error {
	if c == nil {
		return errors.New("collect item is nil")
	}

	if c.Type == CollectionItemTypeDir {
		return errors.New("collect item type is dir")
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NAME_HTTP_REQUEST:
			req, err := node.ToHTTPRequestNode()
			if err != nil {
				return err
			}
			req.Parameters.Del(in, id)
		}
	}
	return nil
}

func (c *Collection) WithoutRef(global *Global, definition *Definitions) error {
	if c == nil {
		return errors.New("collection is nil")
	}

	// 先解 Response 再解 Schema 否则 Response 中引用的 Schema 不会被解
	if err := c.DerefResponse(definition.Responses...); err != nil {
		return err
	}
	if err := c.DerefSchema(definition.Schemas...); err != nil {
		return err
	}
	if len(global.Parameters.Header) > 0 {
		exceptIDs, err := c.GetGlobalExcept("header")
		if err != nil {
			return err
		}
		if len(exceptIDs) > 0 {
			idMap := make(map[int64]bool)
			for _, id := range exceptIDs {
				idMap[id] = true
			}
			for _, v := range global.Parameters.Header {
				if _, exist := idMap[v.ID]; !exist {
					if err := c.AddParameter("header", v); err != nil {
						return err
					}
				}
			}
		} else {
			for _, v := range global.Parameters.Header {
				if err := c.AddParameter("header", v); err != nil {
					return err
				}
			}
		}
	}
	if len(global.Parameters.Query) > 0 {
		exceptIDs, err := c.GetGlobalExcept("query")
		if err != nil {
			return err
		}
		if len(exceptIDs) > 0 {
			idMap := make(map[int64]bool)
			for _, id := range exceptIDs {
				idMap[id] = true
			}
			for _, v := range global.Parameters.Query {
				if _, exist := idMap[v.ID]; !exist {
					if err := c.AddParameter("query", v); err != nil {
						return err
					}
				}
			}
		} else {
			for _, v := range global.Parameters.Query {
				if err := c.AddParameter("query", v); err != nil {
					return err
				}
			}
		}
	}
	if len(global.Parameters.Cookie) > 0 {
		exceptIDs, err := c.GetGlobalExcept("cookie")
		if err != nil {
			return err
		}
		if len(exceptIDs) > 0 {
			idMap := make(map[int64]bool)
			for _, id := range exceptIDs {
				idMap[id] = true
			}
			for _, v := range global.Parameters.Cookie {
				if _, exist := idMap[v.ID]; !exist {
					if err := c.AddParameter("cookie", v); err != nil {
						return err
					}
				}
			}
		} else {
			for _, v := range global.Parameters.Cookie {
				if err := c.AddParameter("cookie", v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *Collection) ItemsTreeToList() Collections {
	list := make(Collections, 0)
	if c.Type != CollectionItemTypeDir {
		return append(list, c)
	}
	c.itemsTreeToList(&list)
	return list
}
func (c *Collection) itemsTreeToList(list *Collections) {
	if c.Items == nil || len(c.Items) == 0 {
		return
	}

	for _, item := range c.Items {
		if item.Type == CollectionItemTypeDir {
			item.itemsTreeToList(list)
		} else {
			*list = append(*list, item)
		}
	}
}
