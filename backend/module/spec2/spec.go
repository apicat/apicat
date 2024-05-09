package spec2

const TYPE_CATEGORY = "category"

type Spec struct {
	ApiCat      string      `json:"apicat" yaml:"apicat"`
	Info        Info        `json:"info" yaml:"info"`
	Servers     []Server    `json:"servers" yaml:"servers"`
	Globals     Globals     `json:"globals" yaml:"globals"`
	Definitions Definitions `json:"definitions,omitempty" yaml:"definitions,omitempty"`
	Collections Collections `json:"collections,omitempty" yaml:"collections,omitempty"`
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

func NewSpec() *Spec {
	return &Spec{}
}

func NewEmptySpec() *Spec {
	return &Spec{}
}

func NewSpecFromJson() *Spec {
	return &Spec{}
}

func (s *Spec) ToJSON() string {
	return ""
}
