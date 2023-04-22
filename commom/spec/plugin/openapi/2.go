package openapi

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/apicat/apicat/commom/spec"
	"github.com/apicat/apicat/commom/spec/jsonschema"
	"github.com/apicat/apicat/commom/spec/markdown"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v2 "github.com/pb33f/libopenapi/datamodel/high/v2"
)

type Swagger struct{}

func (s *Swagger) parseInfo(info *base.Info) *spec.Info {
	return &spec.Info{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}
}

func (s *Swagger) parseServers(in *v2.Swagger) []*spec.Server {
	srvs := make([]*spec.Server, len(in.Schemes))
	for k, v := range in.Schemes {
		srvs[k] = &spec.Server{
			URL:         fmt.Sprintf("%s://%s%s", v, in.Host, in.BasePath),
			Description: v,
		}
	}
	return srvs
}

func (s *Swagger) parseDefinetions(defs *v2.Definitions) spec.Schemas {
	if defs == nil {
		return make(spec.Schemas, 0)
	}
	si := 0
	defines := make(spec.Schemas, len(defs.Definitions))
	for k, v := range defs.Definitions {
		js, err := jsonSchemaConverter(v.Schema())
		if err != nil {
			panic(err)
		}
		defines[si] = &spec.Schema{
			Name:        k,
			Description: k,
			Schema:      js,
		}
		si++
	}
	return defines
}

func (s *Swagger) parseParamtersCommon(in *v2.Swagger) (spec.Schemas, map[string]string) {
	mapping := make(map[string]string)
	if in.Parameters == nil {
		return nil, mapping
	}
	ps := make(spec.Schemas, 0)
	repeat := map[string]struct {
		SrcKey string
		Count  int
	}{}
	for key, v := range in.Parameters.Definitions {
		x := repeat[v.Name]
		x.SrcKey = key
		x.Count++
		repeat[v.Name] = x
	}
	for _, v := range in.Parameters.Definitions {
		if repeat[v.Name].Count != 1 {
			continue
		}
		ps = append(ps, &spec.Schema{
			Name:        v.Name,
			Description: v.Description,
			Required:    *v.Required,
			Schema: &jsonschema.Schema{
				Type:   jsonschema.CreateSliceOrOne(v.Type),
				Format: v.Format,
			},
		})
	}
	for k, v := range repeat {
		if v.Count == 1 {
			mapping[v.SrcKey] = k
		}
	}
	return ps, mapping
}

// 主要处理$ref引用问题
func (s *Swagger) parseContent(b *base.SchemaProxy) *jsonschema.Schema {
	if g := b.GoLow(); g != nil {
		if g.IsSchemaReference() {
			ref := g.GetSchemaReference()
			if strings.HasPrefix(ref, "#/definitions/") {
				r := strings.ReplaceAll(ref, "#/definitions/", "#/definitions/schemas")
				return &jsonschema.Schema{
					Reference: &r,
				}
			}
		}
	}
	js, err := jsonSchemaConverter(b.Schema())
	if err != nil {
		panic(err)
	}
	return js
}

