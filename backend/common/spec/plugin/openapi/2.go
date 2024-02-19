package openapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/apicat/apicat/backend/common/spec"
	"github.com/apicat/apicat/backend/common/spec/jsonschema"
	"github.com/apicat/apicat/backend/common/spec/markdown"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v2 "github.com/pb33f/libopenapi/datamodel/high/v2"
)

type fromSwagger struct {
	schemaMapping          map[string]int64
	parametersMapping      map[string]*spec.Schema
	globalParamtersMappint map[string]struct{}
}

func (s *fromSwagger) parseInfo(info *base.Info) *spec.Info {
	return &spec.Info{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}
}

func (s *fromSwagger) parseServers(in *v2.Swagger) []*spec.Server {
	srvs := make([]*spec.Server, len(in.Schemes))
	if in.BasePath == "/" {
		in.BasePath = ""
	}
	for k, v := range in.Schemes {
		srvs[k] = &spec.Server{
			URL:         fmt.Sprintf("%s://%s%s", v, in.Host, in.BasePath),
			Description: v,
		}
	}
	return srvs
}

func (s *fromSwagger) parseDefinetions(defs *v2.Definitions) spec.Schemas {
	s.schemaMapping = make(map[string]int64)
	defines := make(spec.Schemas, 0)
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
		defines = append(defines, &spec.Schema{
			ID:          id,
			Name:        k,
			Description: k,
			Schema:      js,
		})
	}

	return defines
}

func (s *fromSwagger) parseParametersDefine(in *v2.Swagger) (res spec.HTTPParameters) {
	s.parametersMapping = make(map[string]*spec.Schema)
	// mapping key:swagger paranmters key value:apicat paramter id
	res.Fill()
	if in.Parameters == nil {
		return res
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
		ss := &spec.Schema{
			ID:          id,
			Name:        v.Name,
			Description: v.Description,
			Required:    v.Required != nil && *v.Required,
			Schema: &jsonschema.Schema{
				Type:   jsonschema.CreateSliceOrOne(v.Type),
				Format: v.Format,
			},
		}
		s.parametersMapping[key] = ss
		res.Add(v.In, ss)
	}
	return res
}

func (s *fromSwagger) parseGlobal(inp map[string]any) (res spec.Global) {
	res.Parameters.Fill()
	if inp == nil {
		return res
	}
	global, ok := inp["x-apicat-global-parameters"]
	if !ok {
		return res
	}

	s.globalParamtersMappint = make(map[string]struct{})
	for k, v := range global.(map[string]any) {

		nb, err := json.Marshal(v)
		if err != nil {
			continue
		}

		sc := &spec.Schema{}
		json.Unmarshal(nb, sc)
		in := strings.Index(k, "-")
		if in == -1 {
			continue
		}
		res.Parameters.Add(k[:in], sc)
		s.globalParamtersMappint[sc.Name] = struct{}{}
	}
	return res
}

// 主要处理$ref引用问题
func (s *fromSwagger) parseContent(b *base.SchemaProxy) *jsonschema.Schema {
	js, err := jsonSchemaConverter(b)
	if err != nil {
		panic(err)
	}
	return js
}

