package spec

import (
	"strconv"
	"strings"
)

type HTTPParameters struct {
	Query  Schemas `json:"query"`
	Path   Schemas `json:"path"`
	Cookie Schemas `json:"cookie"`
	Header Schemas `json:"header"`
}

var HttpParameter = []string{"query", "path", "cookie", "header"}

// test func
func (h *HTTPParameters) LookupID(id int64) (s *Schema) {
	s = h.Query.LookupID(id)
	if s != nil {
		return s
	}
	s = h.Path.LookupID(id)
	if s != nil {
		return s
	}
	s = h.Cookie.LookupID(id)
	if s != nil {
		return s
	}
	s = h.Header.LookupID(id)
	if s != nil {
		return s
	}
	return nil
}

func (h *HTTPParameters) Fill() {
	if h.Query == nil {
		h.Query = make(Schemas, 0)
	}
	if h.Path == nil {
		h.Path = make(Schemas, 0)
	}
	if h.Cookie == nil {
		h.Cookie = make(Schemas, 0)
	}
	if h.Header == nil {
		h.Header = make(Schemas, 0)
	}
}

func (h *HTTPParameters) Add(in string, v *Schema) {
	switch in {
	case "query":
		if h.Query.LookupID(v.ID) != nil {
			return
		}
		h.Query = append(h.Query, v)
	case "path":
		if h.Path.LookupID(v.ID) != nil {
			return
		}
		h.Path = append(h.Path, v)
	case "cookie":
		if h.Cookie.LookupID(v.ID) != nil {
			return
		}
		h.Cookie = append(h.Cookie, v)
	case "header":
		if h.Header.LookupID(v.ID) != nil {
			return
		}
		h.Header = append(h.Header, v)
	}
}

func (h *HTTPParameters) Remove(in string, id int64) {
	switch in {
	case "query":
		h.Query.RemoveId(id)
	case "path":
		h.Path.RemoveId(id)
	case "cookie":
		h.Cookie.RemoveId(id)
	case "header":
		h.Header.RemoveId(id)
	}
}

func (h *HTTPParameters) Map() map[string]Schemas {
	m := make(map[string]Schemas)
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
	Path   string  `json:"path"`
	Method string  `json:"method"`
	XDiff  *string `json:"x-apicat-diff,omitempty"`
}

func (HTTPURLNode) Name() string {
	return "apicat-http-url"
}

type HTTPBody map[string]*Schema

type HTTPRequestNode struct {
	GlobalExcepts map[string][]int64 `json:"globalExcepts,omitempty"`
	Parameters    HTTPParameters     `json:"parameters,omitempty"`
	Content       HTTPBody           `json:"content,omitempty"`
}

func (HTTPRequestNode) Name() string {
	return "apicat-http-request"
}

// if content is exect, remove it and return true, else return false
func (h *HTTPRequestNode) tryRemoveGlobalExcept(in string, id int64) bool {
	ids := h.GlobalExcepts[in]
	for _, v := range ids {
		if v == id {
			// remove
			h.RemoveGlobalExcept(in, id)
			return true
		}
	}
	return false
}

func (h *HTTPRequestNode) AddGlobalExcept(in string, id int64) {
	if h == nil {
		return
	}
	if findId(h.GlobalExcepts[in], id) {
		return
	}
	h.GlobalExcepts[in] = append(h.GlobalExcepts[in], id)
}

func (h *HTTPRequestNode) RemoveGlobalExcept(in string, id int64) {
	if h == nil {
		return
	}
	h.GlobalExcepts[in] = removeId(h.GlobalExcepts[in], id)
}

func findId(s []int64, id int64) bool {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == id {
			return true
		}
	}
	return false
}

