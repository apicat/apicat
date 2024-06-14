package project

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	protoproject "github.com/apicat/apicat/v2/backend/route/proto/project"
	projectrequest "github.com/apicat/apicat/v2/backend/route/proto/project/request"
	projectresponse "github.com/apicat/apicat/v2/backend/route/proto/project/response"

	"github.com/apicat/apicat/v2/backend/module/spec/diff"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type definitionSchemaHistoryApiImpl struct{}

func NewDefinitionSchemaHistoryApi() protoproject.DefinitionSchemaHistoryAPi {
	return &definitionSchemaHistoryApiImpl{}
}

func (impl *definitionSchemaHistoryApiImpl) List(ctx *gin.Context, opt *projectrequest.GetDefinitionSchemaHistoryListOption) (*projectresponse.DefinitionSchemaHistoryList, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	ds := &definition.DefinitionSchema{
		ID:        opt.SchemaID,
		ProjectID: selfP.ID,
	}
	exist, err := ds.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ds.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.FailedToGetList"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchema.DoesNotExist"))
	}

	startTime := time.Unix(opt.StartTime, 0)
	endTime := time.Unix(opt.EndTime, 0)
	list, err := definition.GetDefinitionSchemaHistories(ctx, ds, startTime, endTime)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetDefinitionSchemaHistories", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.FailedToGetList"))
	}

	res := make(projectresponse.DefinitionSchemaHistoryList, len(list))
	for i, v := range list {
		if userInfo, err := definition.UserInfo(ctx, v.CreatedBy, true); err == nil {
			res[i] = &projectresponse.DefinitionSchemaHistoryItem{
				IdCreateTimeInfo: protobase.IdCreateTimeInfo{
					ID:        v.ID,
					CreatedAt: v.CreatedAt.Unix(),
				},
				CreatedBy: userInfo.Name,
			}
		}
	}

	return &res, nil
}

func (impl *definitionSchemaHistoryApiImpl) Get(ctx *gin.Context, opt *projectrequest.DefinitionSchemaHistoryIDOption) (*projectresponse.DefinitionSchemaHistory, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	ds := &definition.DefinitionSchema{
		ID:        opt.SchemaID,
		ProjectID: selfP.ID,
	}
	exist, err := ds.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ds.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.FailedToGet"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchema.DoesNotExist"))
	}

	h := &definition.DefinitionSchemaHistory{
		ID:       opt.HistoryID,
		SchemaID: ds.ID,
	}
	exist, err = h.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "h.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.FailedToGet"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchemaHistory.DoesNotExist"))
	}

	userInfo, err := definition.UserInfo(ctx, h.CreatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "h.CreatedMember.UserInfo", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.FailedToGet"))

	}

	return convertModelDefinitionSchemaHistory(h, userInfo), nil
}

func (impl *definitionSchemaHistoryApiImpl) Restore(ctx *gin.Context, opt *projectrequest.DefinitionSchemaHistoryIDOption) (*ginrpc.Empty, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	ds := &definition.DefinitionSchema{
		ID:        opt.SchemaID,
		ProjectID: selfPM.ProjectID,
	}
	exist, err := ds.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ds.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.RestoreFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchema.DoesNotExist"))
	}

	h := &definition.DefinitionSchemaHistory{
		ID:       opt.HistoryID,
		SchemaID: ds.ID,
	}
	exist, err = h.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "h.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.RestoreFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchemaHistory.DoesNotExist"))
	}

	if err := h.Restore(ctx, ds, selfTM.ID); err != nil {
		slog.ErrorContext(ctx, "h.Restore", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.RestoreFailed"))
	}

	return &ginrpc.Empty{}, nil
}

func (impl *definitionSchemaHistoryApiImpl) Diff(ctx *gin.Context, opt *projectrequest.DiffDefinitionSchemaHistoriesOption) (*projectresponse.DiffDefinitionSchemaHistories, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	ds := &definition.DefinitionSchema{
		ID:        opt.SchemaID,
		ProjectID: selfP.ID,
	}
	exist, err := ds.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ds.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchema.DoesNotExist"))
	}
	if ds.Type == definition.SchemaCategory {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
	}

	originalDSH := &definition.DefinitionSchemaHistory{
		ID:       opt.OriginalID,
		SchemaID: ds.ID,
	}
	exist, err = originalDSH.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "originalDSH.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchemaHistory.DoesNotExist"))
	}
	originalDSHUserInfo, err := definition.UserInfo(ctx, originalDSH.CreatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "originalDSH.CreatedMember.UserInfo", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
	}

	var targetDSH *definition.DefinitionSchemaHistory
	if opt.TargetID == 0 {
		targetDSH = &definition.DefinitionSchemaHistory{
			SchemaID:    ds.ID,
			Name:        ds.Name,
			Description: ds.Description,
			Schema:      ds.Schema,
			CreatedBy:   ds.UpdatedBy,
			TimeModel: model.TimeModel{
				CreatedAt: ds.UpdatedAt,
				UpdatedAt: ds.UpdatedAt,
			},
		}
	} else {
		targetDSH = &definition.DefinitionSchemaHistory{
			ID:       opt.TargetID,
			SchemaID: ds.ID,
		}
		exist, err = targetDSH.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "targetDSH.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchemaHistory.DoesNotExist"))
		}
	}
	targetDSHUserInfo, err := definition.UserInfo(ctx, targetDSH.CreatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "targetDSH.CreatedMember.UserInfo", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
	}

	// 对模型进行解引用
	ds.Schema = originalDSH.Schema
	originalSchema, err := dsDerefWithSpec(ctx, ds)
	if err != nil {
		slog.ErrorContext(ctx, "original.dsDerefWithSpec", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
	}

	ds.Schema = targetDSH.Schema
	targetSchema, err := dsDerefWithSpec(ctx, ds)
	if err != nil {
		slog.ErrorContext(ctx, "target.dsDerefWithSpec", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
	}

	if err := diff.DiffModel(originalSchema, targetSchema); err != nil {
		slog.ErrorContext(ctx, "diff.DiffSchema", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
	}

	originalSchemaStr, err := json.Marshal(originalSchema.Schema)
	if err != nil {
		slog.ErrorContext(ctx, "originalSchema.json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
	}
	originalDSH.Schema = string(originalSchemaStr)

	targetSchemaStr, err := json.Marshal(targetSchema.Schema)
	if err != nil {
		slog.ErrorContext(ctx, "targetSchema.json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchemaHistory.DiffFailed"))
	}
	targetDSH.Schema = string(targetSchemaStr)

	return &projectresponse.DiffDefinitionSchemaHistories{
		Schema1: convertModelDefinitionSchemaHistory(originalDSH, originalDSHUserInfo),
		Schema2: convertModelDefinitionSchemaHistory(targetDSH, targetDSHUserInfo),
	}, nil
}
