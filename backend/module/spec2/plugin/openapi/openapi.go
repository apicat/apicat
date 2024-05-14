package openapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec2"
	"github.com/apicat/apicat/v2/backend/module/spec2/jsonschema"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v2"
)

type tagObject struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type openAPIParamter struct {
	Name        string `json:"name,omitempty"`
	In          string `json:"in,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
	// in body
	Schema *jsonschema.Schema `json:"schema,omitempty"`
	// in not body
	Type      string `json:"type,omitempty"`
	Format    string `json:"format,omitempty"`
	Default   any    `json:"default,omitempty"`
	Reference string `json:"$ref,omitempty"`
	Example   any    `json:"example,omitempty"`
}

type specPathItem struct {
	Title       string
	Description string
	OperatorID  string
	Tags        []string
	Req         spec2.CollectionHttpRequest
	Res         spec2.CollectionHttpResponse
}

func Parse(data []byte) (out *spec2.Spec, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	docment, docerr := libopenapi.NewDocument(data)
	if docerr != nil {
		err = docerr
		return
	}
	t := docment.GetSpecInfo().SpecType
	switch t {
	case utils.OpenApi2:
		out, err = parseSwagger(docment)
	case utils.OpenApi3:
		out, err = parseOpenAPI3(docment)
	default:
		err = fmt.Errorf("not support %s", t)
	}
	return
}

func Generator(in *spec2.Spec, version, typ string) ([]byte, error) {
	switch version {
	case "2.0":
		generator := &swaggerGenerator{}
		sp := generator.generateBase(in)
		paths, tags := generator.generatePaths(in)
		sp.Paths = paths
		sp.Tags = tags

		globalParam := in.Globals.Parameters
		m := globalParam.ToMap()
		sp.GlobalParameters = make(map[string]openAPIParamter)
		for in, ps := range m {
			for _, p := range ps {
				sp.GlobalParameters[fmt.Sprintf("%s-%s", in, p.Name)] = toParameter(p, in, version)
			}
		}

		if typ == "yaml" {
			return yaml.Marshal(sp)
		} else {
			return json.MarshalIndent(sp, "", "  ")
		}
	default:
		if strings.HasPrefix(version, "3.") && len(strings.Split(version, ".")) == 3 {
			op := &openapiGenerator{}
			sp := op.generateBase(in, version)
			paths, tag := op.generatePaths(version, in)
			sp.Paths = paths
			sp.Tags = tag

			if typ == "yaml" {
				return yaml.Marshal(sp)
			} else {
				return json.MarshalIndent(sp, "", "  ")
			}
		}
	}
	return nil, fmt.Errorf("openapi %s not support", version)
}

func parseSwagger(document libopenapi.Document) (*spec2.Spec, error) {
	model, errors := document.BuildV2Model()
	if len(errors) > 0 {
		return nil, fmt.Errorf("swagger version:%s parse faild", document.GetVersion())
	}
	sw := &swaggerParser{}
	definitionModels, err := sw.parseDefinitionModels(model.Model.Definitions)
	if err != nil {
		return nil, err
	}

	definitionResponses, err := sw.parseDefinitionResponses(&model.Model)
	if err != nil {
		return nil, err
	}

	return &spec2.Spec{
		ApiCat:      "2.0.1",
		Info:        sw.parseInfo(model.Model.Info),
		Servers:     sw.parseServers(&model.Model),
		Definitions: &spec2.Definitions{Schemas: definitionModels, Responses: definitionResponses},
		Globals:     &spec2.Globals{Parameters: sw.parseGlobalParameters(model.Model.Extensions)},
		Collections: sw.parseCollections(&model.Model, model.Model.Paths),
	}, nil
}

func parseOpenAPI3(document libopenapi.Document) (*spec2.Spec, error) {
	model, errors := document.BuildV3Model()
	if len(errors) > 0 {
		return nil, fmt.Errorf("openapi version:%s parse faild", document.GetVersion())
	}
	o := &openapiParser{
		parametersMapping: make(map[string]*spec2.Parameter),
	}

	definitions, err := o.parseDefinetions(model.Model.Components)
	if err != nil {
		return nil, err
	}

	return &spec2.Spec{
		ApiCat:      "2.0.1",
		Info:        o.parseInfo(model.Model.Info),
		Servers:     o.parseServers(model.Model.Servers),
		Globals:     &spec2.Globals{Parameters: o.parseGlobalParameters(model.Model.Components)},
		Definitions: definitions,
		Collections: o.parseCollections(model.Model.Paths),
	}, nil
}
