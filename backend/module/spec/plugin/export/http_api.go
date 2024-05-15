package export

import (
	"github.com/apicat/apicat/v2/backend/module/spec"
)

type HttpApi struct {
	ID    int64
	Dir   string
	Title string
	spec.HttpRequestAttrs
	Responses spec.Responses
}

type HttpApis map[string]map[string]HttpApi

func NewHttpApis(s *spec.Spec) HttpApis {
	apis := HttpApis{}

	collections := spec.Collections{}
	for _, v := range s.Collections {
		if v.Type == spec.TYPE_CATEGORY {
			collections = append(collections, v.ItemsTreeToList()...)
		} else if v.Type == spec.TYPE_HTTP {
			collections = append(collections, v)
		}
	}

	definitions := &spec.Definitions{}

	for _, v := range s.Definitions.Schemas {
		if v.Type == spec.TYPE_CATEGORY {
			definitions.Schemas = append(definitions.Schemas, v.ItemsTreeToList()...)
		} else {
			definitions.Schemas = append(definitions.Schemas, v)
		}
	}

	for _, v := range s.Definitions.Responses {
		if v.Type == spec.TYPE_CATEGORY {
			definitions.Responses = append(definitions.Responses, v.ItemsTreeToList()...)
		} else {
			definitions.Responses = append(definitions.Responses, v)
		}
	}

	collections.DeepDerefAll(s.Globals.Parameters, definitions)
	for _, c := range collections {
		path, method := "", ""
		http := HttpApi{
			ID:    int64(c.ID),
			Title: c.Title,
		}
		for _, node := range c.Content {
			switch node.NodeType() {
			case spec.NODE_HTTP_URL:
				url := node.ToHttpUrl()
				if url.Attrs.Path == "" && url.Attrs.Method == "" {
					continue
				}

				path = url.Attrs.Path
				method = url.Attrs.Method
			case spec.NODE_HTTP_REQUEST:
				http.HttpRequestAttrs = *node.ToHttpRequest().Attrs
			case spec.NODE_HTTP_RESPONSE:
				http.Responses = node.ToHttpResponse().Attrs.List
			}
		}
		apis[path] = map[string]HttpApi{method: http}
	}
	return apis
}
