package openapi

import (
	"fmt"
	spec2 "github.com/apicat/apicat/backend/module/spec"
	jsonschema2 "github.com/apicat/apicat/backend/module/spec/jsonschema"
	"github.com/apicat/apicat/backend/module/spec/markdown"
	"net/url"
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v2 "github.com/pb33f/libopenapi/datamodel/high/v2"
)

type fromSwagger struct {
	schemaMapping     map[string]int64
	parametersMapping map[string]int64
}

func (s *fromSwagger) parseInfo(info *base.Info) *spec2.Info {
	return &spec2.Info{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}
}

func (s *fromSwagger) parseServers(in *v2.Swagger) []*spec2.Server {
	srvs := make([]*spec2.Server, len(in.Schemes))
	for k, v := range in.Schemes {
		srvs[k] = &spec2.Server{
			URL:         fmt.Sprintf("%s://%s%s", v, in.Host, in.BasePath),
			Description: v,
		}
	}
	return srvs
}

func (s *fromSwagger) parseDefinetions(defs *v2.Definitions) spec2.Schemas {
	s.schemaMapping = make(map[string]int64)
	defines := make(spec2.Schemas, 0)
	if defs == nil {
		return defines
	}

	for k, v := range defs.Definitions {
		js, err := jsonSchemaConverter(v)
		if err != nil {
			panic(err)
		}
		id := stringToUnid(k)
		s.schemaMapping[k] = id
		defines = append(defines, &spec2.Schema{
			ID:          id,
			Name:        k,
			Description: k,
			Schema:      js,
		})
	}

	return defines
}

func (s *fromSwagger) parseParametersDefine(in *v2.Swagger) spec2.Schemas {
	s.parametersMapping = make(map[string]int64)
	ps := make(spec2.Schemas, 0)
	// mapping key:swagger paranmters key value:apicat paramter id
	if in.Parameters == nil {
		return ps
	}
	repeat := map[string]int{}
	// 因为swagger参数是name+in apicat没有in 所以name不能重复 这里只处理不重复的
	for _, v := range in.Parameters.Definitions {
		repeat[v.Name]++
	}
	for key, v := range in.Parameters.Definitions {
		if repeat[v.Name] > 1 {
			continue
		}
		id := stringToUnid(key)
		s.parametersMapping[key] = id
		ps = append(ps, &spec2.Schema{
			ID:          id,
			Name:        v.Name,
			Description: v.Description,
			Required:    v.Required != nil && *v.Required,
			Schema: &jsonschema2.Schema{
				Type:   jsonschema2.CreateSliceOrOne(v.Type),
				Format: v.Format,
			},
		})
	}
	return ps
}

// 主要处理$ref引用问题
func (s *fromSwagger) parseContent(b *base.SchemaProxy) *jsonschema2.Schema {
	js, err := jsonSchemaConverter(b)
	if err != nil {
		panic(err)
	}
	return js
}

func (s *fromSwagger) parseRequest(in *v2.Swagger, info *v2.Operation) spec2.HTTPRequestNode {
	// parameters := &spec.HttpParameters{}
	request := spec2.HTTPRequestNode{
		Content: make(spec2.HTTPBody),
	}
	request.Parameters.Fill()
	var body *jsonschema2.Schema
	// 有效载荷application/x-www-form-urlencoded和multipart/form-data请求是通过使用form参数来描述，而不是body参数。
	formData := &jsonschema2.Schema{
		Type:       jsonschema2.CreateSliceOrOne("object"),
		Properties: make(map[string]*jsonschema2.Schema),
	}

	for _, v := range info.Parameters {
		// 这里引用 #/parameters 暂时无法获取
		// 直接展开
		required := v.Required != nil && *v.Required
		switch v.In {
		case "query", "header", "path":
			request.Parameters.Add(v.In,
				&spec2.Schema{
					Name:        v.Name,
					Description: v.Description,
					Required:    required,
					Schema: &jsonschema2.Schema{
						Type:   jsonschema2.CreateSliceOrOne(v.Type),
						Format: v.Format,
					},
				},
			)
		case "formData":
			formData.Properties[v.Name] = &jsonschema2.Schema{
				Type:        jsonschema2.CreateSliceOrOne(v.Type),
				Description: v.Description,
				Format:      v.Format,
				Default:     v.Default,
			}
			if required {
				formData.Required = append(formData.Required, v.Name)
			}
		case "body":
			body = s.parseContent(v.Schema)
		}
	}

	consumes := info.Consumes
	if len(info.Consumes) == 0 {
		// 从global获取
		consumes = in.Consumes
	}
	// 有些文件没有consunmer 给个默认 否则body不知道什么是mine
	if len(consumes) == 0 && body != nil {
		consumes = []string{defaultSwaggerConsumerProduce}
	}

	for _, v := range consumes {
		if strings.Contains(v, "form") {
			request.Content[v] = &spec2.Schema{Schema: formData}
		} else {
			if body != nil {
				request.Content[v] = &spec2.Schema{Schema: body}
			}
		}
	}
	return request
}

