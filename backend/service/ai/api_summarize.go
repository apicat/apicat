package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
	"github.com/gin-gonic/gin"
)

func APISummarize(ctx *gin.Context, collection *spec.Collection) (string, error) {
	if collection.Content == nil {
		return "", errors.New("collection content is nil")
	}

	summary := make([]string, 0)
	if collection.Title != "" {
		summary = append(summary, fmt.Sprintf("This is an API about \"%s\"", collection.Title))
	}

	var (
		err      error
		url      *spec.HTTPURLNode
		request  *spec.HTTPRequestNode
		response *spec.HTTPResponsesNode
	)
	for _, node := range collection.Content {
		switch node.NodeType() {
		case spec.NAME_HTTP_URL:
			url, err = node.ToHTTPURLNode()
		case spec.NAME_HTTP_REQUEST:
			request, err = node.ToHTTPRequestNode()
		case spec.NAME_HTTP_RESPONSES:
			response, err = node.ToHTTPResponsesNode()
		}
		if err != nil {
			return "", err
		}
	}
	if url == nil || request == nil || response == nil {
		return "", errors.New("incomplete API")
	}

	requestBaseInfo, err := apiRequestSummarize(url, request)
	if err != nil {
		return "", err
	}
	summary = append(summary, requestBaseInfo...)
	summary = append(summary, apiBodySummarize(&request.Content, true)...)
	summary = append(summary, apiResponseSummarize(response)...)

	return strings.Join(summary, "\n"), nil
}

func apiRequestSummarize(url *spec.HTTPURLNode, request *spec.HTTPRequestNode) ([]string, error) {
	summary := make([]string, 0)

	if url.Path == "" || url.Method == "" {
		return summary, errors.New("incomplete url")
	}
	summary = append(summary, fmt.Sprintf("The request path of the API is \"%s\", and the request method is \"%s\".", url.Path, url.Method))

	if request.Parameters.Path != nil && len(request.Parameters.Path) > 0 {
		summary = append(summary, fmt.Sprintf("There are %d Path request parameters, namely: ", len(request.Parameters.Path)))
		summary = append(summary, apiParameterSummarize(&request.Parameters.Path)...)
	}
	if request.Parameters.Header != nil && len(request.Parameters.Header) > 0 {
		summary = append(summary, fmt.Sprintf("There are %d Header request parameters, namely: ", len(request.Parameters.Header)))
		summary = append(summary, apiParameterSummarize(&request.Parameters.Header)...)
	}
	if request.Parameters.Cookie != nil && len(request.Parameters.Cookie) > 0 {
		summary = append(summary, fmt.Sprintf("There are %d Cookie request parameters, namely: ", len(request.Parameters.Cookie)))
		summary = append(summary, apiParameterSummarize(&request.Parameters.Cookie)...)
	}
	if request.Parameters.Query != nil && len(request.Parameters.Query) > 0 {
		summary = append(summary, fmt.Sprintf("There are %d Query request parameters, namely: ", len(request.Parameters.Query)))
		summary = append(summary, apiParameterSummarize(&request.Parameters.Query)...)
	}

	return summary, nil
}

func apiParameterSummarize(parameters *spec.ParameterList) []string {
	summary := make([]string, 0)
	for i, p := range *parameters {
		if p.Name == "" {
			continue
		}

		content := fmt.Sprintf("%d. %s: %s", i+1, p.Name, p.Schema.Type.Value()[0])
		if p.Required {
			content += ", required"
		} else {
			content += ", optional"
		}
		if p.Schema.Description != "" {
			content += fmt.Sprintf(", %s", p.Schema.Description)
		}
		summary = append(summary, content)
	}

	return summary
}

