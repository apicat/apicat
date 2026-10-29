package converter

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/apicat/apicat/v3/internal/model"
)

func TestEndpointRouteIDs(t *testing.T) {
	hash := sha256.Sum256([]byte("GET /missing"))
	reserved := fmt.Sprintf("op-%x", hash[:8])
	endpoints := []model.Endpoint{
		{Method: "GET", Path: "/articles/{id}", OperationID: "getArticleById"},
		{Method: "GET", Path: "/missing"},
		{Method: "GET", Path: "/duplicate-a", OperationID: "duplicate"},
		{Method: "POST", Path: "/duplicate-b", OperationID: "duplicate"},
		{Method: "GET", Path: "/reserved", OperationID: reserved},
		{Method: "GET", Path: "/unsafe", OperationID: "../bad/id%20"},
	}
	doc := &model.APIDoc{TagGroups: []model.TagGroup{{Endpoints: endpoints}, {Endpoints: []model.Endpoint{endpoints[0]}}}}
	assignRouteIDs(doc)
	seen := map[string]bool{}
	byPath := map[string]string{}
	for _, ep := range doc.TagGroups[0].Endpoints {
		if !routeIDPattern.MatchString(ep.RouteID) || seen[ep.RouteID] {
			t.Fatalf("invalid or duplicate route: %+v", ep)
		}
		seen[ep.RouteID] = true
		byPath[ep.Path] = ep.RouteID
	}
	if byPath["/articles/{id}"] != "getArticleById" || byPath["/reserved"] != reserved {
		t.Fatal("unique operationIds were not preserved", byPath)
	}
	if doc.TagGroups[1].Endpoints[0].RouteID != "getArticleById" {
		t.Fatal("multiple tags changed endpoint identity")
	}
	if byPath["/duplicate-a"] == "duplicate" || byPath["/duplicate-b"] == "duplicate" {
		t.Fatal("ambiguous operationId was used")
	}
	// Reordering the document must not change generated links.
	for i, j := 0, len(endpoints)-1; i < j; i, j = i+1, j-1 {
		endpoints[i], endpoints[j] = endpoints[j], endpoints[i]
	}
	doc.TagGroups = doc.TagGroups[:1]
	assignRouteIDs(doc)
	for _, ep := range doc.TagGroups[0].Endpoints {
		if ep.RouteID != byPath[ep.Path] {
			t.Fatalf("route changed after reordering: %+v", ep)
		}
	}
}
