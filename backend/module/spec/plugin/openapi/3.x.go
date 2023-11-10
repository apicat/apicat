package openapi

import (
	"fmt"
	spec2 "github.com/apicat/apicat/backend/module/spec"
	jsonschema2 "github.com/apicat/apicat/backend/module/spec/jsonschema"
	"github.com/apicat/apicat/backend/module/spec/markdown"
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type fromOpenapi struct {
	schemaMapping     map[string]int64
	parametersMapping map[string]int64
}

func (o *fromOpenapi) parseInfo(info *base.Info) *spec2.Info {
	return &spec2.Info{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}
}

func (o *fromOpenapi) parseServers(servs []*v3.Server) []*spec2.Server {
	srvs := make([]*spec2.Server, len(servs))
	for k, v := range servs {
		srvs[k] = &spec2.Server{
			URL:         v.URL,
			Description: v.Description,
		}
	}
	return srvs
}

func (o *fromOpenapi) parseParametersDefine(comp *v3.Components) spec2.Schemas {
	o.parametersMapping = map[string]int64{}
	ps := make(spec2.Schemas, 0)
	if comp == nil {
		return ps
	}
	repeat := map[string]int{}
	for _, v := range comp.Parameters {
		repeat[v.Name]++
	}
	for k, v := range comp.Parameters {
		if repeat[v.Name] > 1 {
			continue
		}
		id := stringToUnid(k)
		o.parametersMapping[k] = id
		var sp = &spec2.Schema{
			ID:       id,
			Name:     v.Name,
			Required: v.Required,
		}
		sp.Schema = &jsonschema2.Schema{}
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
	return ps
}

func (o *fromOpenapi) parseDefinetions(comp *v3.Components) spec2.Definitions {
	if comp == nil {
		return spec2.Definitions{
			Schemas:    make(spec2.Schemas, 0),
			Parameters: make(spec2.Schemas, 0),
			Responses:  make(spec2.HTTPResponseDefines, 0),
		}
	}
	o.schemaMapping = map[string]int64{}
	schemas := make(spec2.Schemas, 0)
	for k, v := range comp.Schemas {
		js, err := jsonSchemaConverter(v)
		if err != nil {
			panic(err)
		}
		schemas = append(schemas, &spec2.Schema{
			ID:          stringToUnid(k),
			Name:        k,
			Description: k,
			Schema:      js,
		})
	}
	rets := []spec2.HTTPResponseDefine{}
	for k, v := range comp.Responses {
		id := stringToUnid(k)
		def := spec2.HTTPResponseDefine{
			Header: make(spec2.Schemas, 0),
			Name:   k,
			ID:     id,
		}
		if v.Headers != nil {
			for k, v := range v.Headers {
				js, err := jsonSchemaConverter(v.Schema)
				if err != nil {
					panic(err)
				}
				js.Description = v.Description
				def.Header = append(def.Header, &spec2.Schema{
					Name:   k,
					Schema: js,
				})
			}
		}
		if v.Content != nil {
			def.Content = o.parseContent(v.Content)
		}
		rets = append(rets, def)
	}
	return spec2.Definitions{
		Schemas:    schemas,
		Responses:  rets,
		Parameters: o.parseParametersDefine(comp),
	}
}

func (o *fromOpenapi) parseParameters(inp []*v3.Parameter) spec2.HTTPParameters {
	var rawparamter spec2.HTTPParameters
	rawparamter.Fill()
	for _, v := range inp {
		if g := v.GoLow(); g.IsReference() {
			id, ok := o.parametersMapping[getRefName(g.GetReference())]
			if ok {
				r := fmt.Sprintf("#/definitions/parameters/%d", id)
				rawparamter.Add(v.In, &spec2.Schema{
					Reference: &r,
				})
				continue
			}
		}
		var sp = &spec2.Schema{
			Name:     v.Name,
			Required: v.Required,
		}
		sp.Schema = &jsonschema2.Schema{}
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

func (o *fromOpenapi) parseContent(mts map[string]*v3.MediaType) spec2.HTTPBody {
	if mts == nil {
		return nil
	}
	content := make(spec2.HTTPBody)
	for contentType, mt := range mts {
		sh := &spec2.Schema{}
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

func (o *fromOpenapi) parseeResoponse(responses map[string]*v3.Response) spec2.HTTPResponsesNode {
	var outresponses spec2.HTTPResponsesNode
	for code, res := range responses {
		c, _ := strconv.Atoi(code)
		resp := spec2.HTTPResponse{
			Code: c,
		}
		resp.Name = res.Description
		resp.Description = res.Description
		resp.Header = make(spec2.Schemas, 0)
		if res.Headers != nil {
			for k, v := range res.Headers {
				js, err := jsonSchemaConverter(v.Schema)
				if err != nil {
					panic(err)
				}
				js.Description = v.Description
				resp.Header = append(resp.Header, &spec2.Schema{
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

func (o *fromOpenapi) parseCollections(paths *v3.Paths) []*spec2.CollectItem {
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
			var req spec2.HTTPRequestNode
			req.Parameters = o.parseParameters(info.Parameters)
			if info.RequestBody != nil {
				req.Content = o.parseContent(info.RequestBody.Content)
			}
			content = append(content, spec2.MuseCreateNodeProxy(spec2.WarpHTTPNode(req)))
			// response
			if info.Responses != nil {
				res := o.parseeResoponse(info.Responses.Codes)
				content = append(content, spec2.MuseCreateNodeProxy(spec2.WarpHTTPNode(res)))
			}

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

///

type openapiSpec struct {
	Openapi    string                                `json:"openapi"`
	Info       *spec2.Info                           `json:"info"`
	Servers    []*spec2.Server                       `json:"servers,omitempty"`
	Components map[string]any                        `json:"components,omitempty"`
	Paths      map[string]map[string]openapiPathItem `json:"paths"`
	Tags       []tagObject                           `json:"tags,omitempty"`
}

type toOpenapi struct {
	schemaMapping map[int64]string
}

func (o *toOpenapi) toBase(in *spec2.Spec, ver string) *openapiSpec {
	return &openapiSpec{
		Openapi: ver,
		Info: &spec2.Info{
			Title:       in.Info.Title,
			Description: in.Info.Description,
			Version:     in.Info.Version,
		},
		Servers:    in.Servers,
		Components: o.toComponents(ver, in),
	}
}

type openapiRequestbody struct {
	Content spec2.HTTPBody `json:"content,omitempty"`
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
func (o *toOpenapi) convertJSONSchema(ver string, in *jsonschema2.Schema) *jsonschema2.Schema {
	if in == nil {
		return nil
	}
	p := toConvertJSONSchemaRef(in, ver, o.schemaMapping)
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
		t := p.Type.Value()
		if len(t) > 0 && t[0] == "file" {
			// jsonschema 没有file
			p.Type.SetValue("array")
			p.Items = &jsonschema2.ValueOrBoolean[*jsonschema2.Schema]{}
			p.Items.SetValue(&jsonschema2.Schema{})
		}
	}
	return p
}

func (o *toOpenapi) toReqParameters(spe *spec2.Spec, ps spec2.HTTPRequestNode, ver string) []openAPIParamter {
	// var out []openAPIParamter
	out := toParameterGlobal(spe.Globals.Parameters, false, ps.GlobalExcepts)
	for in, params := range ps.Parameters.Map() {
		for _, param := range params {
			p := *param
			if p.Reference != nil {
				if defp := spe.Definitions.Parameters.LookupID(toInt64(getRefName(*p.Reference))); defp != nil {
					p = *defp
				}
			}
			item := openAPIParamter{
				Name:     p.Name,
				Required: p.Required,
				Schema:   o.convertJSONSchema(ver, p.Schema),
				In:       in,
			}
			out = append(out, item)
		}
	}
	return out
}

func (o *toOpenapi) toPaths(ver string, in *spec2.Spec) (
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
				Responses:   make(map[string]any),
			}
			for _, v := range op.Tags {
				tags[v] = struct{}{}
			}
			for k, v := range op.Req.Content {
				sp := &spec2.Schema{
					Schema:      o.convertJSONSchema(ver, v.Schema),
					Description: v.Description,
				}
				if item.RequestBody == nil {
					item.RequestBody = &openapiRequestbody{
						Content: make(spec2.HTTPBody),
					}
				}
				item.RequestBody.Content[k] = sp
			}
			for _, v := range op.Res.List {
				res := o.toResponse(in, v.HTTPResponseDefine, ver)
				item.Responses[strconv.Itoa(v.Code)] = res
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

func (o *toOpenapi) toResponse(in *spec2.Spec, def spec2.HTTPResponseDefine, ver string) map[string]any {
	res := map[string]any{}
	v := def
	if def.Reference != nil {
		if strings.HasPrefix(*v.Reference, "#/definitions/responses/") {
			if x := in.Definitions.Responses.LookupID(
				toInt64(getRefName(*v.Reference)),
			); x != nil {
				return map[string]any{
					"$ref": "#/components/responses/" + x.Name,
				}
			}
		}
	}
	if v.Content != nil {
		c := make(map[string]*spec2.Schema)
		for k, xx := range v.Content {
			x := *xx
			x.Schema = o.convertJSONSchema(ver, x.Schema)
			c[k] = &x
		}
		res["content"] = c
	}
	if len(v.Header) > 0 {
		headers := make(map[string]any)
		for _, h := range v.Header {
			headers[h.Name] = map[string]any{
				"description": h.Description,
				"schema":      o.convertJSONSchema(ver, h.Schema),
			}
		}
		res["headers"] = headers
	}
	res["description"] = v.Description
	return res
}

func (o *toOpenapi) toComponents(ver string, in *spec2.Spec) map[string]any {
	schemas := make(map[string]jsonschema2.Schema)
	o.schemaMapping = map[int64]string{}
	for _, v := range in.Definitions.Schemas {
		o.schemaMapping[v.ID] = v.Name
	}
	for _, v := range in.Definitions.Schemas {
		s := o.convertJSONSchema(ver, v.Schema)
		schemas[v.Name] = *s
	}
	respons := make(map[string]any)
	for _, v := range in.Definitions.Responses {
		res := o.toResponse(in, v, ver)
		respons[v.Name] = res
	}

	globalParam := in.Globals.Parameters
	m := globalParam.Map()
	parameters := make(map[string]openAPIParamter)
	for in, ps := range m {
		for _, p := range ps {
			parameters[fmt.Sprintf("%s-%s", in, p.Name)] = toParameter(p, in)
		}
	}

	return map[string]any{
		"schemas":    schemas,
		"responses":  respons,
		"parameters": parameters,
	}
}
