package runner

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/apicat/apicat/v3/internal/workspace"
)

func TestExecuteEnvironmentAndParameters(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/v1/pets/a%2Fb" {
			t.Errorf("path = %s", r.URL.EscapedPath())
		}
		if got := r.URL.Query()["tag"]; len(got) != 2 || got[0] != "a b" || got[1] != "second" {
			t.Errorf("query: %v", got)
		}
		if r.Header.Get("Authorization") != "Bearer local" {
			t.Errorf("auth = %s", r.Header.Get("Authorization"))
		}
		b, _ := io.ReadAll(r.Body)
		if string(b) != `{"name":"local"}` {
			t.Errorf("body = %s", b)
		}
		w.Header().Add("X-Value", "one")
		w.Header().Add("X-Value", "two")
		w.WriteHeader(201)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer target.Close()
	state := workspace.State{Variables: []workspace.Variable{{Key: "token", Value: "global", Enabled: true}}, Environments: []workspace.Environment{{ID: "test", BaseURL: target.URL + "/v1", Variables: []workspace.Variable{{Key: "token", Value: "local", Enabled: true}}}}}
	input := Request{Method: "POST", URL: "/pets/{id}", EnvironmentID: "test", Path: []Pair{{"id", "a/b", true}}, Query: []Pair{{"tag", "a b", true}, {"tag", "second", true}, {"ignored", "yes", false}}, Headers: []Pair{{"Authorization", "Bearer {{token}}", true}}, Body: `{"name":"{{token}}"}`}
	response, err := Execute(context.Background(), input, state)
	if err != nil {
		t.Fatal(err)
	}
	if response.Status != 201 || response.Body != `{"ok":true}` || len(response.Headers["X-Value"]) != 2 {
		t.Fatalf("response: %+v", response)
	}
	if input.Body != `{"name":"{{token}}"}` {
		t.Fatal("mutated input")
	}
}
func TestInvalidRequests(t *testing.T) {
	for _, tc := range []Request{{Method: "GET", URL: "file:///etc/passwd"}, {Method: "GET", URL: "http://example.test/{{missing}}"}, {Method: "GET", URL: "http://example.test/{id}"}, {Method: "CONNECT", URL: "http://example.test"}, {Method: "GET", URL: "/relative"}, {Method: "GET", URL: "http://example.test", EnvironmentID: "missing"}, {Method: "GET", URL: "http://example.test", Headers: []Pair{{"X-Header", "a\nb", true}}}} {
		if _, err := Build(context.Background(), tc, workspace.State{}); err == nil {
			t.Errorf("accepted %+v", tc)
		}
	}
}
func TestRedirectLimitBinaryAndCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/redirect":
			http.Redirect(w, r, "/never", 302)
		case "/never":
			t.Error("followed redirect")
		case "/large":
			w.Write([]byte(strings.Repeat("x", MaxResponseBytes+10)))
		case "/binary":
			w.Write([]byte{0, 255, 1})
		case "/slow":
			<-r.Context().Done()
		}
	}))
	defer server.Close()
	for _, path := range []string{"/redirect", "/large", "/binary"} {
		response, err := Execute(context.Background(), Request{Method: "GET", URL: server.URL + path}, workspace.State{})
		if err != nil {
			t.Fatal(err)
		}
		if path == "/redirect" && response.Status != 302 {
			t.Fatal(response.Status)
		}
		if path == "/large" && (!response.Truncated || response.Size != MaxResponseBytes) {
			t.Fatal("response limit")
		}
		if path == "/binary" && (response.Encoding != "base64" || response.Body != "AP8B") {
			t.Fatalf("binary: %+v", response)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if _, err := Execute(ctx, Request{Method: "GET", URL: server.URL + "/slow"}, workspace.State{}); err == nil {
		t.Fatal("ignored cancellation")
	}
}

func TestResolvedQueryIsNotExpandedAgain(t *testing.T) {
	input := Request{Method: "GET", URL: "http://example.test", QueryResolved: true, Query: []Pair{{Key: "value", Value: "{{literal}}", Enabled: true}}}
	req, err := Build(context.Background(), input, workspace.State{})
	if err != nil || req.URL.Query().Get("value") != "{{literal}}" {
		t.Fatalf("request %v, error %v", req, err)
	}
}
