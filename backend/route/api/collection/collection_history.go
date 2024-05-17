package collection

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/module/spec/diff"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	protocollection "github.com/apicat/apicat/v2/backend/route/proto/collection"
	collectionbase "github.com/apicat/apicat/v2/backend/route/proto/collection/base"
	collectionrequest "github.com/apicat/apicat/v2/backend/route/proto/collection/request"
	collectionresponse "github.com/apicat/apicat/v2/backend/route/proto/collection/response"
	"github.com/apicat/apicat/v2/backend/service/relations"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type collectionHistoryApiImpl struct{}

func NewCollectionHistoryApi() protocollection.CollectionHistoryApi {
	return &collectionHistoryApiImpl{}
}

func (srv *collectionHistoryApiImpl) List(ctx *gin.Context, opt *collectionrequest.GetCollectionHistoryListOption) (*collectionresponse.CollectionHistoryList, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{
		ID:        opt.CollectionID,
		ProjectID: selfP.ID,
	}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.FailedToGetList"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	startTime := time.Unix(opt.StartTime, 0)
	endTime := time.Unix(opt.EndTime, 0)
	list, err := collection.GetCollectionHistories(ctx, c, startTime, endTime)
	if err != nil {
		slog.ErrorContext(ctx, "collection.GetCollectionHistories", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.FailedToGetList"))
	}

	protoList := make(collectionresponse.CollectionHistoryList, len(list))
	for i, ch := range list {
		if userInfo, err := collection.UserInfo(ctx, ch.CreatedBy, true); err == nil {
			protoList[i] = &collectionresponse.CollectionHistoryItem{
				IdCreateTimeInfo: protobase.IdCreateTimeInfo{
					ID:        ch.ID,
					CreatedAt: ch.CreatedAt.Unix(),
				},
				CollectionTypeOption: collectionbase.CollectionTypeOption{
					Type: c.Type,
				},
				CreatedBy: userInfo.Name,
			}
		}
	}

	return &protoList, nil
}

func (srv *collectionHistoryApiImpl) Get(ctx *gin.Context, opt *collectionrequest.CollectionHistoryIDOption) (*collectionresponse.CollectionHistory, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{
		ID:        opt.CollectionID,
		ProjectID: selfP.ID,
	}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.FailedToGet"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	ch := &collection.CollectionHistory{
		ID:           opt.HistoryID,
		CollectionID: c.ID,
	}
	exist, err = ch.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ch.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.FailedToGet"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collectionHistory.DoesNotExist"))
	}

	userInfo, err := collection.UserInfo(ctx, ch.CreatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "ch.CreatedMember.UserInfo", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.FailedToGet"))
	}

	return convertModelCollectionHistory(ch, userInfo), nil
}

func (srv *collectionHistoryApiImpl) Restore(ctx *gin.Context, opt *collectionrequest.CollectionHistoryIDOption) (*ginrpc.Empty, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{
		ID:        opt.CollectionID,
		ProjectID: selfPM.ProjectID,
	}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.RestoreFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	ch := &collection.CollectionHistory{
		ID:           opt.HistoryID,
		CollectionID: c.ID,
	}
	exist, err = ch.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ch.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.RestoreFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collectionHistory.DoesNotExist"))
	}

	if err := ch.Restore(ctx, c, selfTM); err != nil {
		slog.ErrorContext(ctx, "ch.Restore", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.RestoreFailed"))
	}

	return &ginrpc.Empty{}, nil
}

func (srv *collectionHistoryApiImpl) Diff(ctx *gin.Context, opt *collectionrequest.DiffCollectionHistoriesOption) (*collectionresponse.DiffCollectionHistories, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{
		ID:        opt.CollectionID,
		ProjectID: selfP.ID,
	}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.DiffFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}
	if c.Type == collection.CategoryType {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("collectionHistory.DiffFailed"))
	}

	originalCH := &collection.CollectionHistory{
		ID:           opt.OriginalID,
		CollectionID: c.ID,
	}
	exist, err = originalCH.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ch1.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.DiffFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collectionHistory.DoesNotExist"))
	}

	originalCHUserInfo, err := collection.UserInfo(ctx, originalCH.CreatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "originalCH.CreatedMember.UserInfo", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.DiffFailed"))
	}

	var targetCH *collection.CollectionHistory
	if opt.TargetID == 0 {
		targetCH = &collection.CollectionHistory{
			CollectionID: c.ID,
			Title:        c.Title,
			Content:      c.Content,
			CreatedBy:    c.UpdatedBy,
			TimeModel: model.TimeModel{
				CreatedAt: c.UpdatedAt,
				UpdatedAt: c.UpdatedAt,
			},
		}
	} else {
		targetCH = &collection.CollectionHistory{
			ID:           opt.TargetID,
			CollectionID: c.ID,
		}
		exist, err = targetCH.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "ch2.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.DiffFailed"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collectionHistory.DoesNotExist"))
		}
	}

	targetCHUserInfo, err := collection.UserInfo(ctx, targetCH.CreatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "targetCH.CreatedMember.UserInfo", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.DiffFailed"))
	}

	// 对文档进行解引用
	c.Content = originalCH.Content
	originalDoc, err := relations.CollectionDerefWithSpec(ctx, c)
	if err != nil {
		slog.ErrorContext(ctx, "original.collectionDerefWithSpec", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.DiffFailed"))
	}
	c.Content = targetCH.Content
	targetDoc, err := relations.CollectionDerefWithSpec(ctx, c)
	if err != nil {
		slog.ErrorContext(ctx, "target.collectionDerefWithSpec", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.DiffFailed"))
	}

	if err := diff.Diff(originalDoc, targetDoc); err != nil {
		slog.ErrorContext(ctx, "diff.Diff", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.DiffFailed"))
	}

	originalContentStr, err := json.Marshal(originalDoc.Content)
	if err != nil {
		slog.ErrorContext(ctx, "originalDoc.json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.DiffFailed"))
	}
	originalCH.Content = string(originalContentStr)

	targetContentStr, err := json.Marshal(targetDoc.Content)
	if err != nil {
		slog.ErrorContext(ctx, "targetDoc.json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.DiffFailed"))
	}
	targetCH.Content = string(targetContentStr)

	return &collectionresponse.DiffCollectionHistories{
		Doc1: convertModelCollectionHistory(originalCH, originalCHUserInfo),
		Doc2: convertModelCollectionHistory(targetCH, targetCHUserInfo),
	}, nil
}
