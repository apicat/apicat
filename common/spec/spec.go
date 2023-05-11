package spec

import (
	"bytes"
	"encoding/json"
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

// type Common struct {
// 	Parameters Schemas       `json:"parameters,omitempty"`
// 	Responses  HTTPResponses `json:"response,omitempty"`
// }

type Definitions struct {
	Schemas    Schemas             `json:"schemas"`
	Parameters Schemas             `json:"parameters"`
	Responses  HTTPResponseDefines `json:"responses"`
}
