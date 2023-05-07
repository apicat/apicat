package openapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apicat/apicat/common/spec"
	"github.com/apicat/apicat/common/spec/jsonschema"
	"github.com/apicat/apicat/common/spec/markdown"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/utils"
)

// Decode 将openapi解码为spec对象
// 支持swagger，openapi3/3.1
func Decode(data []byte) (out *spec.Spec, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	docment, docerr := libopenapi.NewDocument(data)
	if err != nil {
		err = docerr
		return
	}
	t := docment.GetSpecInfo().SpecType
	switch t {
	case utils.OpenApi2:
		out, err = parseSwagger(docment)
	case utils.OpenApi3:
		out, err = parseOpenAPI3(docment)
	default:
		err = fmt.Errorf("not support %s", t)
	}
	return
}

const defaultSwaggerConsumerProduce = "application/json"

func parseSwagger(document libopenapi.Document) (*spec.Spec, error) {
	model, errors := document.BuildV2Model()
	if len(errors) > 0 {
		return nil, fmt.Errorf("swagger version:%s parse faild", document.GetVersion())
	}
	sw := &Swagger{}
	schemas := sw.parseDefinetions(model.Model.Definitions)
	responseDefinitions := sw.parseResponsesDefine(&model.Model)
	paramters, paramemapping := sw.parseParamtersCommon(&model.Model)
	return &spec.Spec{
		ApiCat:      "2.0.1",
		Info:        sw.parseInfo(model.Model.Info),
		Servers:     sw.parseServers(&model.Model),
		Collections: sw.parseCollections(&model.Model, model.Model.Paths, paramemapping),
		Definitions: spec.Definitions{Schemas: schemas, Responses: responseDefinitions},
		Common:      spec.Common{Parameters: paramters},
	}, nil
}

func parseOpenAPI3(document libopenapi.Document) (*spec.Spec, error) {
	model, errors := document.BuildV3Model()
	if len(errors) > 0 {
		return nil, fmt.Errorf("openapi version:%s parse faild", document.GetVersion())
	}
	o := &OpenAPI{}
	paramters, paramemapping := o.parseParamtersCommon(model.Model.Components)
	return &spec.Spec{
		ApiCat:      "2.0.1",
		Info:        o.parseInfo(model.Model.Info),
		Servers:     o.parseServers(model.Model.Servers),
		Definitions: o.parseDefinetions(model.Model.Components),
		Collections: o.parseCollections(model.Model.Paths, paramemapping),
		Common: spec.Common{
			Parameters: paramters,
		},
	}, nil
}

// Encode 将spec编码为openapi协议
// version 2.0/3.0.0/3.1.0
func Encode(in *spec.Spec, version string) ([]byte, error) {
	switch version {
	case "2.0":
		sw := &Swagger{}
		sp := sw.toBase(in)
		paths, tags := sw.toPaths(in)
		sp.Paths = paths
		sp.Tags = tags
		return json.MarshalIndent(sp, "", "  ")
	default:
		if strings.HasPrefix(version, "3.") && len(strings.Split(version, ".")) == 3 {
			op := &OpenAPI{}
			sp := op.toBase(in, version)
			paths, tag := op.toPaths(version, in)
			sp.Paths = paths
			sp.Tags = tag
			return json.MarshalIndent(sp, "", "  ")
		}
	}
	return nil, fmt.Errorf("openapi %s not support", version)
}

// swagger/open3.x
type openAPIParamter struct {
	Name        string `json:"name,omitempty"`
	In          string `json:"in,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
	// in body
	Schema *jsonschema.Schema `json:"schema,omitempty"`
	// in not body
	Type      string  `json:"type,omitempty"`
	Format    string  `json:"format,omitempty"`
	Default   any     `json:"default,omitempty"`
	Reference *string `json:"$ref,omitempty"`
}

func toParameter(p *spec.Schema, in string) openAPIParamter {
	tp := "string"
	if n := len(p.Schema.Type.Value()); n > 0 {
		tp = p.Schema.Type.Value()[0]
	}
	return openAPIParamter{
		In:       in,
		Type:     tp,
		Name:     p.Name,
		Required: p.Required,
		Format:   p.Schema.Format,
		Default:  p.Schema.Default,
	}
}

func toParameterGlobal(globalsParmaters spec.HTTPParameters, isSwagger bool, skip map[string][]string) []openAPIParamter {
	var outs []openAPIParamter
	skips := make(map[string]bool)
	if skip != nil {
		for k, v := range skip {
			for _, x := range v {
				skips[k+"|"+x] = true
			}
		}
	}
	for in, ps := range globalsParmaters.Map() {
		for _, v := range ps {
			if skips[in+"|"+v.Name] {
				continue
			}
			ref := fmt.Sprintf("global-%s-%s", in, v.Name)
			if isSwagger {
				ref = "#/parameters/" + ref
			} else {
				ref = "#/components/parameters/" + ref
			}
			outs = append(outs, openAPIParamter{
				Reference: &ref,
			})
		}
	}
	return outs
}

type tagObject struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type specPathItem struct {
	Title       string
	Description string
	OperatorID  string
	Tags        []string
	Req         spec.HTTPRequestNode
	Res         spec.HTTPResponsesNode
}

func walkHttpCollection(doc *spec.Spec) map[string]map[string]specPathItem {
	paths := make(map[string]map[string]specPathItem)
	doc.WalkCollections(
		func(v *spec.CollectItem) bool {
			if v.Type != spec.ContentItemTypeHttp {
				return true
			}
			var (
				info    spec.HTTPURLNode
				docRoot spec.Document
			)
			item := specPathItem{
				Title:      v.Title,
				OperatorID: fmt.Sprintf("%d", v.ID),
				Tags:       v.Tags,
			}
			for _, n := range v.Content {
				switch nx := n.Node.(type) {
				case *spec.HTTPNode[spec.HTTPURLNode]:
					info = nx.Attrs
				case *spec.HTTPNode[spec.HTTPRequestNode]:
					item.Req = nx.Attrs
				case *spec.HTTPNode[spec.HTTPResponsesNode]:
					item.Res = nx.Attrs
				case *spec.DocNode:
					docRoot.Items = append(docRoot.Items, nx)
				}
			}
			if len(docRoot.Items) > 0 {
				if raw, err := markdown.ToMarkdown(&docRoot); err == nil {
					item.Description = string(raw)
				}
			}
			if _, ok := paths[info.Path]; !ok {
				paths[info.Path] = map[string]specPathItem{
					info.Method: item,
				}
			} else {
				paths[info.Path][info.Method] = item
			}
			return true
		},
	)
	return paths
}

// 将jsonschema 转为对应的 openaapi版本 主要是引用
func toConvertJSONSchemaRef(v *jsonschema.Schema, ver string) {
	if v.Reference != nil {
		ref := getRefName(*v.Reference)
		if ver[0] == '2' {
			ref = fmt.Sprintf("#/definitions/%s", ref)
		} else {
			ref = fmt.Sprintf("#/components/schemas/%s", ref)
		}
		v.Reference = &ref
		return
	}
	if v.Properties != nil {
		for _, v := range v.Properties {
			toConvertJSONSchemaRef(v, ver)
		}
	}
	if v.Items != nil {
		if !v.Items.IsBool() {
			toConvertJSONSchemaRef(v.Items.Value(), ver)
		}
	}
	if v.AdditionalProperties != nil {
		if !v.AdditionalProperties.IsBool() {
			toConvertJSONSchemaRef(v.AdditionalProperties.Value(), ver)
		}
	}
}