func (s *Swagger) parseRequest(in *v2.Swagger, info *v2.Operation, paramtags map[string]string) spec.HTTPRequestNode {
	// paramters := &spec.HttpParameters{}
	request := spec.HTTPRequestNode{
		Content: make(spec.HTTPBody),
	}
	var body *jsonschema.Schema
	// 有效载荷application/x-www-form-urlencoded和multipart/form-data请求是通过使用form参数来描述，而不是body参数。
	formData := &jsonschema.Schema{
		Type:       jsonschema.CreateSliceOrOne("object"),
		Properties: make(map[string]*jsonschema.Schema),
	}

	for _, v := range info.Parameters {
		// 这里引用 #/parameters 暂时无法获取
		// 直接展开
		switch v.In {
		case "query", "header", "path":
			request.Parameters.Add(v.In,
				&spec.Schema{
					Name:        v.Name,
					Description: v.Description,
					Required:    *v.Required,
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
			if *v.Required {
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
		if strings.Index(v, "form") != -1 {
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
func (s *Swagger) parseResponsesDefine(in *v2.Swagger) []spec.HTTPResponseDefine {
	if in.Responses == nil {
		return nil
	}
	list := make([]spec.HTTPResponseDefine, 0)
	for key, res := range in.Responses.Definitions {
		var header []*spec.Schema
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
			} else {
				for _, v := range in.Produces {
					content[v] = sh
				}
			}
		}
		list = append(list, spec.HTTPResponseDefine{
			Name:    key,
			Header:  header,
			Content: content,
		})
	}
	return nil
}

func (s *Swagger) parseResponse(info *v2.Operation) *spec.HTTPResponsesNode {
	if info.Responses == nil {
		return nil
	}
	var outresponses spec.HTTPResponsesNode
	if info.Responses.Default != nil {
		// 我们没有default
		// todo
	}
	for code, res := range info.Responses.Codes {
		// res github.com/pb33f/libopenapi 不支持response ref 所以无法获取
		// 这里的common无法转换
		c, err := strconv.Atoi(code)
		if err != nil {
			continue
		}
		resp := spec.HTTPResponse{
			Code:        c,
			Description: res.Description,
		}
		resp.Content = make(spec.HTTPBody)
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
				resp.Content[v] = &spec.Schema{
					Schema: js,
				}
			}
		}
		outresponses.List = append(outresponses.List, resp)
	}
	return &outresponses
}

func (s *Swagger) parseCollections(in *v2.Swagger, paths *v2.Paths, paramtags map[string]string) []*spec.CollectItem {
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
			req := spec.WarpHTTPNode(s.parseRequest(in, info, paramtags))
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
	Swagger     string                                `json:"swagger"`
	Info        *spec.Info                            `json:"info"`
	Tags        []tagObject                           `json:"tags,omitempty"`
	Host        string                                `json:"host,omitempty"`
	BasePath    string                                `json:"basePath"`
	Schemas     []string                              `json:"schemas,omitempty"`
	Definitions map[string]jsonschema.Schema          `json:"definitions"`
	Parameters  map[string]openAPIParamter            `json:"parameters,omitempty"`
	Responses   map[string]any                        `json:"responses,omitempty"`
	Paths       map[string]map[string]swaggerPathItem `json:"paths"`
}

func (s *Swagger) toBase(in *spec.Spec) *swaggerSpec {
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
	for _, v := range in.Definitions.Schemas {
		out.Definitions[v.Name] = *v.Schema
	}

	globalParam := in.Globals.Parameters
	m := globalParam.Map()
	out.Parameters = make(map[string]openAPIParamter)
	for in, ps := range m {
		for _, p := range ps {
			out.Parameters["global"+strings.Title(p.Name)] = toParameter(p, in)
		}
	}

	if out.BasePath == "" {
		out.BasePath = "/"
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

// func (s *Swagger) toParameter(p *spec.Schema, in string) openAPIParamter {
// 	tp := "string"
// 	if n := len(p.Schema.Type.Value()); n > 0 {
// 		tp = p.Schema.Type.Value()[0]
// 	}
// 	return openAPIParamter{
// 		In:       in,
// 		Type:     tp,
// 		Name:     p.Name,
// 		Required: p.Required,
// 		Format:   p.Schema.Format,
// 		Default:  p.Schema.Default,
// 	}
// }

func (s *Swagger) toReqParameters(ps spec.HTTPRequestNode, spe *spec.Spec) []openAPIParamter {
	out := toParameterGlobal(spe.Globals.Parameters, ps.GlobalExcepts)
	for in, params := range ps.Parameters.Map() {
		switch in {
		case "header", "query", "path":
			for _, v := range params {
				if v.Reference != nil {
					if !strings.HasPrefix(*v.Reference, "#/commons/parameters/") {
						continue
					}
					p := spe.Common.Parameters.Lookup(
						getRefName(*v.Reference),
					)
					if p != nil {
						out = append(out, toParameter(p, in))
					}
				} else {
					out = append(out, toParameter(v, in))
				}
			}
		}
	}
	if ps.Content != nil {
		var hasBody bool
		for ct, c := range ps.Content {
			if strings.Index(ct, "form") != -1 {
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
						Schema:      v,
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
					Schema:      c.Schema,
					In:          "body",
					Required:    true,
				})
				hasBody = true
			}
		}
	}
	return out
}

func (s *Swagger) toPathResponse(in *spec.Spec, resp []spec.HTTPResponse) (map[string]any, []string) {
	product := map[string]struct{}{}
	reslist := make(map[string]any)
	for _, r := range resp {
		v := r
		if v.Reference != nil {
			switch {
			case strings.HasPrefix(*v.Reference, "#/definitions/responses/"):
				x := in.Definitions.Responses.Lookup(
					getRefName(*v.Reference),
				)
				if x != nil {
					v.HTTPResponseDefine = *x
					v.Reference = nil
				}
			case strings.HasPrefix(*v.Reference, "#/commons/responses/"):
				x := in.Common.Responses.Lookup(
					getRefName(*v.Reference),
				)
				v = *x
			default:
				// not support
				continue
			}
		}
		var res *spec.Schema
		for k := range v.Content {
			if _, ok := product[k]; !ok {
				product[k] = struct{}{}
			}
			if res != nil {
				break
			}
			res = &spec.Schema{
				Schema:      v.Content[k].Schema,
				Description: v.Description,
			}
		}
		if res == nil {
			res = &spec.Schema{Description: v.Description}
		}
		reslist[strconv.Itoa(v.Code)] = res
	}
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

func (s *Swagger) toPaths(in *spec.Spec) (map[string]map[string]swaggerPathItem, []tagObject) {
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
