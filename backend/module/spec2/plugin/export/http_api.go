package export

import (
	"github.com/apicat/apicat/v2/backend/module/spec2"
)

type HttpApi struct {
	ID    int64
	Dir   string
	Title string
	spec2.HttpRequestAttrs
	Responses spec2.Responses
}

type HttpApis map[string]map[string]HttpApi

func NewHttpApis(s *spec2.Spec) HttpApis {
	apis := HttpApis{}

	collections := spec2.Collections{}
	for _, v := range s.Collections {
		if v.Type == spec2.TYPE_CATEGORY {
			collections = append(collections, v.ItemsTreeToList()...)
		} else if v.Type == spec2.TYPE_HTTP {
			collections = append(collections, v)
		}
	}

	models := spec2.DefinitionModels{}
	for _, v := range s.Definitions.Schemas {
		if v.Type == spec2.TYPE_CATEGORY {
			models = append(models, v.ItemsTreeToList()...)
		} else {
			models = append(models, v)
		}
	}

	responses := spec2.DefinitionResponses{}
	for _, v := range s.Definitions.Responses {
		if v.Type == spec2.TYPE_CATEGORY {
			responses = append(responses, v.ItemsTreeToList()...)
		} else {
			responses = append(responses, v)
		}
	}

	collections.DeepDerefAll(s.Globals.Parameters, models, responses)
	for _, c := range collections {
		path, method := "", ""
		http := HttpApi{
			ID:    int64(c.ID),
			Title: c.Title,
		}
		for _, node := range c.Content {
			switch node.NodeType() {
			case spec2.NODE_HTTP_URL:
				url := node.ToHttpUrl()
				if url.Attrs.Path == "" && url.Attrs.Method == "" {
					continue
				}

				path = url.Attrs.Path
				method = url.Attrs.Method
			case spec2.NODE_HTTP_REQUEST:
				http.HttpRequestAttrs = *node.ToHttpRequest().Attrs
			case spec2.NODE_HTTP_RESPONSE:
				http.Responses = node.ToHttpResponse().Attrs.List
			}
		}
		apis[path] = map[string]HttpApi{method: http}
	}
	return apis
}
