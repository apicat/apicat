package spec2

import "strconv"

const NODE_HTTP_REQUEST = "apicat-http-request"

type CollectionHttpRequest struct {
	Type  string            `json:"type" yaml:"type"`
	Attrs *HttpRequestAttrs `json:"attr" yaml:"attrs"`
}

type HttpRequestAttrs struct {
	GlobalExcepts *HttpRequestGlobalExcepts `json:"globalExcepts,omitempty" yaml:"globalExcepts,omitempty"`
	Parameters    *HTTPParameters           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Content       HTTPBody                  `json:"content,omitempty" yaml:"content,omitempty"`
}

type HttpRequestGlobalExcepts struct {
	Header []int64 `json:"header" yaml:"header"`
	Cookie []int64 `json:"cookie" yaml:"cookie"`
	Query  []int64 `json:"query" yaml:"query"`
}

func NewCollectionHttpRequest() *CollectionHttpRequest {
	return &CollectionHttpRequest{
		Type: NODE_HTTP_REQUEST,
		Attrs: &HttpRequestAttrs{
			GlobalExcepts: NewHttpRequestGlobalExcepts(),
			Parameters:    NewHTTPParameters(),
			Content:       HTTPBody{},
		},
	}
}

func NewHttpRequestGlobalExcepts() *HttpRequestGlobalExcepts {
	return &HttpRequestGlobalExcepts{
		Header: []int64{},
		Cookie: []int64{},
		Query:  []int64{},
	}
}

func (r *CollectionHttpRequest) NodeType() string {
	return r.Type
}

func (r *CollectionHttpRequest) GetGlobalExcept(in string) []int64 {
	if r == nil || r.Attrs == nil || r.Attrs.GlobalExcepts == nil {
		return nil
	}

	switch in {
	case "header":
		return r.Attrs.GlobalExcepts.Header
	case "cookie":
		return r.Attrs.GlobalExcepts.Cookie
	case "query":
		return r.Attrs.GlobalExcepts.Query
	}
	return nil
}

func (r *CollectionHttpRequest) GetGlobalExceptAll() map[string][]int64 {
	if r == nil || r.Attrs == nil || r.Attrs.GlobalExcepts == nil {
		return nil
	}

	return map[string][]int64{
		"header": r.Attrs.GlobalExcepts.Header,
		"cookie": r.Attrs.GlobalExcepts.Cookie,
		"query":  r.Attrs.GlobalExcepts.Query,
	}
}

func (r *CollectionHttpRequest) AddGlobalExcept(in string, id int64) {
	if r == nil || r.Attrs == nil || r.Attrs.GlobalExcepts == nil {
		return
	}

	switch in {
	case "header":
		if len(r.Attrs.GlobalExcepts.Header) == 0 {
			r.Attrs.GlobalExcepts.Header = append(r.Attrs.GlobalExcepts.Header, id)
		} else {
			for _, v := range r.Attrs.GlobalExcepts.Header {
				if v == id {
					return
				}
			}
			r.Attrs.GlobalExcepts.Header = append(r.Attrs.GlobalExcepts.Header, id)
		}
	case "cookie":
		if len(r.Attrs.GlobalExcepts.Cookie) == 0 {
			r.Attrs.GlobalExcepts.Cookie = append(r.Attrs.GlobalExcepts.Cookie, id)
		} else {
			for _, v := range r.Attrs.GlobalExcepts.Cookie {
				if v == id {
					return
				}
			}
			r.Attrs.GlobalExcepts.Cookie = append(r.Attrs.GlobalExcepts.Cookie, id)
		}
	case "query":
		if len(r.Attrs.GlobalExcepts.Query) == 0 {
			r.Attrs.GlobalExcepts.Query = append(r.Attrs.GlobalExcepts.Query, id)
		} else {
			for _, v := range r.Attrs.GlobalExcepts.Query {
				if v == id {
					return
				}
			}
			r.Attrs.GlobalExcepts.Query = append(r.Attrs.GlobalExcepts.Query, id)
		}
	}
}

func (r *CollectionHttpRequest) DelGlobalExcept(in string, id int64) {
	if r == nil || r.Attrs == nil || r.Attrs.GlobalExcepts == nil {
		return
	}

	switch in {
	case "header":
		if len(r.Attrs.GlobalExcepts.Header) == 0 {
			return
		}
		for i, v := range r.Attrs.GlobalExcepts.Header {
			if v == id {
				r.Attrs.GlobalExcepts.Header = append(r.Attrs.GlobalExcepts.Header[:i], r.Attrs.GlobalExcepts.Header[i+1:]...)
				return
			}
		}
	case "cookie":
		if len(r.Attrs.GlobalExcepts.Cookie) == 0 {
			return
		}
		for i, v := range r.Attrs.GlobalExcepts.Cookie {
			if v == id {
				r.Attrs.GlobalExcepts.Cookie = append(r.Attrs.GlobalExcepts.Cookie[:i], r.Attrs.GlobalExcepts.Cookie[i+1:]...)
				return
			}
		}
	case "query":
		if len(r.Attrs.GlobalExcepts.Query) == 0 {
			return
		}
		for i, v := range r.Attrs.GlobalExcepts.Query {
			if v == id {
				r.Attrs.GlobalExcepts.Query = append(r.Attrs.GlobalExcepts.Query[:i], r.Attrs.GlobalExcepts.Query[i+1:]...)
				return
			}
		}
	}
}

func (r *CollectionHttpRequest) Deref(ref *Model) {
	if r == nil || r.Attrs == nil || r.Attrs.Content == nil || ref == nil {
		return
	}
	ref.Schema.ID = ref.ID

	for _, v := range r.Attrs.Content {
		if v.Schema != nil {
			refSchemas := v.Schema.DeepFindRefById(strconv.FormatInt(ref.ID, 10))
			if len(refSchemas) > 0 {
				for _, schema := range refSchemas {
					schema.ReplaceRef(ref.Schema)
				}
			}
		}
	}
}

func (r *CollectionHttpRequest) DelRef(ref *Model) {
	if r == nil || r.Attrs == nil || r.Attrs.Content == nil || ref == nil {
		return
	}
	ref.Schema.ID = ref.ID

	for _, v := range r.Attrs.Content {
		if v.Schema != nil {
			if v.Schema.Ref() {
				v.Schema.DelRootRef(ref.Schema)
			}
			v.Schema.DelChildrenRef(ref.Schema)
		}
	}
}
