package openapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

type openapiParser struct {
	modelMapping      map[string]int64
	parametersMapping map[string]*spec.Parameter
}

type openapiGenerator struct {
	modelMapping map[int64]string
}

type openapiSpec struct {
	Openapi    string                                `json:"openapi"`
	Info       spec.Info                             `json:"info"`
	Servers    []spec.Server                         `json:"servers,omitempty"`
	Components map[string]any                        `json:"components,omitempty"`
	Paths      map[string]map[string]openapiPathItem `json:"paths"`
	Tags       []tagObject                           `json:"tags,omitempty"`
}

type openapiPathItem struct {
	Summary     string              `json:"summary"`
	Description string              `json:"description,omitempty"`
	OperationId string              `json:"operationId"`
	Tags        []string            `json:"tags,omitempty"`
	Parameters  []openAPIParamter   `json:"parameters,omitempty"`
	RequestBody *openapiRequestbody `json:"requestBody,omitempty"`
	Responses   map[string]any      `json:"responses,omitempty"`
}

type openapiRequestbody struct {
	Content map[string]*jsonschema.Schema `json:"content,omitempty"`
}

func (o *openapiParser) parseInfo(info *base.Info) spec.Info {
	return spec.Info{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}
}

func (o *openapiParser) parseServers(servs []*v3.Server) []spec.Server {
	srvs := make([]spec.Server, len(servs))
	for k, v := range servs {
		srvs[k] = spec.Server{
			URL:         v.URL,
			Description: v.Description,
		}
	}
	return srvs
}

func (o *openapiParser) parseContent(mts *orderedmap.Map[string, *v3.MediaType]) (spec.HTTPBody, error) {
	if mts == nil {
		return nil, errors.New("no content")
	}

	content := make(spec.HTTPBody)
	for pair := range orderedmap.Iterate(context.Background(), mts) {
		contentType := pair.Key()
		mediaType := pair.Value()
		body := &spec.Body{}

		js, err := jsonSchemaConverter(mediaType.Schema)
		if err != nil {
			return nil, err
		}

		if mediaType.Example != nil {
			js.Examples = mediaType.Example.Value
		}

		if orderedmap.Len(mediaType.Examples) > 0 {
			i := 0
			body.Examples = make(map[string]spec.Example)
			for examplePair := range orderedmap.Iterate(context.Background(), mediaType.Examples) {
				v := examplePair.Value()
				if example, err := json.Marshal(v); err == nil {
					body.Examples[strconv.Itoa(i)] = spec.Example{
						Summary: v.Summary,
						Value:   string(example),
					}
					i++
				}
			}
		}
		body.Schema = js
		content[contentType] = body
	}
	return content, nil
}

