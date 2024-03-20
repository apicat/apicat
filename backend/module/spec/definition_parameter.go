package spec

import "github.com/apicat/apicat/v2/backend/module/spec/jsonschema"

var HttpParameterType = []string{"query", "path", "cookie", "header"}

type Parameter struct {
	ID          int64              `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string             `json:"name,omitempty" yaml:"name,omitempty"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool               `json:"required,omitempty" yaml:"required,omitempty"`
	Schema      *jsonschema.Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	XDiff       *string            `json:"x-apicat-diff,omitempty" yaml:"x-apicat-diff,omitempty"`
}

func (p *Parameter) EqualNomal(o *Parameter, OnlyCore bool) (b bool) {
	b = true
	if OnlyCore {
		if p.Required != o.Required {
			b = false
		}
		return b
	} else {
		if p.Description != o.Description || p.Required != o.Required {
			b = false
		}
	}
	return b
}

func (p *Parameter) SetXDiff(x *string) {
	if p.Schema != nil {
		p.Schema.SetXDiff(x)
	}
	p.XDiff = x
}

type ParameterList []*Parameter

func (p *ParameterList) LookupByID(id int64) *Parameter {
	if p == nil {
		return nil
	}
	for _, v := range *p {
		if id == v.ID {
			return v
		}
	}
	return nil
}

func (p *ParameterList) LookupByName(name string) *Parameter {
	if p == nil {
		return nil
	}
	for _, v := range *p {
		if name == v.Name {
			return v
		}
	}
	return nil
}

func (p *ParameterList) DelByID(id int64) {
	if p == nil {
		return
	}
	for i, v := range *p {
		if v.ID == id {
			*p = append((*p)[:i], (*p)[i+1:]...)
		}
	}
}

func (p *ParameterList) SetXDiff(x *string) {
	for _, v := range *p {
		v.SetXDiff(x)
	}
}

type HTTPParameters struct {
	Query  ParameterList `json:"query" yaml:"query"`
	Path   ParameterList `json:"path" yaml:"path"`
	Cookie ParameterList `json:"cookie" yaml:"cookie"`
	Header ParameterList `json:"header" yaml:"header"`
}

func (h *HTTPParameters) LookupByID(id int64) (p *Parameter) {
	p = h.Query.LookupByID(id)
	if p != nil {
		return p
	}
	p = h.Path.LookupByID(id)
	if p != nil {
		return p
	}
	p = h.Cookie.LookupByID(id)
	if p != nil {
		return p
	}
	p = h.Header.LookupByID(id)
	if p != nil {
		return p
	}
	return nil
}

func (h *HTTPParameters) Fill() {
	if h.Query == nil {
		h.Query = make(ParameterList, 0)
	}
	if h.Path == nil {
		h.Path = make(ParameterList, 0)
	}
	if h.Cookie == nil {
		h.Cookie = make(ParameterList, 0)
	}
	if h.Header == nil {
		h.Header = make(ParameterList, 0)
	}
}

func (h *HTTPParameters) Add(in string, v *Parameter) {
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

func (h *HTTPParameters) Del(in string, id int64) {
	switch in {
	case "query":
		h.Query.DelByID(id)
	case "path":
		h.Path.DelByID(id)
	case "cookie":
		h.Cookie.DelByID(id)
	case "header":
		h.Header.DelByID(id)
	}
}

func (h *HTTPParameters) Map() map[string]ParameterList {
	m := make(map[string]ParameterList)
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