// parseResponsesDefine 因为swagger response 没有code 所以这个只能放到definition里
func (s *fromSwagger) parseResponsesDefine(in *v2.Swagger) []spec2.HTTPResponseDefine {
	list := make([]spec2.HTTPResponseDefine, 0)
	if in.Responses == nil {
		return list
	}
	for key, res := range in.Responses.Definitions {
		header := make([]*spec2.Schema, 0)
		content := make(spec2.HTTPBody)
		if res.Headers != nil {
			for k, v := range res.Headers {
				header = append(header, &spec2.Schema{
					Name: k,
					Schema: &jsonschema2.Schema{
						Type:        jsonschema2.CreateSliceOrOne(v.Type),
						Format:      v.Format,
						Description: v.Description,
					},
				})
			}
		}
		if res.Schema != nil {
			js := s.parseContent(res.Schema)
			sh := &spec2.Schema{
				Schema: js,
			}
			if len(in.Produces) == 0 {
				content[defaultSwaggerConsumerProduce] = sh
			} else {
				for _, v := range in.Produces {
					content[v] = sh
				}
			}
		}
		list = append(list, spec2.HTTPResponseDefine{
			Name:    key,
			Header:  header,
			Content: content,
		})
	}
	return list
}

func (s *fromSwagger) parseResponse(info *v2.Operation) *spec2.HTTPResponsesNode {
	if info.Responses == nil {
		return nil
	}
	var outresponses spec2.HTTPResponsesNode
	// if info.Responses.Default != nil {
	// 	// 我们没有default
	// 	// todo
	// }
	for code, res := range info.Responses.Codes {
		// res github.com/pb33f/libopenapi 不支持response ref 所以无法获取
		// 这里的common无法转换
		c, err := strconv.Atoi(code)
		if err != nil {
			continue
		}
		resp := spec2.HTTPResponse{
			Code: c,
		}
		resp.Name = res.Description
		resp.Description = res.Description
		resp.Content = make(spec2.HTTPBody)
		resp.Header = make(spec2.Schemas, 0)
		if res.Headers != nil {
			for k, v := range res.Headers {
				resp.Header = append(resp.Header, &spec2.Schema{
					Name: k,
					Schema: &jsonschema2.Schema{
						Type:        jsonschema2.CreateSliceOrOne(v.Type),
						Format:      v.Format,
						Description: v.Description,
					},
				})
			}
		}
		if res.Schema != nil {
			js := s.parseContent(res.Schema)
			for _, v := range info.Produces {
				resp.Content[v] = &spec2.Schema{
					Schema: js,
				}
			}
		}
		outresponses.List = append(outresponses.List, resp)
	}
	if len(outresponses.List) == 0 {
		outresponses.List = append(outresponses.List, spec2.HTTPResponse{
			Code:               200,
			HTTPResponseDefine: spec2.HTTPResponseDefine{Description: "success"},
		})
	}
	return &outresponses
}

