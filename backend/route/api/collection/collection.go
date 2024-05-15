package collection

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/iteration"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	collection_proto "github.com/apicat/apicat/v2/backend/route/proto/collection"
	collectionbase "github.com/apicat/apicat/v2/backend/route/proto/collection/base"
	collectionrequest "github.com/apicat/apicat/v2/backend/route/proto/collection/request"
	collectionresponse "github.com/apicat/apicat/v2/backend/route/proto/collection/response"
	"github.com/apicat/apicat/v2/backend/service/ai"
	"github.com/apicat/apicat/v2/backend/service/reference"
	"github.com/apicat/apicat/v2/backend/service/relations"
	"github.com/apicat/apicat/v2/backend/utils/onetime_token"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/plugin/export"
	"github.com/apicat/apicat/v2/backend/module/spec/plugin/openapi"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type collectionApiImpl struct{}

func NewCollectionApi() collection_proto.CollectionApi {
	return &collectionApiImpl{}
}

func (cai *collectionApiImpl) Create(ctx *gin.Context, opt *collectionrequest.CreateCollectionOption) (*collectionresponse.Collection, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if opt.ParentID != 0 {
		parentC := &collection.Collection{ID: opt.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := parentC.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "parentC.Get", "err", err)
			if opt.Type == collection.CategoryType {
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.CreationFailed"))
			}
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CreationFailed"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
	}

	c := &collection.Collection{
		ProjectID: selfPM.ProjectID,
		ParentID:  opt.ParentID,
		Title:     opt.Title,
		Type:      opt.Type,
		Content:   opt.Content,
	}
	if err := c.Create(ctx, selfTM); err != nil {
		slog.ErrorContext(ctx, "c.Create", "err", err)
		if c.Type == collection.CategoryType {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.CreationFailed"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CreationFailed"))
	}

	if opt.IterationID != "" {
		i := &iteration.Iteration{ID: opt.IterationID}
		exist, err := i.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "i.Get", "err", err)
			if c.Type == collection.CategoryType {
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.CreationFailed"))
			}
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CreationFailed"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("iteration.DoesNotExist"))
		}

		if err := i.BatchCreateCollection(ctx, []*collection.Collection{c}); err != nil {
			slog.ErrorContext(ctx, "i.BatchCreateCollection", "err", err)
			if c.Type == collection.CategoryType {
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.CreationFailed"))
			}
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CreationFailed"))
		}
	}

	userInfo := jwt.GetUser(ctx)
	return convertModelCollection(ctx, c, userInfo, userInfo), nil
}

func (cai *collectionApiImpl) List(ctx *gin.Context, opt *collectionrequest.GetCollectionListOption) (*collectionresponse.CollectionTree, error) {
	selfP := access.GetSelfProject(ctx)
	collections, err := collection.GetCollections(ctx, selfP.ID)
	if err != nil {
		slog.ErrorContext(ctx, "collection.GetCollections", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToGetList"))
	}

	var tree collectionresponse.CollectionTree
	if opt.IterationID == "" {
		tree = buildTree(0, collections, nil)
	} else {
		i := &iteration.Iteration{ID: opt.IterationID}
		exist, err := i.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "i.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToGetList"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("iteration.DoesNotExist"))
		}

		cIDs, err := i.GetCollectionIDs(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "i.GetCollectionIDs", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToGetList"))
		}

		tree = buildTree(0, collections, cIDs)
	}

	return &tree, nil
}

func (cai *collectionApiImpl) Get(ctx *gin.Context, opt *collectionbase.ProjectCollectionIDOption) (*collectionresponse.Collection, error) {
	selfP := access.GetSelfProject(ctx)

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: selfP.ID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		if c.Type == collection.CategoryType {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToGet"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToGet"))
	}
	if !exist {
		if c.Type == collection.CategoryType {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.DoesNotExist"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.DoesNotExist"))
	}

	cUserInfo, err := collection.UserInfo(ctx, c.CreatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "c.CreatedMember.UserInfo", "err", err)
		if c.Type == collection.CategoryType {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToGet"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToGet"))
	}
	uUserInfo, err := collection.UserInfo(ctx, c.UpdatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "c.UpdatedMember.UserInfo", "err", err)
		if c.Type == collection.CategoryType {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToGet"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToGet"))
	}

	return convertModelCollection(ctx, c, cUserInfo, uUserInfo), nil
}

