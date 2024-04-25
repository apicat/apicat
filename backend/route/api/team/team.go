package team

import (
	"log/slog"

	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	prototeam "github.com/apicat/apicat/v2/backend/route/proto/team"
	prototeambase "github.com/apicat/apicat/v2/backend/route/proto/team/base"
	prototeamrequest "github.com/apicat/apicat/v2/backend/route/proto/team/request"
	prototeamresponse "github.com/apicat/apicat/v2/backend/route/proto/team/response"
	"github.com/apicat/apicat/v2/backend/service/team_relations"

	"net/http"

	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type teamApiImpl struct{}

func NewTeamApi() prototeam.TeamApi {
	return &teamApiImpl{}
}

// Create 创建团队
func (t *teamApiImpl) Create(ctx *gin.Context, opt *prototeambase.TeamDataOption) (*prototeamresponse.Team, error) {
	tm, err := team.Create(ctx, jwt.GetUser(ctx), opt.Name)
	if err != nil {
		slog.ErrorContext(ctx, "team.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("team.CreationFailed"))
	}
	if err := switchTeam(ctx, tm.ID); err != nil {
		slog.ErrorContext(ctx, "switchTeam", "err", err)
	}
	return team_relations.ConvertModelTeam(ctx, tm), nil
}

// TeamList 团队列表
func (t *teamApiImpl) TeamList(ctx *gin.Context, opt *prototeamrequest.RolesOption) (*prototeamresponse.TeamList, error) {
	list, err := team.GetUserTeams(ctx, jwt.GetUser(ctx).ID, opt.Roles...)
	if err != nil {
		slog.ErrorContext(ctx, "team.GetUserTeams", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("team.FailedToGetList"))
	}
	var resp prototeamresponse.TeamList
	resp.Items = make([]*prototeamresponse.Team, len(list))
	for k, v := range list {
		resp.Items[k] = team_relations.ConvertModelTeam(ctx, v)
	}
	return &resp, nil
}

// Current 当前团队
func (t *teamApiImpl) Current(ctx *gin.Context, opt *ginrpc.Empty) (*prototeamresponse.CurrentTeamRes, error) {
	selfUser := jwt.GetUser(ctx)

	tm, err := team.GetLastActiveTeam(ctx, selfUser.ID)
	if err != nil {
		slog.ErrorContext(ctx, "team.GetLastActiveTeam", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("team.FailedToGetCurrentTeam"),
		)
	}
	if tm == nil {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("team.NoTeam"),
		)
	}

	teamInstance := &team.Team{ID: tm.TeamID}
	exist, err := teamInstance.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "team.Get", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("team.FailedToGetCurrentTeam"),
		)
	}
	if !exist {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("team.CurrentTeamDoesNotExist"),
		)
	}

	return &prototeamresponse.CurrentTeamRes{
		Team: *team_relations.ConvertModelTeam(ctx, teamInstance),
		Role: tm.Role,
	}, nil
}

// Get 团队详情
func (t *teamApiImpl) Get(ctx *gin.Context, opt *protobase.TeamIdOption) (*prototeamresponse.Team, error) {
	teamRecord, err := team.GetTeam(ctx, opt.TeamID)
	if err != nil {
		slog.ErrorContext(ctx, "team.GetTeam", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("team.DoesNotExist"),
		)
	}
	return team_relations.ConvertModelTeam(ctx, teamRecord), nil
}

// Switch 切换当前团队
func (t *teamApiImpl) Switch(ctx *gin.Context, opt *protobase.TeamIdOption) (*ginrpc.Empty, error) {
	if err := switchTeam(ctx, opt.TeamID); err != nil {
		slog.ErrorContext(ctx, "switchTeam", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("team.FailedToSwitch"))
	}
	return &ginrpc.Empty{}, nil
}

// GetInvitationToken 获取邀请token
func (t *teamApiImpl) GetInvitationToken(ctx *gin.Context, opt *protobase.TeamIdOption) (*protobase.InvitationTokenOption, error) {
	selfMember := access.GetSelfTeamMember(ctx)
	if selfMember.Role.Lower(team.RoleAdmin) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	return &protobase.InvitationTokenOption{
		InvitationToken: selfMember.InvitationToken,
	}, nil
}

// ResetInvitationToken 重置团队邀请token
func (t *teamApiImpl) ResetInvitationToken(ctx *gin.Context, opt *protobase.TeamIdOption) (*protobase.InvitationTokenOption, error) {
	selfMember := access.GetSelfTeamMember(ctx)
	if selfMember.Role.Lower(team.RoleAdmin) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if err := selfMember.ResetInvitationToken(ctx); err != nil {
		slog.ErrorContext(ctx, "team.GetTeam", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("team.InvitationTokenResetFailed"))
	}
	return &protobase.InvitationTokenOption{
		InvitationToken: selfMember.InvitationToken,
	}, nil
}

