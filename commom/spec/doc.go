package spec

import (
	"fmt"
)

// DocNode 文档节点
// https://prosemirror.net/docs/guide/#doc
type DocNode struct {
	Type    string         `json:"type"`
	Content []*DocNode     `json:"content,omitempty"`
	Text    string         `json:"text,omitempty"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Mark    []*DocNode     `json:"mark,omitempty"`
}

func (*DocNode) NodeType() string {
	return "doc"
}

// Docment node节点的集合文档
type Document struct {
	Items []*DocNode
}

func (d *DocNode) getAttr(k string) any {
	if d.Attrs != nil {
		return d.Attrs[k]
	}
	return nil
}

func (d *DocNode) LookupAttrString(k string) string {
	v := d.getAttr(k)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func (d *DocNode) LookupAttrNumber(k string) float64 {
	v := d.getAttr(k)
	if v != nil {
		switch x := v.(type) {
		case float64:
			return x
		case float32:
			return float64(x)
		case int:
			return float64(x)
		case int32:
			return float64(x)
		case int64:
			return float64(x)
		}
	}
	return 0
}
