package suggestion

import (
	suggestionreq "github.com/apicat/apicat/v2/backend/route/proto/suggestion/request"
	suggestionres "github.com/apicat/apicat/v2/backend/route/proto/suggestion/response"
	"github.com/gin-gonic/gin"
)

type SuggestionApi interface {
	// Collection suggestion
	// @route POST /projects/{projectID}/suggestion/collection
	GenCollection(*gin.Context, *suggestionreq.CollectionOption) (*suggestionres.CollectionSuggestion, error)

	// Model suggestion
	// @route POST /projects/{projectID}/suggestion/model
	GenModel(*gin.Context, *suggestionreq.ModelOption) (*suggestionres.ModelSuggestion, error)

	// Schema suggestion
	// @route POST /projects/{projectID}/suggestion/schema
	GenSchema(*gin.Context, *suggestionreq.SchemaOption) (*suggestionres.ModelSuggestion, error)

	// Reference suggestion
	// @route POST /projects/{projectID}/suggestion/reference
	GenReference(*gin.Context, *suggestionreq.RefOption) (*suggestionres.ModelSuggestion, error)
}
