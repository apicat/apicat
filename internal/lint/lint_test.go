package lint

import (
	"strings"
	"testing"

	"github.com/apicat/apicat/v3/internal/model"
)

func TestAnnotateCleanSpec(t *testing.T) {
	doc := &model.APIDoc{
		Description: "documented",
		TagGroups: []model.TagGroup{{
			Endpoints: []model.Endpoint{{
				Method:      "GET",
				Path:        "/pets",
				Summary:     "List pets",
				OperationID: "listPets",
				Parameters: []model.Parameter{
					{Name: "limit", Description: "max results"},
				},
				Responses: []model.Response{
					{StatusCode: "200", Content: []model.MediaTypeContent{{
						MediaType: "application/json",
						Schema: &model.SchemaView{Type: "object", Properties: []model.SchemaView{
							{Name: "id", Description: "pet id"},
						}},
					}}},
					{StatusCode: "default"},
				},
			}},
		}},
	}

	Annotate(doc)

	if doc.LintTotal != 0 {
		t.Errorf("LintTotal = %d, want 0; counts = %+v", doc.LintTotal, doc.LintCounts)
	}
	if hints := doc.TagGroups[0].Endpoints[0].LintHints; len(hints) != 0 {
		t.Errorf("LintHints = %v, want none", hints)
	}
}

func TestAnnotateFindsProblems(t *testing.T) {
	doc := &model.APIDoc{
		// no API description
		TagGroups: []model.TagGroup{{
			Endpoints: []model.Endpoint{{
				Method: "POST",
				Path:   "/pets",
				// no summary/description, no operationId
				Parameters: []model.Parameter{
					{Name: "verbose"}, // no description
					{Name: "limit", Description: "ok"},
				},
				RequestBody: &model.RequestBody{Content: []model.MediaTypeContent{{
					MediaType: "application/json",
					Schema: &model.SchemaView{Type: "object", Properties: []model.SchemaView{
						{Name: "id"},                            // no description
						{Name: "name", Description: "the name"}, // ok
						{Name: "owner", Properties: []model.SchemaView{{Name: "email"}}}, // nested: owner + email
					}},
				}}},
				Responses: []model.Response{{StatusCode: "200"}}, // no error responses
			}},
		}},
	}

	Annotate(doc)

	ep := doc.TagGroups[0].Endpoints[0]
	joined := strings.Join(ep.LintHints, "; ")
	for _, want := range []string{
		"no summary or description",
		"no operationId",
		"no error responses",
		`parameter "verbose" has no description`,
		"3 schema fields without description",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("LintHints missing %q; got: %s", want, joined)
		}
	}
	if len(ep.LintHints) != 5 {
		t.Errorf("len(LintHints) = %d, want 5: %v", len(ep.LintHints), ep.LintHints)
	}

	// api(1) + summary(1) + operationId(1) + errors(1) + param(1) + fields(3)
	if doc.LintTotal != 8 {
		t.Errorf("LintTotal = %d, want 8; counts = %+v", doc.LintTotal, doc.LintCounts)
	}

	counts := map[string]int{}
	for _, c := range doc.LintCounts {
		counts[c.Rule] = c.Count
	}
	if counts["field-no-description"] != 3 {
		t.Errorf("field-no-description = %d, want 3", counts["field-no-description"])
	}
	if counts["param-no-description"] != 1 {
		t.Errorf("param-no-description = %d, want 1", counts["param-no-description"])
	}
}

func TestAnnotateIsIdempotent(t *testing.T) {
	doc := &model.APIDoc{
		TagGroups: []model.TagGroup{{
			Endpoints: []model.Endpoint{{Method: "GET", Path: "/x"}},
		}},
	}
	Annotate(doc)
	first := doc.LintTotal
	Annotate(doc)
	if doc.LintTotal != first {
		t.Errorf("LintTotal changed on re-annotate: %d -> %d", first, doc.LintTotal)
	}
	if n := len(doc.TagGroups[0].Endpoints[0].LintHints); n != 3 {
		t.Errorf("LintHints accumulated on re-annotate: %d", n)
	}
}
