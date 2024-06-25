package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v2 "github.com/pb33f/libopenapi/datamodel/high/v2"
	"github.com/pb33f/libopenapi/orderedmap"
	"gopkg.in/yaml.v3"
)

type swaggerParser struct {
	modelMapping           map[string]int64
	parametersMapping      map[string]*spec.Parameter
	globalParamtersMapping map[string]struct{}
}

type swaggerGenerator struct {
	modelNames map[int64]string
}

type swaggerSpec struct {
	Swagger          string                                `json:"swagger"`
	Info             *spec.Info                            `json:"info"`
	Tags             []tagObject                           `json:"tags,omitempty"`
	Host             string                                `json:"host,omitempty"`
	BasePath         string                                `json:"basePath"`
	Schemas          []string                              `json:"schemes,omitempty"`
	Definitions      map[string]jsonschema.Schema          `json:"definitions"`
	Parameters       map[string]openAPIParamter            `json:"parameters,omitempty"`
	Responses        map[string]any                        `json:"responses,omitempty"`
	Paths            map[string]map[string]swaggerPathItem `json:"paths"`
	GlobalParameters map[string]openAPIParamter            `json:"x-apicat-global-parameters,omitempty"`
}

type swaggerPathItem struct {
	Summary     string            `json:"summary"`
	Tags        []string          `json:"tags,omitempty"`
	Description string            `json:"description,omitempty"`
	OperationId string            `json:"operationId"`
	Consumes    []string          `json:"consumes,omitempty"`
	Produces    []string          `json:"produces,omitempty"`
	Parameters  []openAPIParamter `json:"parameters,omitempty"`
	Responses   map[string]any    `json:"responses,omitempty"`
}

func (s *swaggerParser) parseInfo(info *base.Info) spec.Info {
	return spec.Info{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}
}

func (s *swaggerParser) parseServers(in *v2.Swagger) []spec.Server {
	servers := make([]spec.Server, len(in.Schemes))
	if in.BasePath == "/" {
		in.BasePath = ""
	}
	for k, v := range in.Schemes {
		servers[k] = spec.Server{
			URL:         fmt.Sprintf("%s://%s%s", v, in.Host, in.BasePath),
			Description: v,
		}
	}
	return servers
}

func (s *swaggerParser) parseDefinitionModels(defs *v2.Definitions) (spec.DefinitionModels, error) {
	s.modelMapping = make(map[string]int64)
	models := make(spec.DefinitionModels, 0)
	if defs == nil {
		return models, nil
	}

	// orderedmap.Iterate(context.Context(), defs.Definitions)
	for pair := range orderedmap.Iterate(context.Background(), defs.Definitions) {
		k := pair.Key()
		v := pair.Value()

		js, err := jsonSchemaConverter(v)
		if err != nil {
			return nil, err
		}

		id := stringToUnid(k)
		s.modelMapping[k] = id

		if js.Type.First() != jsonschema.T_OBJ && js.Type.First() != jsonschema.T_ARR && js.AllOf == nil {
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
				ID:          id,
				Name:        k,
				Description: k,
				Schema:      parentJS,
			})
		} else {
			models = append(models, &spec.DefinitionModel{
				ID:          id,
				Name:        k,
				Description: k,
				Schema:      js,
			})
		}
	}

	return models, nil
}

func (s *swaggerParser) parseDefinitionParameters(in *v2.Swagger) error {
	if in.Parameters == nil {
		return nil
	}
	for pair := range orderedmap.Iterate(context.Background(), in.Parameters.Definitions) {
		parameter := pair.Value()
		if parameter.Schema != nil {
			js, err := jsonSchemaConverter(parameter.Schema)
			if err != nil {
				return err
			}
			k := fmt.Sprintf("%s-%s", parameter.In, parameter.Name)
			s.parametersMapping[k] = &spec.Parameter{
				ID:          stringToUnid(parameter.Name),
				Name:        parameter.Name,
				Description: parameter.Description,
				Required:    parameter.Required != nil && *parameter.Required,
				Schema:      js,
			}
		}
	}
	return nil
}

func (s *swaggerParser) parseJsonSchema(b *base.SchemaProxy) (*jsonschema.Schema, error) {
	js, err := jsonSchemaConverter(b)
	if err != nil {
		return nil, err
	}
	return js, nil
}

