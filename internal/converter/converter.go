package converter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/apicat/apicat/v3/internal/model"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"go.yaml.in/yaml/v4"
)

const maxSchemaDepth = 10

func Convert(docModel *libopenapi.DocumentModel[v3high.Document]) (*model.APIDoc, error) {
	doc := &docModel.Model

	apiDoc := &model.APIDoc{
		Title:   doc.Info.Title,
		Version: doc.Info.Version,
	}
	if doc.Info.Description != "" {
		apiDoc.Description = doc.Info.Description
	}

	for _, s := range doc.Servers {
		apiDoc.Servers = append(apiDoc.Servers, model.Server{
			URL:         s.URL,
			Description: s.Description,
		})
	}

	tagDescMap := make(map[string]string)
	for _, t := range doc.Tags {
		tag := model.Tag{Name: t.Name, Description: t.Description}
		apiDoc.Tags = append(apiDoc.Tags, tag)
		tagDescMap[t.Name] = t.Description
	}

	apiDoc.Security = extractSecuritySchemes(doc)

	endpointsByTag := make(map[string][]model.Endpoint)
	if doc.Paths != nil && doc.Paths.PathItems != nil {
		for pair := doc.Paths.PathItems.First(); pair != nil; pair = pair.Next() {
			pathStr := pair.Key()
			pathItem := pair.Value()
			ops := pathItem.GetOperations()
			if ops == nil {
				continue
			}
			for opPair := ops.First(); opPair != nil; opPair = opPair.Next() {
				method := strings.ToUpper(opPair.Key())
				op := opPair.Value()
				endpoint := convertOperation(method, pathStr, op)
				tags := op.Tags
				if len(tags) == 0 {
					tags = []string{"default"}
				}
				for _, tag := range tags {
					endpointsByTag[tag] = append(endpointsByTag[tag], endpoint)
				}
			}
		}
	}

	for _, t := range apiDoc.Tags {
		if eps, ok := endpointsByTag[t.Name]; ok {
			apiDoc.TagGroups = append(apiDoc.TagGroups, model.TagGroup{
				Tag:       t,
				Endpoints: eps,
			})
			delete(endpointsByTag, t.Name)
		}
	}
	for tagName, eps := range endpointsByTag {
		apiDoc.TagGroups = append(apiDoc.TagGroups, model.TagGroup{
			Tag:       model.Tag{Name: tagName, Description: tagDescMap[tagName]},
			Endpoints: eps,
		})
	}

	return apiDoc, nil
}

func extractSecuritySchemes(doc *v3high.Document) []model.SecurityScheme {
	var schemes []model.SecurityScheme
	if doc.Components == nil || doc.Components.SecuritySchemes == nil {
		return schemes
	}
	for pair := doc.Components.SecuritySchemes.First(); pair != nil; pair = pair.Next() {
		ss := pair.Value()
		scheme := model.SecurityScheme{
			Name: pair.Key(),
			Type: ss.Type,
		}
		if ss.Scheme != "" {
			scheme.Scheme = ss.Scheme
		}
		schemes = append(schemes, scheme)
	}
	return schemes
}

func convertOperation(method, path string, op *v3high.Operation) model.Endpoint {
	ep := model.Endpoint{
		Method:      method,
		Path:        path,
		Summary:     op.Summary,
		Description: op.Description,
		OperationID: op.OperationId,
		Anchor:      makeAnchor(method, path),
	}
	if op.Deprecated != nil && *op.Deprecated {
		ep.Deprecated = true
	}

	for _, p := range op.Parameters {
		ep.Parameters = append(ep.Parameters, convertParameter(p))
	}

	if op.RequestBody != nil {
		ep.RequestBody = convertRequestBody(op.RequestBody)
	}

	if op.Responses != nil {
		ep.Responses = convertResponses(op.Responses)
	}

	return ep
}

func convertParameter(p *v3high.Parameter) model.Parameter {
	param := model.Parameter{
		Name:        p.Name,
		In:          p.In,
		Description: p.Description,
		Deprecated:  p.Deprecated,
	}
	if p.Required != nil && *p.Required {
		param.Required = true
	}
	if p.Schema != nil {
		param.Schema = convertSchemaProxy(p.Schema, "", false, 0)
	}
	return param
}

