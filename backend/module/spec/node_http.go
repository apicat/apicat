package spec

import (
	"errors"
	"strconv"
)

type HTTPNoder interface {
	Name() string
}

type HTTPNode[T HTTPNoder] struct {
	Type  string `json:"type"`
	Attrs T      `json:"attrs"`
}

func (n *HTTPNode[T]) NodeType() string {
	return n.Type
}

type HTTPURLNode struct {
	Path   string  `json:"path"`
	Method string  `json:"method"`
	XDiff  *string `json:"x-apicat-diff,omitempty"`
}

const NAME_HTTP_URL string = "apicat-http-url"

func (HTTPURLNode) Name() string {
	return NAME_HTTP_URL
}

type HTTPRequestNode struct {
	GlobalExcepts map[string][]int64 `json:"globalExcepts"`
	Parameters    HTTPParameters     `json:"parameters,omitempty"`
	Content       HTTPBody           `json:"content"`
}

func (h *HTTPRequestNode) FillGlobalExcepts() {
	if h.GlobalExcepts == nil {
		h.GlobalExcepts = make(map[string][]int64)
		for _, v := range HttpParameterType {
			if _, ok := h.GlobalExcepts[v]; !ok {
				h.GlobalExcepts[v] = []int64{}
			}
		}
	}
}

func (h *HTTPRequestNode) InitContent() {
	if h.Content == nil {
		h.Content = make(HTTPBody)
	}
}

const NAME_HTTP_REQUEST string = "apicat-http-request"

func (HTTPRequestNode) Name() string {
	return NAME_HTTP_REQUEST
}

func (h *HTTPRequestNode) Deref(wantToDeref ...*Schema) {
	for _, v := range h.Content {
		v.Deref(wantToDeref...)
	}
}

func (h *HTTPRequestNode) DelRef(schema *Schema) {
	for _, v := range h.Content {
		v.DelRef(schema)
	}
}

// if content is exect, remove it and return true, else return false
func (h *HTTPRequestNode) DelGlobalExceptID(in string, id int64) {
	if h == nil || len(h.GlobalExcepts[in]) == 0 {
		return
	}

	for i, v := range h.GlobalExcepts[in] {
		if v == id {
			h.GlobalExcepts[in] = append(h.GlobalExcepts[in][:i], h.GlobalExcepts[in][i+1:]...)
			return
		}
	}
}

func (h *HTTPRequestNode) AddGlobalExcept(in string, id int64) {
	if h == nil {
		return
	}
	if len(h.GlobalExcepts[in]) == 0 {
		h.GlobalExcepts[in] = append(h.GlobalExcepts[in], id)
	} else {
		for _, v := range h.GlobalExcepts[in] {
			if v == id {
				return
			}
		}
		h.GlobalExcepts[in] = append(h.GlobalExcepts[in], id)
	}
}

type HTTPResponse struct {
	Code int `json:"code"`
	HTTPResponseDefine
}

func (h *HTTPResponse) SetXDiff(x *string) {
	h.Header.SetXDiff(x)
	h.Content.SetXDiff(x)
	h.HTTPResponseDefine.SetXDiff(x)
}

type HTTPResponses []*HTTPResponse

func (h *HTTPResponses) Map() map[int]HTTPResponseDefine {
	m := make(map[int]HTTPResponseDefine)
	for _, v := range *h {
		m[v.Code] = v.HTTPResponseDefine
	}
	return m
}

func (h *HTTPResponses) Merge(code int, hrd *HTTPResponseDefine) {
	for _, v := range *h {
		if v.Code == code {
			// want replace
			// 在 diff 时，可能会出现重复的 code
			// 如果重复，做并集处理
			v.HTTPResponseDefine = *hrd
			return
		}
	}
	*h = append(*h, &HTTPResponse{
		Code:               code,
		HTTPResponseDefine: *hrd,
	})
}

func (h *HTTPResponses) LookupCode(code int) *HTTPResponse {
	for _, v := range *h {
		if v.Code == code {
			return v
		}
	}
	return nil
}

type HTTPResponsesNode struct {
	List HTTPResponses `json:"list,omitempty"`
}

const NAME_HTTP_RESPONSES string = "apicat-http-response"

func (HTTPResponsesNode) Name() string {
	return NAME_HTTP_RESPONSES
}

// range responses list to dereference sub response
// @title Deref
// @description 遍历 responses list，将入参中想要解引用的 Definition Response List 替换为非 $ref 的方式
// @param sub *HTTPResponseDefine 引用的 Definition Response
// @return error
func (resp *HTTPResponsesNode) Deref(sub []*HTTPResponseDefine) (err error) {
	for _, r := range resp.List {
		for _, subr := range sub {
			err = r.HTTPResponseDefine.Deref(subr)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// 将引用了对应 id 的 response 删除
func (resp *HTTPResponsesNode) DelResponse(id int64) error {
	if resp == nil {
		return errors.New("responses is nil")
	}

	i := 0
	for i < len(resp.List) {
		// just todo remove
		if resp.List[i].IsRefID(strconv.Itoa(int(id))) {
			resp.List = append(resp.List[:i], resp.List[i+1:]...)
			continue
		}
		i++
	}
	return nil
}

// range responses list to dereference sub response
func (resp *HTTPResponsesNode) DerefSchema(wantToDeref ...*Schema) {
	for _, r := range resp.List {
		r.HTTPResponseDefine.DerefSchema(wantToDeref...)
	}
}

// // range responses list to remove sub response
func (resp *HTTPResponsesNode) DelRefSchema(schema *Schema) {
	for _, r := range resp.List {
		r.HTTPResponseDefine.DelRefSchema(schema)
	}
}

func WarpHTTPNode[T HTTPNoder](n T) Node {
	return &HTTPNode[T]{
		Type:  n.Name(),
		Attrs: n,
	}
}

type HTTPPart struct {
	Title string
	ID    int64
	Dir   string
	HTTPRequestNode
	Responses HTTPResponses `json:"responses,omitempty"`
}

func (h *HTTPPart) ToCollectItem(urlnode HTTPURLNode) *Collection {
	item := &Collection{
		Title: h.Title,
		Type:  CollectionItemTypeHttp,
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
