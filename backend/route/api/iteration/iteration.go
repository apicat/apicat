package iteration

import (
	"log/slog"
	"math"
	"net/http"

	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/iteration"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/route/middleware/access"
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	protoiteration "github.com/apicat/apicat/backend/route/proto/iteration"
	iterationbase "github.com/apicat/apicat/backend/route/proto/iteration/base"
	iterationrequest "github.com/apicat/apicat/backend/route/proto/iteration/request"
	iterationresponse "github.com/apicat/apicat/backend/route/proto/iteration/response"
	projectbase "github.com/apicat/apicat/backend/route/proto/project/base"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type iterationApiImpl struct{}

func NewIterationApi() protoiteration.IterationApi {
	return &iterationApiImpl{}
}

func (iai *iterationApiImpl) Create(ctx *gin.Context, opt *iterationrequest.CreateIterationOption) (*ginrpc.Empty, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	p := &project.Project{ID: opt.ProjectID}
	exist, err := p.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "p.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.CreationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("project.DoesNotExist"))
	}
	pm := &project.ProjectMember{ProjectID: p.ID, MemberID: selfTM.ID}
	exist, err = pm.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "pm.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.CreationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("projectMember.DoesNotExist"))
	}
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	it := &iteration.Iteration{
		ProjectID:   p.ID,
		Title:       opt.Title,
		Description: opt.Description,
	}
	if err = it.Create(ctx, access.GetSelfTeamMember(ctx)); err != nil {
		slog.ErrorContext(ctx, "it.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.CreationFailed"))
	}

	if len(opt.CollectionIDs) > 0 {
		collections, err := collection.GetCollections(ctx, p, opt.CollectionIDs...)
		if err != nil {
			slog.ErrorContext(ctx, "collection.GetCollections", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.CreationFailed"))
		}
		if err := it.PlanningIterationApi(ctx, collections); err != nil {
			slog.ErrorContext(ctx, "it.PlanningIterationApi", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.CreationFailed"))
		}
	}
	return &ginrpc.Empty{}, nil
}

func (iai *iterationApiImpl) List(ctx *gin.Context, opt *iterationrequest.GetIterationListOption) (*iterationresponse.IterationList, error) {
	selfTeam := access.GetSelfTeam(ctx)
	selfTM := access.GetSelfTeamMember(ctx)
	if opt.Page <= 0 {
		opt.Page = 1
	}
	if opt.PageSize <= 0 {
		opt.PageSize = 15
	}

	projectIDs := make([]string, 0)
	if opt.ProjectID != "" {
		p := &project.Project{ID: opt.ProjectID}
		exist, err := p.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "p.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.FailedToGetList"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("project.DoesNotExist"))
		}
		projectIDs = append(projectIDs, p.ID)
	} else {
		projects, err := project.GetProjects(ctx, access.GetSelfTeamMember(ctx))
		if err != nil {
			slog.ErrorContext(ctx, "project.GetProjects", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.FailedToGetList"))
		}
		for _, p := range projects {
			projectIDs = append(projectIDs, p.ID)
		}
	}
	pms, err := project.GetInvolvedProjectMembers(ctx, selfTM)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetInvolvedProjectMembers", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.FailedToGetList"))
	}

	pmDict := map[string]*project.ProjectMember{}
	for _, v := range pms {
		pmDict[v.ProjectID] = v
	}

	iterations, err := iteration.GetIterations(ctx, selfTeam.ID, opt.Page, opt.PageSize, projectIDs...)
	if err != nil {
		slog.ErrorContext(ctx, "iteration.GetIterations", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.FailedToGetList"))
	}
	count, err := iteration.GetIterationsCount(ctx, selfTeam.ID, projectIDs...)
	if err != nil {
		slog.ErrorContext(ctx, "iteration.GetIterationsCount", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.FailedToGetList"))
	}

	list := &iterationresponse.IterationList{
		PaginationInfo: protobase.PaginationInfo{
			Count:       int(count),
			TotalPage:   int(math.Ceil(float64(count) / float64(opt.PaginationOption.PageSize))),
			CurrentPage: opt.PaginationOption.Page,
		},
		Items: make([]*iterationresponse.IterationListItem, len(iterations)),
	}
	for index, i := range iterations {
		if p, err := i.ProjectInfo(ctx); err == nil {
			count, err := i.GetIterationApiCount(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "i.GetIterationApiCount", "err", err)
			}

			permission := project.ProjectMemberNone
			if pm, ok := pmDict[i.ProjectID]; ok {
				permission = pm.Permission
			}

			list.Items[index] = &iterationresponse.IterationListItem{
				IdCreateTimeInfo: protobase.IdCreateTimeInfo{
					ID:        i.ID,
					CreatedAt: i.CreatedAt.Unix(),
				},
				IterationData: iterationbase.IterationData{
					Title:       i.Title,
					Description: i.Description,
				},
				ApisCount: count,
				Project: &iterationresponse.IterationListProject{
					ID:    p.ID,
					Title: p.Title,
					SelfMember: protobase.ProjectMemberPermission{
						Permission: permission,
					},
				},
			}
		}
	}

	return list, nil
}