func (cai *collectionApiImpl) Update(ctx *gin.Context, opt *collectionrequest.UpdateCollectionOption) (*ginrpc.Empty, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: selfPM.ProjectID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	if err := c.Update(ctx, opt.Title, opt.Content, selfTM.ID); err != nil {
		slog.ErrorContext(ctx, "c.Update", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	// 编辑文档时更新文档引用关系
	if c.Type != collection.CategoryType {
		if err := reference.UpdateCollectionRef(ctx, c); err != nil {
			slog.ErrorContext(ctx, "collectionrelations.UpdateCollectionRef", "err", err)
		}
	}

	return &ginrpc.Empty{}, nil
}

func (cai *collectionApiImpl) Delete(ctx *gin.Context, opt *collectionrequest.DeleteCollectionOption) (*ginrpc.Empty, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	selfTM := access.GetSelfTeamMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: selfPM.ProjectID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToDelete"))
	}
	if !exist {
		if c.Type == collection.CategoryType {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	if err := relations.DeleteCollections(ctx, selfPM.ProjectID, c, selfTM); err != nil {
		slog.ErrorContext(ctx, "relations.DeleteCollections", "err", err)
		if c.Type == collection.CategoryType {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToDelete"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToDelete"))
	}

	return &ginrpc.Empty{}, nil
}

func (cai *collectionApiImpl) Move(ctx *gin.Context, opt *collectionrequest.MoveCollectionOption) (*ginrpc.Empty, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if opt.Target.ParentID != 0 {
		targetParentC := &collection.Collection{ID: opt.Target.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := targetParentC.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "targetParentC.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToMove"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
		if targetParentC.Type != collection.CategoryType {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("category.IsNotCategory"))
		}
	}

	if opt.Origin.ParentID != opt.Target.ParentID && opt.Origin.ParentID != 0 {
		originParentC := &collection.Collection{ID: opt.Origin.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := originParentC.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "originParentC.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToMove"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
		if originParentC.Type != collection.CategoryType {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("category.IsNotCategory"))
		}
	}

	for i, cID := range opt.Target.IDs {
		c := &collection.Collection{ID: cID, ProjectID: selfPM.ProjectID}
		exist, err := c.Get(ctx)
		if err != nil || !exist {
			slog.ErrorContext(ctx, "c.Get", "err", err)
			continue
		}

		if err := c.Sort(ctx, opt.Target.ParentID, i+1); err != nil {
			slog.ErrorContext(ctx, "c.Sort", "err", err)
			continue
		}
	}

	if opt.Origin.ParentID != opt.Target.ParentID {
		for i, cID := range opt.Origin.IDs {
			c := &collection.Collection{ID: cID, ProjectID: selfPM.ProjectID}
			exist, err := c.Get(ctx)
			if err != nil || !exist {
				slog.ErrorContext(ctx, "c.Get", "err", err)
				continue
			}

			if err := c.Sort(ctx, opt.Origin.ParentID, i+1); err != nil {
				slog.ErrorContext(ctx, "c.Sort", "err", err)
				continue
			}
		}
	}

	return &ginrpc.Empty{}, nil
}

func (cai *collectionApiImpl) Copy(ctx *gin.Context, opt *collectionrequest.CopyCollectionOption) (*collectionresponse.Collection, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: selfPM.ProjectID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CopyFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	if c.Type == collection.CategoryType {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.CannotBeCopied"))
	}

	newC := &collection.Collection{
		ProjectID:    selfPM.ProjectID,
		ParentID:     c.ParentID,
		Title:        fmt.Sprintf("%s (copy)", c.Title),
		Type:         c.Type,
		Content:      c.Content,
		DisplayOrder: c.DisplayOrder,
	}
	if err := newC.Create(ctx, selfTM); err != nil {
		slog.ErrorContext(ctx, "newC.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CopyFailed"))
	}

	if opt.IterationID != "" {
		i := &iteration.Iteration{ID: opt.IterationID}
		exist, err = i.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "i.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CopyFailed"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("iteration.DoesNotExist"))
		}

		if err := i.BatchCreateCollection(ctx, []*collection.Collection{newC}); err != nil {
			slog.ErrorContext(ctx, "i.BatchCreateCollection", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CopyFailed"))
		}
	}

	userInfo := jwt.GetUser(ctx)
	return convertModelCollection(ctx, newC, userInfo, userInfo), nil
}

