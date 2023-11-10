package postman

import (
	"encoding/json"
	"github.com/apicat/apicat/backend/module/spec"
	"github.com/apicat/apicat/backend/module/spec/jsonschema"

	"golang.org/x/exp/slices"
)

// https://schema.postman.com/collection/json/v2.1.0/draft-07/docs/index.html
type Spec struct {
	Info  Info   `json:"info"`
	Items []Item `json:"item"`
}
type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
}

type Item struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Request     *Request   `json:"request,omitempty"`
	Response    []Response `json:"response"`
	Items       []Item     `json:"item"`
}

type Request struct {
	Method      string     `json:"method"`
	Headers     []Variable `json:"header"`
	Url         URL        `json:"url"`
	Description string     `json:"description"`
	Body        *Body
}

type Variable struct {
	Key         string  `json:"key"`
	Value       string  `json:"value"`
	Description string  `json:"description"`
	Type        *string `json:"type"`
	Disabled    bool
}

func (v *Variable) toJSONSchema() *jsonschema.Schema {
	t := "string"
	typelist := []string{"string", "boolean", "number"}
	if v.Type != nil && slices.Contains(typelist, *v.Type) {
		t = *v.Type
	}
	sh := jsonschema.Create(t)
	sh.Description = v.Description
	sh.Example = v.Value
	return sh
}

func (v *Variable) toSchema() *spec.Schema {
	return &spec.Schema{
		Name:        v.Key,
		Description: v.Description,
		Schema:      v.toJSONSchema(),
	}
}

type Cookie struct{}

type URL struct {
	Raw       string     `json:"raw"`
	Protocol  string     `json:"protocol"`
	Host      []string   `json:"host"`
	Path      []string   `json:"path"`
	Queries   []Variable `json:"query"`
	Variables []Variable `json:"variable"`
}

type Response struct {
	Name                    string
	OriginalRequest         OriginalRequest
	Status                  string
	Code                    int
	PostmanePreviewLanguage string `json:"_postman_previewlanguage"`
	Header                  []Variable
	Cookie                  []Cookie
	Body                    string
}

type OriginalRequest struct {
	Method string
	Header []Variable
	URL    URL
}

type Body struct {
	Mode       string
	Raw        string
	Urlencoded []Variable
	Formdata   []Variable
	File       struct {
		Src     *string
		Content string
	}
	Options struct {
		Raw struct {
			Language string
		}
	}
	Disabled bool
}

func jsonToSchema(b string) *jsonschema.Schema {
	var m any
	if err := json.Unmarshal([]byte(b), &m); err != nil {
		m = make(map[string]any)
	}
	ret := tojsonschema(m)
	ret.Example = b
	return ret
}

func tojsonschema(a any) *jsonschema.Schema {
	var ret *jsonschema.Schema
	switch x := a.(type) {
	case []any:
		ret = jsonschema.Create("array")
		for _, v := range x {
			var items jsonschema.ValueOrBoolean[*jsonschema.Schema]
			items.SetValue(tojsonschema(v))
			ret.Items = &items
			break
		}
	case map[string]any:
		ret = jsonschema.Create("object")
		ret.Properties = make(map[string]*jsonschema.Schema)
		for k, v := range x {
			ret.Properties[k] = tojsonschema(v)
		}
	case float64:
		ret = jsonschema.Create("number")
		ret.Example = x
	case bool:
		ret = jsonschema.Create("boolean")
	default:
		ret = jsonschema.Create("string")
	}
	return ret
}
