package openapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apicat/apicat/commom/spec"
	"github.com/apicat/apicat/commom/spec/jsonschema"
	"github.com/apicat/apicat/commom/spec/markdown"
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

func parseSwagger(document libopenapi.Document) (*spec.Spec, error) {
	model, errors := document.BuildV2Model()
	if len(errors) > 0 {
		return nil, fmt.Errorf("swagger version:%s parse faild", document.GetVersion())
	}
	sw := &Swagger{}
	return &spec.Spec{
		ApiCat:      "2.0",
		Info:        sw.parseInfo(model.Model.Info),
		Servers:     sw.parseServers(&model.Model),
		Definitions: sw.parseDefinetions(model.Model.Definitions),
		Collections: sw.parseCollections(&model.Model, model.Model.Paths),
	}, nil
}

func parseOpenAPI3(document libopenapi.Document) (*spec.Spec, error) {
	model, errors := document.BuildV3Model()
	if len(errors) > 0 {
		return nil, fmt.Errorf("openapi version:%s parse faild", document.GetVersion())
	}
	o := &OpenAPI{}
	return &spec.Spec{
		Info:        o.parseInfo(model.Model.Info),
		Servers:     o.parseServers(model.Model.Servers),
		Definitions: o.parseDefinetions(model.Model.Components.Schemas),
		Collections: o.parseCollections(model.Model.Paths),
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
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required"`
	Description string `json:"description,omitempty"`
	// in body
	Schema *jsonschema.Schema `json:"schema,omitempty"`
	// in not body
	Type    string `json:"type,omitempty"`
	Format  string `json:"format,omitempty"`
	Default any    `json:"default,omitempty"`
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