func (cai *collectionApiImpl) Trashes(ctx *gin.Context, opt *protobase.ProjectIdOption) (*collectionresponse.TrashList, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	// 获取回收站列表
	collections, err := collection.GetDeletedCollections(ctx, pm.ProjectID)
	if err != nil {
		slog.ErrorContext(ctx, "collection.GetDeletedCollections", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.FailedToGetList"))
	}

	list := make(collectionresponse.TrashList, len(collections))
	for i, c := range collections {
		if userInfo, err := collection.UserInfo(ctx, c.DeletedBy, true); err == nil {
			list[i] = &collectionresponse.Trash{
				CollectionIDOption: collectionbase.CollectionIDOption{
					CollectionID: c.ID,
				},
				CollectionTitle: c.Title,
				DeleteInfo: collectionresponse.DeleteInfo{
					DeletedAt: c.DeletedAt.Time,
					DeletedBy: userInfo.Name,
				},
			}
		}
	}
	return &list, nil
}

func (cai *collectionApiImpl) Restore(ctx *gin.Context, opt *collectionrequest.RestoreOption) (*collectionresponse.RestoreNum, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)

	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	deletedCollections, err := collection.GetDeletedCollectionsByIDs(ctx, selfPM.ProjectID, opt.CollectionIDs)
	if err != nil {
		slog.ErrorContext(ctx, "collection.GetDeletedCollectionsByIDs", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.RestoreFailed"))
	}
	if len(deletedCollections) == 0 {
		return &collectionresponse.RestoreNum{
			Num: 0,
		}, nil
	}

	restoreIDs, err := collection.RestoreCollections(ctx, selfTM, deletedCollections)
	if err != nil {
		slog.ErrorContext(ctx, "collection.RestoreCollections", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collectionHistory.RestoreFailed"))
	}

	if err := iteration.RestoreIterationApi(ctx, restoreIDs); err != nil {
		slog.ErrorContext(ctx, "iteration.RestoreIterationApi", "err", err)
	}

	return &collectionresponse.RestoreNum{
		Num: len(restoreIDs),
	}, nil
}

func (cai *collectionApiImpl) GetExportPath(ctx *gin.Context, opt *collectionrequest.GetExportPathOption) (*collectionresponse.ExportCollection, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	selfTM := access.GetSelfTeamMember(ctx)
	if selfPM.Permission.Equal(project.ProjectMemberRead) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: selfPM.ProjectID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.ExportFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}
	if c.Type == collection.CategoryType {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.CannotBeExported"))
	}

	tokenKey := fmt.Sprintf(
		"ExportCollection-%d-%d",
		selfTM.ID,
		time.Now().Unix(),
	)
	ca, err := cache.NewCache(config.Get().Cache.ToCfg())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.ExportFailed"))
	}
	token, err := onetime_token.NewTokenHelper(ca).GenerateToken(tokenKey, opt, time.Minute)
	if err != nil {
		slog.ErrorContext(ctx, "onetime_token.GenerateToken", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.ExportFailed"))
	}

	return &collectionresponse.ExportCollection{
		Path: fmt.Sprintf("/api/projects/%s/collections/%d/export/%s", selfPM.ProjectID, c.ID, token),
	}, nil
}

func (cai *collectionApiImpl) AIGenerate(ctx *gin.Context, opt *collectionrequest.AIGenerateCollectionOption) (*collectionresponse.Collection, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if opt.ParentID != 0 {
		parentC := &collection.Collection{ID: opt.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := parentC.Get(ctx)
		if err != nil {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.GenerationFailed"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
	}

	c, err := ai.DocGenerate(ctx, opt.Prompt)
	if err != nil {
		slog.ErrorContext(ctx, "ai.CreateAPI", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.GenerationFailed"))
	}

	c.ProjectID = selfPM.ProjectID
	c.ParentID = opt.ParentID
	if err := c.Create(ctx, selfTM); err != nil {
		slog.ErrorContext(ctx, "c.Create", "err", err)
		if c.Type == collection.CategoryType {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.GenerationFailed"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CreationFailed"))
	}

	if opt.IterationID != "" {
		i := &iteration.Iteration{ID: opt.IterationID}
		exist, err := i.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "i.Get", "err", err)
			if c.Type == collection.CategoryType {
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.CreationFailed"))
			}
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CreationFailed"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("iteration.DoesNotExist"))
		}

		if err := i.BatchCreateCollection(ctx, []*collection.Collection{c}); err != nil {
			slog.ErrorContext(ctx, "i.BatchCreateCollection", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.CreationFailed"))
		}
	}

	userInfo := jwt.GetUser(ctx)
	return convertModelCollection(ctx, c, userInfo, userInfo), nil
}

func Export(ctx *gin.Context) {
	// 解析和校验 URI 中的参数
	opt := &collectionrequest.ExportCodeOption{}
	if err := ctx.ShouldBindUri(&opt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ca, err := cache.NewCache(config.Get().Cache.ToCfg())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.NewErr("collection.ExportFailed").Error(),
		})
		return
	}
	tokenHelper := onetime_token.NewTokenHelper(ca)

	t := collectionrequest.GetExportPathOption{}
	if !tokenHelper.CheckToken(opt.Code, &t) {
		slog.ErrorContext(ctx, "onetime_token.CheckToken", "err", i18n.NewErr("collection.ExportFailed"))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.NewErr("collection.ExportFailed").Error(),
		})
		return
	}

	p := &project.Project{ID: ctx.Param("projectID")}
	exist, err := p.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "p.Get", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.NewErr("collection.ExportFailed").Error(),
		})
		return
	}
	if !exist {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": i18n.NewErr("collection.DoesNotExist").Error(),
		})
		return
	}

	if t.ProjectID != opt.ProjectID || t.CollectionID != opt.CollectionID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.NewErr("collection.ExportFailed").Error(),
		})
		return
	}

	c := &collection.Collection{ID: t.CollectionID, ProjectID: t.ProjectID}
	exist, err = c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.NewErr("collection.ExportFailed").Error(),
		})
		return
	}
	if !exist {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": i18n.NewErr("collection.DoesNotExist").Error(),
		})
		return
	}

	apicatData, err := relations.CollectionDerefWithApiCatSpec(ctx, c)
	if err != nil {
		slog.ErrorContext(ctx, "collectionDerefWithApiCatSpec", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.NewErr("collection.ExportFailed").Error(),
		})
		return
	}
	relations.SpecFillInfo(ctx, apicatData, p)
	relations.SpecFillServers(ctx, apicatData, p.ID)

	if apicatDataJson, err := json.Marshal(apicatData); err != nil {
		slog.ErrorContext(ctx, "export", "marshalErr", err)
	} else {
		slog.InfoContext(ctx, "export", "apicatData", apicatDataJson)
	}

	var (
		content []byte
	)
	switch t.Type {
	case "swagger":
		content, err = openapi.Generate(apicatData, "2.0", "json")
	case "openapi3.0.0":
		content, err = openapi.Generate(apicatData, "3.0.0", "json")
	case "openapi3.0.1":
		content, err = openapi.Generate(apicatData, "3.0.1", "json")
	case "openapi3.0.2":
		content, err = openapi.Generate(apicatData, "3.0.2", "json")
	case "openapi3.1.0":
		content, err = openapi.Generate(apicatData, "3.1.0", "json")
	case "HTML":
		content, err = export.HTML(apicatData)
	case "md":
		content, err = export.Markdown(apicatData)
	case "apicat":
		content, err = apicatData.ToJSON(spec.JSONOption{Indent: "  "})
	default:
		content, err = apicatData.ToJSON(spec.JSONOption{Indent: "  "})
	}
	if err != nil {
		slog.ErrorContext(ctx, "export", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.NewErr("collection.ExportFailed").Error(),
		})
		return
	}

	slog.InfoContext(ctx, "export", t.Type, content)

	switch t.Download {
	case true:
		filename := fmt.Sprintf("%s-%s", p.Title, t.Type)
		switch t.Type {
		case "HTML":
			ctx.Header("Content-Disposition", "attachment; filename="+filename+".html")
		case "md":
			ctx.Header("Content-Disposition", "attachment; filename="+filename+".md")
		default:
			ctx.Header("Content-Disposition", "attachment; filename="+filename+".json")
		}
		ctx.Data(http.StatusOK, "application/octet-stream", content)
	default:
		switch t.Type {
		case "HTML":
			ctx.Data(http.StatusOK, "text/html; charset=utf-8", content)
		case "md":
			ctx.Data(http.StatusOK, "text/markdown; charset=utf-8", content)
		default:
			ctx.Data(http.StatusOK, "application/json", content)
		}
	}

	tokenHelper.DelToken(opt.Code)
}
