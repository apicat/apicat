package converter

import (
	"testing"

	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/apicat/apicat/v3/internal/model"
)

func buildModel(t *testing.T, spec string) *libopenapi.DocumentModel[v3high.Document] {
	t.Helper()
	doc, err := libopenapi.NewDocument([]byte(spec))
	if err != nil {
		t.Fatalf("NewDocument: %v", err)
	}
	m, err := doc.BuildV3Model()
	if err != nil {
		t.Fatalf("BuildV3Model: %v", err)
	}
	return m
}

func mustConvert(t *testing.T, spec string) *model.APIDoc {
	t.Helper()
	apiDoc, err := Convert(buildModel(t, spec))
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	return apiDoc
}

func TestMakeAnchor(t *testing.T) {
	cases := []struct {
		method, path, want string
	}{
		{"GET", "/users", "get-users"},
		{"GET", "/users/{id}", "get-users-id"},
		{"POST", "/users/{id}/pets", "post-users-id-pets"},
		{"DELETE", "/a/b_c.d", "delete-a-b-c-d"},
		{"GET", "/", "get"},
	}
	for _, c := range cases {
		if got := makeAnchor(c.method, c.path); got != c.want {
			t.Errorf("makeAnchor(%q, %q) = %q, want %q", c.method, c.path, got, c.want)
		}
	}
}

const petStoreSpec = `openapi: 3.0.3
info:
  title: Pet Store
  version: 1.2.3
  description: A **test** spec.
servers:
  - url: https://api.example.com
    description: prod
tags:
  - name: pets
    description: Pet operations
paths:
  /pets:
    get:
      tags: [pets]
      summary: List pets
      operationId: listPets
      parameters:
        - name: limit
          in: query
          required: true
          schema:
            type: integer
            format: int32
      responses:
        '200':
          description: ok
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Pet'
    post:
      tags: [pets]
      summary: Create pet
      deprecated: true
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Pet'
      responses:
        '201':
          description: created
        default:
          description: error
  /health:
    get:
      summary: Health
      responses:
        '200':
          description: ok
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
  schemas:
    Pet:
      type: object
      required: [id, name]
      properties:
        id:
          type: integer
          format: int64
        name:
          type: string
        status:
          type: string
          enum: [available, sold]
        nickname:
          type: string
          nullable: true
`

func TestConvertDocumentInfo(t *testing.T) {
	apiDoc := mustConvert(t, petStoreSpec)

	if apiDoc.Title != "Pet Store" {
		t.Errorf("Title = %q, want %q", apiDoc.Title, "Pet Store")
	}
	if apiDoc.Version != "1.2.3" {
		t.Errorf("Version = %q, want %q", apiDoc.Version, "1.2.3")
	}
	if apiDoc.Description != "A **test** spec." {
		t.Errorf("Description = %q", apiDoc.Description)
	}

	if len(apiDoc.Servers) != 1 {
		t.Fatalf("len(Servers) = %d, want 1", len(apiDoc.Servers))
	}
	if apiDoc.Servers[0].URL != "https://api.example.com" || apiDoc.Servers[0].Description != "prod" {
		t.Errorf("Servers[0] = %+v", apiDoc.Servers[0])
	}

	if len(apiDoc.Security) != 1 {
		t.Fatalf("len(Security) = %d, want 1", len(apiDoc.Security))
	}
	sec := apiDoc.Security[0]
	if sec.Name != "bearerAuth" || sec.Type != "http" || sec.Scheme != "bearer" {
		t.Errorf("Security[0] = %+v", sec)
	}
}

