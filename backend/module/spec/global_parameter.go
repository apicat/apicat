package spec

type GlobalParameters struct {
	Query  ParameterList `json:"query" yaml:"query"`
	Cookie ParameterList `json:"cookie" yaml:"cookie"`
	Header ParameterList `json:"header" yaml:"header"`
}

type GlobalResponses []*Response

type Globals struct {
	Parameters *GlobalParameters `json:"parameters" yaml:"parameters"`
	Responses  GlobalResponses   `json:"responses" yaml:"responses"`
}

func NewGlobalParameters() *GlobalParameters {
	return &GlobalParameters{
		Query:  ParameterList{},
		Cookie: ParameterList{},
		Header: ParameterList{},
	}
}

func (p *GlobalParameters) Add(in string, v *Parameter) {
	switch in {
	case "query":
		p.Query = append(p.Query, v)
	case "cookie":
		p.Cookie = append(p.Cookie, v)
	case "header":
		p.Header = append(p.Header, v)
	}
}

func (p *GlobalParameters) ToMap() map[string]ParameterList {
	return map[string]ParameterList{
		"query":  p.Query,
		"cookie": p.Cookie,
		"header": p.Header,
	}
}
