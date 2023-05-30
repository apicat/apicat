package spec

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/apicat/apicat/common/spec/jsonschema"
)

// Content 文档类型
type ContentType string

const (
	ContentItemTypeDir  ContentType = "dir"
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
	ApiCat      string         `json:"apicat"`
	Info        *Info          `json:"info"`
	Servers     []*Server      `json:"servers"`
	Globals     Global         `json:"globals"`
	Definitions Definitions    `json:"definitions"`
	Collections []*CollectItem `json:"collections"`
}

// WalkCollections 遍历集合
// 如果回调函数返回false等退出 否则继续遍历
func (s *Spec) WalkCollections(f func(v *CollectItem) bool) {
	walkCollectionHandle(s.Collections, f)
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

func (s *Spec) expendRef(v Referencer, max int, parentRef ...string) {
	if v == nil {
		return
	}
	switch x := v.(type) {
	case *jsonschema.Schema:
		if x.Ref() {
			// 重复出现的递归次数
			count := refDept(*x.Reference, parentRef...)
			if count > max {
				return
			}
			parentRef = append(parentRef, *x.Reference)
			*x = *(s.Definitions.Schemas.LookupID(mustGetRefID(*x.Reference)).Schema)
		}
		if x.Properties != nil {
			for _, v := range x.Properties {
				s.expendRef(v, max, parentRef...)
			}
		}
		if x.Items != nil && !x.Items.IsBool() {
			s.expendRef(x.Items.Value(), max, parentRef...)
		}
	case *Schema:
		if x.Ref() {
			ps := strings.Split(*x.Reference, "/")
			if len(ps) == 4 {
				id, _ := strconv.ParseInt(ps[3], 10, 64)
				switch ps[2] {
				case "schemas":
					*x = *(s.Definitions.Schemas.LookupID(id))
				case "parameters":
					*x = *(s.Definitions.Parameters.LookupID(id))
				}
			}
		}
		s.expendRef(x.Schema, max, parentRef...)
	case *HTTPResponseDefine:
		if x.Ref() {
			*x = *(s.Definitions.Responses.LookupID(mustGetRefID(*x.Reference)))
		}
		for _, v := range x.Header {
			s.expendRef(v, max, parentRef...)
		}
		for _, v := range x.Content {
			s.expendRef(v, max, parentRef...)
		}
	}
}

// CollectionsMap 返回map结构的api路由定义
// expend 是否展开内部引用
// refexpendMaxCount 引用最大解开次数 防止递归引用死循环
func (s *Spec) CollectionsMap(expend bool, refexpendMaxCount int) map[string]map[string]HTTPPart {
	paths := map[string]map[string]HTTPPart{}
	s.WalkCollections(func(v *CollectItem) bool {
		if v.Type != ContentItemTypeHttp {
			return true
		}
		var (
			method string
			path   string
			part   HTTPPart
		)
		for _, item := range v.Content {
			switch nx := item.Node.(type) {
			case *HTTPNode[HTTPURLNode]:
				method, path = nx.Attrs.Method, nx.Attrs.Path
			case *HTTPNode[HTTPRequestNode]:
				if expend {
					mp := nx.Attrs.Parameters.Map()
					for _, v := range mp {
						for _, a := range v {
							s.expendRef(a, refexpendMaxCount)
						}
					}
					if nx.Attrs.Content != nil {
						for _, a := range nx.Attrs.Content {
							s.expendRef(a, refexpendMaxCount)
						}
					}
				}
				part.HTTPRequestNode = nx.Attrs
			case *HTTPNode[HTTPResponsesNode]:
				res := nx.Attrs.List
				if expend {
					for _, v := range res {
						s.expendRef(&v.HTTPResponseDefine, refexpendMaxCount)
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
}

func walkCollectionHandle(list []*CollectItem,
	f func(v *CollectItem) bool) bool {
	for _, v := range list {
		switch v.Type {
		case ContentItemTypeDir:
			if !walkCollectionHandle(v.Items, f) {
				return false
			}
		default:
			if !f(v) {
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

type Definitions struct {
	Schemas    Schemas             `json:"schemas"`
	Parameters Schemas             `json:"parameters"`
	Responses  HTTPResponseDefines `json:"responses"`
}

func mustGetRefID(v string) int64 {
	ps := strings.Split(v, "/")
	id, _ := strconv.ParseInt(ps[len(ps)-1], 10, 64)
	return id
}