func TestConvertTagGrouping(t *testing.T) {
	apiDoc := mustConvert(t, petStoreSpec)

	if len(apiDoc.TagGroups) != 2 {
		t.Fatalf("len(TagGroups) = %d, want 2", len(apiDoc.TagGroups))
	}

	// Declared tags come first, in declaration order.
	pets := apiDoc.TagGroups[0]
	if pets.Tag.Name != "pets" || pets.Tag.Description != "Pet operations" {
		t.Errorf("TagGroups[0].Tag = %+v", pets.Tag)
	}
	if len(pets.Endpoints) != 2 {
		t.Fatalf("pets endpoints = %d, want 2", len(pets.Endpoints))
	}
	if pets.Endpoints[0].Method != "GET" || pets.Endpoints[1].Method != "POST" {
		t.Errorf("pets endpoint order = %s, %s", pets.Endpoints[0].Method, pets.Endpoints[1].Method)
	}

	// Untagged operations land in "default".
	def := apiDoc.TagGroups[1]
	if def.Tag.Name != "default" {
		t.Errorf("TagGroups[1].Tag.Name = %q, want %q", def.Tag.Name, "default")
	}
	if len(def.Endpoints) != 1 || def.Endpoints[0].Path != "/health" {
		t.Errorf("default group endpoints = %+v", def.Endpoints)
	}
}

func TestConvertEndpointFields(t *testing.T) {
	apiDoc := mustConvert(t, petStoreSpec)
	listPets := apiDoc.TagGroups[0].Endpoints[0]

	if listPets.Anchor != "get-pets" {
		t.Errorf("Anchor = %q, want %q", listPets.Anchor, "get-pets")
	}
	if listPets.Summary != "List pets" || listPets.OperationID != "listPets" {
		t.Errorf("Summary/OperationID = %q / %q", listPets.Summary, listPets.OperationID)
	}
	if listPets.Deprecated {
		t.Error("GET /pets should not be deprecated")
	}

	if len(listPets.Parameters) != 1 {
		t.Fatalf("len(Parameters) = %d, want 1", len(listPets.Parameters))
	}
	p := listPets.Parameters[0]
	if p.Name != "limit" || p.In != "query" || !p.Required {
		t.Errorf("Parameters[0] = %+v", p)
	}
	if p.Schema == nil || p.Schema.Type != "integer" || p.Schema.Format != "int32" {
		t.Errorf("Parameters[0].Schema = %+v", p.Schema)
	}

	createPet := apiDoc.TagGroups[0].Endpoints[1]
	if !createPet.Deprecated {
		t.Error("POST /pets should be deprecated")
	}
	if createPet.RequestBody == nil || !createPet.RequestBody.Required {
		t.Fatalf("RequestBody = %+v", createPet.RequestBody)
	}
	if len(createPet.RequestBody.Content) != 1 || createPet.RequestBody.Content[0].MediaType != "application/json" {
		t.Errorf("RequestBody.Content = %+v", createPet.RequestBody.Content)
	}

	// Response codes in declaration order, default last.
	if len(createPet.Responses) != 2 {
		t.Fatalf("len(Responses) = %d, want 2", len(createPet.Responses))
	}
	if createPet.Responses[0].StatusCode != "201" || createPet.Responses[1].StatusCode != "default" {
		t.Errorf("response order = %s, %s", createPet.Responses[0].StatusCode, createPet.Responses[1].StatusCode)
	}
}

