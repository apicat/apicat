package spec2

type Response struct {
	ID          int64         `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string        `json:"name,omitempty" yaml:"name,omitempty"`
	Code        int           `json:"code" yaml:"code"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Header      ParameterList `json:"header,omitempty" yaml:"header,omitempty"`
	Content     HTTPBody      `json:"content" yaml:"content"`
}
