package spec2

import "fmt"

const NODE_DOC = "doc"

type CollectionDoc struct {
	Type    string           `json:"type"`
	Content []*CollectionDoc `json:"content,omitempty"`
	Text    string           `json:"text,omitempty"`
	Attrs   map[string]any   `json:"attrs,omitempty"`
	Mark    []*CollectionDoc `json:"mark,omitempty"`
}

type Document struct {
	Items []*CollectionDoc
}

func (d *CollectionDoc) NodeType() string {
	return d.Type
}

func (d *CollectionDoc) getAttr(k string) any {
	if d.Attrs != nil {
		return d.Attrs[k]
	}
	return nil
}

func (d *CollectionDoc) GetAttrString(k string) string {
	v := d.getAttr(k)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func (d *CollectionDoc) GetAttrNumber(k string) float64 {
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
