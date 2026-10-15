package converter

import (
	"context"
	"fmt"
	"strings"

	"github.com/apicat/apicat/v3/internal/model"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"go.yaml.in/yaml/v4"
)

func convertComponents(components *v3high.Components) ([]model.ComponentGroup, error) {
	if components == nil {
		return nil, nil
	}
	var groups []model.ComponentGroup
	var firstErr error
	add := func(group model.ComponentGroup, err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
		if len(group.Definitions) > 0 {
			groups = append(groups, group)
		}
	}
	add(componentGroup("schemas", "Schemas", components.Schemas))
	add(componentGroup("parameters", "Parameters", components.Parameters))
	add(componentGroup("requestBodies", "Request bodies", components.RequestBodies))
	add(componentGroup("responses", "Responses", components.Responses))
	add(componentGroup("headers", "Headers", components.Headers))
	add(componentGroup("securitySchemes", "Security schemes", components.SecuritySchemes))
	add(componentGroup("examples", "Examples", components.Examples))
	add(componentGroup("links", "Links", components.Links))
	add(componentGroup("callbacks", "Callbacks", components.Callbacks))
	add(componentGroup("pathItems", "Path items", components.PathItems))
	add(componentGroup("mediaTypes", "Media types", components.MediaTypes))
	return groups, firstErr
}

type componentRenderer interface {
	Render() ([]byte, error)
}

func componentDefinition(value componentRenderer) ([]byte, error) {
	var resolved any = value
	if proxy, ok := value.(*base.SchemaProxy); ok {
		if schema := proxy.Schema(); schema != nil {
			resolved = schema
		}
	}
	// The resolved source node contains the definition, while Render can return
	// only the registration's external $ref. Serialize one node without inlining
	// nested references, preserving recursive refs and all original metadata.
	if high, ok := resolved.(interface{ GoLowUntyped() any }); ok {
		if low, ok := high.GoLowUntyped().(interface {
			GetRootNode() *yaml.Node
			GetIndex() *index.SpecIndex
		}); ok && low.GetRootNode() != nil {
			node, idx := low.GetRootNode(), low.GetIndex()
			seen := make(map[*yaml.Node]bool)
			for idx != nil && !seen[node] {
				seen[node] = true
				ref := ""
				if node.Kind == yaml.MappingNode {
					for i := 0; i+1 < len(node.Content); i += 2 {
						if node.Content[i].Value == "$ref" {
							ref = node.Content[i+1].Value
							break
						}
					}
				}
				if ref == "" {
					break
				}
				target := idx.FindComponent(context.Background(), ref)
				if target == nil || target.Node == nil {
					break
				}
				node, idx = target.Node, target.Index
			}
			return yaml.Marshal(node)
		}
	}
	return value.Render()
}

func componentGroup[T componentRenderer](kind, title string, values *orderedmap.Map[string, T]) (model.ComponentGroup, error) {
	group := model.ComponentGroup{Kind: kind, Title: title}
	if values == nil {
		return group, nil
	}
	for pair := values.First(); pair != nil; pair = pair.Next() {
		value := pair.Value()
		// Keep the complete definition available for metadata outside the table view.
		definition, err := componentDefinition(value)
		if err != nil {
			return group, fmt.Errorf("cannot render component %s/%s: %w", kind, pair.Key(), err)
		}
		component := model.Component{Name: pair.Key(), Definition: strings.TrimSpace(string(definition))}
		switch v := any(value).(type) {
		case *base.SchemaProxy:
			component.Schema = convertSchemaProxy(v, "", false, 0)
			if component.Schema != nil {
				component.Description = component.Schema.Description
			}
		case *v3high.Parameter:
			component.Description = v.Description
			component.Parameters = []model.Parameter{convertParameter(v)}
		case *v3high.RequestBody:
			component.Description = v.Description
			component.RequestBody = convertRequestBody(v)
		case *v3high.Response:
			component.Description = v.Description
			response := convertResponse("", v)
			component.Response = &response
		case *v3high.Header:
			component.Description = v.Description
			component.Parameters = []model.Parameter{{
				Name: pair.Key(), In: "header", Description: v.Description,
				Required: v.Required, Deprecated: v.Deprecated,
				Schema: convertSchemaProxy(v.Schema, "", false, 0),
			}}
		case *v3high.SecurityScheme:
			component.Description = v.Description
		case *base.Example:
			component.Description = v.Summary
			if component.Description == "" {
				component.Description = v.Description
			}
		case *v3high.Link:
			component.Description = v.Description
		case *v3high.PathItem:
			component.Description = v.Description
		case *v3high.MediaType:
			component.Schema = convertSchemaProxy(v.Schema, "", false, 0)
		}
		group.Definitions = append(group.Definitions, component)
	}
	return group, nil
}