func convertRequestBody(rb *v3high.RequestBody) *model.RequestBody {
	result := &model.RequestBody{
		Description: rb.Description,
	}
	if rb.Required != nil && *rb.Required {
		result.Required = true
	}
	if rb.Content != nil {
		for pair := rb.Content.First(); pair != nil; pair = pair.Next() {
			mt := model.MediaTypeContent{MediaType: pair.Key()}
			if pair.Value().Schema != nil {
				mt.Schema = convertSchemaProxy(pair.Value().Schema, "", false, 0)
			}
			result.Content = append(result.Content, mt)
		}
	}
	return result
}

func convertResponses(responses *v3high.Responses) []model.Response {
	var result []model.Response
	if responses.Codes != nil {
		for pair := responses.Codes.First(); pair != nil; pair = pair.Next() {
			result = append(result, convertResponse(pair.Key(), pair.Value()))
		}
	}
	if responses.Default != nil {
		result = append(result, convertResponse("default", responses.Default))
	}
	return result
}

func convertResponse(code string, resp *v3high.Response) model.Response {
	r := model.Response{
		StatusCode:  code,
		Description: resp.Description,
	}
	if resp.Content != nil {
		for pair := resp.Content.First(); pair != nil; pair = pair.Next() {
			mt := model.MediaTypeContent{MediaType: pair.Key()}
			if pair.Value().Schema != nil {
				mt.Schema = convertSchemaProxy(pair.Value().Schema, "", false, 0)
			}
			r.Content = append(r.Content, mt)
		}
	}
	return r
}

func convertSchemaProxy(proxy *base.SchemaProxy, name string, required bool, depth int) *model.SchemaView {
	if proxy == nil || depth > maxSchemaDepth {
		return nil
	}

	sv := &model.SchemaView{
		Name:     name,
		Required: required,
		Depth:    depth,
		IsRef:    proxy.IsReference(),
	}

	if proxy.IsReference() {
		sv.RefName = extractRefName(proxy.GetReference())
	}

	schema := proxy.Schema()
	if schema == nil {
		return sv
	}

	if len(schema.Type) > 0 {
		sv.Type = schema.Type[0]
	}
	sv.Format = schema.Format
	sv.Description = schema.Description
	sv.Pattern = schema.Pattern

	if schema.Nullable != nil && *schema.Nullable {
		sv.Nullable = true
	}
	if schema.Deprecated != nil && *schema.Deprecated {
		sv.Deprecated = true
	}

	sv.MinLength = schema.MinLength
	sv.MaxLength = schema.MaxLength
	sv.Minimum = schema.Minimum
	sv.Maximum = schema.Maximum
	sv.MinItems = schema.MinItems
	sv.MaxItems = schema.MaxItems

	sv.Enum = extractEnum(schema.Enum)

	if schema.Default != nil {
		sv.Default = yamlNodeValue(schema.Default)
	}
	if schema.Example != nil {
		sv.Example = yamlNodeValue(schema.Example)
	}

	if schema.Properties != nil {
		requiredSet := toSet(schema.Required)
		for pair := schema.Properties.First(); pair != nil; pair = pair.Next() {
			propView := convertSchemaProxy(pair.Value(), pair.Key(), requiredSet[pair.Key()], depth+1)
			if propView != nil {
				sv.Properties = append(sv.Properties, *propView)
			}
		}
	}

	if schema.Items != nil && schema.Items.IsA() {
		sv.Items = convertSchemaProxy(schema.Items.A, "", false, depth+1)
	}

	for _, ap := range schema.AllOf {
		child := convertSchemaProxy(ap, "", false, depth+1)
		if child != nil {
			sv.AllOf = append(sv.AllOf, *child)
		}
	}
	for _, op := range schema.OneOf {
		child := convertSchemaProxy(op, "", false, depth+1)
		if child != nil {
			sv.OneOf = append(sv.OneOf, *child)
		}
	}
	for _, ap := range schema.AnyOf {
		child := convertSchemaProxy(ap, "", false, depth+1)
		if child != nil {
			sv.AnyOf = append(sv.AnyOf, *child)
		}
	}

	return sv
}

func extractRefName(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}

func extractEnum(nodes []*yaml.Node) []string {
	var values []string
	for _, n := range nodes {
		if n != nil {
			values = append(values, n.Value)
		}
	}
	return values
}

func yamlNodeValue(n *yaml.Node) string {
	if n == nil {
		return ""
	}
	return n.Value
}

func toSet(items []string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, item := range items {
		m[item] = true
	}
	return m
}

var anchorReplacer = regexp.MustCompile(`[^a-z0-9]+`)

func makeAnchor(method, path string) string {
	s := strings.ToLower(fmt.Sprintf("%s-%s", method, path))
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")
	s = anchorReplacer.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
