package spec

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

type BasicResponse struct {
	ID          int64         `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Header      ParameterList `json:"header,omitempty" yaml:"header,omitempty"`
	Content     HTTPBody      `json:"content,omitempty" yaml:"content,omitempty"`
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

func (r *Response) GetRefID() (int64, error) {
	if !r.Ref() {
		return 0, errors.New("no reference")
	}

	i := strings.LastIndex(r.Reference, "/")
	if i != -1 {
		id, _ := strconv.ParseInt(r.Reference[i+1:], 10, 64)
		return id, nil
	}
	return 0, errors.New("no reference")
}

func (r *Response) ReplaceRef(ref *BasicResponse) error {
	if !r.Ref() || ref == nil {
		return errors.New("response is not a reference or ref is nil")
	}

	if refID, err := r.GetRefID(); err != nil || refID != ref.ID {
		return errors.New("ref id does not match")
	}

	r.BasicResponse = *ref
	r.Reference = ""
	return nil
}

func (r *Response) SetXDiff(x string) {
	r.Header.SetXDiff(x)
	r.Content.SetXDiff(x)
}

func (r *Response) IsEmpty() bool {
	if len(r.Header) > 0 {
		return false
	}
	if len(r.Content) > 0 {
		for _, body := range r.Content {
			if len(body.Schema.Properties) > 0 {
				return false
			}
			if len(body.Schema.AllOf) > 0 {
				return false
			}
			if len(body.Schema.OneOf) > 0 {
				return false
			}
			if len(body.Schema.AnyOf) > 0 {
				return false
			}
			if body.Schema.Reference != nil {
				return false
			}
			if body.Schema.Items != nil {
				return false
			}
		}
	}
	return true
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

func (r *Responses) Sort() {
	m := make(map[int]*Response)
	l := make([]int, 0)
	for _, v := range *r {
		m[v.Code] = v
		l = append(l, v.Code)
	}

	new := make(Responses, 0)
	sort.Ints(l)
	for _, v := range l {
		new = append(new, m[v])
	}
	*r = new
}
