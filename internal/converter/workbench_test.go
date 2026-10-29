package converter

import "testing"

func TestWorkbenchInheritanceAndExamples(t *testing.T) {
	doc := mustConvert(t, `openapi: 3.1.2
info: {title: Workbench, version: '1'}
servers:
  - url: https://{region}.example.com
    variables:
      region: {default: eu}
paths:
  /pets/{id}:
    servers: [{url: 'https://path.example.com'}]
    parameters:
      - {name: id, in: path, required: true, schema: {type: string}, example: abc}
      - {name: limit, in: query, schema: {type: integer}, example: 5}
    post:
      parameters:
        - {name: limit, in: query, schema: {type: integer}, example: 10}
      servers: [{url: 'https://operation.example.com'}]
      requestBody:
        content:
          application/json:
            example: {name: cat, count: 2}
            schema:
              type: object
              example: {name: cat}
      responses:
        '200': {description: OK}
`)
	if doc.Servers[0].URL != "https://eu.example.com" {
		t.Fatal(doc.Servers)
	}
	ep := doc.TagGroups[0].Endpoints[0]
	if len(ep.Parameters) != 2 || ep.Parameters[0].Name != "id" || ep.Parameters[1].Example != "10" {
		t.Fatalf("parameters: %+v", ep.Parameters)
	}
	if ep.Servers[0].URL != "https://operation.example.com" {
		t.Fatal(ep.Servers)
	}
	media := ep.RequestBody.Content[0]
	if media.Example != `{"count":2,"name":"cat"}` || media.Schema.JSONExample != `{"name":"cat"}` {
		t.Fatalf("example: %+v", media)
	}
}