func (iai *iterationApiImpl) Get(ctx *gin.Context, opt *iterationbase.IterationIDOption) (*iterationresponse.Iteration, error) {
	i := &iteration.Iteration{ID: opt.IterationID}
	exist, err := i.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "i.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.FailedToGet"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("iteration.DoesNotExist"))
	}

	pm := access.GetSelfProjectMember(ctx)
	p := access.GetSelfProject(ctx)
	return &iterationresponse.Iteration{
		IdCreateTimeInfo: protobase.IdCreateTimeInfo{
			ID:        i.ID,
			CreatedAt: i.CreatedAt.Unix(),
		},
		IterationData: iterationbase.IterationData{
			Title:       i.Title,
			Description: i.Description,
		},
		Project: &iterationresponse.IterationProject{
			OnlyIdInfo: protobase.OnlyIdInfo{
				ID: p.ID,
			},
			ProjectDataOption: projectbase.ProjectDataOption{
				Title:       p.Title,
				Cover:       p.Cover,
				Description: p.Description,
			},
			SelfMember: protobase.ProjectMemberPermission{
				Permission: pm.Permission,
			},
		},
	}, nil
}

func (iai *iterationApiImpl) Update(ctx *gin.Context, opt *iterationrequest.UpdateIterationOption) (*ginrpc.Empty, error) {
	i := &iteration.Iteration{ID: opt.IterationID}
	exist, err := i.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "i.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.FailedToGet"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("iteration.DoesNotExist"))
	}

	p := access.GetSelfProject(ctx)
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	i.Title = opt.Title
	i.Description = opt.Description
	if err := i.Update(ctx, access.GetSelfTeamMember(ctx)); err != nil {
		slog.ErrorContext(ctx, "i.Update", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	collections := make([]*collection.Collection, 0)
	if len(opt.CollectionIDs) > 0 {
		collections, err = collection.GetCollections(ctx, p, opt.CollectionIDs...)
		if err != nil {
			slog.ErrorContext(ctx, "collection.GetCollections", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
		}
	}

	if err := i.PlanningIterationApi(ctx, collections); err != nil {
		slog.ErrorContext(ctx, "selfI.PlanningIterationApi", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	return &ginrpc.Empty{}, nil
}

func (iai *iterationApiImpl) Delete(ctx *gin.Context, opt *iterationbase.IterationIDOption) (*ginrpc.Empty, error) {
	i := &iteration.Iteration{ID: opt.IterationID}
	exist, err := i.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "i.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.FailedToGet"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("iteration.DoesNotExist"))
	}

	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if err := i.Delete(ctx, access.GetSelfTeamMember(ctx)); err != nil {
		slog.ErrorContext(ctx, "selfI.Delete", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("iteration.FailedToDelete"))
	}
	return &ginrpc.Empty{}, nil
}
