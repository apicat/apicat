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

func (c *Collection) DerefModel(ref *DefinitionModel) {
	if c == nil {
		return
	}

	for _, node := range c.Content {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DerefModel(ref)
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

func (c *Collection) DelRefModel(ref *DefinitionModel) {
	if c == nil {
		return
	}

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

func (c *Collection) DerefGlobalParameters(params *GlobalParameters) {
	if c == nil || params == nil {
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

// Make sure the order of the content is: URL, Request, Response
func (c *Collection) HttpFormat() {
	if c == nil || c.Type != TYPE_HTTP || len(c.Content) < 3 {
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

func (c *Collection) DeepDerefAllDefinitions(refModels DefinitionModels, refResponses DefinitionResponses) {
	if c == nil {
		return
	}
}
