package spec2

type Spec struct {
	ApiCat  string   `json:"apicat" yaml:"apicat"`
	Info    Info     `json:"info" yaml:"info"`
	Servers []Server `json:"servers" yaml:"servers"`
	Globals Globals  `json:"globals" yaml:"globals"`
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
