package spec2

import (
	"strconv"
	"strings"
)

type BasicResponse struct {
	ID          int64         `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Header      ParameterList `json:"header,omitempty" yaml:"header,omitempty"`
	Content     HTTPBody      `json:"content" yaml:"content"`
	XDiff       string        `json:"x-apicat-diff,omitempty" yaml:"x-apicat-diff,omitempty"`
}

type Response struct {
	BasicResponse
	Code      int    `json:"code" yaml:"code"`
	Reference string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

type Responses []*Response

func (r *Response) Ref() bool { return r != nil && r.Reference != "" }

func (r *Response) IsRefID(id string) bool {
	if r == nil || r.Reference == "" {
		return false
	}

	i := strings.LastIndex(r.Reference, "/")
	if i != -1 {
		if id == (r.Reference)[i+1:] {
			return true
		}
	}
	return false
}

func (r *Response) GetRefID() int64 {
	if !r.Ref() {
		return 0
	}

	i := strings.LastIndex(r.Reference, "/")
	if i != -1 {
		id, _ := strconv.ParseInt(r.Reference[i+1:], 10, 64)
		return id
	}
	return 0
}

func (r *Response) ReplaceRef(ref *BasicResponse) {
	if !r.Ref() || ref == nil {
		return
	}

	refID := r.GetRefID()
	if refID != ref.ID {
		return
	}

	r.BasicResponse = *ref
}

func (r *Response) SetXDiff(x string) {
	r.Header.SetXDiff(x)
	r.Content.SetXDiff(x)
}

func (r *Responses) FindByCode(code int) *Response {
	for _, v := range *r {
		if v.Code == code {
			return v
		}
	}
	return nil
}

func (r *Responses) ToMap() map[int]*Response {
	m := make(map[int]*Response)
	for _, v := range *r {
		m[v.Code] = v
	}
	return m
}

func (r *Responses) AddOrUpdate(res *Response) {
	for _, v := range *r {
		if v.Code == res.Code {
			*v = *res
			return
		}
	}
	*r = append(*r, res)
}
