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

func (c *Collection) DerefModel(ref *DefinitionModel) error {
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			return node.ToHttpRequest().DerefModel(ref)
		case NODE_HTTP_RESPONSE:
			return node.ToHttpResponse().DerefModel(ref)
		}
	}

	return nil
}

func (c *Collection) DerefResponse(ref *DefinitionResponse) error {
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			return node.ToHttpResponse().DerefResponse(ref)
		}
	}
	return nil
}

func (c *Collection) DeepDerefAll(params *GlobalParameters, definitions *Definitions) error {
	helper := jsonschema.NewDerefHelper(definitions.Schemas.ToJsonSchemaMap())

	c.DerefGlobalParameters(params)

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
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
	return nil
}

func (c *Collection) DelRefModel(ref *DefinitionModel) {
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DelRefModel(ref)
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().DelRefModel(ref)
		}
	}
}

func (c *Collection) DelRefResponse(ref *DefinitionResponse) {
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().DelRefResponse(ref)
		}
	}
}

func (c *Collection) DerefGlobalParameter(in string, param *Parameter) {
	if param == nil {
		return
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().Attrs.Parameters.Add(in, param)
		}
	}
}

func (c *Collection) DerefGlobalParameters(params *GlobalParameters) {
	if params == nil {
		return
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DerefGlobalParameters(params)
		}
	}
}

func (c *Collection) DelGlobalExcept(in string, id int64) {
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DelGlobalExcept(in, id)
		}
	}
}

func (c *Collection) GetGlobalExcept(in string) []int64 {
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			return node.ToHttpRequest().GetGlobalExcept(in)
		}
	}
	return nil
}

func (c *Collection) GetGlobalExceptAll() map[string][]int64 {
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			return node.ToHttpRequest().GetGlobalExceptAll()
		}
	}
	return nil
}

func (c *Collection) GetRefModelIDs() []int64 {
	ids := make([]int64, 0)
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			ids = append(ids, node.ToHttpRequest().GetRefModelIDs()...)
		case NODE_HTTP_RESPONSE:
			ids = append(ids, node.ToHttpResponse().GetRefModelIDs()...)
		}
	}
	return ids
}

func (c *Collection) GetRefResponseIDs() []int64 {
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			return node.ToHttpResponse().GetRefResponseIDs()
		}
	}
	return nil
}

func (c *Collection) AddGlobalExcept(in string, id int64) {
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().AddGlobalExcept(in, id)
		}
	}
}

func (c *Collection) AddReqParameter(in string, p *Parameter) {
	if p == nil {
		return
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().Attrs.Parameters.Add(in, p)
		}
	}
}

func (c *Collection) SortResponses() {
	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().Sort()
		}
	}
}

func (v *Collection) HasTag(tag string) bool {
	for _, t := range v.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Make sure the order of the content is: URL, Request, Response
func (c *Collection) HttpFormat() {
	if c.Type != TYPE_HTTP || len(c.Content) < 3 {
		return
	}

	for i, node := range c.Content {
		if node.NodeType() == NODE_HTTP_URL {
			c.Content[0], c.Content[i] = c.Content[i], c.Content[0]
		}
		if node.NodeType() == NODE_HTTP_REQUEST {
			c.Content[1], c.Content[i] = c.Content[i], c.Content[1]
		}
		if node.NodeType() == NODE_HTTP_RESPONSE {
			c.Content[2], c.Content[i] = c.Content[i], c.Content[2]
		}
	}
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

		c.DerefGlobalParameters(params)

		for _, node := range c.Content {
			switch node.NodeType() {
			case NODE_HTTP_REQUEST:
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
