package project

import (
	"log/slog"
	"math"
	"net/http"

	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	protoproject "github.com/apicat/apicat/v2/backend/route/proto/project"
	projectrequest "github.com/apicat/apicat/v2/backend/route/proto/project/request"
	projectresponse "github.com/apicat/apicat/v2/backend/route/proto/project/response"
	prototeamresponse "github.com/apicat/apicat/v2/backend/route/proto/team/response"
	"github.com/apicat/apicat/v2/backend/service/team_relations"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type projectMemberApiImpl struct{}

func NewProjectMemberApi() protoproject.ProjectMemberApi {
	return &projectMemberApiImpl{}
}

// Create 创建项目成员
func (pgai *projectMemberApiImpl) Create(ctx *gin.Context, opt *projectrequest.CreateProjectMemberOption) (*ginrpc.Empty, error) {
	t := access.GetSelfTeam(ctx)
	p := access.GetSelfProject(ctx)
	pm := access.GetSelfProjectMember(ctx)

	if pm.Permission.Lower(project.ProjectMemberManage) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	if opt.Permission == project.ProjectMemberManage {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("projectMember.CanNotAddProjectManager"),
		)
	}

	tms, err := team.GetMembers(ctx, t.ID, 0, 0, "")
	if err != nil {
		slog.ErrorContext(ctx, "team.GetMembers", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectMember.FailedToAddProjectMember"))
	}
	var teamMembers []*team.TeamMember
	for _, tmID := range opt.MemberIDs {
		for _, tm := range tms {
			if tm.ID == tmID {
				teamMembers = append(teamMembers, tm)
			}
		}
	}

	project.BatchCreateMember(ctx, p.ID, teamMembers, opt.Permission)
	return &ginrpc.Empty{}, nil
}

// Members 获取项目成员列表
func (pgai *projectMemberApiImpl) Members(ctx *gin.Context, opt *projectrequest.GetProjectMemberListOption) (*projectresponse.GetProjectMemberListResponse, error) {
	p := access.GetSelfProject(ctx)
	if opt.PaginationOption.Page <= 0 {
		opt.PaginationOption.Page = 1
	}
	if opt.PaginationOption.PageSize <= 0 {
		opt.PaginationOption.PageSize = 15
	}
	pms, err := project.GetProjectMembers(ctx, p.ID, opt.PaginationOption.Page, opt.PaginationOption.PageSize)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetProjectMembers", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectMember.FailedToGetList"))
	}

	count, err := project.GetProjectMembersCount(ctx, p.ID)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetProjectMembersCount", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectMember.FailedToGetList"))
	}

	list := &projectresponse.GetProjectMemberListResponse{
		PaginationInfo: protobase.PaginationInfo{
			Count:       int(count),
			TotalPage:   int(math.Ceil(float64(count) / float64(opt.PaginationOption.PageSize))),
			CurrentPage: opt.PaginationOption.Page,
		},
		Items: make(projectresponse.ProjectMembers, len(pms)),
	}
	for i, pm := range pms {
		if memberInfo, err := pm.MemberInfo(ctx, false); err == nil {
			if userInfo, err := memberInfo.UserInfo(ctx, false); err == nil {
				list.Items[i] = convertModelProjectMember(pm, memberInfo, userInfo)
			}
		}
	}

	return list, nil
}

// NotInProjectMembers 不在此项目的成员列表
func (pgai *projectMemberApiImpl) NotInProjectMembers(ctx *gin.Context, opt *protobase.ProjectIdOption) (*prototeamresponse.TeamMembers, error) {
	p := access.GetSelfProject(ctx)
	t := access.GetSelfTeam(ctx)
	tms, err := team.GetMembers(ctx, t.ID, 0, 0, team.MemberStatusActive)
	if err != nil {
		slog.ErrorContext(ctx, "team.GetMembers", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectMember.FailedToGetList"))
	}

	pms, err := project.GetProjectMembers(ctx, p.ID, 0, 0)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetProjectMembers", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectMember.FailedToGetList"))
	}
	projectMemberMap := make(map[uint]*project.ProjectMember)
	for _, pm := range pms {
		projectMemberMap[pm.MemberID] = pm
	}

	res := make(prototeamresponse.TeamMembers, 0)
	for _, tm := range tms {
		if _, exist := projectMemberMap[tm.ID]; !exist {
			if userInfo, err := tm.UserInfo(ctx, false); err == nil {
				res = append(res, team_relations.ConvertModelTeamMember(ctx, tm, userInfo))
			}
		}
	}

	return &res, nil
}

// Update 更新项目成员
func (pgai *projectMemberApiImpl) Update(ctx *gin.Context, opt *projectrequest.UpdateProjectMemberOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberManage) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	targetPm := project.ProjectMember{ID: opt.MemberID}
	exist, err := targetPm.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "targetPm.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("projectMember.NotInTheProject"),
		)
	}

	if targetPm.ProjectID != pm.ProjectID {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("projectMember.NotInTheProject"),
		)
	}

	if opt.Permission == project.ProjectMemberManage {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("projectMember.CanNotAddProjectManager"),
		)
	}

	targetPm.Permission = opt.Permission
	if err := targetPm.Update(ctx); err != nil {
		slog.ErrorContext(ctx, "targetPm.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	return &ginrpc.Empty{}, nil
}

// Delete 删除项目成员
func (pgai *projectMemberApiImpl) Delete(ctx *gin.Context, opt *projectrequest.ProjectMemberIDOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberManage) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	targetPm := project.ProjectMember{ID: opt.MemberID}
	exist, err := targetPm.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "targetPm.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("projectMember.RemoveFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("projectMember.NotInTheProject"),
		)
	}

	if targetPm.ProjectID != pm.ProjectID {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("projectMember.NotInTheProject"),
		)
	}

	if targetPm.Permission == project.ProjectMemberManage {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	if err := targetPm.Delete(ctx); err != nil {
		slog.ErrorContext(ctx, "targetPm.Delete", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("projectMember.RemoveFailed"))
	}

	return &ginrpc.Empty{}, nil
}