func (o *openapiParser) parseDefinetions(comp *v3.Components) (*spec.Definitions, error) {
	if comp == nil {
		return &spec.Definitions{
			Schemas:   make(spec.DefinitionModels, 0),
			Responses: make(spec.DefinitionResponses, 0),
		}, nil
	}

	o.modelMapping = map[string]int64{}
	models := make(spec.DefinitionModels, 0)

	for pair := range orderedmap.Iterate(context.Background(), comp.Schemas) {
		k := pair.Key()
		v := pair.Value()

		js, err := jsonSchemaConverter(v)
		if err != nil {
			return nil, err
		}

		if js.Type.First() != jsonschema.T_OBJ && js.Type.First() != jsonschema.T_ARR {
			parentJS := jsonschema.NewSchema(jsonschema.T_OBJ)
			if js.AnyOf != nil || js.OneOf != nil {
				parentJS.Properties = map[string]*jsonschema.Schema{
					k: js,
				}
			} else {
				parentJS.Properties = map[string]*jsonschema.Schema{
					k: js,
				}
			}
			models = append(models, &spec.DefinitionModel{
				ID:          stringToUnid(k),
				Name:        k,
				Description: js.Description,
				Schema:      parentJS,
			})
		} else {
			models = append(models, &spec.DefinitionModel{
				ID:          stringToUnid(k),
				Name:        k,
				Description: js.Description,
				Schema:      js,
			})
		}
	}

	responses := make(spec.DefinitionResponses, 0)
	for pair := range orderedmap.Iterate(context.Background(), comp.Responses) {
		responseName := pair.Key()
		response := pair.Value()
		id := stringToUnid(responseName)

		def := &spec.DefinitionResponse{
			BasicResponse: spec.BasicResponse{
				Header:      make(spec.ParameterList, 0),
				Name:        responseName,
				ID:          id,
				Description: response.Description,
			},
		}

		if response.Headers != nil {
			for headerPair := range orderedmap.Iterate(context.Background(), response.Headers) {
				v := headerPair.Value()
				js, err := jsonSchemaConverter(v.Schema)
				if err != nil {
					return nil, err
				}
				js.Description = v.Description
				def.Header = append(def.Header, &spec.Parameter{
					Name:   headerPair.Key(),
					Schema: js,
				})
			}
		}

		if response.Content != nil {
			content, err := o.parseContent(response.Content)
			if err != nil {
				return nil, err
			}
			def.Content = content
		}
		responses = append(responses, def)
	}

	if comp.Parameters != nil {
		for pair := range orderedmap.Iterate(context.Background(), comp.Parameters) {
			parameter := pair.Value()
			if parameter.Schema != nil {
				js, err := jsonSchemaConverter(parameter.Schema)
				if err != nil {
					return nil, err
				}
				k := fmt.Sprintf("%s-%s", parameter.In, parameter.Name)
				o.parametersMapping[k] = &spec.Parameter{
					ID:          stringToUnid(parameter.Name),
					Name:        parameter.Name,
					Description: parameter.Description,
					Required:    parameter.Required != nil && *parameter.Required,
					Schema:      js,
				}
			}
		}
	}

	return &spec.Definitions{
		Schemas:   models,
		Responses: responses,
	}, nil
}

func (o *openapiParser) parseGlobalParameters(com *v3.Components) *spec.GlobalParameters {
	res := spec.NewGlobalParameters()
	if com == nil {
		return res
	}

	inp := com.Extensions
	if inp == nil {
		return res
	}

	global, ok := inp.Get("x-apicat-global-parameters")
	if !ok {
		return res
	}

	var globals map[string]any
	if err := global.Decode(&globals); err != nil {
		return res
	}

	for k, v := range globals {
		nb, err := json.Marshal(v)
		if err != nil {
			continue
		}

		s := &spec.Parameter{}
		json.Unmarshal(nb, s)
		in := strings.Index(k, "-")
		if in == -1 {
			continue
		}
		res.Add(k[:in], s)
	}
	return res
}

func (o *openapiParser) parseParameters(inp []*v3.Parameter) (*spec.HTTPParameters, error) {
	rawparamter := spec.NewHTTPParameters()
	for _, v := range inp {
		if g := v.GoLow(); g.IsReference() {
			// if this parameter is a global parameter
			if isGlobalParameter(g.Reference.GetReference()) {
				continue
			}

			name := fmt.Sprintf("%s-%s", v.In, v.Name)
			sc, ok := o.parametersMapping[name]
			if ok {
				rawparamter.Add(v.In, sc)
				continue
				// r := fmt.Sprintf("#/definitions/parameters/%d", id)
				// rawparamter.Add(v.In, &spec.Schema{
				// 	Reference: &r,
				// })
				// continue
			}
		}

		sp := &spec.Parameter{
			Name:     v.Name,
			Required: v.Required != nil && *v.Required,
		}

		sp.Schema = &jsonschema.Schema{}
		if v.Schema != nil {
			js, err := jsonSchemaConverter(v.Schema)
			if err != nil {
				return nil, err
			}
			sp.Schema = js
		}
		sp.Schema.Description = v.Description
		if v.Example != nil {
			sp.Schema.Examples = v.Example
		}
		sp.Schema.Deprecated = &v.Deprecated
		rawparamter.Add(v.In, sp)
	}
	return rawparamter, nil
}