func (s *fromSwagger) parseRequest(in *v2.Swagger, info *v2.Operation) spec.HTTPRequestNode {
	// parameters := &spec.HttpParameters{}
	request := spec.HTTPRequestNode{
		Content: make(spec.HTTPBody),
	}
	request.Parameters.Fill()
	body := &jsonschema.Schema{}
	// 有效载荷application/x-www-form-urlencoded和multipart/form-data请求是通过使用form参数来描述，而不是body参数。
	formData := &jsonschema.Schema{
		Type:       jsonschema.CreateSliceOrOne("object"),
		Properties: make(map[string]*jsonschema.Schema),
	}

	for _, v := range info.Parameters {
		// 这里引用 #/parameters 暂时无法获取
		// 直接展开
		_, ok := s.parametersMapping[v.Name]
		if ok {
			continue
		}
		_, ok = s.globalParamtersMappint[v.Name]
		if ok {
			continue
		}
		required := v.Required != nil && *v.Required
		switch v.In {
		case "query", "header", "path", "cookie":
			request.Parameters.Add(v.In,
				&spec.Schema{
					Name:        v.Name,
					Description: v.Description,
					Required:    required,
					Schema: &jsonschema.Schema{
						Type:   jsonschema.CreateSliceOrOne(v.Type),
						Format: v.Format,
					},
				},
			)
		case "formData":
			formData.Properties[v.Name] = &jsonschema.Schema{
				Type:        jsonschema.CreateSliceOrOne(v.Type),
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
			request.Content[v] = &spec.Schema{Schema: formData}
		} else {
			if body != nil {
				request.Content[v] = &spec.Schema{Schema: body}
			}
		}
	}
	return request
}

// parseResponsesDefine 因为swagger response 没有code 所以这个只能放到definition里
func (s *fromSwagger) parseResponsesDefine(in *v2.Swagger) []spec.HTTPResponseDefine {
	list := make([]spec.HTTPResponseDefine, 0)
	if in.Responses == nil {
		return list
	}
	for key, res := range in.Responses.Definitions {
		header := make([]*spec.Schema, 0)
		content := make(spec.HTTPBody)
		if res.Headers != nil {
			for k, v := range res.Headers {
				header = append(header, &spec.Schema{
					Name: k,
					Schema: &jsonschema.Schema{
						Type:        jsonschema.CreateSliceOrOne(v.Type),
						Format:      v.Format,
						Description: v.Description,
					},
				})
			}
		}
		if res.Schema != nil {
			js := s.parseContent(res.Schema)
			sh := &spec.Schema{
				Schema: js,
			}
			if len(in.Produces) == 0 {
				content[defaultSwaggerConsumerProduce] = sh
				if res.Examples != nil {
					sh.Examples = make(map[string]*spec.Example)
					for k, v := range res.Examples.Values {
						sh.Examples[k] = &spec.Example{
							Summary: "example",
							Value:   v,
						}
					}
				}
			} else {
				for _, v := range in.Produces {
					content[v] = sh
					if res.Examples != nil {
						emp, ok := res.Examples.Values[v]
						if ok {
							sh.Examples = map[string]*spec.Example{
								v: {
									Summary: "example",
									Value:   emp,
								},
							}
						}
					}
				}
			}
		}
		list = append(list, spec.HTTPResponseDefine{
			ID:          stringToUnid(key),
			Name:        key,
			Header:      header,
			Content:     content,
			Description: res.Description,
		})
	}
	return list
}

func (s *fromSwagger) parseResponse(info *v2.Operation) *spec.HTTPResponsesNode {
	if info.Responses == nil {
		return nil
	}
	var outresponses spec.HTTPResponsesNode
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
		resp := spec.HTTPResponse{
			Code: c,
		}
		if _, ok := res.Extensions["x-apicat-response-name"]; ok {
			resp.Name = res.Extensions["x-apicat-response-name"].(string)
		}
		resp.Description = res.Description
		resp.Content = make(spec.HTTPBody)
		resp.Header = make(spec.Schemas, 0)
		if res.Schema.GetReference() != "" {
			ref := res.Schema.GetReference()
			refs := fmt.Sprintf("#/definitions/responses/%d", stringToUnid(ref[strings.LastIndex(ref, "/")+1:]))
			resp.Reference = &refs
			outresponses.List = append(outresponses.List, &resp)
			continue
		}
		if res.Headers != nil {
			for k, v := range res.Headers {
				resp.Header = append(resp.Header, &spec.Schema{
					Name: k,
					Schema: &jsonschema.Schema{
						Type:        jsonschema.CreateSliceOrOne(v.Type),
						Format:      v.Format,
						Description: v.Description,
					},
				})
			}
		}
		if res.Schema != nil {
			js := s.parseContent(res.Schema)
			for _, v := range info.Produces {
				sh := &spec.Schema{
					Schema: js,
				}
				if res.Examples != nil {
					mp, ok := res.Examples.Values[v]
					if ok {
						sh.Examples = map[string]*spec.Example{
							v: {
								Summary: "example",
								Value:   mp,
							},
						}
					}
				}
				resp.Content[v] = sh
			}
		}
		outresponses.List = append(outresponses.List, &resp)
	}
	if len(outresponses.List) == 0 {
		outresponses.List = append(outresponses.List, &spec.HTTPResponse{
			Code:               200,
			HTTPResponseDefine: spec.HTTPResponseDefine{Description: "success"},
		})
	}
	return &outresponses
}

