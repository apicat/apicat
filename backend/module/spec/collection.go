package spec

import (
	"encoding/json"

	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

const (
	TYPE_HTTP = "http"
	TYPE_DOC  = "doc"
)

type Collection struct {
	ID       int64           `json:"id,omitempty" yaml:"id,omitempty"`
	ParentID int64           `json:"parentid,omitempty" yaml:"parentid,omitempty"`
	Title    string          `json:"title" yaml:"title"`
	Type     string          `json:"type" yaml:"type"`
	Content  CollectionNodes `json:"content,omitempty" yaml:"content,omitempty"`
	Tags     []string        `json:"tag,omitempty" yaml:"tag,omitempty"`
	Items    Collections     `json:"items,omitempty" yaml:"items,omitempty"`
}

type Collections []*Collection

func NewCollection(title, typ string) *Collection {
	return &Collection{
		Title: title,
		Type:  typ,
	}
}

func NewCollectionFromJson(c string) (*Collection, error) {
	var collection Collection
	if err := json.Unmarshal([]byte(c), &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

func (v *Collection) HasTag(tag string) bool {
	for _, t := range v.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (c *Collection) ItemsTreeToList() Collections {
	list := make(Collections, 0)
	if c.Type != TYPE_CATEGORY {
		return append(list, c)
	}

	if c.Items != nil && len(c.Items) > 0 {
		for _, item := range c.Items {
			list = append(list, item.ItemsTreeToList()...)
		}
	}
	return list
}

func (c *Collection) ToJson() (string, error) {
	res, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (cs *Collections) DeepDerefAll(params *GlobalParameters, definitions *Definitions) error {
	helper := jsonschema.NewDerefHelper(definitions.Schemas.ToJsonSchemaMap())

	for _, c := range *cs {
		if c.Type == TYPE_CATEGORY {
			continue
		}

		for _, node := range c.Content {
			switch node.NodeType() {
			case NODE_HTTP_REQUEST:
				node.ToHttpRequest().DerefGlobalParameters(params)
				if err := node.ToHttpRequest().DeepDerefModelByHelper(helper); err != nil {
					return err
				}
			case NODE_HTTP_RESPONSE:
				res := node.ToHttpResponse()
				if err := res.DerefAllResponses(definitions.Responses); err != nil {
					return err
				}
				if err := res.DeepDerefModelByHelper(helper); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