func (o *openapiParser) parseResponses(responses *orderedmap.Map[string, *v3.Response]) (*spec.CollectionHttpResponse, error) {
	outresponses := spec.NewCollectionHttpResponse()

	for pair := range orderedmap.Iterate(context.Background(), responses) {
		code := pair.Key()
		res := pair.Value()

		c, err := strconv.Atoi(code)
		if err != nil {
			continue
		}

		resp := spec.Response{
			BasicResponse: spec.BasicResponse{
				Name: fmt.Sprintf("response%s", code),
			},
			Code: c,
		}

		s := res.GoLow().Reference.GetReference()
		if s != "" {
			refs := fmt.Sprintf("#/definitions/responses/%d", stringToUnid(s[strings.LastIndex(s, "/")+1:]))
			resp.Reference = refs
			outresponses.Attrs.List = append(outresponses.Attrs.List, &resp)
			continue
		}

		if v, ok := res.Extensions.Get("x-apicat-response-name"); ok {
			resp.Name = v.Value
		}

		resp.Description = res.Description
		resp.Header = make(spec.ParameterList, 0)
		if res.Headers != nil {
			for pair := range orderedmap.Iterate(context.Background(), res.Headers) {
				v := pair.Value()
				js, err := jsonSchemaConverter(v.Schema)
				if err != nil {
					return nil, err
				}

				js.Description = v.Description
				resp.Header = append(resp.Header, &spec.Parameter{
					Name:   pair.Key(),
					Schema: js,
				})
			}
		}

		if res.Content != nil {
			content, err := o.parseContent(res.Content)
			if err != nil {
				return nil, err
			}
			resp.Content = content
		}

		outresponses.Attrs.List = append(outresponses.Attrs.List, &resp)
	}
	return outresponses, nil
}

func (o *openapiParser) parseCollections(paths *v3.Paths) spec.Collections {
	collections := make(spec.Collections, 0)
	if paths == nil {
		return collections
	}

	var err error
	for pair := range orderedmap.Iterate(context.Background(), paths.PathItems) {
		path := pair.Key()
		operations := pair.Value().GetOperations()
		for operation := range orderedmap.Iterate(context.Background(), operations) {
			method := operation.Key()
			info := operation.Value()

			content := spec.CollectionNodes{
				spec.NewCollectionHttpUrl(path, method).ToCollectionNode(),
			}

			// parse markdown to doc
			// doctree := markdown.ToDocment([]byte(info.Description))
			// for _, v := range doctree.Items {
			// 	content = append(content, spec.MuseCreateNodeProxy(v))
			// }

			// request
			req := spec.NewCollectionHttpRequest()
			if req.Attrs.Parameters, err = o.parseParameters(info.Parameters); err != nil {
				continue
			}
			if info.RequestBody != nil {
				if req.Attrs.Content, err = o.parseContent(info.RequestBody.Content); err != nil {
					continue
				}
			}
			content = append(content, req.ToCollectionNode())

			// response
			res := spec.NewCollectionHttpResponse()
			if info.Responses != nil {
				if res, err = o.parseResponses(info.Responses.Codes); err != nil {
					continue
				}
			}
			content = append(content, res.ToCollectionNode())

			title := info.Summary
			if title == "" {
				title = path
			}

			c := &spec.Collection{
				Type:    spec.TYPE_HTTP,
				Title:   title,
				Tags:    info.Tags,
				Content: content,
			}
			c.Content.SortResponses()
			collections = append(collections, c)
		}
	}
	return collections
}

