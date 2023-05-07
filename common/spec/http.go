package spec

type HTTPParameters struct {
	Query  []*Schema `json:"query,omitempty"`
	Path   []*Schema `json:"path,omitempty"`
	Cookie []*Schema `json:"cookie,omitempty"`
	Header []*Schema `json:"header,omitempty"`
}

func (h *HTTPParameters) Add(in string, v *Schema) {
	switch in {
	case "query":
		h.Query = append(h.Query, v)
	case "path":
		h.Path = append(h.Path, v)
	case "cookie":
		h.Cookie = append(h.Cookie, v)
	case "header":
		h.Header = append(h.Header, v)
	}
}

func (h *HTTPParameters) Map() map[string][]*Schema {
	m := make(map[string][]*Schema)
	if h.Query != nil {
		m["query"] = h.Query
	}
	if h.Path != nil {
		m["path"] = h.Path
	}
	if h.Header != nil {
		m["header"] = h.Header
	}
	if h.Cookie != nil {
		m["cookie"] = h.Cookie
	}
	return m
}

type HTTPNode[T HTTPNoder] struct {
	Type  string `json:"type"`
	Attrs T      `json:"attrs"`
}

type HTTPNoder interface {
	Name() string
}

func (n *HTTPNode[T]) NodeType() string {
	return n.Type
}

func WarpHTTPNode[T HTTPNoder](n T) Node {
	return &HTTPNode[T]{
		Type:  n.Name(),
		Attrs: n,
	}
}

type HTTPURLNode struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

func (HTTPURLNode) Name() string {
	return "apicat-http-url"
}

type HTTPBody map[string]*Schema

type HTTPRequestNode struct {
	GlobalExcepts map[string][]string `json:"globalExcepts,omitempty"`
	Parameters    HTTPParameters      `json:"parameters,omitempty"`
	Content       HTTPBody            `json:"content,omitempty"`
}

func (HTTPRequestNode) Name() string {
	return "apicat-http-request"
}

type HTTPResponsesNode struct {
	List HTTPResponses `json:"list,omitempty"`
}

func (HTTPResponsesNode) Name() string {
	return "apicat-http-response"
}

type HTTPResponse struct {
	Name        string `json:"name,omitempty"`
	Code        int    `json:"code"`
	Description string `json:"description"`
	HTTPResponseDefine
}

type HTTPResponses []HTTPResponse

func (h HTTPResponses) Lookup(name string) *HTTPResponse {
	for _, v := range h {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

type HTTPResponseDefine struct {
	Name      string   `json:"name,omitempty"`
	Content   HTTPBody `json:"content,omitempty"`
	Header    Schemas  `json:"header,omitempty"`
	Reference *string  `json:"$ref,omitempty"`
}

type HTTPResponseDefines []HTTPResponseDefine

func (h HTTPResponseDefines) Lookup(name string) *HTTPResponseDefine {
	for _, v := range h {
		if v.Name == name {
			return &v
		}
	}
	return nil
}
