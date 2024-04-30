package spec2

type GlobalParameters struct {
	Query  ParameterList `json:"query" yaml:"query"`
	Cookie ParameterList `json:"cookie" yaml:"cookie"`
	Header ParameterList `json:"header" yaml:"header"`
}

type GlobalResponses []*Response

type Globals struct {
	Parameters GlobalParameters `json:"parameters" yaml:"parameters"`
	Responses  GlobalResponses  `json:"responses" yaml:"responses"`
}