func (o *openapiGenerator) generateBase(in *spec.Spec, version string) *openapiSpec {
	return &openapiSpec{
		Openapi: version,
		Info: spec.Info{
			Title:       in.Info.Title,
			Description: in.Info.Description,
			Version:     in.Info.Version,
		},
		Servers:    in.Servers,
		Components: o.generateComponents(version, in),
	}
}

func (o *openapiGenerator) convertJsonSchema(version string, in *jsonschema.Schema) *jsonschema.Schema {
	if in == nil {
		return nil
	}
	p := convertJsonSchemaRef(in, version, o.modelMapping)
	// if p.Reference == nil && strings.HasPrefix(ver, "3.0") {
	// 	if p.Items != nil {
	// 		if !p.Items.IsBool() {
	// 			p.Items.SetValue(&jsonschema.Schema{})
	// 		}
	// 	}
	// 	if p.Properties != nil {
	// 		for _, v := range p.Properties {
	// 			o.convertJSONSchema(ver, v)
	// 		}
	// 	}
	// 	if p.AdditionalProperties != nil {
	// 		if !p.AdditionalProperties.IsBool() {
	// 			o.convertJSONSchema(ver, p.AdditionalProperties.Value())
	// 		}
	// 	}
	// 	p.Type.SetValue(p.Type.Value()[0])
	// }
	if p.Type != nil {
		t := p.Type.List()
		if len(t) > 0 && t[0] == "file" {
			// jsonschema 没有file
			p.Type.Set(jsonschema.T_ARR)
			p.Items = &jsonschema.ValueOrBoolean[*jsonschema.Schema]{}
			p.Items.SetValue(&jsonschema.Schema{})
		}
	}
	return p
}

func (o *openapiGenerator) generateResponseWithoutRef(resp *spec.BasicResponse, version string) map[string]any {
	result := map[string]any{}
	if resp.Content != nil {
		c := make(map[string]*spec.Body)
		for contentType, body := range resp.Content {
			body.Schema = o.convertJsonSchema(version, body.Schema)
			c[contentType] = body
		}
		result["content"] = c
	}

	if len(resp.Header) > 0 {
		headers := make(map[string]any)
		for _, h := range resp.Header {
			headers[h.Name] = map[string]any{
				"description": h.Description,
				"schema":      o.convertJsonSchema(version, h.Schema),
			}
		}
		result["headers"] = headers
	}

	result["description"] = resp.Description
	result["x-apicat-response-name"] = resp.Name
	return result
}

func (o *openapiGenerator) generateReqParams(collectionReq spec.CollectionHttpRequest, globalsParmaters *spec.GlobalParameters, version string) []openAPIParamter {
	// var out []openAPIParamter
	out := globalToLocalParameters(globalsParmaters, false, collectionReq.Attrs.GlobalExcepts.ToMap())

	for in, params := range collectionReq.Attrs.Parameters.ToMap() {
		for _, param := range params {
			p := *param
			// if p.Reference != nil {
			// 	if defp := spe.Definitions.Parameters.LookupID(toInt64(getRefName(*p.Reference))); defp != nil {
			// 		p = *defp
			// 	}
			// }
			item := openAPIParamter{
				Name:        p.Name,
				Required:    p.Required,
				Description: p.Description,
				// Example:     p.Example,
				Schema: o.convertJsonSchema(version, p.Schema),
				In:     in,
			}
			out = append(out, item)
		}
	}
	return out
}

func (o *openapiGenerator) generateResponse(resp *spec.Response, definitionsResps spec.DefinitionResponses, version string) map[string]any {
	if resp.Ref() {
		if strings.HasPrefix(resp.Reference, "#/definitions/responses/") {
			if x := definitionsResps.FindByID(
				toInt64(getRefName(resp.Reference)),
			); x != nil {
				name_id := fmt.Sprintf("%s-%d", x.Name, x.ID)
				return map[string]any{
					"$ref": "#/components/responses/" + name_id,
				}
			}
		}
	}
	return o.generateResponseWithoutRef(&resp.BasicResponse, version)
}