func (s *swaggerParser) parseDefinitionResponses(in *v2.Swagger) (spec.DefinitionResponses, error) {
	list := make(spec.DefinitionResponses, 0)
	if in.Responses == nil {
		return list, nil
	}

	for pair := range orderedmap.Iterate(context.Background(), in.Responses.Definitions) {
		responseName := pair.Key()
		response := pair.Value()

		headers := make([]*spec.Parameter, 0)
		content := make(spec.HTTPBody)

		if response.Headers != nil {
			for headerPair := range orderedmap.Iterate(context.Background(), response.Headers) {
				v := headerPair.Value()
				headers = append(headers, &spec.Parameter{
					Name: headerPair.Key(),
					Schema: &jsonschema.Schema{
						Type:        jsonschema.NewSchemaType(v.Type),
						Format:      v.Format,
						Description: v.Description,
						Examples:    v.Default,
					},
				})
			}
		}

		if response.Schema != nil {
			js, err := s.parseJsonSchema(response.Schema)
			if err != nil {
				return list, err
			}

			body := &spec.Body{Schema: js}
			if len(in.Produces) == 0 {
				content["application/json"] = body
				if response.Examples != nil {
					i := 0
					body.Examples = make(map[string]spec.Example)
					for examplePair := range orderedmap.Iterate(context.Background(), response.Examples.Values) {
						if example, err := json.Marshal(examplePair.Value()); err == nil {
							body.Examples[strconv.Itoa(i)] = spec.Example{
								Summary: examplePair.Key(),
								Value:   string(example),
							}
							i++
						}
					}
				}
			} else {
				i := 0
				body.Examples = make(map[string]spec.Example)
				for _, v := range in.Produces {
					content[v] = body
					if response.Examples != nil {
						emp, ok := response.Examples.Values.Get(v)
						if ok {
							if example, err := json.Marshal(emp); err == nil {
								body.Examples[strconv.Itoa(i)] = spec.Example{
									Summary: v,
									Value:   string(example),
								}
								i++
							}
						}
					}
				}
			}
		}
		list = append(list, &spec.DefinitionResponse{
			BasicResponse: spec.BasicResponse{
				ID:          stringToUnid(responseName),
				Name:        responseName,
				Header:      headers,
				Content:     content,
				Description: response.Description,
			},
		})
	}
	return list, nil
}

func (s *swaggerParser) parseGlobalParameters(inp *orderedmap.Map[string, *yaml.Node]) *spec.GlobalParameters {
	params := spec.NewGlobalParameters()
	if inp == nil {
		return params
	}
	global, ok := inp.Get("x-apicat-global-parameters")
	if !ok {
		return params
	}

	s.globalParamtersMapping = make(map[string]struct{})

	var globals map[string]any
	if err := global.Decode(&globals); err != nil {
		return params
	}

	for k, v := range globals {
		nb, err := json.Marshal(v)
		if err != nil {
			continue
		}

		p := &spec.Parameter{}
		json.Unmarshal(nb, p)
		in := strings.Index(k, "-")
		if in == -1 {
			continue
		}
		params.Add(k[:in], p)
		s.globalParamtersMapping[p.Name] = struct{}{}
	}
	return params
}

func (s *swaggerParser) parseRequest(in *v2.Swagger, info *v2.Operation) (*spec.CollectionHttpRequest, error) {
	request := spec.NewCollectionHttpRequest()

	var err error
	body := &spec.Body{}
	// 有效载荷application/x-www-form-urlencoded和multipart/form-data请求是通过使用form参数来描述，而不是body参数。
	formData := &jsonschema.Schema{
		Type:       jsonschema.NewSchemaType(jsonschema.T_OBJ),
		Properties: make(map[string]*jsonschema.Schema),
	}

	for _, v := range info.Parameters {
		// 这里引用 #/parameters 暂时无法获取
		// 直接展开
		k := fmt.Sprintf("%s-%s", v.In, v.Name)
		if sc, ok := s.parametersMapping[k]; ok {
			request.Attrs.Parameters.Add(v.In, sc)
			continue
		}
		if _, ok := s.globalParamtersMapping[v.Name]; ok {
			continue
		}

		required := v.Required != nil && *v.Required
		switch v.In {
		case "query", "header", "path", "cookie":
			request.Attrs.Parameters.Add(v.In,
				&spec.Parameter{
					Name:        v.Name,
					Description: v.Description,
					Required:    required,
					Schema: &jsonschema.Schema{
						Type:   jsonschema.NewSchemaType(v.Type),
						Format: v.Format,
					},
				},
			)
		case "formData":
			formData.Properties[v.Name] = &jsonschema.Schema{
				Type:        jsonschema.NewSchemaType(v.Type),
				Description: v.Description,
				Format:      v.Format,
				Default:     v.Default,
			}
			if required {
				formData.Required = append(formData.Required, v.Name)
			}
		case "body":
			body.Schema, err = s.parseJsonSchema(v.Schema)
			if err != nil {
				return nil, err
			}
		}
	}

	consumes := info.Consumes
	if len(info.Consumes) == 0 {
		// 从global获取
		consumes = in.Consumes
	}
	if len(consumes) == 0 && body.Schema != nil {
		consumes = []string{"application/json"}
	}

	for _, v := range consumes {
		if strings.Contains(v, "form") {
			request.Attrs.Content[v] = &spec.Body{Schema: formData}
		} else {
			if body.Schema != nil {
				request.Attrs.Content[v] = body
			}
		}
	}
	return request, nil
}

