package openapi

import (
	"fmt"

	"github.com/apicat/apicat/v2/backend/module/spec2"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/utils"
)

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
		// out, err = decodeOpenAPI3(docment)
	default:
		err = fmt.Errorf("not support %s", t)
	}
	return
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

	globalParameters := sw.parseGlobalParameters(model.Model.Extensions)

	return &spec2.Spec{
		ApiCat:      "2.0.1",
		Info:        sw.parseInfo(model.Model.Info),
		Servers:     sw.parseServers(&model.Model),
		Definitions: &spec2.Definitions{Schemas: definitionModels, Responses: definitionResponses},
		Globals:     &spec2.Globals{Parameters: globalParameters},
		Collections: sw.parseCollections(&model.Model, model.Model.Paths),
	}, nil
}
