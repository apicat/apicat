package loader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/apicat/apicat/v3/internal/converter"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadMultiFileRef(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "openapi.yaml"), `openapi: 3.0.3
info: {title: Multi, version: '1'}
paths:
  /pets:
    get:
      responses:
        '200':
          description: ok
          content:
            application/json:
              schema:
                $ref: './schemas/pet.yaml#/Pet'
`)
	writeFile(t, filepath.Join(dir, "schemas", "pet.yaml"), `Pet:
  type: object
  properties:
    name: {type: string}
`)

	docModel, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	apiDoc, err := converter.Convert(docModel)
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}

	sv := apiDoc.TagGroups[0].Endpoints[0].Responses[0].Content[0].Schema
	if sv == nil {
		t.Fatal("schema is nil")
	}
	if sv.Type != "object" {
		t.Errorf("Type = %q, want object (cross-file $ref not resolved)", sv.Type)
	}
	if len(sv.Properties) != 1 || sv.Properties[0].Name != "name" {
		t.Errorf("Properties = %+v", sv.Properties)
	}
}

func TestLoadSpecFilePriority(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "openapi.json"), `{"openapi":"3.0.3","info":{"title":"json","version":"1"},"paths":{}}`)
	writeFile(t, filepath.Join(dir, "openapi.yaml"), "openapi: 3.0.3\ninfo: {title: yaml, version: '1'}\npaths: {}\n")

	docModel, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := docModel.Model.Info.Title; got != "yaml" {
		t.Errorf("Title = %q, want %q (openapi.yaml should win over openapi.json)", got, "yaml")
	}
}

func TestLoadMissingSpec(t *testing.T) {
	if _, err := Load(t.TempDir()); err == nil {
		t.Fatal("Load on empty dir should fail")
	}
}