func (s *fromSwagger) parseCollections(in *v2.Swagger, paths *v2.Paths) []*spec2.CollectItem {
	collects := make([]*spec2.CollectItem, 0)
	for path, p := range paths.PathItems {
		op := p.GetOperations()
		for method, info := range op {
			content := []*spec2.NodeProxy{
				spec2.MuseCreateNodeProxy(
					spec2.WarpHTTPNode(spec2.HTTPURLNode{
						Path:   path,
						Method: method,
					}),
				),
			}

			// parse markdown to doc
			doctree := markdown.ToDocment([]byte(info.Description))
			for _, v := range doctree.Items {
				content = append(content, spec2.MuseCreateNodeProxy(v))
			}

			// request
			req := spec2.WarpHTTPNode(s.parseRequest(in, info))
			content = append(content, spec2.MuseCreateNodeProxy(req))
			// response
			res := spec2.WarpHTTPNode(s.parseResponse(info))
			content = append(content, spec2.MuseCreateNodeProxy(res))

			title := info.Summary
			if title == "" {
				title = path
			}

			collects = append(collects, &spec2.CollectItem{
				Type:    spec2.ContentItemTypeHttp,
				Title:   title,
				Tags:    info.Tags,
				Content: content,
			})
		}
	}
	return collects
}

//

type swaggerSpec struct {
	Swagger     string                                `json:"swagger"`
	Info        *spec2.Info                           `json:"info"`
	Tags        []tagObject                           `json:"tags,omitempty"`
	Host        string                                `json:"host,omitempty"`
	BasePath    string                                `json:"basePath"`
	Schemas     []string                              `json:"schemas,omitempty"`
	Definitions map[string]jsonschema2.Schema         `json:"definitions"`
	Parameters  map[string]openAPIParamter            `json:"parameters,omitempty"`
	Responses   map[string]any                        `json:"responses,omitempty"`
	Paths       map[string]map[string]swaggerPathItem `json:"paths"`
}

type toSwagger struct {
	schemas map[int64]string
}

