package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apicat/apicat/v3/internal/model"
)

func TestWorkbenchRoutesAndDocumentReload(t *testing.T) {
	doc := &model.APIDoc{Title: "First", ComponentGroups: []model.ComponentGroup{{Kind: "schemas", Definitions: []model.Component{{Name: "名字 & value", Definition: "<script>alert(1)</script>"}}}}}
	handler, err := newHandler(func() *model.APIDoc { return doc }, t.TempDir(), "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/", "/endpoints", "/endpoints/getArticleById", "/components?type=schemas&name=User", "/environments"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest("GET", "http://localhost"+path, nil))
		if response.Code != 200 || !strings.Contains(response.Body.String(), `id="root"`) {
			t.Fatalf("%s: %d %s", path, response.Code, response.Body.String())
		}
		if strings.Contains(response.Body.String(), "bulma") {
			t.Fatal("legacy UI shipped")
		}
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "http://localhost/api/document", nil))
	var received model.APIDoc
	if err := json.Unmarshal(response.Body.Bytes(), &received); err != nil {
		t.Fatal(err)
	}
	if received.Title != "First" || received.ComponentGroups[0].Definitions[0].Name != "名字 & value" {
		t.Fatalf("document = %+v", received)
	}
	doc = &model.APIDoc{Title: "Updated"}
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "http://localhost/api/document", nil))
	if !strings.Contains(response.Body.String(), "Updated") {
		t.Fatal("stale document")
	}
	for _, path := range []string{"/api/missing", "/assets/missing.js", "/missing"} {
		response = httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest("GET", "http://localhost"+path, nil))
		if response.Code != 404 {
			t.Errorf("%s = %d", path, response.Code)
		}
	}
}
func TestLocalAPIProtectionAndSettings(t *testing.T) {
	handler, err := newHandler(func() *model.APIDoc { return &model.APIDoc{} }, t.TempDir(), "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	bootstrap := httptest.NewRecorder()
	handler.ServeHTTP(bootstrap, httptest.NewRequest("GET", "http://localhost/api/session", nil))
	var session map[string]string
	json.Unmarshal(bootstrap.Body.Bytes(), &session)
	body := `{"revision":0,"activeEnvironment":"dev","variables":[],"environments":[{"id":"dev","name":"开发环境","baseUrl":"http://localhost:3000","variables":[]}]}`
	for _, tc := range []struct {
		name, host, origin, token string
		code                      int
	}{
		{"no token", "localhost", "", "", 403}, {"cross origin", "localhost", "https://evil.example", session["token"], 403}, {"rebound host", "evil.example", "", session["token"], 403}, {"save", "localhost", "http://localhost", session["token"], 200}, {"stale revision", "localhost", "", session["token"], 409},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "http://"+tc.host+"/api/workspace", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Apicat-Token", tc.token)
			req.Header.Set("Origin", tc.origin)
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			if res.Code != tc.code {
				t.Fatalf("%d: %s", res.Code, res.Body.String())
			}
		})
	}
	for _, path := range []string{"/api/session", "/api/document", "/api/workspace"} {
		req := httptest.NewRequest(http.MethodGet, "http://evil.example"+path, nil)
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		if res.Code != 403 {
			t.Fatalf("rebound read %s: %d", path, res.Code)
		}
	}
}

func TestErrorMessagesCarryLocalizationValues(t *testing.T) {
	handler, err := newHandler(func() *model.APIDoc { return &model.APIDoc{} }, t.TempDir(), "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	sessionResponse := httptest.NewRecorder()
	handler.ServeHTTP(sessionResponse, httptest.NewRequest("GET", "http://localhost/api/session", nil))
	var session map[string]string
	if err := json.Unmarshal(sessionResponse.Body.Bytes(), &session); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("POST", "http://localhost/api/requests/execute", strings.NewReader(`{"method":"GET","url":"http://example.test/{{中文_token}}"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Apicat-Token", session["token"])
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	var payload struct {
		MessageKey string            `json:"messageKey"`
		Values     map[string]string `json:"values"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if response.Code != 400 || payload.MessageKey != "变量 {name} 未定义或未启用" || payload.Values["name"] != "中文_token" {
		t.Fatalf("unexpected localized error: %d %s", response.Code, response.Body.String())
	}
}
