package openapi

import (
	"fmt"
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

func (o *OpenAPI) parseParamtersCommon(comp *v3.Components) (spec.Schemas, map[string]string) {
	mapping := make(map[string]string)
	ps := make(spec.Schemas, 0)
	if comp == nil {
		return ps, mapping
	}
	repeat := map[string]struct {
		SrcKey string
		Count  int
	}{}
	for key, v := range comp.Parameters {
		x := repeat[v.Name]
		x.SrcKey = key
		x.Count++
		repeat[v.Name] = x
	}
	for _, v := range comp.Parameters {
		if repeat[v.Name].Count != 1 {
			continue
		}
		var sp = &spec.Schema{
			Name:     v.Name,
			Required: v.Required,
		}
		sp.Schema = &jsonschema.Schema{}
		if v.Schema != nil {
			js, err := jsonSchemaConverter(v.Schema)
			if err != nil {
				panic(err)
			}
			sp.Schema = js
		}
		sp.Schema.Description = v.Description
		sp.Schema.Example = v.Example
		sp.Schema.Deprecated = v.Deprecated
		ps = append(ps, sp)
	}
	for k, v := range repeat {
		if v.Count == 1 {
			mapping[v.SrcKey] = k
		}
	}
	return ps, mapping
}

func (o *OpenAPI) parseDefinetions(comp *v3.Components) spec.Definitions {
	if comp == nil {
		return spec.Definitions{}
	}
	si := 0
	schemas := make(spec.Schemas, len(comp.Schemas))
	for k, v := range comp.Schemas {
		js, err := jsonSchemaConverter(v)
		if err != nil {
			panic(err)
		}
		schemas[si] = &spec.Schema{
			Name:        k,
			Description: k,
			Schema:      js,
		}
		si++
	}
	si = 0
	rets := make([]spec.HTTPResponseDefine, len(comp.Responses))
	for k, v := range comp.Responses {
		if v.Content == nil {
			continue
		}
		var def spec.HTTPResponseDefine
		def.Name = k
		if v.Headers != nil {
			for k, v := range v.Headers {
				js, err := jsonSchemaConverter(v.Schema)
				if err != nil {
					panic(err)
				}
				js.Description = v.Description
				def.Header = append(def.Header, &spec.Schema{
					Name:   k,
					Schema: js,
				})
			}
		}
		def.Content = o.parseContent(v.Content)
		rets[si] = def
		si++
	}
	return spec.Definitions{
		Schemas:   schemas,
		Responses: rets,
	}
}

func (o *OpenAPI) parseParameters(inp []*v3.Parameter, commParamters map[string]string) spec.HTTPParameters {
	var rawparamter spec.HTTPParameters
	for _, v := range inp {
		if g := v.GoLow(); g.IsReference() && commParamters != nil {
			x, ok := commParamters[getRefName(g.GetReference())]
			if ok {
				r := "#/commons/parameters/" + x
				rawparamter.Add(v.In, &spec.Schema{
					Reference: &r,
				})
				continue
			}
		}
		var sp = &spec.Schema{
			Name:     v.Name,
			Required: v.Required,
		}

		sp.Schema = &jsonschema.Schema{}
		if v.Schema != nil {
			js, err := jsonSchemaConverter(v.Schema)
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
	return rawparamter
}

func (o *OpenAPI) parseContent(mts map[string]*v3.MediaType) spec.HTTPBody {
	if mts == nil {
		return nil
	}
	content := make(spec.HTTPBody)
	for contentType, mt := range mts {
		sh := &spec.Schema{}
		js, err := jsonSchemaConverter(mt.Schema)
		if err != nil {
			panic(err)
		}
		js.Example = mt.Example
		sh.Schema = js
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
				js, err := jsonSchemaConverter(v.Schema)
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

func (o *OpenAPI) parseCollections(paths *v3.Paths, commParamters map[string]string) []*spec.CollectItem {
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
			req.Parameters = o.parseParameters(info.Parameters, commParamters)
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
	Servers    []*spec.Server                        `json:"servers,omitempty"`
	Components map[string]any                        `json:"components,omitempty"`
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
	Parameters  []openAPIParamter   `json:"parameters,omitempty"`
	RequestBody *openapiRequestbody `json:"requestBody,omitempty"`
	Responses   map[string]any      `json:"responses,omitempty"`
}

// 3.0/3.1使用的jsonschema标准不太一样 3.1偏标准
func (o *OpenAPI) convertJSONSchema(ver string, in *jsonschema.Schema) {
	if in == nil {
		return
	}
	toConvertJSONSchemaRef(in, ver)
	if in.Reference == nil && strings.HasPrefix(ver, "3.0") {
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
			in.Items.SetValue(&jsonschema.Schema{})
		}
	}
}

func (o *OpenAPI) toReqParameters(spe *spec.Spec, ps spec.HTTPRequestNode, ver string) []openAPIParamter {
	// var out []openAPIParamter
	out := toParameterGlobal(spe.Globals.Parameters, false, ps.GlobalExcepts)
	for in, params := range ps.Parameters.Map() {
		for _, param := range params {
			if param.Reference != nil {
				param = spe.Common.Parameters.Lookup(getRefName(*param.Reference))
			}
			o.convertJSONSchema(ver, param.Schema)
			item := openAPIParamter{
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
				Parameters:  o.toReqParameters(in, op.Req, ver),
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
			for _, v := range op.Res.List {
				code, res := o.toResponse(in, v, ver)
				if item.Responses != nil {
					item.Responses[code] = res
				} else {
					item.Responses = map[string]any{
						code: res,
					}
				}
			}
			if len(op.Res.List) == 0 {
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

func (o *OpenAPI) toResponse(in *spec.Spec, def spec.HTTPResponse, ver string) (string, map[string]any) {
	res := map[string]any{}
	v := def
	if def.Reference != nil {
		switch {
		case strings.HasPrefix(*v.Reference, "#/definitions/responses/"):
			// openapi3 并没有define response
			x := in.Definitions.Responses.Lookup(
				getRefName(*v.Reference),
			)
			v.HTTPResponseDefine = *x
			v.Reference = nil
		case strings.HasPrefix(*v.Reference, "#/commons/responses/"):
			x := in.Common.Responses.Lookup(
				getRefName(*v.Reference),
			)
			v = *x
		default:
			panic(fmt.Sprintf("error response ref:%s", *def.Reference))
		}
	}
	if v.Content != nil {
		for _, xx := range v.Content {
			o.convertJSONSchema(ver, xx.Schema)
		}
		res["content"] = v.Content
	}
	if len(v.Header) > 0 {
		headers := make(map[string]any)
		for _, h := range v.Header {
			o.convertJSONSchema(ver, h.Schema)
			headers[h.Name] = map[string]any{
				"description": h.Description,
				"schema":      h.Schema,
			}
		}
		res["headers"] = headers
	}
	res["description"] = v.Description
	return strconv.Itoa(v.Code), res
}

func (o *OpenAPI) toComponents(ver string, in *spec.Spec) map[string]any {
	schemas := make(map[string]jsonschema.Schema)
	for _, v := range in.Definitions.Schemas {
		o.convertJSONSchema(ver, v.Schema)
		schemas[v.Name] = *v.Schema
	}
	return map[string]any{
		"schemas": schemas,
	}
}