func (s *fromSwagger) parseCollections(in *v2.Swagger, paths *v2.Paths) []*spec.CollectItem {
	collects := make([]*spec.CollectItem, 0)
	for path, p := range paths.PathItems {
		op := p.GetOperations()
		for method, info := range op {
			content := []*spec.NodeProxy{
				spec.MuseCreateNodeProxy(
					spec.WarpHTTPNode(spec.HTTPURLNode{
						Path:   path,
						Method: method,
					}),
				),
			}

			// parse markdown to doc
			doctree := markdown.ToDocment([]byte(info.Description))
			for _, v := range doctree.Items {
				content = append(content, spec.MuseCreateNodeProxy(v))
			}

			// request
			req := spec.WarpHTTPNode(s.parseRequest(in, info))
			content = append(content, spec.MuseCreateNodeProxy(req))
			// response
			res := spec.WarpHTTPNode(s.parseResponse(info))
			content = append(content, spec.MuseCreateNodeProxy(res))

			title := info.Summary
			if title == "" {
				title = path
			}

			collects = append(collects, &spec.CollectItem{
				Type:    spec.ContentItemTypeHttp,
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

type toSwagger struct {
	schemas map[int64]string
}

func (s *toSwagger) toBase(in *spec.Spec) *swaggerSpec {
	s.schemas = map[int64]string{}
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
	}

	ss := spec.Schemas{}
	// fill s.schemas
	for _, v := range in.Definitions.Schemas {
		if v.Type == string(spec.ContentItemTypeDir) {
			items := v.ItemsTreeToList()
			for _, item := range items {
				s.schemas[item.ID] = item.Name
			}
			ss = append(ss, items...)
		} else {
			s.schemas[v.ID] = v.Name
			ss = append(ss, v)
		}
	}
	for _, v := range ss {
		name_id := fmt.Sprintf("%s-%d", v.Name, v.ID)
		v.Schema.Description = v.Description
		out.Definitions[name_id] = *s.convertJSONSchema(v.Schema)
	}

	globalParam := in.Globals.Parameters
	m := globalParam.Map()
	out.GlobalParameters = make(map[string]openAPIParamter)
	for in, ps := range m {
		for _, p := range ps {
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
			if v.Type == string(spec.ContentItemTypeDir) {
				items := v.ItemsTreeToList()
				for _, item := range items {
					name_id := fmt.Sprintf("%s-%d", item.Name, item.ID)
					out.Responses[name_id] = s.parseResponse(in, item)
				}
			} else {
				name_id := fmt.Sprintf("%s-%d", v.Name, v.ID)
				out.Responses[name_id] = s.parseResponse(in, v)
			}
		}
	}

	definitionParameters := in.Definitions.Parameters
	dp := definitionParameters.Map()
	parameters := make(map[string]openAPIParamter)
	for in, ps := range dp {
		for _, p := range ps {
			opp := toParameter(p, in, "2.0")
			// remove format from parameter, because it's not support in openapi3.components.parameters.item
			opp.Format = ""
			parameters[fmt.Sprintf("%s-%s", in, p.Name)] = opp
		}
	}
	// add Definitions parameters
	out.Parameters = parameters
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

func (s *toSwagger) toReqParameters(ps spec.HTTPRequestNode, spe *spec.Spec) []openAPIParamter {
	// 添加启用的全局参数
	out := toParameterGlobal(spe.Globals.Parameters, true, ps.GlobalExcepts)
	for in, params := range ps.Parameters.Map() {
		switch in {
		case "header", "query", "path", "cookie":
			for _, v := range params {
				if v.Reference != nil {
					// 解开公共参数
					if id := toInt64(getRefName(*v.Reference)); id != 0 {
						v = spe.Definitions.Parameters.LookupID(id)
					}
				}
				newv := *v
				newv.Schema = s.convertJSONSchema(v.Schema)
				out = append(out, toParameter(&newv, in, "2.0"))
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
					Type:        v.Type.Value()[0],
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

func (s *toSwagger) convertJSONSchema(v *jsonschema.Schema) *jsonschema.Schema {
	if v == nil {
		return v
	}
	return toConvertJSONSchemaRef(v, "2.0", s.schemas)
}

func (s *toSwagger) parseResponse(in *spec.Spec, res spec.HTTPResponseDefine) map[string]any {
	if res.Reference != nil {
		ref := *res.Reference
		if strings.HasPrefix(ref, "#/definitions/responses/") {
			x := in.Definitions.Responses.LookupID(
				toInt64(getRefName(ref)),
			)
			if x != nil {
				name := fmt.Sprintf("%s-%d", x.Name, x.ID)
				return map[string]any{
					"$ref": "#/responses/" + name,
				}
			}
		}
		return nil
	}
	resp := map[string]any{
		"description":            res.Description,
		"x-apicat-response-name": res.Name,
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
		for k, v := range res.Content {
			resp["schema"] = s.convertJSONSchema(v.Schema)
			if v.Examples != nil {
				for _, v := range v.Examples {
					resp["examples"] = map[string]any{
						k: v,
					}
					break
				}
			}
			break
		}
	}
	return resp
}

func (s *toSwagger) toPathResponse(in *spec.Spec, resp []*spec.HTTPResponse) (map[string]any, []string) {
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

func (s *toSwagger) toPaths(in *spec.Spec) (map[string]map[string]swaggerPathItem, []tagObject) {
	out := make(map[string]map[string]swaggerPathItem)
	tags := make(map[string]struct{})
	for path, ops := range walkHttpCollection(in) {
		if path == "" {
			continue
		}
		for method, op := range ops {
			reslist, product := s.toPathResponse(in, op.Res.List)
			if len(reslist) == 0 {
				reslist["default"] = &spec.Schema{Description: "success"}
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