func (s *swaggerParser) parseResponse(info *v2.Operation) (*spec.CollectionHttpResponse, error) {
	if info.Responses == nil {
		return nil, nil
	}
	outresponses := spec.NewCollectionHttpResponse()
	// if info.Responses.Default != nil {
	// 	// 我们没有default
	// 	// todo
	// }
	for pair := range orderedmap.Iterate(context.Background(), info.Responses.Codes) {
		code := pair.Key()
		res := pair.Value()

		// res github.com/pb33f/libopenapi 不支持response ref 所以无法获取
		// 这里的common无法转换
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
		if v, ok := res.Extensions.Get("x-apicat-response-name"); ok {
			resp.Name = v.Value
		}
		resp.Description = res.Description
		resp.Content = make(spec.HTTPBody)
		resp.Header = make(spec.ParameterList, 0)

		// libopenapi not support response ref, in swagger 2.0
		// it's like dereference
		if res.GoLow().Schema.GetReference() != "" {
			ref := res.GoLow().Schema.GetReference()
			refs := fmt.Sprintf("#/definitions/responses/%d", stringToUnid(ref[strings.LastIndex(ref, "/")+1:]))
			resp.Reference = refs
			outresponses.Attrs.List = append(outresponses.Attrs.List, &resp)
			continue
		}

		if res.Headers != nil {
			for headerPair := range orderedmap.Iterate(context.Background(), res.Headers) {
				v := headerPair.Value()
				resp.Header = append(resp.Header, &spec.Parameter{
					Name: headerPair.Key(),
					Schema: &jsonschema.Schema{
						Type:        jsonschema.NewSchemaType(v.Type),
						Format:      v.Format,
						Description: v.Description,
						Examples:    v.Default,
					},
				})
			}
		}

		if res.Schema != nil {
			js, err := s.parseJsonSchema(res.Schema)
			if err != nil {
				return nil, err
			}

			for _, v := range info.Produces {
				body := &spec.Body{
					Schema:   js,
					Examples: make(map[string]spec.Example),
				}
				if res.Examples != nil {
					mp, ok := res.Examples.Values.Get(v)
					if ok {
						if example, err := json.Marshal(mp); err == nil {
							body.Examples["0"] = spec.Example{
								Summary: v,
								Value:   string(example),
							}
						}
					}
				}
				resp.Content[v] = body
			}
		}
		outresponses.Attrs.List = append(outresponses.Attrs.List, &resp)
	}
	if len(outresponses.Attrs.List) == 0 {
		outresponses.Attrs.List = append(outresponses.Attrs.List, &spec.Response{
			Code:          200,
			BasicResponse: spec.BasicResponse{Description: "success"},
		})
	}
	return outresponses, nil
}

