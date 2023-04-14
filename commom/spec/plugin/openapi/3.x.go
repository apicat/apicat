package openapi

import (
	"strconv"
	"strings"

	"github.com/apicat/apicat/commom/spec"
	"github.com/apicat/apicat/commom/spec/jsonschema"
	"github.com/apicat/apicat/commom/spec/markdown"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type OpenAPI struct{}

func (o *OpenAPI) parseInfo(info *base.Info) *spec.Info {
	return &spec.Info{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}
}

func (o *OpenAPI) parseServers(servs []*v3.Server) []*spec.Server {
	srvs := make([]*spec.Server, len(servs))
	for k, v := range servs {
		srvs[k] = &spec.Server{
			URL:         v.URL,
			Description: v.Description,
		}
	}
	return srvs
}

func (o *OpenAPI) parseDefinetions(shs map[string]*base.SchemaProxy) spec.Schemas {
	si := 0
	defines := make(spec.Schemas, len(shs))
	for k, v := range shs {
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

func (o *OpenAPI) parseParameters(inp []*v3.Parameter) *spec.HTTPParameters {
	var rawparamter spec.HTTPParameters
	for _, v := range inp {
		var sp = &spec.Schema{
			Name:     v.Name,
			Required: v.Required,
		}
		sp.Schema = &jsonschema.Schema{}
		if v.Schema != nil {
			js, err := jsonSchemaConverter(v.Schema.Schema())
			if err != nil {
				panic(err)
			}
			sp.Schema = js
		}
		sp.Schema.Description = v.Description
		sp.Schema.Example = v.Example
		sp.Schema.Deprecated = v.Deprecated
		rawparamter.Add(v.In, sp)
	}
	return &rawparamter
}

func (o *OpenAPI) parseContent(mts map[string]*v3.MediaType) spec.HTTPBody {
	if mts == nil {
		return nil
	}
	content := make(spec.HTTPBody)
	for contentType, mt := range mts {
		sh := &spec.Schema{}
		if g := mt.Schema.GoLow(); g.IsSchemaReference() {
			ref := strings.ReplaceAll(
				g.GetSchemaReference(),
				"#/components/schemas",
				"#/definitions",
			)
			sh.Schema = &jsonschema.Schema{
				Reference: &ref,
			}
		} else {
			js, err := jsonSchemaConverter(mt.Schema.Schema())
			if err != nil {
				panic(err)
			}
			js.Example = mt.Example
			sh.Schema = js
		}
		content[contentType] = sh
	}
	return content
}

func (o *OpenAPI) parseeResoponse(responses map[string]*v3.Response) spec.HTTPResponsesNode {
	var outresponses spec.HTTPResponsesNode
	for code, res := range responses {
		c, _ := strconv.Atoi(code)
		resp := spec.HTTPResponse{
			Code:        c,
			Description: res.Description,
		}
		if res.Headers != nil {
			for k, v := range res.Headers {
				js, err := jsonSchemaConverter(v.Schema.Schema())
				if err != nil {
					panic(err)
				}
				js.Description = v.Description
				resp.Header = append(resp.Header, &spec.Schema{
					Name:   k,
					Schema: js,
				})
			}
		}
		resp.Content = o.parseContent(res.Content)
		outresponses.List = append(outresponses.List, resp)
	}
	return outresponses
}

func (o *OpenAPI) parseCollections(paths *v3.Paths) []*spec.CollectItem {
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
			var req spec.HTTPRequestNode
			req.Parameters = o.parseParameters(info.Parameters)
			if info.RequestBody != nil {
				req.Content = o.parseContent(info.RequestBody.Content)
			}
			content = append(content, spec.MuseCreateNodeProxy(spec.WarpHTTPNode(req)))
			// response
			if info.Responses != nil {
				res := o.parseeResoponse(info.Responses.Codes)
				content = append(content, spec.MuseCreateNodeProxy(spec.WarpHTTPNode(res)))
			}

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

///

type openapiSpec struct {
	Openapi    string                                `json:"openapi"`
	Info       *spec.Info                            `json:"info"`
	Servers    []*spec.Server                        `json:"servers"`
	Components map[string]any                        `json:"components"`
	Paths      map[string]map[string]openapiPathItem `json:"paths"`
	Tags       []tagObject                           `json:"tags,omitempty"`
}

func (o *OpenAPI) toBase(in *spec.Spec, ver string) *openapiSpec {
	return &openapiSpec{
		Openapi: ver,
		Info: &spec.Info{
			Title:       in.Info.Title,
			Description: in.Info.Description,
			Version:     in.Info.Version,
		},
		Servers:    in.Servers,
		Components: o.toComponents(ver, in),
	}
}

type openapiRequestbody struct {
	Content spec.HTTPBody `json:"content,omitempty"`
}
type openapiPathItem struct {
	Summary     string              `json:"summary"`
	Description string              `json:"description,omitempty"`
	OperationId string              `json:"operationId"`
	Tags        []string            `json:"tags,omitempty"`
	Parameters  []*openAPIParamter  `json:"parameters,omitempty"`
	RequestBody *openapiRequestbody `json:"requestBody,omitempty"`
	Responses   map[string]any      `json:"responses,omitempty"`
}

// 3.0/3.1使用的jsonschema标准不太一样 3.1偏标准
func (o *OpenAPI) convertJSONSchema(ver string, in *jsonschema.Schema) {
	if in == nil {
		return
	}
	if in.Reference != nil {
		*in.Reference = strings.ReplaceAll(
			*in.Reference, "#/definitions", "#/components/schemas",
		)
	} else if strings.HasPrefix(ver, "3.0") {
		if in.Items != nil {
			if !in.Items.IsBool() {
				in.Items.SetValue(&jsonschema.Schema{})
			}
		}
		if in.Properties != nil {
			for _, v := range in.Properties {
				o.convertJSONSchema(ver, v)
			}
		}
		if in.AdditionalProperties != nil {
			if !in.AdditionalProperties.IsBool() {
				o.convertJSONSchema(ver, in.AdditionalProperties.Value())
			}
		}
		in.Type.SetValue(in.Type.Value()[0])
	}
	if in.Type != nil {
		t := in.Type.Value()
		if len(t) > 0 && t[0] == "file" {
			// jsonschema 没有file
			in.Type.SetValue("array")
			in.Items = &jsonschema.ValueOrBoolean[*jsonschema.Schema]{}
		}
	}
}

func (o *OpenAPI) toReqParameters(ps *spec.HTTPParameters, ver string) []*openAPIParamter {
	if ps == nil {
		return nil
	}
	var out []*openAPIParamter
	for in, params := range ps.Map() {
		for _, param := range params {
			o.convertJSONSchema(ver, param.Schema)
			item := &openAPIParamter{
				Name:     param.Name,
				Required: param.Required,
				Schema:   param.Schema,
				In:       in,
			}
			out = append(out, item)
		}
	}
	return out
}

func (o *OpenAPI) toPaths(ver string, in *spec.Spec) (
	map[string]map[string]openapiPathItem, []tagObject) {
	var (
		out  = make(map[string]map[string]openapiPathItem)
		tags = make(map[string]struct{})
	)
	for path, ops := range walkHttpCollection(in) {
		if path == "" {
			continue
		}
		for method, op := range ops {
			item := openapiPathItem{
				Summary:     op.Title,
				Description: op.Description,
				OperationId: op.OperatorID,
				Tags:        op.Tags,
				Parameters:  o.toReqParameters(op.Req.Parameters, ver),
			}
			for _, v := range op.Tags {
				tags[v] = struct{}{}
			}
			for k, v := range op.Req.Content {
				o.convertJSONSchema(ver, v.Schema)
				sp := &spec.Schema{
					Schema:      v.Schema,
					Description: v.Description,
				}
				if item.RequestBody == nil {
					item.RequestBody = &openapiRequestbody{
						Content: make(spec.HTTPBody),
					}
				}
				item.RequestBody.Content[k] = sp
			}

			if len(op.Res.List) == 0 {
				if in.Common != nil && len(in.Common.Responses) > 0 {
					op.Res.List = in.Common.Responses
				}
			}
			for _, v := range op.Res.List {
				res := map[string]any{
					"description": v.Description,
				}
				if v.Content != nil {
					for _, xx := range v.Content {
						o.convertJSONSchema(ver, xx.Schema)
					}
					res["content"] = v.Content
				}
				if len(v.Header) > 0 {
					headers := make(map[string]any)
					for _, v := range v.Header {
						o.convertJSONSchema(ver, v.Schema)
						headers[v.Name] = map[string]any{
							"description": v.Description,
							"schema":      v.Schema,
						}
					}
					res["headers"] = headers
				}
				if item.Responses != nil {
					item.Responses[strconv.Itoa(v.Code)] = res
				} else {
					item.Responses = map[string]any{
						strconv.Itoa(v.Code): res,
					}
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

func (o *OpenAPI) toComponents(ver string, in *spec.Spec) map[string]any {
	schemas := make(map[string]jsonschema.Schema)
	for _, v := range in.Definitions {
		o.convertJSONSchema(ver, v.Schema)
		schemas[v.Name] = *v.Schema
	}
	return map[string]any{
		"schemas": schemas,
	}
}
