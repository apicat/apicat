package spec

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/apicat/apicat/backend/common/spec/jsonschema"
)

// Content 文档类型
type ContentType string

const (
	ContentItemTypeDir  ContentType = "category"
	ContentItemTypeHttp             = "http"
	ContentItemTypeDoc              = "doc"
)

func init() {
	RegisterNode(&DocNode{})
	RegisterNode(WarpHTTPNode(HTTPURLNode{}))
	RegisterNode(WarpHTTPNode(HTTPRequestNode{}))
	RegisterNode(WarpHTTPNode(HTTPResponsesNode{}))
}

// Spec 是apicat的协议的整体结构
type Spec struct {
	// spec schema版本 当前固定2.0
	ApiCat      string      `json:"apicat"`
	Info        *Info       `json:"info"`
	Servers     []*Server   `json:"servers"`
	Globals     Global      `json:"globals"`
	Definitions Definitions `json:"definitions"`
	Collections Collections `json:"collections"`
}

type Collections []*CollectItem

// WalkCollections 遍历集合
// 如果回调函数返回false等退出 否则继续遍历
func (s *Spec) WalkCollections(f func(*CollectItem, []string) bool) {
	walkCollectionHandle(s.Collections, []string{}, f)
}

func refDept(r string, refs ...string) int {
	var n int
	for _, v := range refs {
		if v == r {
			n++
		}
	}
	return n
}

func (s *Spec) ExpendRef(v Referencer, max int, parentRef ...string) {
	s.expendRef(v, max, parentRef...)
}

func (s *Spec) expendRef(v Referencer, max int, parentRef ...string) {
	if v == nil {
		return
	}
	switch x := v.(type) {
	case *jsonschema.Schema:
		if x == nil {
			return
		}
		if x.Ref() {
			// this feature is important
			// 重复出现的递归次数
			count := refDept(*x.Reference, parentRef...)
			if count > max {
				return
			}
			refv := s.Definitions.Schemas.LookupID(mustGetRefID(*x.Reference))
			if refv == nil || refv.Schema == nil {
				return
			}
			parentRef = append(parentRef, *x.Reference)
			b, _ := json.Marshal(*refv.Schema)
			var sc jsonschema.Schema
			json.Unmarshal(b, &sc)
			*x = sc
		}
		if x.Properties != nil {
			for _, p := range x.Properties {
				s.expendRef(p, max, parentRef...)
			}
		}
		if x.Items != nil && !x.Items.IsBool() {
			v := *x.Items.Value()
			s.expendRef(&v, max, parentRef...)
			x.Items.SetValue(&v)
		}
	case *Schema:
		if x.Ref() {
			ps := strings.Split(*x.Reference, "/")
			if len(ps) == 4 {
				id, _ := strconv.ParseInt(ps[3], 10, 64)
				switch ps[2] {
				case "schemas":
					if refv := s.Definitions.Schemas.LookupID(id); refv != nil {
						*x = *refv
					}
				case "parameters":
					if refv := s.Definitions.Parameters.LookupID(id); refv != nil {
						*x = *refv
					}
				}
			}
		}
		s.expendRef(x.Schema, max, parentRef...)
	case *HTTPResponseDefine:
		if x.Ref() {
			refv := s.Definitions.Responses.LookupID(mustGetRefID(*x.Reference))
			if refv == nil {
				*x = HTTPResponseDefine{}
			} else {
				*x = *refv
			}
		}
		for k := range x.Header {
			s.expendRef(x.Header[k], max, parentRef...)
		}
		for k := range x.Content {
			s.expendRef(x.Content[k], max, parentRef...)
		}
	}
}

// CollectionsMap 返回map结构的api路由定义
// expend 是否展开内部引用
// refexpendMaxCount 引用最大解开次数 防止递归引用死循环
func (s *Spec) CollectionsMap(expend bool, refexpendMaxCount int) map[string]map[string]HTTPPart {
	paths := map[string]map[string]HTTPPart{}
	s.WalkCollections(func(v *CollectItem, p []string) bool {
		if v.Type != ContentItemTypeHttp {
			return true
		}
		var (
			method string
			path   string
		)
		part := HTTPPart{
			Title: v.Title,
			ID:    v.ID,
			Dir:   strings.Join(p, "/"),
		}
		for _, item := range v.Content {
			switch nx := item.Node.(type) {
			case *HTTPNode[HTTPURLNode]:
				method, path = nx.Attrs.Method, nx.Attrs.Path
			case *HTTPNode[HTTPRequestNode]:
				nx.Attrs.Parameters.Fill()
				if expend {
					mp := nx.Attrs.Parameters.Map()
					for _, v := range mp {
						for k := range v {
							s.expendRef(v[k], refexpendMaxCount)
						}
					}
					if nx.Attrs.Content != nil {
						for k := range nx.Attrs.Content {
							s.expendRef(nx.Attrs.Content[k], refexpendMaxCount)
						}
					}
				}
				part.HTTPRequestNode = nx.Attrs
			case *HTTPNode[HTTPResponsesNode]:
				res := nx.Attrs.List
				if expend {
					for i, v := range res {
						s.expendRef(&v.HTTPResponseDefine, refexpendMaxCount)
						res[i] = v
					}
				}
				part.Responses = res
			}
		}
		subs, ok := paths[path]
		if !ok {
			subs = make(map[string]HTTPPart)
		}
		subs[method] = part
		paths[path] = subs
		return true
	})
	return paths
}

func (s *Spec) Valid() error {
	return nil
}

// JSONOption Spec进行json输出时的配置
type JSONOption struct {
	EscapeHTML bool
	Indent     string
}