// CheckInvitationToken 检查邀请链接有效性
func (t *teamApiImpl) CheckInvitationToken(ctx *gin.Context, opt *protobase.InvitationTokenOption) (*prototeamresponse.TeamInviteContent, error) {
	errResp := ginrpc.NewError(
		http.StatusBadRequest,
		i18n.NewErr("common.LinkExpired"),
	)

	if opt.InvitationToken == "" {
		errResp.Attrs = map[string]any{
			"emoji":       "😳",
			"title":       i18n.NewTran("common.LinkExpiredTitle").Translate(ctx),
			"description": i18n.NewTran("team.InvitationTokenNotFound").Translate(ctx),
		}
		return nil, errResp
	}

	tm, err := team.GetMemberByToken(ctx, opt.InvitationToken)
	if err != nil {
		errResp.Attrs = map[string]any{
			"emoji":       "😳",
			"title":       i18n.NewTran("common.LinkExpiredTitle").Translate(ctx),
			"description": i18n.NewTran("team.InvalidInvitationToken").Translate(ctx),
		}
		return nil, errResp
	}
	userInfo, err := tm.UserInfo(ctx, false)
	if err != nil {
		errResp.Attrs = map[string]any{
			"emoji":       "😳",
			"title":       i18n.NewTran("common.LinkExpiredTitle").Translate(ctx),
			"description": i18n.NewTran("team.InvalidInvitationToken").Translate(ctx),
		}
		return nil, errResp
	}

	currentTeam := &team.Team{ID: tm.TeamID}
	exist, err := currentTeam.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "team.Get", "err", err)
		errResp.Code = http.StatusInternalServerError
		errResp.Attrs = map[string]any{
			"emoji":       "😳",
			"title":       i18n.NewTran("common.LinkExpiredTitle").Translate(ctx),
			"description": i18n.NewTran("team.InvalidInvitationToken").Translate(ctx),
		}
		return nil, errResp
	}
	if !exist {
		errResp.Attrs = map[string]any{
			"emoji":       "😳",
			"title":       i18n.NewTran("common.LinkExpiredTitle").Translate(ctx),
			"description": i18n.NewTran("team.InvalidInvitationToken").Translate(ctx),
		}
		return nil, errResp
	}

	return &prototeamresponse.TeamInviteContent{
		Inviter: userInfo.Name,
		Team:    currentTeam.Name,
	}, nil
}

// Join 加入团队 需要邀请码
func (t *teamApiImpl) Join(ctx *gin.Context, opt *protobase.InvitationTokenOption) (*ginrpc.Empty, error) {
	self := jwt.GetUser(ctx)
	if err := team_relations.JoinTeam(ctx, opt.InvitationToken, self); err != nil {
		return nil, ginrpc.NewError(http.StatusBadRequest, err)
	}
	return &ginrpc.Empty{}, nil
}

// Setting 设置团队信息
func (t *teamApiImpl) Setting(ctx *gin.Context, opt *prototeamrequest.SettingOption) (*ginrpc.Empty, error) {
	selfTeam := access.GetSelfTeam(ctx)
	selfMember := access.GetSelfTeamMember(ctx)
	if !selfMember.Role.Equal(team.RoleOwner) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	// 目前只能修改名称
	selfTeam.Name = opt.Name
	if err := selfTeam.Update(ctx); err != nil {
		slog.ErrorContext(ctx, "team.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	return &ginrpc.Empty{}, nil
}

// Transfer 团队转让
func (t *teamApiImpl) Transfer(ctx *gin.Context, opt *prototeamrequest.GetTeamMemberOption) (*ginrpc.Empty, error) {
	selfTeam := access.GetSelfTeam(ctx)
	selfMember := access.GetSelfTeamMember(ctx)
	if !selfMember.Role.Equal(team.RoleOwner) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}
	targetMember, err := team.GetMember(ctx, opt.MemberID)
	if err != nil {
		slog.ErrorContext(ctx, "team.GetMember", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("teamMember.NotInTheTeam"))
	}

	if targetMember.TeamID != opt.TeamID || selfMember.TeamID != targetMember.TeamID {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("teamMember.TeamTransferFailed"),
		)
	}

	if !targetMember.Role.Equal(team.RoleAdmin) {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("teamMember.TeamTransferInvalidMember"),
		)
	}

	if targetMember.Status != team.MemberStatusActive {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("teamMember.Deactivated"),
		)
	}

	if err := selfTeam.Transfer(ctx, selfMember, targetMember); err != nil {
		slog.ErrorContext(ctx, "selfTeam.Transfer", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("teamMember.TeamTransferFailed"),
		)
	}

	return &ginrpc.Empty{}, nil
}

// Delete 删除团队
func (t *teamApiImpl) Delete(ctx *gin.Context, opt *protobase.TeamIdOption) (*ginrpc.Empty, error) {
	selfTeam := access.GetSelfTeam(ctx)
	selfMember := access.GetSelfTeamMember(ctx)
	if !selfMember.Role.Equal(team.RoleOwner) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}
	if err := selfTeam.Delete(ctx); err != nil {
		slog.ErrorContext(ctx, "selfTeam.Delete", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("common.DeletionFailed"),
		)
	}
	return &ginrpc.Empty{}, nil
}
