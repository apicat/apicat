package postman

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apicat/apicat/backend/common/spec"
	"github.com/apicat/apicat/backend/common/spec/jsonschema"
)

func Import(data []byte) (*spec.Spec, error) {
	var pm Spec
	if err := json.Unmarshal(data, &pm); err != nil {
		return nil, err
	}

	p := &spec.Spec{
		ApiCat: "2.0.1",
		Info: &spec.Info{
			Title:       pm.Info.Name,
			Description: pm.Info.Description,
		},
		Servers: func() []*spec.Server {
			for _, v := range pm.Items {
				if v.Request != nil {
					return []*spec.Server{{
						URL: fmt.Sprintf("%s://%s",
							v.Request.Url.Protocol,
							strings.Join(v.Request.Url.Host, "/"),
						),
						Description: "default",
					}}
				}
			}
			return []*spec.Server{}
		}(),
		Globals: func() spec.Global {
			var parmts spec.HTTPParameters
			parmts.Fill()
			return spec.Global{
				Parameters: parmts,
			}
		}(),
		Definitions: spec.Definitions{
			Schemas:    make(spec.Schemas, 0),
			Parameters: make(spec.Schemas, 0),
			Responses:  make(spec.HTTPResponseDefines, 0),
		},
		Collections: walkCpllection(pm.Items, 1000),
	}
	return p, nil
}

func walkCpllection(items []Item, parentid int64) []*spec.CollectItem {
	cs := make([]*spec.CollectItem, 0)
	for i, v := range items {
		// http request
		id := parentid*1024 + int64(i) + 1
		if v.Request != nil {
			specItem := &spec.CollectItem{
				ID:       id,
				ParentID: parentid,
				Type:     spec.ContentItemTypeHttp,
				Title:    v.Name,
				Content:  convertContent(v),
			}
			cs = append(cs, specItem)
		}
		// dir
		if len(v.Items) > 0 {
			specItem := &spec.CollectItem{
				ID:       id,
				ParentID: parentid,
				Type:     spec.ContentItemTypeDir,
				Title:    v.Name,
				Items:    walkCpllection(v.Items, id),
			}
			cs = append(cs, specItem)
		}
	}
	return cs
}

func convertContent(item Item) []*spec.NodeProxy {
	req := spec.HTTPRequestNode{
		GlobalExcepts: make(map[string][]int64),
	}
	req.Parameters.Fill()
	for k, v := range item.Request.Url.Path {
		if !strings.HasPrefix(v, ":") {
			continue
		}
		for _, x := range item.Request.Url.Variables {
			if x.Key == v[1:] {
				item.Request.Url.Path[k] = "{" + x.Key + "}"
				req.Parameters.Path = append(req.Parameters.Path, x.toSchema())
				break
			}
		}
	}

	nodes := []*spec.NodeProxy{
		spec.MuseCreateNodeProxy(spec.WarpHTTPNode(spec.HTTPURLNode{
			Path:   "/" + strings.Join(item.Request.Url.Path, "/"),
			Method: item.Request.Method,
		})),
	}

	for _, v := range item.Request.Url.Queries {
		if v.Disabled {
			continue
		}
		req.Parameters.Query = append(req.Parameters.Query, v.toSchema())
	}
	for _, v := range item.Request.Headers {
		req.Parameters.Header = append(req.Parameters.Header, v.toSchema())
	}

	if body := encodeRequestBody(item.Request.Body); body != nil {
		req.Content = body
	}
	nodes = append(nodes, spec.MuseCreateNodeProxy(spec.WarpHTTPNode(req)))
	nodes = append(nodes, spec.MuseCreateNodeProxy(spec.WarpHTTPNode(encodeResponseBody(item.Response))))
	return nodes
}

var contenttypemapp = map[string]string{
	"json":      "application/json",
	"urlencode": "application/x-www-form-urlencoded",
	"formdata":  "multipart/form-data",
	"plain":     "text/plain",
}

func encodeRequestBody(body *Body) map[string]*spec.Schema {
	if body == nil || body.Disabled {
		return nil
	}
	switch body.Mode {
	case "raw":
		if body.Options.Raw.Language == "json" {
			b := jsonToSchema(body.Raw)
			return map[string]*spec.Schema{
				contenttypemapp["json"]: {
					Schema: b,
				},
			}
		}
	case "formdata", "urlencode":
		b := jsonschema.Create("object")
		b.Properties = make(map[string]*jsonschema.Schema)
		for _, v := range body.Formdata {
			if v.Disabled {
				continue
			}
			b.Properties[v.Key] = v.toJSONSchema()
		}
		return map[string]*spec.Schema{
			contenttypemapp[body.Mode]: {
				Schema: b,
			},
		}
	case "file":
	case "graphql":
	default:
	}

	return nil
}

func encodeResponseBody(res []Response) *spec.HTTPResponsesNode {
	response := &spec.HTTPResponsesNode{
		List: make(spec.HTTPResponses, 0),
	}
	for _, v := range res {
		r := spec.HTTPResponse{Code: v.Code}
		r.Description = v.Name
		switch v.PostmanePreviewLanguage {
		case "json":

			// fmt.Println("json.........")
			b := jsonToSchema(v.Body)
			b.Example = v.Body
			r.Content = map[string]*spec.Schema{
				contenttypemapp["json"]: {
					Schema: b,
				},
			}

		case "plain":
			b := jsonschema.Create("string")
			b.Example = v.Body
			r.Content = map[string]*spec.Schema{
				contenttypemapp["plain"]: {
					Schema: b,
				},
			}
		default:

		}

		for _, v := range v.Header {
			if v.Disabled {
				continue
			}
			b := jsonschema.Create("string")
			b.Example = v.Value
			r.Header = append(r.Header, &spec.Schema{
				Name:        v.Key,
				Description: v.Description,
				Schema:      b,
			})

		}

		// fmt.Println(r)

		response.List = append(response.List, &r)
	}

	if len(response.List) == 0 {
		defaultres := spec.HTTPResponse{Code: 200}
		defaultres.Name = "success"
		response.List = append(response.List, &defaultres)
	}
	return response
}