func (s *swaggerParser) parseCollections(in *v2.Swagger, paths *v2.Paths) spec.Collections {
	collections := make(spec.Collections, 0)
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
			// 	content = append(content, spec.NewCollectionDoc(v).ToCollectionNode())
			// }

			// request
			if req, err := s.parseRequest(in, info); err != nil {
				continue
			} else {
				content = append(content, req.ToCollectionNode())
			}

			// response
			if res, err := s.parseResponse(info); err != nil {
				continue
			} else {
				content = append(content, res.ToCollectionNode())
			}

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

func (s *swaggerGenerator) convertJsonSchema(v *jsonschema.Schema) *jsonschema.Schema {
	if v == nil {
		return v
	}
	return convertJsonSchemaRef(v, "2.0", s.modelNames)
}

func (s *swaggerGenerator) generateBase(in *spec.Spec) *swaggerSpec {
	s.modelNames = map[int64]string{}
	out := &swaggerSpec{
		Swagger: "2.0",
		Info: &spec.Info{
			Title:       in.Info.Title,
			Description: in.Info.Description,
			Version:     in.Info.Version,
		},
		Definitions: make(map[string]jsonschema.Schema),
	}

	for _, v := range in.Servers {
		u, err := url.Parse(v.URL)
		if err != nil {
			continue
		}

		if out.Host == "" {
			out.Host = u.Host
			out.BasePath = u.Path
		}
		out.Schemas = append(out.Schemas, u.Scheme)
		// just need fist one
		break
	}

	definitionModels := spec.DefinitionModels{}
	for _, v := range in.Definitions.Schemas {
		if v.Type == string(spec.TYPE_CATEGORY) {
			items := v.ItemsTreeToList()
			for _, item := range items {
				s.modelNames[item.ID] = strings.ReplaceAll(item.Name, " ", "")
			}
			definitionModels = append(definitionModels, items...)
		} else {
			s.modelNames[v.ID] = strings.ReplaceAll(v.Name, " ", "")
			definitionModels = append(definitionModels, v)
		}
	}

	for _, v := range definitionModels {
		name_id := fmt.Sprintf("%s-%d", strings.ReplaceAll(v.Name, " ", ""), v.ID)
		out.Definitions[name_id] = *s.convertJsonSchema(v.Schema)
	}

	globalParams := in.Globals.Parameters.ToMap()
	out.GlobalParameters = make(map[string]openAPIParamter)
	for in, paramList := range globalParams {
		for _, p := range paramList {
			name_id := fmt.Sprintf("%s-%d", p.Name, p.ID)
			out.GlobalParameters[name_id] = toParameter(p, in, "2.0")
		}
	}

	if out.BasePath == "" {
		out.BasePath = "/"
	}
	if len(in.Definitions.Responses) > 0 {
		out.Responses = make(map[string]any)
		for _, v := range in.Definitions.Responses {
			if v.Type == string(spec.TYPE_CATEGORY) {
				items := v.ItemsTreeToList()
				for _, item := range items {
					name_id := fmt.Sprintf("%s-%d", strings.ReplaceAll(item.Name, " ", ""), item.ID)
					out.Responses[name_id] = s.generateResponseWithoutRef(&item.BasicResponse)
				}
			} else {
				name_id := fmt.Sprintf("%s-%d", strings.ReplaceAll(v.Name, " ", ""), v.ID)
				out.Responses[name_id] = s.generateResponseWithoutRef(&v.BasicResponse)
			}
		}
	}

	return out
}

func (s *swaggerGenerator) generateReqParams(collectionReq spec.CollectionHttpRequest, globalsParmaters *spec.GlobalParameters) []openAPIParamter {
	// 添加启用的全局参数
	out := globalToLocalParameters(globalsParmaters, true, collectionReq.Attrs.GlobalExcepts.ToMap())

	for in, params := range collectionReq.Attrs.Parameters.ToMap() {
		switch in {
		case "header", "query", "path", "cookie":
			for _, v := range params {
				newv := *v
				newv.Schema = s.convertJsonSchema(v.Schema)
				out = append(out, toParameter(&newv, in, "2.0"))
			}
		}
	}

	if collectionReq.Attrs.Content == nil {
		return out
	}

	var hasBody bool
	for contentType, body := range collectionReq.Attrs.Content {
		// contentType incloud form use parameters in
		if strings.Contains(contentType, "form") {
			if body.Schema == nil {
				continue
			}

			if num := len(body.Schema.Type.List()); num == 0 {
				continue
			}

			typ := body.Schema.Type.First()
			if typ != jsonschema.T_OBJ || body.Schema.Properties == nil {
				continue
			}

			for k, v := range body.Schema.Properties {
				content := openAPIParamter{
					Name:        k,
					In:          "formData",
					Type:        v.Type.First(),
					Description: v.Description,
					Schema:      s.convertJsonSchema(v),
					Required: func() bool {
						for _, r := range v.Required {
							if r == k {
								return true
							}
						}
						return false
					}(),
				}
				if v != nil {
					t := v.Type.List()
					if len(t) > 0 && t[0] == "file" {
						content.Type = t[0]
					}
				}
				out = append(out, content)
			}
		} else {
			if hasBody {
				continue
			}

			out = append(out, openAPIParamter{
				Name:     "body",
				Schema:   s.convertJsonSchema(body.Schema),
				In:       "body",
				Required: true,
			})
			hasBody = true
		}
	}
	return out
}

func (s *swaggerGenerator) generateResponseWithoutRef(resp *spec.BasicResponse) map[string]any {
	response := map[string]any{
		"x-apicat-response-name": resp.Name,
		"description":            resp.Description,
	}

	if len(resp.Header) > 0 {
		header := make(map[string]any)
		for _, v := range resp.Header {
			if v.Schema.Description == "" {
				v.Schema.Description = v.Description
			}
			v.Schema.Default = v.Schema.Examples
			v.Schema.Examples = nil
			header[v.Name] = v.Schema
		}
		response["headers"] = header
	}

	if resp.Content != nil {
		for k, v := range resp.Content {
			response["schema"] = s.convertJsonSchema(v.Schema)
			if v.Examples != nil {
				for _, v := range v.Examples {
					response["examples"] = map[string]any{
						k: v,
					}
					break
				}
			}
			break
		}
	}
	return response
}

func (s *swaggerGenerator) generateResponse(resp *spec.Response, definitionsResps spec.DefinitionResponses) map[string]any {
	if resp.Reference != "" {
		if strings.HasPrefix(resp.Reference, "#/definitions/responses/") {
			x := definitionsResps.FindByID(
				toInt64(getRefName(resp.Reference)),
			)
			if x != nil {
				name_id := fmt.Sprintf("%s-%d", strings.ReplaceAll(x.Name, " ", ""), x.ID)
				return map[string]any{
					"$ref": "#/responses/" + name_id,
				}
			}
		}
		return nil
	}
	return s.generateResponseWithoutRef(&resp.BasicResponse)
}

func (s *swaggerGenerator) getRefResponseContentType(resp *spec.Response, definitionsResps spec.DefinitionResponses) string {
	x := definitionsResps.FindByID(
		toInt64(getRefName(resp.Reference)),
	)
	if x != nil {
		for k := range x.Content {
			return k
		}
	}
	return ""
}

func (s *swaggerGenerator) generatePathResponse(resp spec.CollectionHttpResponse, definitionsResps spec.DefinitionResponses) (map[string]any, []string) {
	product := map[string]struct{}{}
	result := make(map[string]any)

	for _, r := range resp.Attrs.List {
		result[strconv.Itoa(r.Code)] = s.generateResponse(r, definitionsResps)
		if r.Reference != "" {
			if ct := s.getRefResponseContentType(r, definitionsResps); ct != "" {
				product[ct] = struct{}{}
			}
		} else {
			for k := range r.Content {
				if _, ok := product[k]; !ok {
					product[k] = struct{}{}
					continue
				}
			}
		}
	}
	if len(result) == 0 {
		result["default"] = map[string]string{
			"description": "success",
		}
	}
	return result, func() (ret []string) {
		if len(product) == 0 {
			return []string{"application/json"}
		}
		for k := range product {
			ret = append(ret, k)
		}
		return
	}()
}

func (s *swaggerGenerator) generatePaths(in *spec.Spec) (map[string]map[string]swaggerPathItem, []tagObject) {
	out := make(map[string]map[string]swaggerPathItem)
	tags := make(map[string]struct{})

	for path, methods := range deepGetHttpCollection(&in.Collections) {
		if path == "" {
			continue
		}
		for method, op := range methods {
			reslist, product := s.generatePathResponse(op.Res, in.Definitions.Responses)
			if len(reslist) == 0 {
				reslist["default"] = &jsonschema.Schema{Description: "success"}
			}
			item := swaggerPathItem{
				Summary:     op.Title,
				Description: op.Description,
				OperationId: op.OperatorID,
				Parameters:  s.generateReqParams(op.Req, in.Globals.Parameters),
				Produces:    product,
				Responses:   reslist,
				Tags:        op.Tags,
			}
			for k := range op.Req.Attrs.Content {
				item.Consumes = append(item.Consumes, k)
			}
			if _, ok := out[path]; !ok {
				out[path] = make(map[string]swaggerPathItem)
			}
			for _, v := range op.Tags {
				tags[v] = struct{}{}
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
