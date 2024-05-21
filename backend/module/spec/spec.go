package spec

import (
	"bytes"
	"encoding/json"
)

const TYPE_CATEGORY = "category"

type Spec struct {
	ApiCat      string       `json:"apicat" yaml:"apicat"`
	Info        Info         `json:"info" yaml:"info"`
	Servers     []Server     `json:"servers" yaml:"servers"`
	Globals     *Globals     `json:"globals" yaml:"globals"`
	Definitions *Definitions `json:"definitions,omitempty" yaml:"definitions,omitempty"`
	Collections Collections  `json:"collections,omitempty" yaml:"collections,omitempty"`
}

type Info struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty"`
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string `json:"version" yaml:"version"`
}

type Server struct {
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type Definitions struct {
	Schemas   DefinitionModels    `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses DefinitionResponses `json:"responses,omitempty" yaml:"responses,omitempty"`
}

type JSONOption struct {
	EscapeHTML bool
	Indent     string
}

func NewSpec() *Spec {
	return &Spec{
		ApiCat: "2.0",
	}
}

func NewEmptySpec() *Spec {
	return &Spec{
		ApiCat:  "2.0",
		Info:    Info{},
		Servers: make([]Server, 0),
		Globals: &Globals{
			Parameters: NewGlobalParameters(),
		},
		Definitions: &Definitions{
			Schemas:   make(DefinitionModels, 0),
			Responses: make(DefinitionResponses, 0),
		},
		Collections: make(Collections, 0),
	}
}

func NewSpecFromJson(raw []byte) (*Spec, error) {
	var spec Spec
	if err := json.Unmarshal(raw, &spec); err != nil {
		return nil, err
	}
	if spec.ApiCat == "" {
		spec.ApiCat = "2.0"
	}
	return &spec, nil
}

func (s *Spec) ToJSON(opt JSONOption) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(opt.EscapeHTML)
	if opt.Indent != "" {
		enc.SetIndent("", opt.Indent)
	}
	if err := enc.Encode(s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
