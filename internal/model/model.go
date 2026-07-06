package model

type APIDoc struct {
	Title       string
	Description string
	Version     string
	Servers     []Server
	Tags        []Tag
	TagGroups   []TagGroup
	Security    []SecurityScheme

	// Filled by lint.Annotate after conversion.
	LintTotal  int
	LintCounts []LintCount
}

type LintCount struct {
	Rule    string
	Message string
	Count   int
}

type Server struct {
	URL         string
	Description string
}

type Tag struct {
	Name        string
	Description string
}

type SecurityScheme struct {
	Name   string
	Type   string
	Scheme string
}

type TagGroup struct {
	Tag       Tag
	Endpoints []Endpoint
}

type Endpoint struct {
	Method      string
	Path        string
	Summary     string
	Description string
	OperationID string
	Deprecated  bool
	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   []Response
	Anchor      string

	// Filled by lint.Annotate after conversion.
	LintHints []string
}

type Parameter struct {
	Name        string
	In          string
	Description string
	Required    bool
	Deprecated  bool
	Schema      *SchemaView
}

type RequestBody struct {
	Description string
	Required    bool
	Content     []MediaTypeContent
}

type MediaTypeContent struct {
	MediaType string
	Schema    *SchemaView
}

type Response struct {
	StatusCode  string
	Description string
	Content     []MediaTypeContent
}

type SchemaView struct {
	Name        string
	Type        string
	Format      string
	Description string
	Required    bool
	Nullable    bool
	Deprecated  bool
	Enum        []string
	Default     string
	Example     string

	MinLength *int64
	MaxLength *int64
	Minimum   *float64
	Maximum   *float64
	MinItems  *int64
	MaxItems  *int64
	Pattern   string

	Properties []SchemaView
	Items      *SchemaView

	AllOf []SchemaView
	OneOf []SchemaView
	AnyOf []SchemaView

	RefName string
	IsRef   bool
	Depth   int
}