func (o *openapiGenerator) generateComponents(version string, in *spec.Spec) map[string]any {
	schemas := make(map[string]jsonschema.Schema)
	respons := make(map[string]any)
	o.modelMapping = map[int64]string{}

	definitionModels := spec.DefinitionModels{}
	for _, v := range in.Definitions.Schemas {
		if v.Type == string(spec.TYPE_CATEGORY) {
			items := v.ItemsTreeToList()
			for _, item := range items {
				o.modelMapping[item.ID] = item.Name
			}
			definitionModels = append(definitionModels, items...)
		} else {
			o.modelMapping[v.ID] = v.Name
			definitionModels = append(definitionModels, v)
		}
	}
	for _, v := range definitionModels {
		name_id := fmt.Sprintf("%s-%d", strings.ReplaceAll(v.Name, " ", ""), v.ID)
		schemas[name_id] = *o.convertJsonSchema(version, v.Schema)
	}

	for _, v := range in.Definitions.Responses {
		if v.Type == string(spec.TYPE_CATEGORY) {
			resps := v.ItemsTreeToList()
			for _, resp := range resps {
				name_id := fmt.Sprintf("%s-%d", resp.Name, resp.ID)
				respons[name_id] = o.generateResponseWithoutRef(&resp.BasicResponse, version)
			}
		} else {
			name_id := fmt.Sprintf("%s-%d", strings.ReplaceAll(v.Name, " ", ""), v.ID)
			respons[name_id] = o.generateResponseWithoutRef(&v.BasicResponse, version)
		}
	}

	globalParam := in.Globals.Parameters
	m := globalParam.ToMap()
	globals := make(map[string]openAPIParamter)
	for in, ps := range m {
		for _, p := range ps {
			globals[fmt.Sprintf("%s-%s", in, strings.ReplaceAll(p.Name, " ", ""))] = toParameter(p, in, version)
		}
	}

	return map[string]any{
		"schemas":                    schemas,
		"responses":                  respons,
		"x-apicat-global-parameters": globals,
	}
}

func (o *openapiGenerator) generatePaths(version string, in *spec.Spec) (map[string]map[string]openapiPathItem, []tagObject) {
	var (
		out  = make(map[string]map[string]openapiPathItem)
		tags = make(map[string]struct{})
	)

	for path, ops := range deepGetHttpCollection(&in.Collections) {
		if path == "" {
			continue
		}

		for method, op := range ops {
			item := openapiPathItem{
				Summary:     op.Title,
				Description: op.Description,
				OperationId: op.OperatorID,
				Tags:        op.Tags,
				Parameters:  o.generateReqParams(op.Req, in.Globals.Parameters, version),
				Responses:   make(map[string]any),
			}

			for _, v := range op.Tags {
				tags[v] = struct{}{}
			}

			for contentType, body := range op.Req.Attrs.Content {
				if contentType == "none" {
					continue
				}

				sp := o.convertJsonSchema(version, body.Schema)
				sp.Examples = body.Examples
				if item.RequestBody == nil {
					item.RequestBody = &openapiRequestbody{
						Content: make(map[string]*jsonschema.Schema),
					}
				}
				item.RequestBody.Content[contentType] = sp
			}

			for _, v := range op.Res.Attrs.List {
				res := o.generateResponse(v, in.Definitions.Responses, version)
				item.Responses[strconv.Itoa(v.Code)] = res
			}

			if len(op.Res.Attrs.List) == 0 {
				item.Responses["200"] = map[string]any{
					"description": "success",
				}
			}

			if _, ok := out[path]; !ok {
				out[path] = make(map[string]openapiPathItem)
			}

			out[path][method] = item
		}
	}

	return out, func() (list []tagObject) {
		for k := range tags {
			list = append(list, tagObject{Name: k})
		}
		return
	}()
}
