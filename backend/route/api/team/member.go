package team

import (
	"log/slog"
	"math"
	"net/http"

	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	prototeam "github.com/apicat/apicat/v2/backend/route/proto/team"
	prototeamrequest "github.com/apicat/apicat/v2/backend/route/proto/team/request"
	prototeamresponse "github.com/apicat/apicat/v2/backend/route/proto/team/response"
	"github.com/apicat/apicat/v2/backend/service/team_relations"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type teamMemberApiImpl struct{}

func NewTeamMemberApi() prototeam.TeamMemberApi {
	return &teamMemberApiImpl{}
}

// ListMembers 团队成员
func (t *teamMemberApiImpl) MemberList(ctx *gin.Context, opt *prototeamrequest.MembersOption) (*prototeamresponse.TeamMemberList, error) {
	if opt.ListOption.Page <= 0 {
		opt.ListOption.Page = 1
	}
	if opt.ListOption.PageSize <= 0 {
		opt.ListOption.PageSize = 15
	}

	items, err := team.GetMembers(ctx, opt.TeamID, opt.ListOption.Page, opt.ListOption.PageSize, opt.Status, opt.Roles...)
	if err != nil {
		slog.ErrorContext(ctx, "team.GetMembers", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("teamMember.FailedToGetList"))
	}

	count, err := team.GetMembersCount(ctx, opt.TeamID, opt.Roles...)
	if err != nil {
		slog.ErrorContext(ctx, "team.GetMembersCount", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("teamMember.FailedToGetList"))
	}

	var list = &prototeamresponse.TeamMemberList{
		PaginationInfo: protobase.PaginationInfo{
			Count:       int(count),
			TotalPage:   int(math.Ceil(float64(count) / float64(opt.ListOption.PageSize))),
			CurrentPage: opt.ListOption.Page,
		},
		Items: make([]*prototeamresponse.TeamMember, len(items)),
	}
	for k, v := range items {
		if userInfo, err := v.UserInfo(ctx, false); err == nil {
			list.Items[k] = team_relations.ConvertModelTeamMember(ctx, v, userInfo)
		}
	}
	return list, nil
}

// UpdateMember 编辑团队成员
func (t *teamMemberApiImpl) UpdateMember(ctx *gin.Context, opt *prototeamrequest.UpdateTeamMemberOption) (*prototeamresponse.TeamMember, error) {
	if opt.Role == "" && opt.Status == "" {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.OperationFailed"))
	}

	selfMember := access.GetSelfTeamMember(ctx)
	if selfMember.ID == opt.MemberID {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.DoItToSelf"))
	}
	targetMember, err := team.GetMember(ctx, opt.MemberID)
	if err != nil {
		slog.ErrorContext(ctx, "team.GetMember", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.OperationFailed"))
	}
	if targetMember.TeamID != opt.TeamID || selfMember.TeamID != targetMember.TeamID {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.OperationFailed"))
	}
	userInfo, err := targetMember.UserInfo(ctx, false)
	if err != nil {
		slog.ErrorContext(ctx, "targetMember.UserInfo", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.OperationFailed"))
	}

	var needModify bool
	// 如果目标权限和新权限不一致才会更新权限
	if opt.Role != "" && targetMember.Role != opt.Role {
		if !selfMember.Role.Equal(team.RoleOwner) {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.PermissionDenied"))
		}
		targetMember.Role = opt.Role
		needModify = true
	}
	if opt.Status != "" && targetMember.Status != opt.Status {
		if selfMember.Role.LowerOrEqual(targetMember.Role) {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.PermissionDenied"))
		}
		targetMember.Status = opt.Status
		needModify = true
	}

	if needModify {
		if err := targetMember.Update(ctx); err != nil {
			slog.ErrorContext(ctx, "targetMember.Update", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.OperationFailed"))
		}
	}

	return team_relations.ConvertModelTeamMember(ctx, targetMember, userInfo), nil
}

// DeleteMember 删除团队成员
func (t *teamMemberApiImpl) DeleteMember(ctx *gin.Context, opt *prototeamrequest.GetTeamMemberOption) (*ginrpc.Empty, error) {
	selfMember := access.GetSelfTeamMember(ctx)
	if selfMember.ID == opt.MemberID {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("teamMember.RemoveFailed"),
		)
	}

	target, err := team.GetMember(ctx, opt.MemberID)
	if err != nil {
		slog.ErrorContext(ctx, "team.GetMember", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("teamMember.DoesNotExist"))
	}

	if target.TeamID != opt.TeamID || selfMember.TeamID != target.TeamID {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("teamMember.RemoveFailed"),
		)
	}

	if selfMember.Role.LowerOrEqual(target.Role) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("teamMember.RemoveFailed"),
		)
	}

	if err := target.Quit(ctx); err != nil {
		slog.ErrorContext(ctx, "target.Quit", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("teamMember.RemoveFailed"))
	}

	return &ginrpc.Empty{}, nil
}

// Quit 退出团队
func (t *teamMemberApiImpl) Quit(ctx *gin.Context, opt *protobase.TeamIdOption) (*ginrpc.Empty, error) {
	selfMember := access.GetSelfTeamMember(ctx)
	if selfMember.Role == team.RoleOwner {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("teamMember.CanNotQuitOwnTeam"),
		)
	}

	if err := selfMember.Quit(ctx); err != nil {
		slog.ErrorContext(ctx, "selfMember.Quit", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("teamMember.TeamQuitFailed"))
	}

	return &ginrpc.Empty{}, nil
}
