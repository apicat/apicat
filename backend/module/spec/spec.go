package spec

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

// Collection 集合类型
type CollectionType string

const (
	CollectionItemTypeDir  CollectionType = "category"
	CollectionItemTypeHttp CollectionType = "http"
	CollectionItemTypeDoc  CollectionType = "doc"
)

type Info struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty"`
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string `json:"version" yaml:"version"`
}

type Server struct {
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type Global struct {
	Parameters *HTTPParameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func NewGlobal() *Global {
	return &Global{
		Parameters: &HTTPParameters{
			Query:  make(ParameterList, 0),
			Header: make(ParameterList, 0),
			Cookie: make(ParameterList, 0),
			Path:   make(ParameterList, 0),
		},
	}
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
	Schemas    Schemas             `json:"schemas" yaml:"schemas"`
	Parameters *HTTPParameters     `json:"parameters" yaml:"parameters"`
	Responses  HTTPResponseDefines `json:"responses" yaml:"responses"`
}

func NewDefinitions() *Definitions {
	return &Definitions{
		Schemas:   make(Schemas, 0),
		Responses: make(HTTPResponseDefines, 0),
		Parameters: &HTTPParameters{
			Query:  make(ParameterList, 0),
			Header: make(ParameterList, 0),
			Cookie: make(ParameterList, 0),
			Path:   make(ParameterList, 0),
		},
	}
}

// Spec 是 apicat 的协议的整体结构
type Spec struct {
	ApiCat      string       `json:"apicat"`
	Info        *Info        `json:"info"`
	Servers     []*Server    `json:"servers"`
	Globals     *Global      `json:"globals"`
	Definitions *Definitions `json:"definitions"`
	Collections Collections  `json:"collections"`
}

type Collections []*Collection

func init() {
	RegisterNode(&DocNode{})
	RegisterNode(WarpHTTPNode(HTTPURLNode{}))
	RegisterNode(WarpHTTPNode(HTTPRequestNode{}))
	RegisterNode(WarpHTTPNode(HTTPResponsesNode{}))
}

func NewSpec() *Spec {
	s := &Spec{
		ApiCat:      "2.0",
		Info:        &Info{},
		Servers:     []*Server{},
		Collections: make(Collections, 0),
	}

	s.Globals = NewGlobal()
	s.Definitions = NewDefinitions()

	return s
}

// WalkCollections 遍历集合
// 如果回调函数返回false等退出 否则继续遍历
func (s *Spec) WalkCollections(f func(*Collection, []string) bool) {
	walkCollectionHandle(s.Collections, []string{}, f)
}

func walkCollectionHandle(list []*Collection, p []string,
	f func(*Collection, []string) bool) bool {
	for _, v := range list {
		switch v.Type {
		case CollectionItemTypeDir:
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

// func refDept(r string, refs ...string) int {
// 	var n int
// 	for _, v := range refs {
// 		if v == r {
// 			n++
// 		}
// 	}
// 	return n
// }

// func (s *Spec) ExpendRef(v Referencer, max int, parentRef ...string) {
// 	s.expendRef(v, max, parentRef...)
// }

// func (s *Spec) expendRef(v Referencer, max int, parentRef ...string) {
// 	if v == nil {
// 		return
// 	}
// 	switch x := v.(type) {
// 	case *jsonschema.Schema:
// 		if x == nil {
// 			return
// 		}
// 		if x.Ref() {
// 			// this feature is important
// 			// 重复出现的递归次数
// 			count := refDept(*x.Reference, parentRef...)
// 			if count > max {
// 				return
// 			}
// 			refv := s.Definitions.Schemas.LookupID(mustGetRefID(*x.Reference))
// 			if refv == nil || refv.Schema == nil {
// 				return
// 			}
// 			parentRef = append(parentRef, *x.Reference)
// 			b, _ := json.Marshal(*refv.Schema)
// 			var sc jsonschema.Schema
// 			json.Unmarshal(b, &sc)
// 			*x = sc
// 		}
// 		if x.Properties != nil {
// 			for _, p := range x.Properties {
// 				s.expendRef(p, max, parentRef...)
// 			}
// 		}
// 		if x.Items != nil && !x.Items.IsBool() {
// 			v := *x.Items.Value()
// 			s.expendRef(&v, max, parentRef...)
// 			x.Items.SetValue(&v)
// 		}
// 	case *Schema:
// 		if x.Ref() {
// 			ps := strings.Split(*x.Reference, "/")
// 			if len(ps) == 4 {
// 				id, _ := strconv.ParseInt(ps[3], 10, 64)
// 				switch ps[2] {
// 				case "schemas":
// 					if refv := s.Definitions.Schemas.LookupID(id); refv != nil {
// 						*x = *refv
// 					}
// 				case "parameters":
// 					if refv := s.Definitions.Parameters.LookupID(id); refv != nil {
// 						*x = *refv
// 					}
// 				}
// 			}
// 		}
// 		s.expendRef(x.Schema, max, parentRef...)
// 	case *HTTPResponseDefine:
// 		if x.Ref() {
// 			refv := s.Definitions.Responses.LookupID(mustGetRefID(*x.Reference))
// 			if refv == nil {
// 				*x = HTTPResponseDefine{}
// 			} else {
// 				*x = *refv
// 			}
// 		}
// 		for k := range x.Header {
// 			s.expendRef(x.Header[k], max, parentRef...)
// 		}
// 		for k := range x.Content {
// 			s.expendRef(x.Content[k], max, parentRef...)
// 		}
// 	}
// }

// CollectionsMap 返回map结构的api路由定义
// expend 是否展开内部引用
// refexpendMaxCount 引用最大解开次数 防止递归引用死循环
func (s *Spec) CollectionsMap(expend bool, refexpendMaxCount int) map[string]map[string]HTTPPart {
	paths := map[string]map[string]HTTPPart{}
	// s.WalkCollections(func(v *CollectItem, p []string) bool {
	// 	if v.Type != CollectionItemTypeHttp {
	// 		return true
	// 	}
	// 	var (
	// 		method string
	// 		path   string
	// 	)
	// 	part := HTTPPart{
	// 		Title: v.Title,
	// 		ID:    v.ID,
	// 		Dir:   strings.Join(p, "/"),
	// 	}
	// 	for _, item := range v.Content {
	// 		switch nx := item.Node.(type) {
	// 		case *HTTPNode[HTTPURLNode]:
	// 			method, path = nx.Attrs.Method, nx.Attrs.Path
	// 		case *HTTPNode[HTTPRequestNode]:
	// 			nx.Attrs.Parameters.Fill()
	// 			if expend {
	// 				mp := nx.Attrs.Parameters.Map()
	// 				for _, v := range mp {
	// 					for k := range v {
	// 						s.expendRef(v[k], refexpendMaxCount)
	// 					}
	// 				}
	// 				if nx.Attrs.Content != nil {
	// 					for k := range nx.Attrs.Content {
	// 						s.expendRef(nx.Attrs.Content[k], refexpendMaxCount)
	// 					}
	// 				}
	// 			}
	// 			part.HTTPRequestNode = nx.Attrs
	// 		case *HTTPNode[HTTPResponsesNode]:
	// 			res := nx.Attrs.List
	// 			if expend {
	// 				for i, v := range res {
	// 					s.expendRef(&v.HTTPResponseDefine, refexpendMaxCount)
	// 					res[i] = v
	// 				}
	// 			}
	// 			part.Responses = res
	// 		}
	// 	}
	// 	subs, ok := paths[path]
	// 	if !ok {
	// 		subs = make(map[string]HTTPPart)
	// 	}
	// 	subs[method] = part
	// 	paths[path] = subs
	// 	return true
	// })
	return paths
}

func (s *Spec) GetPaths() map[string]map[string]HTTPPart {
	paths := map[string]map[string]HTTPPart{}
	cs := Collections{}
	for _, v := range s.Collections {
		if v.Type == CollectionItemTypeDir {
			cs = append(cs, v.ItemsTreeToList()...)
		} else {
			cs = append(cs, v)
		}
	}

	// need tree to list
	ss := make(Schemas, 0)
	for _, v := range s.Definitions.Schemas {
		if v.Type == string(CollectionItemTypeDir) {
			ss = append(ss, v.ItemsTreeToList()...)
		} else {
			ss = append(ss, v)
		}
	}
	resps := []*HTTPResponseDefine{}
	for _, v := range s.Definitions.Responses {
		if v.Type == string(CollectionItemTypeDir) {
			resps = append(resps, v.ItemsTreeToList()...)
		} else {
			resps = append(resps, v)
		}
	}

	definition := &Definitions{
		Responses: resps,
		Schemas:   ss,
	}
	for _, c := range cs {
		c.WithoutRef(s.Globals, definition)
		sortCollectionNodeType(c.Content)
		path, method := "", ""
		part := HTTPPart{
			Title: c.Title,
			ID:    int64(c.ID),
		}
		for _, node := range c.Content {
			switch node.NodeType() {
			case NAME_HTTP_URL:
				url, err := node.ToHTTPURLNode()
				if err != nil {
					return nil
				}
				path = url.Path
				method = url.Method
				if url.Path == "" && url.Method == "" && c.Type == CollectionItemTypeHttp {
					goto e
				}
			case NAME_HTTP_REQUEST:
				req, err := node.ToHTTPRequestNode()
				if err != nil {
					return nil
				}
				part.HTTPRequestNode = *req
			case NAME_HTTP_RESPONSES:
				resps, err := node.ToHTTPResponsesNode()
				if err != nil {
					return nil
				}
				part.Responses = resps.List
			}
		}
		paths[path] = map[string]HTTPPart{method: part}
	e:
	}
	return paths
}

// this func will be sort collection.content by node type, result is [url,request,response,....]
func sortCollectionNodeType(c []*NodeProxy) {
	if len(c) < 3 {
		return
	}
	for i, v := range c {
		if v.NodeType() == NAME_HTTP_URL {
			c[0], c[i] = c[i], c[0]
		}
		if v.NodeType() == NAME_HTTP_REQUEST {
			c[1], c[i] = c[i], c[1]
		}
		if v.NodeType() == NAME_HTTP_RESPONSES {
			c[2], c[i] = c[i], c[2]
		}
	}
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

func mustGetRefID(v string) int64 {
	ps := strings.Split(v, "/")
	id, _ := strconv.ParseInt(ps[len(ps)-1], 10, 64)
	return id
}