// ToJSON 转为json格式
func (s *Spec) ToJSON(opt JSONOption) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(opt.EscapeHTML)
	if opt.Indent != "" {
		enc.SetIndent("", opt.Indent)
	}
	if err := enc.Encode(s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ParseJSON 将json数据转为spec
func ParseJSON(raw []byte) (*Spec, error) {
	var spec Spec
	if err := json.Unmarshal(raw, &spec); err != nil {
		return nil, err
	}
	if spec.ApiCat == "" {
		spec.ApiCat = "2.0"
	}
	return &spec, nil
}

// CollectItem 集合中的每一项结构定义
type CollectItem struct {
	Type     ContentType    `json:"type"`
	ID       int64          `json:"id,omitempty"`
	ParentID int64          `json:"parentid,omitempty"`
	Title    string         `json:"title"`
	Tags     []string       `json:"tag,omitempty"`
	Content  []*NodeProxy   `json:"content,omitempty"`
	Items    []*CollectItem `json:"items,omitempty"`
	XDiff    *string        `json:"x-apicat-diff,omitempty"`
}

func walkCollectionHandle(list []*CollectItem, p []string,
	f func(*CollectItem, []string) bool) bool {
	for _, v := range list {
		switch v.Type {
		case ContentItemTypeDir:
			if !walkCollectionHandle(v.Items, append(p, v.Title), f) {
				return false
			}
		default:
			if !f(v, p) {
				return false
			}
		}
	}
	return true
}

func (v *CollectItem) HasTag(tag string) bool {
	for _, t := range v.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (c *CollectItem) DereferenceResponse(sub *HTTPResponseDefine) error {
	if c == nil {
		return nil
	}

	// if it type is "category", just return nil
	if c.Type == ContentItemTypeDir {
		return nil
	}

	// just range 3 times
	for _, node := range c.Content {
		switch node.NodeType() {
		// just this type to reference response
		case "apicat-http-response":
			resps, err := node.ToHTTPResponsesNode()
			if err != nil {
				return err
			}
			resps.DereferenceResponses(sub)
		}
	}
	return nil
}

func (c *CollectItem) RemoveResponse(s_id int64) error {
	if c == nil {
		return nil
	}

	if c.Type == ContentItemTypeDir {
		return nil
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case "apicat-http-response":
			resps, err := node.ToHTTPResponsesNode()
			if err != nil {
				return err
			}
			resps.RemoveResponse(s_id)
		}
	}
	return nil
}

func (c *CollectItem) DereferenceSchema(sub *Schema) error {
	if c == nil {
		return nil
	}

	if c.Type == ContentItemTypeDir {
		return nil
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case "apicat-http-response":
			resps, err := node.ToHTTPResponsesNode()
			if err != nil {
				return err
			}
			resps.DereferenceSchema(sub)
		}
	}
	return nil
}

func (c *CollectItem) RemoveSchema(s_id int64) error {
	if c == nil {
		return nil
	}

	if c.Type == ContentItemTypeDir {
		return nil
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case "apicat-http-response":
			resps, err := node.ToHTTPResponsesNode()
			if err != nil {
				return err
			}
			resps.RemoveSchema(s_id)
		}

	}
	return nil
}

func (c *CollectItem) DereferenceGlobalParameters(in string, sub *Schema) error {
	if c == nil {
		return nil
	}

	if c.Type == ContentItemTypeDir {
		return nil
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case "apicat-http-request":
			req, err := node.ToHTTPRequestNode()
			if err != nil {
				return err
			}

			ok := req.tryRemoveGlobalExcept(in, sub.ID)
			if ok {
				return nil
			} else {
				req.Parameters.Add(in, sub)
			}
		}
	}
	return nil
}

func (c *CollectItem) OpenGlobalParameters(in string, s_id int64) error {
	if c == nil {
		return nil
	}

	if c.Type == ContentItemTypeDir {
		return nil
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case "apicat-http-request":
			req, err := node.ToHTTPRequestNode()
			if err != nil {
				return err
			}
			req.RemoveGlobalExcept(in, s_id)
		}
	}
	return nil
}

func (c *CollectItem) AddParameters(in string, sub *Schema) error {
	if c == nil {
		return nil
	}

	if c.Type == ContentItemTypeDir {
		return nil
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case "apicat-http-request":
			req, err := node.ToHTTPRequestNode()
			if err != nil {
				return err
			}
			req.Parameters.Add(in, sub)
		}
	}
	return nil
}

func (c *CollectItem) RemoveParameters(in string, s_id int64) error {
	if c == nil {
		return nil
	}

	if c.Type == ContentItemTypeDir {
		return nil
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case "apicat-http-request":
			req, err := node.ToHTTPRequestNode()
			if err != nil {
				return err
			}
			req.Parameters.Remove(in, s_id)
		}
	}
	return nil
}

type Info struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Global struct {
	Parameters HTTPParameters `json:"parameters,omitempty"`
}

// this type is only used to represent Definition item
type DefinitionType string

var (
	// if item's type is "response", then it is a definition response
	TypeResponse DefinitionType = "response"
	// if item's type is "category", then it is a category, it's items property is a list of response or schema, like it self
	TypeCategory DefinitionType = "category"
	// if item's type is "schema", then it is a definition schema
	TypeSchema DefinitionType = "schema"
)

type Definitions struct {
	Schemas    Schemas             `json:"schemas"`
	Parameters HTTPParameters      `json:"parameters"`
	Responses  HTTPResponseDefines `json:"responses"`
}

func mustGetRefID(v string) int64 {
	ps := strings.Split(v, "/")
	id, _ := strconv.ParseInt(ps[len(ps)-1], 10, 64)
	return id
}
