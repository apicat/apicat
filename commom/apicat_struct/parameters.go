package apicat_struct

type ParametersObject struct {
	Path   []SchemaObject `json:"path"`
	Query  []SchemaObject `json:"query"`
	Header []SchemaObject `json:"header"`
	Cookie []SchemaObject `json:"cookie"`
}

func (p *ParametersObject) CheckPathRef(schema SchemaObject) {
	p.Path = append([]SchemaObject{schema}, p.Path...)
}

func (p *ParametersObject) CheckQueryRef(schema SchemaObject) {
	p.Query = append([]SchemaObject{schema}, p.Query...)
}

func (p *ParametersObject) CheckHeaderRef(schema SchemaObject) {
	p.Header = append([]SchemaObject{schema}, p.Header...)
}

func (p *ParametersObject) CheckCookieRef(schema SchemaObject) {
	p.Cookie = append([]SchemaObject{schema}, p.Cookie...)
}
