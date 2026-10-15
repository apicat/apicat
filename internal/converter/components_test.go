package converter

import (
	"strings"
	"testing"
)

func TestConvertComponents(t *testing.T) {
	doc := mustConvert(t, `openapi: 3.1.0
info: {title: Components, version: '1'}
paths: {}
components:
  schemas:
    Node:
      type: object
      description: Tree node
      properties:
        child: {$ref: '#/components/schemas/Node'}
    Unused: {type: string, enum: [a, b]}
  parameters:
    Limit:
      name: limit
      in: query
      description: Result limit
      schema: {type: integer, default: 10}
  requestBodies:
    CreateNode:
      description: New node
      required: true
      content:
        application/json:
          schema: {$ref: '#/components/schemas/Node'}
  responses:
    NodeResult:
      description: A node
      content:
        application/json:
          schema: {$ref: '#/components/schemas/Node'}
  headers:
    RateLimit:
      description: Remaining requests
      schema: {type: integer}
  securitySchemes:
    ApiKey:
      type: apiKey
      name: X-API-Key
      in: header
      description: API key authentication
  examples:
    NodeExample:
      summary: Empty node
      value: {child: null}
  links:
    Next:
      operationId: getNode
      description: Fetch next node
  callbacks:
    OnUpdate:
      '{$request.body#/callbackUrl}':
        post:
          responses:
            '200': {description: received}
  pathItems:
    SharedPath:
      get:
        responses:
          '200': {description: ok}
`)
	wantKinds := []string{"schemas", "parameters", "requestBodies", "responses", "headers", "securitySchemes", "examples", "links", "callbacks", "pathItems"}
	if len(doc.ComponentGroups) != len(wantKinds) {
		t.Fatalf("groups = %d, want %d", len(doc.ComponentGroups), len(wantKinds))
	}
	for i, group := range doc.ComponentGroups {
		if group.Kind != wantKinds[i] || len(group.Definitions) == 0 {
			t.Fatalf("group %d = %+v", i, group)
		}
		for _, component := range group.Definitions {
			if component.Name == "" || component.Definition == "" {
				t.Errorf("missing name or complete definition: %+v", component)
			}
		}
	}
	schemas := doc.ComponentGroups[0].Definitions
	if len(schemas) != 2 || schemas[1].Name != "Unused" {
		t.Fatal("unreferenced schemas must be included in declaration order")
	}
	if !schemas[0].Schema.Properties[0].Recursive || schemas[0].Description != "Tree node" {
		t.Fatalf("recursive schema = %+v", schemas[0])
	}
	if !strings.Contains(schemas[0].Definition, "$ref:") || len(schemas[0].Definition) > 1000 {
		t.Fatal("source must preserve recursive refs without unbounded expansion")
	}
	if p := doc.ComponentGroups[1].Definitions[0].Parameters[0]; p.Name != "limit" || p.Schema.Default != "10" {
		t.Fatalf("parameter = %+v", p)
	}
	if body := doc.ComponentGroups[2].Definitions[0].RequestBody; !body.Required || body.Content[0].Schema.Type != "object" {
		t.Fatalf("body = %+v", body)
	}
	if response := doc.ComponentGroups[3].Definitions[0].Response; response.StatusCode != "" || response.Content[0].Schema.Type != "object" {
		t.Fatalf("response = %+v", response)
	}
	if header := doc.ComponentGroups[4].Definitions[0].Parameters[0]; header.In != "header" || header.Schema.Type != "integer" {
		t.Fatalf("header = %+v", header)
	}
	if !strings.Contains(doc.ComponentGroups[5].Definitions[0].Definition, "X-API-Key") {
		t.Fatal("security details must remain visible")
	}
}

func TestConvertWithoutComponents(t *testing.T) {
	doc := mustConvert(t, "openapi: 3.0.3\ninfo: {title: Empty, version: '1'}\npaths: {}\n")
	if len(doc.ComponentGroups) != 0 {
		t.Fatalf("unexpected groups: %+v", doc.ComponentGroups)
	}
}
