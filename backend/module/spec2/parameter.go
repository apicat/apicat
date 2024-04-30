package spec2

import "github.com/apicat/apicat/v2/backend/module/spec2/jsonschema"

type Parameter struct {
	ID          int64              `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string             `json:"name,omitempty" yaml:"name,omitempty"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool               `json:"required,omitempty" yaml:"required,omitempty"`
	Schema      *jsonschema.Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	XDiff       string             `json:"x-apicat-diff,omitempty" yaml:"x-apicat-diff,omitempty"`
}

type ParameterList []*Parameter

type HTTPParameters struct {
	Query  ParameterList `json:"query" yaml:"query"`
	Path   ParameterList `json:"path" yaml:"path"`
	Cookie ParameterList `json:"cookie" yaml:"cookie"`
	Header ParameterList `json:"header" yaml:"header"`
}

func (p *Parameter) SetXDiff(x string) {
	if p.Schema != nil {
		p.Schema.SetXDiff(x)
	}
	p.XDiff = x
}

func (pl *ParameterList) FindByID(id int64) *Parameter {
	if pl == nil {
		return nil
	}
	for _, v := range *pl {
		if id == v.ID {
			return v
		}
	}
	return nil
}

func (pl *ParameterList) FindByName(name string) *Parameter {
	if pl == nil {
		return nil
	}
	for _, v := range *pl {
		if name == v.Name {
			return v
		}
	}
	return nil
}

func (pl *ParameterList) DelByID(id int64) {
	if pl == nil {
		return
	}
	for i, v := range *pl {
		if v.ID == id {
			*pl = append((*pl)[:i], (*pl)[i+1:]...)
		}
	}
}

func (hp *HTTPParameters) FindByID(id int64) (p *Parameter) {
	p = hp.Query.FindByID(id)
	if p != nil {
		return p
	}
	p = hp.Path.FindByID(id)
	if p != nil {
		return p
	}
	p = hp.Cookie.FindByID(id)
	if p != nil {
		return p
	}
	p = hp.Header.FindByID(id)
	if p != nil {
		return p
	}
	return nil
}

func (hp *HTTPParameters) Fill() {
	if hp.Query == nil {
		hp.Query = make(ParameterList, 0)
	}
	if hp.Path == nil {
		hp.Path = make(ParameterList, 0)
	}
	if hp.Cookie == nil {
		hp.Cookie = make(ParameterList, 0)
	}
	if hp.Header == nil {
		hp.Header = make(ParameterList, 0)
	}
}

func (hp *HTTPParameters) Add(in string, v *Parameter) {
	switch in {
	case "query":
		hp.Query = append(hp.Query, v)
	case "path":
		hp.Path = append(hp.Path, v)
	case "cookie":
		hp.Cookie = append(hp.Cookie, v)
	case "header":
		hp.Header = append(hp.Header, v)
	}
}

func (hp *HTTPParameters) DelByID(in string, id int64) {
	switch in {
	case "query":
		hp.Query.DelByID(id)
	case "path":
		hp.Path.DelByID(id)
	case "cookie":
		hp.Cookie.DelByID(id)
	case "header":
		hp.Header.DelByID(id)
	}
}

func (hp *HTTPParameters) ToMap() map[string]ParameterList {
	m := make(map[string]ParameterList)
	if hp.Query != nil {
		m["query"] = hp.Query
	}
	if hp.Path != nil {
		m["path"] = hp.Path
	}
	if hp.Header != nil {
		m["header"] = hp.Header
	}
	if hp.Cookie != nil {
		m["cookie"] = hp.Cookie
	}
	return m
}
