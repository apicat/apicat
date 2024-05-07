package spec2

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

func (c *Collection) DerefModel(ref *Model) {
	if c == nil {
		return
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().Deref(ref)
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().DerefModel(ref)
		}
	}
}

func (c *Collection) DerefResponse(ref *DefinitionResponse) {
	if c == nil {
		return
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().DerefResponse(ref)
		}
	}
}

func (c *Collection) DelRefModel(ref *Model) {
	if c == nil {
		return
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DelRef(ref)
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().DelRefModel(ref)
		}
	}
}

func (c *Collection) DelRefResponse(ref *DefinitionResponse) {
	if c == nil {
		return
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().DelRefResponse(ref)
		}
	}
}

func (c *Collection) DelGlobalExcept(in string, id int64) {
	if c == nil {
		return
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DelGlobalExcept(in, id)
		}
	}
}

func (c *Collection) GetGlobalExcept(in string) []int64 {
	if c == nil {
		return nil
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			return node.ToHttpRequest().GetGlobalExcept(in)
		}
	}
	return nil
}

func (c *Collection) GetGlobalExceptAll() map[string][]int64 {
	if c == nil {
		return nil
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			return node.ToHttpRequest().GetGlobalExceptAll()
		}
	}
	return nil
}

func (c *Collection) AddGlobalExcept(in string, id int64) {
	if c == nil {
		return
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().AddGlobalExcept(in, id)
		}
	}
}
