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

type projectGroupApiImpl struct{}

func NewProjectGroupApi() protoproject.ProjectGroupApi {
	return &projectGroupApiImpl{}
}

// Create 创建项目分组
func (pgai *projectGroupApiImpl) Create(ctx *gin.Context, opt *projectrequest.CreateProjectGroupOption) (*projectresponse.ProjectGroup, error) {
	tm := access.GetSelfTeamMember(ctx)

	pg := &project.ProjectGroup{
		MemberID: tm.ID,
		Name:     opt.Name,
	}
	exist, err := pg.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "pg.Get", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("projectGroup.CreationFailed"),
		)
	}
	if exist {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("projectGroup.NameHasBeenUsed"),
		)
	}

	pg, err = pg.Create(ctx, opt.Name, tm)
	if err != nil {
		slog.ErrorContext(ctx, "pg.Create", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("projectGroup.CreationFailed"),
		)
	}

	return convertModelProjectGroup(pg), nil
}

// List 获取项目分组列表
func (pgai *projectGroupApiImpl) List(ctx *gin.Context, opt *protobase.TeamIdOption) (*projectresponse.ProjectGroups, error) {
	tm := access.GetSelfTeamMember(ctx)

	pgs, err := project.GetProjectGroups(ctx, tm.ID)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetProjectGroups", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("projectGroup.FailedToGetList"),
		)
	}

	list := make(projectresponse.ProjectGroups, len(pgs))
	for i, v := range pgs {
		list[i] = convertModelProjectGroup(v)
	}

	return &list, nil
}

// Delete 删除项目分组
func (pgai *projectGroupApiImpl) Delete(ctx *gin.Context, opt *projectrequest.ProjectGroupIdOption) (*ginrpc.Empty, error) {
	tm := access.GetSelfTeamMember(ctx)
	pg := &project.ProjectGroup{ID: opt.GroupID}
	exist, err := pg.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "pg.Get", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("projectGroup.FailedToDelete"),
		)
	}
	if !exist {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("projectGroup.DoesNotExist"),
		)
	}

	if pg.MemberID != tm.ID {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("projectGroup.DoesNotExist"),
		)
	}

	if err := pg.Delete(ctx); err != nil {
		slog.ErrorContext(ctx, "pg.Delete", "err", err)
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("projectGroup.FailedToDelete"),
		)
	}

	return &ginrpc.Empty{}, nil
}

// Rename 重命名项目分组
func (pgai *projectGroupApiImpl) Rename(ctx *gin.Context, opt *projectrequest.RenameProjectGroupOption) (*ginrpc.Empty, error) {
	tm := access.GetSelfTeamMember(ctx)
	pg := &project.ProjectGroup{ID: opt.GroupID}
	exist, err := pg.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "pg.Get", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("projectGroup.FailedToDelete"),
		)
	}
	if !exist {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("projectGroup.DoesNotExist"),
		)
	}

	if pg.MemberID != tm.ID {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("projectGroup.DoesNotExist"),
		)
	}

	newPg := &project.ProjectGroup{MemberID: tm.ID, Name: opt.Name}
	exits, err := newPg.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "pg.Get", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("projectGroup.RenameFailed"),
		)
	}
	if exits {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("projectGroup.NameHasBeenUsed"),
		)
	}

	pg.Name = opt.Name
	if err := pg.Rename(ctx); err != nil {
		slog.ErrorContext(ctx, "pg.Rename", "err", err)
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("projectGroup.RenameFailed"),
		)
	}

	return &ginrpc.Empty{}, nil
}

// Sort 项目URL排序
func (pgai *projectGroupApiImpl) Sort(ctx *gin.Context, opt *projectrequest.SortProjectGroupOption) (*ginrpc.Empty, error) {
	tm := access.GetSelfTeamMember(ctx)

	if err := project.GroupSort(ctx, tm.ID, opt.GroupIDs); err != nil {
		slog.ErrorContext(ctx, "project.GroupSort", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("projectGroup.SortingFailed"),
		)
	}

	return &ginrpc.Empty{}, nil
}
