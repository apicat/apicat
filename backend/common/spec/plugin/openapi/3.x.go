package openapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/apicat/apicat/backend/common/spec"
	"github.com/apicat/apicat/backend/common/spec/jsonschema"
	"github.com/apicat/apicat/backend/common/spec/markdown"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type fromOpenapi struct {
	schemaMapping     map[string]int64
	parametersMapping map[string]*spec.Schema
}

func (o *fromOpenapi) parseInfo(info *base.Info) *spec.Info {
	return &spec.Info{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}
}

func (o *fromOpenapi) parseServers(servs []*v3.Server) []*spec.Server {
	srvs := make([]*spec.Server, len(servs))
	for k, v := range servs {
		srvs[k] = &spec.Server{
			URL:         v.URL,
			Description: v.Description,
		}
	}
	return srvs
}

func (o *fromOpenapi) parseParametersDefine(comp *v3.Components) (res spec.HTTPParameters) {
	res.Fill()
	if comp == nil {
		return res
	}
	repeat := map[string]int{}
	for _, v := range comp.Parameters {
		repeat[v.Name]++
	}
	for k, v := range comp.Parameters {
		if repeat[v.Name] > 1 {
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
		o.parametersMapping[k] = sp

		res.Add(v.In, sp)
	}
	return res
}

func (o *fromOpenapi) parseDefinetions(comp *v3.Components) spec.Definitions {
	if comp == nil {
		return spec.Definitions{
			Schemas:    make(spec.Schemas, 0),
			Parameters: spec.HTTPParameters{},
			Responses:  make(spec.HTTPResponseDefines, 0),
		}
	}
	o.schemaMapping = map[string]int64{}
	schemas := make(spec.Schemas, 0)
	for k, v := range comp.Schemas {
		js, err := jsonSchemaConverter(v)
		if err != nil {
			panic(err)
		}
		schemas = append(schemas, &spec.Schema{
			ID:          stringToUnid(k),
			Name:        k,
			Description: js.Description,
			Schema:      js,
		})
	}
	rets := []spec.HTTPResponseDefine{}
	for k, v := range comp.Responses {
		id := stringToUnid(k)
		def := spec.HTTPResponseDefine{
			Header: make(spec.Schemas, 0),
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
				def.Header = append(def.Header, &spec.Schema{
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
	return spec.Definitions{
		Schemas:    schemas,
		Responses:  rets,
		Parameters: o.parseParametersDefine(comp),
	}
}

func (o *fromOpenapi) parseParameters(inp []*v3.Parameter) spec.HTTPParameters {
	var rawparamter spec.HTTPParameters
	rawparamter.Fill()
	for _, v := range inp {
		if g := v.GoLow(); g.IsReference() {
			// if this parameter is a global parameter
			if isGlobalParameter(g.Reference.Reference) {
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

func (o *fromOpenapi) parseGlobal(com *v3.Components) (res spec.Global) {
	res.Parameters.Fill()
	if com == nil {
		return res
	}
	inp := com.Extensions
	if inp == nil {
		return res
	}
	global, ok := inp["x-apicat-globals"]
	if !ok {
		return res
	}
	var rawparamter spec.HTTPParameters
	rawparamter.Fill()

	for k, v := range global.(map[string]any) {
		nb, err := json.Marshal(v)
		if err != nil {
			continue
		}

		s := &spec.Schema{}
		json.Unmarshal(nb, s)
		in := strings.Index(k, "-")
		if in == -1 {
			continue
		}
		rawparamter.Add(k[:in], s)
	}
	res.Parameters = rawparamter
	return res
}

func (o *fromOpenapi) parseContent(mts map[string]*v3.MediaType) spec.HTTPBody {
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
		// sh.Example = mt.Example
		if len(mt.Examples) > 0 {
			sh.Examples = make(map[string]*spec.Example)
			for k, v := range mt.Examples {
				sh.Examples[k] = &spec.Example{
					Summary: v.Summary,
					Value:   v.Value,
				}
			}
		}
		sh.Schema = js
		content[contentType] = sh
	}
	return content
}

func (o *fromOpenapi) parseeResoponse(responses map[string]*v3.Response) spec.HTTPResponsesNode {
	var outresponses spec.HTTPResponsesNode
	for code, res := range responses {
		c, _ := strconv.Atoi(code)
		resp := spec.HTTPResponse{
			Code: c,
		}
		s := res.GoLow().Reference.Reference
		if s != "" {
			resp.Reference = &s
			outresponses.List = append(outresponses.List, &resp)
			continue
		}
		resp.Name = res.Description
		resp.Description = res.Description
		resp.Header = make(spec.Schemas, 0)
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
		outresponses.List = append(outresponses.List, &resp)
	}
	return outresponses
}

func (o *fromOpenapi) parseCollections(paths *v3.Paths) []*spec.CollectItem {
	collects := make([]*spec.CollectItem, 0)
	if paths == nil {
		return collects
	}
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
			req.FillGlobalExcepts()
			req.InitContent()
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
	Servers    []*spec.Server                        `json:"servers,omitempty"`
	Components map[string]any                        `json:"components,omitempty"`
	Paths      map[string]map[string]openapiPathItem `json:"paths"`
	Tags       []tagObject                           `json:"tags,omitempty"`
}

type toOpenapi struct {
	// get all schema that type is not category,in toComponents func, used to toConvertJSONSchemaRef func, get schema's name by schema's id
	schemaMapping map[int64]string
}

func (o *toOpenapi) toBase(in *spec.Spec, ver string) *openapiSpec {
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
func (o *toOpenapi) convertJSONSchema(ver string, in *jsonschema.Schema) *jsonschema.Schema {
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
			p.Items = &jsonschema.ValueOrBoolean[*jsonschema.Schema]{}
			p.Items.SetValue(&jsonschema.Schema{})
		}
	}
	return p
}

func (o *toOpenapi) toReqParameters(spe *spec.Spec, ps spec.HTTPRequestNode, ver string) []openAPIParamter {
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
				Name:        p.Name,
				Required:    p.Required,
				Description: p.Description,
				Example:     p.Example,
				Schema:      o.convertJSONSchema(ver, p.Schema),
				In:          in,
			}
			out = append(out, item)
		}
	}
	return out
}

func (o *toOpenapi) toPaths(ver string, in *spec.Spec) (
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
				if k == "none" {
					continue
				}
				sp := &spec.Schema{
					Schema:      o.convertJSONSchema(ver, v.Schema),
					Description: v.Description,
					Examples:    v.Examples,
				}
				if sp.Schema.Example != nil {
					sp.Example = sp.Schema.Example
				}
				if item.RequestBody == nil {
					item.RequestBody = &openapiRequestbody{
						Content: make(spec.HTTPBody),
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

func (o *toOpenapi) toResponse(in *spec.Spec, def spec.HTTPResponseDefine, ver string) map[string]any {
	res := map[string]any{}
	v := def
	if def.Reference != nil {
		if strings.HasPrefix(*v.Reference, "#/definitions/responses/") {
			if x := in.Definitions.Responses.LookupID(
				toInt64(getRefName(*v.Reference)),
			); x != nil {
				name_id := fmt.Sprintf("%s-%d", x.Name, x.ID)
				return map[string]any{
					"$ref": "#/components/responses/" + name_id,
				}
			}
		}
	}
	if v.Content != nil {
		c := make(map[string]*spec.Schema)
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
	res["x-apicat-category"] = v.Category
	res["name"] = v.Name
	return res
}

func (o *toOpenapi) toComponents(ver string, in *spec.Spec) map[string]any {
	schemas := make(map[string]jsonschema.Schema)
	o.schemaMapping = map[int64]string{}
	for _, v := range in.Definitions.Schemas {
		// if type is category, it's not have schema, need to range it's items
		if v.Type == string(spec.ContentItemTypeDir) {
			ss := v.ItemsTreeToList()
			for _, s := range ss {
				name_id := fmt.Sprintf("%s-%d", s.Name, s.ID)
				schemas[name_id] = *o.convertJSONSchema(ver, s.Schema)
				o.schemaMapping[s.ID] = s.Name
			}
		} else {
			name_id := fmt.Sprintf("%s-%d", v.Name, v.ID)
			schemas[name_id] = *o.convertJSONSchema(ver, v.Schema)
			o.schemaMapping[v.ID] = v.Name
		}
	}
	respons := make(map[string]any)
	for _, v := range in.Definitions.Responses {
		if v.Type == string(spec.ContentItemTypeDir) {
			rs := v.ItemsTreeToList()
			for _, resp := range rs {
				name_id := fmt.Sprintf("%s-%d", resp.Name, resp.ID)
				respons[name_id] = o.toResponse(in, resp, ver)
			}
		} else {
			name_id := fmt.Sprintf("%s-%d", v.Name, v.ID)
			respons[name_id] = o.toResponse(in, v, ver)
		}
	}

	globalParam := in.Globals.Parameters
	m := globalParam.Map()
	globals := make(map[string]openAPIParamter)
	for in, ps := range m {
		for _, p := range ps {
			globals[fmt.Sprintf("%s-%s", in, p.Name)] = toParameter(p, in)
		}
	}

	definitionParameters := in.Definitions.Parameters
	dp := definitionParameters.Map()
	parameters := make(map[string]openAPIParamter)
	for in, ps := range dp {
		for _, p := range ps {
			opp := toParameter(p, in)
			// remove format from parameter, because it's not support in openapi3.components.parameters.item
			opp.Format = ""
			parameters[fmt.Sprintf("%s-%s", in, p.Name)] = opp
		}
	}

	return map[string]any{
		"schemas":          schemas,
		"responses":        respons,
		"parameters":       parameters,
		"x-apicat-globals": globals,
	}
}
