package spec2

import "github.com/apicat/apicat/v2/backend/module/spec2/jsonschema"

type Body struct {
	Schema   *jsonschema.Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	Examples []Example          `json:"examples,omitempty" yaml:"examples,omitempty"`
}

type Example struct {
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	Value   string `json:"value,omitempty" yaml:"value,omitempty"`
}

type HTTPBody map[string]*Body
