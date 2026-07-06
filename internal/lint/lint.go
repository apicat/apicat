// Package lint walks the converted view model and annotates it with spec
// quality hints: endpoints get per-endpoint messages (rendered as a badge
// on the detail page), the document gets per-rule counts (rendered as a
// summary on the overview page). Field-level "missing description" marks
// are rendered directly by the templates from emptiness checks; this
// package only aggregates.
package lint

import (
	"fmt"

	"github.com/apicat/apicat/v3/internal/model"
)

const (
	ruleAPINoDescription      = "api-no-description"
	ruleEndpointNoSummary     = "endpoint-no-summary"
	ruleEndpointNoOperationID = "endpoint-no-operation-id"
	ruleEndpointNoErrors      = "endpoint-no-error-response"
	ruleParamNoDescription    = "param-no-description"
	ruleFieldNoDescription    = "field-no-description"
)

// ruleMessages doubles as the fixed display order of the summary.
var ruleOrder = []struct {
	rule    string
	message string
}{
	{ruleAPINoDescription, "API has no description"},
	{ruleEndpointNoSummary, "endpoints without summary or description"},
	{ruleEndpointNoOperationID, "endpoints without operationId"},
	{ruleEndpointNoErrors, "endpoints without error responses (4xx/5xx/default)"},
	{ruleParamNoDescription, "parameters without description"},
	{ruleFieldNoDescription, "schema fields without description"},
}

// Annotate fills Endpoint.LintHints and the document-level LintCounts /
// LintTotal. Safe to call on every reload; it overwrites previous results.
func Annotate(doc *model.APIDoc) {
	counts := map[string]int{}

	if doc.Description == "" {
		counts[ruleAPINoDescription]++
	}

	for gi := range doc.TagGroups {
		eps := doc.TagGroups[gi].Endpoints
		for ei := range eps {
			annotateEndpoint(&eps[ei], counts)
		}
	}

	doc.LintCounts = doc.LintCounts[:0]
	doc.LintTotal = 0
	for _, r := range ruleOrder {
		if n := counts[r.rule]; n > 0 {
			doc.LintCounts = append(doc.LintCounts, model.LintCount{
				Rule:    r.rule,
				Message: r.message,
				Count:   n,
			})
			doc.LintTotal += n
		}
	}
}

func annotateEndpoint(ep *model.Endpoint, counts map[string]int) {
	ep.LintHints = nil

	if ep.Summary == "" && ep.Description == "" {
		ep.LintHints = append(ep.LintHints, "no summary or description")
		counts[ruleEndpointNoSummary]++
	}
	if ep.OperationID == "" {
		ep.LintHints = append(ep.LintHints, "no operationId")
		counts[ruleEndpointNoOperationID]++
	}
	if !hasErrorResponse(ep.Responses) {
		ep.LintHints = append(ep.LintHints, "no error responses (4xx/5xx/default)")
		counts[ruleEndpointNoErrors]++
	}

	for _, p := range ep.Parameters {
		if p.Description == "" {
			ep.LintHints = append(ep.LintHints, fmt.Sprintf("parameter %q has no description", p.Name))
			counts[ruleParamNoDescription]++
		}
	}

	fields := 0
	if ep.RequestBody != nil {
		for _, mt := range ep.RequestBody.Content {
			fields += countUndescribedFields(mt.Schema)
		}
	}
	for _, resp := range ep.Responses {
		for _, mt := range resp.Content {
			fields += countUndescribedFields(mt.Schema)
		}
	}
	if fields > 0 {
		noun := "fields"
		if fields == 1 {
			noun = "field"
		}
		ep.LintHints = append(ep.LintHints, fmt.Sprintf("%d schema %s without description", fields, noun))
		counts[ruleFieldNoDescription] += fields
	}
}

func hasErrorResponse(responses []model.Response) bool {
	for _, r := range responses {
		if r.StatusCode == "default" {
			return true
		}
		if len(r.StatusCode) > 0 && (r.StatusCode[0] == '4' || r.StatusCode[0] == '5') {
			return true
		}
	}
	return false
}

// countUndescribedFields counts named properties without a description,
// recursively. Container nodes (array wrappers, composition branches) are
// unnamed and not counted themselves, only walked.
func countUndescribedFields(sv *model.SchemaView) int {
	if sv == nil {
		return 0
	}
	n := 0
	if sv.Name != "" && sv.Description == "" {
		n++
	}
	for i := range sv.Properties {
		n += countUndescribedFields(&sv.Properties[i])
	}
	n += countUndescribedFields(sv.Items)
	for i := range sv.AllOf {
		n += countUndescribedFields(&sv.AllOf[i])
	}
	for i := range sv.OneOf {
		n += countUndescribedFields(&sv.OneOf[i])
	}
	for i := range sv.AnyOf {
		n += countUndescribedFields(&sv.AnyOf[i])
	}
	return n
}
