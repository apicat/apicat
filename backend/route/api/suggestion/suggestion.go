package suggestion

import (
	"log/slog"
	"net/http"

	"github.com/apicat/apicat/v2/backend/core/content_suggestion"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"
	protosuggestion "github.com/apicat/apicat/v2/backend/route/proto/suggestion"
	suggestionreq "github.com/apicat/apicat/v2/backend/route/proto/suggestion/request"
	suggestionres "github.com/apicat/apicat/v2/backend/route/proto/suggestion/response"
	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type suggestionApiImpl struct{}

func NewSuggestionApi() protosuggestion.SuggestionApi {
	return &suggestionApiImpl{}
}

func (s *suggestionApiImpl) GenCollection(ctx *gin.Context, opt *suggestionreq.CollectionOption) (*suggestionres.CollectionSuggestion, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{ProjectID: opt.ProjectID, Title: opt.Title, Path: opt.Path}
	generator, err := content_suggestion.NewCollectionGenerator(c, jwt.GetUser(ctx).Language)
	if err != nil {
		slog.ErrorContext(ctx, "content_suggestion.NewCollectionGenerator", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.CollectionGenerationFailed"))
	}

	newC, err := generator.Generate()
	if err != nil {
		slog.ErrorContext(ctx, "generator.Generate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.CollectionGenerationFailed"))
	}

	contentSpec, err := newC.ToSpec()
	if err != nil {
		slog.ErrorContext(ctx, "newC.ToSpec", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.CollectionGenerationFailed"))
	}
	content, err := contentSpec.ToJson()
	if err != nil {
		slog.ErrorContext(ctx, "contentSpec.ToJson", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.CollectionGenerationFailed"))
	}
	return &suggestionres.CollectionSuggestion{
		RequestID: opt.RequestID,
		Content:   content,
	}, nil
}

func (s *suggestionApiImpl) GenModel(ctx *gin.Context, opt *suggestionreq.ModelOption) (*suggestionres.ModelSuggestion, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	m := &definition.DefinitionSchema{ProjectID: opt.ProjectID, Name: opt.Title}
	generator, err := content_suggestion.NewModelGenerator(m, jwt.GetUser(ctx).Language)
	if err != nil {
		slog.ErrorContext(ctx, "content_suggestion.NewModelGenerator", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.ModelGenerationFailed"))
	}

	newM, err := generator.Generate()
	if err != nil {
		slog.ErrorContext(ctx, "generator.Generate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.ModelGenerationFailed"))
	}
	schemaSpec, err := newM.ToSpec()
	if err != nil {
		slog.ErrorContext(ctx, "newM.ToSpec", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.ModelGenerationFailed"))
	}
	return &suggestionres.ModelSuggestion{
		RequestID: opt.RequestID,
		Schema:    schemaSpec.Schema.ToJson(),
	}, nil
}

func (s *suggestionApiImpl) GenSchema(ctx *gin.Context, opt *suggestionreq.SchemaOption) (*suggestionres.ModelSuggestion, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	js, err := jsonschema.NewSchemaFromJson(opt.Schema)
	if err != nil {
		slog.ErrorContext(ctx, "jsonschema.NewSchemaFromJson", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("suggestion.SchemaGenerationFailed"))
	}

	generator, err := content_suggestion.NewSchemaGenerator(selfPM.ProjectID, js)
	if err != nil {
		slog.ErrorContext(ctx, "content_suggestion.NewSchemaGenerator", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.SchemaGenerationFailed"))
	}

	newS, err := generator.Generate()
	if err != nil {
		slog.ErrorContext(ctx, "generator.Generate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.SchemaGenerationFailed"))
	}
	if newS == nil {
		return &suggestionres.ModelSuggestion{
			RequestID: opt.RequestID,
			Schema:    "",
		}, nil
	}
	return &suggestionres.ModelSuggestion{
		RequestID: opt.RequestID,
		Schema:    newS.ToJson(),
	}, nil
}

func (s *suggestionApiImpl) GenReference(ctx *gin.Context, opt *suggestionreq.SchemaOption) (*suggestionres.ModelSuggestion, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	js, err := jsonschema.NewSchemaFromJson(opt.Schema)
	if err != nil {
		slog.ErrorContext(ctx, "jsonschema.NewSchemaFromJson", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("suggestion.SchemaGenerationFailed"))
	}

	generator, err := content_suggestion.NewReferenceMatcher(ctx, selfPM.ProjectID)
	if err != nil {
		slog.ErrorContext(ctx, "content_suggestion.NewReferenceMatcher", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.ReferenceGenerationFailed"))
	}

	newS, err := generator.Match(opt.Title, js, opt.ModelID)
	if err != nil {
		slog.ErrorContext(ctx, "generator.Match", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("suggestion.ReferenceGenerationFailed"))
	}
	if newS == nil {
		return &suggestionres.ModelSuggestion{
			RequestID: opt.RequestID,
			Schema:    "",
		}, nil
	}
	return &suggestionres.ModelSuggestion{
		RequestID: opt.RequestID,
		Schema:    newS.ToJson(),
	}, nil
}