func (s *toSwagger) toBase(in *spec2.Spec) *swaggerSpec {
	s.schemas = map[int64]string{}
	out := &swaggerSpec{
		Swagger: "2.0",
		Info: &spec2.Info{
			Title:       in.Info.Title,
			Description: in.Info.Description,
			Version:     in.Info.Version,
		},
		Definitions: make(map[string]jsonschema2.Schema),
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
	}
	for _, v := range in.Definitions.Schemas {
		s.schemas[v.ID] = v.Name
	}
	for _, v := range in.Definitions.Schemas {
		out.Definitions[v.Name] = *s.convertJSONSchema(v.Schema)
	}

	globalParam := in.Globals.Parameters
	m := globalParam.Map()
	out.Parameters = make(map[string]openAPIParamter)
	for in, ps := range m {
		for _, p := range ps {
			out.Parameters[fmt.Sprintf("%s-%s", in, p.Name)] = toParameter(p, in)
		}
	}

	if out.BasePath == "" {
		out.BasePath = "/"
	}
	if len(in.Definitions.Responses) > 0 {
		out.Responses = make(map[string]any)
		for _, v := range in.Definitions.Responses {
			out.Responses[v.Name] = s.parseResponse(in, v)
		}
	}
	return out
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

func (s *toSwagger) toReqParameters(ps spec2.HTTPRequestNode, spe *spec2.Spec) []openAPIParamter {
	// 添加启用的全局参数
	out := toParameterGlobal(spe.Globals.Parameters, true, ps.GlobalExcepts)
	for in, params := range ps.Parameters.Map() {
		switch in {
		case "header", "query", "path":
			for _, v := range params {
				if v.Reference != nil {
					// 解开公共参数
					if id := toInt64(getRefName(*v.Reference)); id != 0 {
						v = spe.Definitions.Parameters.LookupID(id)
					}
				}
				newv := *v
				newv.Schema = s.convertJSONSchema(v.Schema)
				out = append(out, toParameter(&newv, in))
			}
		}
	}
	if ps.Content == nil {
		return out
	}
	var hasBody bool
	for ct, c := range ps.Content {
		// contentType incloud form use parameters in
		if strings.Contains(ct, "form") {
			if c.Schema == nil {
				continue
			}
			if n := len(c.Schema.Type.Value()); n == 0 {
				continue
			}
			tp := c.Schema.Type.Value()[0]
			if tp != "object" || c.Schema.Properties == nil {
				continue
			}

			for k, v := range c.Schema.Properties {
				content := openAPIParamter{
					Name:        k,
					In:          "formData",
					Description: v.Description,
					Schema:      s.convertJSONSchema(v),
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
					t := v.Type.Value()
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
				Name:        "body",
				Description: c.Description,
				Schema:      s.convertJSONSchema(c.Schema),
				In:          "body",
				Required:    true,
			})
			hasBody = true
		}
	}
	return out
}

func (s *toSwagger) convertJSONSchema(v *jsonschema2.Schema) *jsonschema2.Schema {
	if v == nil {
		return v
	}
	return toConvertJSONSchemaRef(v, "2.0", s.schemas)
}

func (s *toSwagger) parseResponse(in *spec2.Spec, res spec2.HTTPResponseDefine) map[string]any {
	if res.Reference != nil {
		ref := *res.Reference
		if strings.HasPrefix(ref, "#/definitions/responses/") {
			x := in.Definitions.Responses.LookupID(
				toInt64(getRefName(ref)),
			)
			if x != nil {
				return map[string]any{
					"$ref": "#/responses/" + x.Name,
				}
			}
		}
		return nil
	}
	resp := map[string]any{
		"description": res.Description,
	}
	if len(res.Header) > 0 {
		h := make(map[string]any)
		for _, v := range res.Header {
			if v.Schema.Description == "" {
				v.Schema.Description = v.Description
			}
			h[v.Name] = v.Schema
		}
		resp["headers"] = h
	}
	if res.Content != nil {
		for _, v := range res.Content {
			resp["schema"] = s.convertJSONSchema(v.Schema)
			break
		}
	}
	return resp
}

func (s *toSwagger) toPathResponse(in *spec2.Spec, resp []spec2.HTTPResponse) (map[string]any, []string) {
	product := map[string]struct{}{}
	reslist := make(map[string]any)
	for _, r := range resp {
		reslist[strconv.Itoa(r.Code)] = s.parseResponse(in, r.HTTPResponseDefine)
		for k := range r.Content {
			if _, ok := product[k]; !ok {
				product[k] = struct{}{}
			}
		}
	}
	if len(reslist) == 0 {
		reslist["default"] = map[string]string{
			"description": "success",
		}
	}
	// for _, r := range resp {
	// 	v := r
	// 	if v.Reference != nil {
	// 		switch {
	// 		case strings.HasPrefix(*v.Reference, "#/definitions/responses/"):
	// 			x := in.Definitions.Responses.Lookup(
	// 				getRefName(*v.Reference),
	// 			)
	// 			if x != nil {
	// 				v.HTTPResponseDefine = *x
	// 				v.Reference = nil
	// 			}
	// 		case strings.HasPrefix(*v.Reference, "#/commons/responses/"):
	// 			x := in.Common.Responses.Lookup(
	// 				getRefName(*v.Reference),
	// 			)
	// 			v = *x
	// 		default:
	// 			// not support
	// 			continue
	// 		}
	// 	}
	// 	var res *spec.Schema
	// 	for k := range v.Content {
	// 		if _, ok := product[k]; !ok {
	// 			product[k] = struct{}{}
	// 		}
	// 		if res != nil {
	// 			break
	// 		}
	// 		res = &spec.Schema{
	// 			Schema:      s.convertJSONSchema(v.Content[k].Schema),
	// 			Description: v.Description,
	// 		}

	// 	}
	// 	if res == nil {
	// 		res = &spec.Schema{Description: v.Description}
	// 	}
	// 	reslist[strconv.Itoa(v.Code)] = res
	// }
	return reslist, func() (ret []string) {
		if len(product) == 0 {
			return []string{defaultSwaggerConsumerProduce}
		}
		for k := range product {
			ret = append(ret, k)
		}
		return
	}()
}

func (s *toSwagger) toPaths(in *spec2.Spec) (map[string]map[string]swaggerPathItem, []tagObject) {
	out := make(map[string]map[string]swaggerPathItem)
	tags := make(map[string]struct{})
	for path, ops := range walkHttpCollection(in) {
		if path == "" {
			continue
		}
		for method, op := range ops {
			reslist, product := s.toPathResponse(in, op.Res.List)
			if len(reslist) == 0 {
				reslist["default"] = &spec2.Schema{Description: "success"}
			}
			item := swaggerPathItem{
				Summary:     op.Title,
				Description: op.Description,
				OperationId: op.OperatorID,
				Parameters:  s.toReqParameters(op.Req, in),
				Produces:    product,
				Responses:   reslist,
				Tags:        op.Tags,
			}
			for k := range op.Req.Content {
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
