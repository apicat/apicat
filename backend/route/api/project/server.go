package project

import (
	"apicat-cloud/backend/i18n"
	"apicat-cloud/backend/model/project"
	"apicat-cloud/backend/route/middleware/access"
	protobase "apicat-cloud/backend/route/proto/base"
	protoproject "apicat-cloud/backend/route/proto/project"
	projectrequest "apicat-cloud/backend/route/proto/project/request"
	projectresponse "apicat-cloud/backend/route/proto/project/response"
	"log/slog"

	"net/http"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type projectServerApiImpl struct{}

func NewProjectServerApi() protoproject.ProjectServerApi {
	return &projectServerApiImpl{}
}

// Create 创建项目URL
func (psai *projectServerApiImpl) Create(ctx *gin.Context, opt *projectrequest.CreateProjectServerOption) (*projectresponse.ProjectServer, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	ps := project.Server{
		ProjectID:   pm.ProjectID,
		URL:         opt.URL,
		Description: opt.Description,
	}
	extis, err := ps.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ps.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectServer.CreationFailed"))
	}
	if extis {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("projectServer.HasBeenUsed"),
		)
	}

	psm, err := ps.Create(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ps.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectServer.CreationFailed"))
	}

	return convertModelProjectServer(psm), nil
}

// List 获取项目URL列表
func (psai *projectServerApiImpl) List(ctx *gin.Context, opt *protobase.ProjectIdOption) (*projectresponse.ProjectServerList, error) {
	p := access.GetSelfProject(ctx)

	list, err := project.GetServers(ctx, p.ID)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetServers", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectServer.FailedToGetList"))
	}

	ret := make(projectresponse.ProjectServerList, len(list))
	for i, v := range list {
		ret[i] = convertModelProjectServer(v)
	}

	return &ret, nil
}

// Update 修改项目URL
func (psai *projectServerApiImpl) Update(ctx *gin.Context, opt *projectrequest.UpdateProjectServerOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	ps := project.Server{
		ID:          opt.ServerID,
		ProjectID:   pm.ProjectID,
		URL:         opt.URL,
		Description: opt.Description,
	}
	exits, err := ps.CheckRepeat(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ps.CheckRepeat", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}
	if exits {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("projectServer.HasBeenUsed"))
	}

	if err := ps.Update(ctx); err != nil {
		slog.ErrorContext(ctx, "ps.Update", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	return &ginrpc.Empty{}, nil
}

// Delete 删除项目URL
func (psai *projectServerApiImpl) Delete(ctx *gin.Context, opt *projectrequest.GetProjectServerOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	ps := &project.Server{ID: opt.ServerID}
	if err := ps.Delete(ctx); err != nil {
		slog.ErrorContext(ctx, "ps.Delete", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectServer.FailedToDelete"))
	}

	return &ginrpc.Empty{}, nil
}

// Sort 项目URL排序
func (psai *projectServerApiImpl) Sort(ctx *gin.Context, opt *projectrequest.SortProjectServerOpt) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	if err := project.ServerSort(ctx, pm.ProjectID, opt.ServerIDs); err != nil {
		slog.ErrorContext(ctx, "project.ServerSort", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectServer.SortingFailed"))
	}

	return &ginrpc.Empty{}, nil
}