func removeId(s []int64, id int64) []int64 {
	if len(s) == 0 {
		return s
	}
	for i, v := range s {
		if v == id {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

type HTTPResponsesNode struct {
	List HTTPResponses `json:"list,omitempty"`
}

func (HTTPResponsesNode) Name() string {
	return "apicat-http-response"
}

// range responses list to dereference sub response
func (resp *HTTPResponsesNode) DereferenceResponses(sub *HTTPResponseDefine) {
	for _, r := range resp.List {
		r.HTTPResponseDefine.DereferenceResponses(sub)
	}
}

func (resp *HTTPResponsesNode) RemoveResponse(s_id int64) {
	id := strconv.Itoa(int(s_id))

	i := 0
	for i < len(resp.List) {
		// just todo remove
		if resp.List[i].IsRefId(id) {
			resp.List = append(resp.List[:i], resp.List[i+1:]...)
			continue
		}
		i++
	}
}

// range responses list to dereference sub response
func (resp *HTTPResponsesNode) DereferenceSchema(sub *Schema) {
	for _, r := range resp.List {
		r.HTTPResponseDefine.DereferenceSchema(sub)
	}
}

// range responses list to remove sub response
func (resp *HTTPResponsesNode) RemoveSchema(s_id int64) {
	for _, r := range resp.List {
		r.HTTPResponseDefine.RemoveSchema(s_id)
	}
}

type HTTPResponse struct {
	Code  int     `json:"code"`
	XDiff *string `json:"x-apicat-diff,omitempty"`
	HTTPResponseDefine
}

type HTTPResponses []*HTTPResponse

func (h *HTTPResponses) Map() map[int]HTTPResponseDefine {
	m := make(map[int]HTTPResponseDefine)
	for _, v := range *h {
		m[v.Code] = v.HTTPResponseDefine
	}
	return m
}

func (h *HTTPResponses) Add(code int, hrd *HTTPResponseDefine) {
	for _, v := range *h {
		if v.Code == code {
			v.HTTPResponseDefine = *hrd
			v.XDiff = hrd.XDiff
			return
		}
	}
	*h = append(*h, &HTTPResponse{
		Code:               code,
		XDiff:              hrd.XDiff,
		HTTPResponseDefine: *hrd,
	})
}

func (h *HTTPResponse) SetXDiff(x *string) {
	h.Header.SetXDiff(x)
	h.Content.SetXDiff(x)
	h.HTTPResponseDefine.SetXDiff(x)
}

type HTTPResponseDefine struct {
	ID          int64               `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Type        string              `json:"type,omitempty"`
	ParentId    uint64              `json:"parentid,omitempty"`
	Description string              `json:"description,omitempty"`
	Content     HTTPBody            `json:"content,omitempty"`
	Header      Schemas             `json:"header,omitempty"`
	Items       HTTPResponseDefines `json:"items,omitempty"`
	Reference   *string             `json:"$ref,omitempty"`
	XDiff       *string             `json:"x-apicat-diff,omitempty"`
}

func (h *HTTPResponseDefine) Ref() bool { return h.Reference != nil }

func (h *HTTPResponseDefine) IsRefId(id string) bool {
	if h == nil {
		return false
	}

	if h.Reference != nil {
		i := strings.LastIndex(*h.Reference, "/")
		if i != -1 {
			if id == (*h.Reference)[i+1:] {
				return true
			}
		}
	}
	return false
}

func (h *HTTPResponseDefine) DereferenceResponses(sub *HTTPResponseDefine) {
	id := strconv.Itoa(int(sub.ID))

	// this response is not reference sub
	if !h.IsRefId(id) {
		return
	}

	*h = *sub
}

func (h *HTTPResponseDefine) DereferenceSchema(sub *Schema) {

	// dereference content
	for _, body := range h.Content {
		body.DereferenceSchema(sub)
	}

}

func (h *HTTPResponseDefine) RemoveSchema(s_id int64) {
	// remove content
	for _, body := range h.Content {
		body.RemoveSchema(s_id)
	}
}

func (h *HTTPResponseDefine) ItemsTreeToList() (res HTTPResponseDefines) {
	if h.Items == nil || len(h.Items) == 0 {
		return res
	}

	for _, item := range h.Items {
		if item.Type == string(ContentItemTypeDir) {
			res = append(res, item.ItemsTreeToList()...)
		} else {
			res = append(res, item)
		}
	}

	return res
}

func (h *HTTPResponseDefine) SetXDiff(x *string) {
	h.Header.SetXDiff(x)
	h.Content.SetXDiff(x)
	h.XDiff = x
}

type HTTPResponseDefines []HTTPResponseDefine

func (h HTTPResponseDefines) Lookup(name string) *HTTPResponseDefine {
	for _, v := range h {
		if v.Type == string(ContentItemTypeDir) {
			if res := v.Items.Lookup(name); res != nil {
				return res
			}
		} else {
			if v.Name == name {
				return &v
			}
		}
	}
	return nil
}

func (h HTTPResponseDefines) LookupID(id int64) *HTTPResponseDefine {
	for _, v := range h {
		if v.Type == string(ContentItemTypeDir) {
			if res := v.Items.LookupID(id); res != nil {
				return res
			}
		} else {
			if v.ID == id {
				return &v
			}
		}
	}
	return nil
}

func (h *HTTPResponseDefines) SetXDiff(x *string) {
	for _, v := range *h {
		v.SetXDiff(x)
	}
}

type HTTPPart struct {
	Title string
	ID    int64
	Dir   string
	HTTPRequestNode
	Responses HTTPResponses `json:"responses,omitempty"`
}

func (h *HTTPPart) ToCollectItem(urlnode HTTPURLNode) *CollectItem {
	item := &CollectItem{
		Title: h.Title,
		Type:  ContentItemTypeHttp,
	}
	content := make([]*NodeProxy, 0)
	content = append(content, MuseCreateNodeProxy(WarpHTTPNode(urlnode)))
	content = append(content, MuseCreateNodeProxy(WarpHTTPNode(h.HTTPRequestNode)))
	content = append(content, MuseCreateNodeProxy(WarpHTTPNode(&HTTPResponsesNode{List: h.Responses})))
	item.Content = content
	return item
}

func (hb *HTTPBody) SetXDiff(x *string) {
	for _, v := range *hb {
		v.SetXDiff(x)
	}
}