func TestConvertSchemaExpansion(t *testing.T) {
	apiDoc := mustConvert(t, petStoreSpec)
	listPets := apiDoc.TagGroups[0].Endpoints[0]

	if len(listPets.Responses) != 1 || len(listPets.Responses[0].Content) != 1 {
		t.Fatalf("unexpected responses: %+v", listPets.Responses)
	}
	arr := listPets.Responses[0].Content[0].Schema
	if arr == nil || arr.Type != "array" || arr.Depth != 0 {
		t.Fatalf("array schema = %+v", arr)
	}

	pet := arr.Items
	if pet == nil {
		t.Fatal("array items schema is nil")
	}
	if !pet.IsRef || pet.RefName != "Pet" {
		t.Errorf("items IsRef/RefName = %v/%q, want true/\"Pet\"", pet.IsRef, pet.RefName)
	}
	if pet.Type != "object" || pet.Depth != 1 {
		t.Errorf("items Type/Depth = %q/%d, want object/1", pet.Type, pet.Depth)
	}
	if len(pet.Properties) != 4 {
		t.Fatalf("len(pet.Properties) = %d, want 4", len(pet.Properties))
	}

	props := map[string]model.SchemaView{}
	for _, prop := range pet.Properties {
		props[prop.Name] = prop
	}
	if !props["id"].Required || !props["name"].Required {
		t.Error("id and name should be required")
	}
	if props["status"].Required {
		t.Error("status should not be required")
	}
	if got := props["status"].Enum; len(got) != 2 || got[0] != "available" || got[1] != "sold" {
		t.Errorf("status enum = %v", got)
	}
	if !props["nickname"].Nullable {
		t.Error("nickname should be nullable")
	}
	if props["id"].Depth != 2 {
		t.Errorf("property depth = %d, want 2", props["id"].Depth)
	}
	if props["id"].Format != "int64" {
		t.Errorf("id format = %q, want int64", props["id"].Format)
	}
}

func TestConvertComposition(t *testing.T) {
	spec := `openapi: 3.0.3
info: {title: T, version: '1'}
paths:
  /things:
    get:
      responses:
        '200':
          description: ok
          content:
            application/json:
              schema:
                allOf:
                  - $ref: '#/components/schemas/Base'
                  - type: object
                    properties:
                      extra: {type: string}
    post:
      requestBody:
        content:
          application/json:
            schema:
              oneOf:
                - {type: string}
                - {type: integer}
              anyOf:
                - {type: boolean}
      responses:
        '200': {description: ok}
components:
  schemas:
    Base:
      type: object
      properties:
        id: {type: integer}
`
	apiDoc := mustConvert(t, spec)
	eps := apiDoc.TagGroups[0].Endpoints

	allOf := eps[0].Responses[0].Content[0].Schema.AllOf
	if len(allOf) != 2 {
		t.Fatalf("len(AllOf) = %d, want 2", len(allOf))
	}
	if !allOf[0].IsRef || allOf[0].RefName != "Base" {
		t.Errorf("AllOf[0] = %+v", allOf[0])
	}
	if len(allOf[0].Properties) != 1 || allOf[0].Properties[0].Name != "id" {
		t.Errorf("AllOf[0].Properties = %+v", allOf[0].Properties)
	}
	if len(allOf[1].Properties) != 1 || allOf[1].Properties[0].Name != "extra" {
		t.Errorf("AllOf[1].Properties = %+v", allOf[1].Properties)
	}

	body := eps[1].RequestBody.Content[0].Schema
	if len(body.OneOf) != 2 || body.OneOf[0].Type != "string" || body.OneOf[1].Type != "integer" {
		t.Errorf("OneOf = %+v", body.OneOf)
	}
	if len(body.AnyOf) != 1 || body.AnyOf[0].Type != "boolean" {
		t.Errorf("AnyOf = %+v", body.AnyOf)
	}
}

func TestConvertRecursiveSchemaDepthLimit(t *testing.T) {
	spec := `openapi: 3.0.3
info: {title: T, version: '1'}
paths:
  /nodes:
    get:
      responses:
        '200':
          description: ok
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Node'
components:
  schemas:
    Node:
      type: object
      properties:
        child:
          $ref: '#/components/schemas/Node'
`
	// Must terminate; without the depth guard this would recurse forever.
	apiDoc := mustConvert(t, spec)

	sv := apiDoc.TagGroups[0].Endpoints[0].Responses[0].Content[0].Schema
	depth := 0
	for sv != nil && len(sv.Properties) > 0 {
		sv = &sv.Properties[0]
		depth++
	}
	if depth > maxSchemaDepth {
		t.Errorf("schema chain depth = %d, exceeds maxSchemaDepth %d", depth, maxSchemaDepth)
	}
	if depth < 3 {
		t.Errorf("schema chain depth = %d, recursion cut off too early", depth)
	}
}