func apiResponseSummarize(responses *spec.HTTPResponsesNode) []string {
	summary := make([]string, 0)

	summary = append(summary, fmt.Sprintf("%d response situations are possible:", len(responses.List)))
	for i, response := range responses.List {
		desc := fmt.Sprintf("%d. Response status code %d: %s", i+1, response.Code, response.Name)
		if response.Description != "" {
			summary = append(summary, fmt.Sprintf("%s, %s", desc, response.Description))
		} else {
			summary = append(summary, desc)
		}

		if response.Header != nil && len(response.Header) > 0 {
			summary = append(summary, fmt.Sprintf("There are %d Header response parameters, namely: ", len(response.Header)))
			summary = append(summary, apiParameterSummarize(&response.Header)...)
		}

		if response.Content != nil && len(response.Content) > 0 {
			summary = append(summary, apiBodySummarize(&response.Content, false)...)
		}
	}

	return summary
}

func apiBodySummarize(body *spec.HTTPBody, isRequest bool) []string {
	summary := make([]string, 0)
	for contentType, schema := range *body {
		if contentType == "none" {
			continue
		}

		var contentTypeDesc string
		if isRequest {
			contentTypeDesc = fmt.Sprintf("The \"Content-Type\" in the request body is: %s", contentType)
		} else {
			contentTypeDesc = fmt.Sprintf("The \"Content-Type\" in the response body is: %s", contentType)
		}

		content := make(map[string]any)
		apiJsonSchemaSummarize("root", true, schema.Schema, content)

		switch contentType {
		case "multipart/form-data", "application/x-www-form-urlencoded":
			if len(content) > 0 {
				summary = append(summary, fmt.Sprintf("%s, containing the following %d parameters: ", contentTypeDesc, len(content)))
				i := 1
				for k, v := range content["root"].(map[string]any) {
					summary = append(summary, fmt.Sprintf("%d. %s: %s", i, k, v.(string)))
					i++
				}
			}
		case "application/json", "application/xml":
			if len(content) > 0 {
				summary = append(summary, fmt.Sprintf("%s, and the structure is as follows: ", contentTypeDesc))
				js, err := json.Marshal(content["root"])
				if err != nil {
					return summary
				}
				summary = append(summary, string(js))
			}
		case "raw", "application/octet-stream":
			if schema.Schema.Example != nil {
				if example, ok := schema.Schema.Example.(string); ok {
					summary = append(summary, fmt.Sprintf("%s. For example: %s", contentTypeDesc, example))
				}
			}
		}
	}

	return summary
}

// @title apiJsonSchemaSummarize
// @description 总结 JSON Schema 要表达的含义
// @param parameterName string 参数的名称
// @param required bool 参数是否必填
// @param schema *jsonschema 参数的 JSON Schema 描述
// @param result map[string]any 带结构的描述结果
// @return nil
func apiJsonSchemaSummarize(parameterName string, required bool, schema *jsonschema.Schema, result map[string]any) {
	if schema.Type == nil {
		return
	}

	switch schema.Type.Value()[0] {
	case "string", "integer", "number", "boolean":
		requiredStr := "optional"
		if required {
			requiredStr = "required"
		}
		desc := ""
		if schema.Description != "" {
			desc = fmt.Sprintf(", %s", schema.Description)
		}
		result[parameterName] = fmt.Sprintf("%s, %s%s", schema.Type.Value()[0], requiredStr, desc)
	case "object":
		if schema.Properties != nil && len(schema.Properties) > 0 {
			requiredParams := make(map[string]bool)
			if schema.Required != nil && len(schema.Required) > 0 {
				for _, v := range schema.Required {
					requiredParams[v] = true
				}
			}

			children := make(map[string]any)
			for k, v := range schema.Properties {
				_, r := requiredParams[k]
				apiJsonSchemaSummarize(k, r, v, children)
			}
			if len(children) > 0 {
				result[parameterName] = children
			}
		}
	case "array":
		if schema.Items != nil && !schema.Items.IsBool() {
			children := make(map[string]any)
			apiJsonSchemaSummarize("items", true, schema.Items.Value(), children)
			if len(children) > 0 {
				for _, v := range children {
					result[parameterName] = v
					break
				}
			}
		}
	}
}
