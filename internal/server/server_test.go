package server

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/apicat/apicat/v3/internal/model"
	"github.com/apicat/apicat/v3/internal/render"
)

func TestComponentPages(t *testing.T) {
	if err := render.Init(); err != nil {
		t.Fatal(err)
	}
	doc := &model.APIDoc{Title: "Test", ComponentGroups: []model.ComponentGroup{
		{Kind: "schemas", Title: "Schemas", Definitions: []model.Component{
			{Name: "User.A", Schema: &model.SchemaView{Type: "string"}, Definition: "type: string"},
			{Name: "User-A", Schema: &model.SchemaView{Type: "integer"}, Definition: "type: integer"},
			{Name: "名字 & value", Schema: &model.SchemaView{Type: "boolean"}, Definition: "type: boolean"},
		}},
		{Kind: "examples", Title: "Examples", Definitions: []model.Component{
			{Name: "User.A", Definition: "value: <script>alert(1)</script>"},
		}},
	}}
	handler, err := newHandler(func() *model.APIDoc { return doc })
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		path, want string
		code       int
	}{
		{"/", "/components?type=schemas", 200},
		{"/endpoints", "Components", 200},
		{"/components", "Reusable definitions", 200},
		{"/components?type=schemas&name=User.A", "type: string", 200},
		{"/components?type=schemas&name=User-A", "type: integer", 200},
		{"/components?type=schemas&name=" + url.QueryEscape("名字 & value"), "type: boolean", 200},
		{"/components?type=examples&name=User.A", "&lt;script&gt;", 200},
		{"/components?type=schemas&name=Unknown", "", 404},
		{"/components?type=unknown&name=User.A", "", 404},
		{"/components?type=schemas", "", 404},
	} {
		t.Run(tc.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest("GET", tc.path, nil))
			if response.Code != tc.code || !strings.Contains(response.Body.String(), tc.want) {
				t.Fatalf("status=%d, want %d; missing %q: %s", response.Code, tc.code, tc.want, response.Body.String())
			}
			if strings.Contains(response.Body.String(), "<script>alert(1)</script>") {
				t.Fatal("unescaped definition")
			}
		})
	}
	// Each request must pick up the latest document, including the component list.
	doc = &model.APIDoc{Title: "Updated"}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/components", nil))
	if response.Code != 200 || !strings.Contains(response.Body.String(), "No reusable definitions") {
		t.Fatal("stale component list")
	}
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/components?type=schemas&name=User.A", nil))
	if response.Code != 404 {
		t.Fatal("removed component is still served")
	}
}
