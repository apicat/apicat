package team

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	teambase "github.com/apicat/apicat/backend/route/proto/team/base"
	teamrequest "github.com/apicat/apicat/backend/route/proto/team/request"
	teamresponse "github.com/apicat/apicat/backend/route/proto/team/response"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

// TeamApi 团队相关
type TeamApi interface {
	// Create 创建团队
	// @route POST /teams
	Create(*gin.Context, *teambase.TeamDataOption) (*teamresponse.Team, error)

	// TeamList 团队列表
	// @route GET /teams
	TeamList(*gin.Context, *teamrequest.RolesOption) (*teamresponse.TeamList, error)

	// Current 当前团队
	// @route GET /teams/current
	Current(*gin.Context, *ginrpc.Empty) (*teamresponse.CurrentTeamRes, error)

	// Get 我的单个团队
	// @route GET /teams/{teamID}
	Get(*gin.Context, *protobase.TeamIdOption) (*teamresponse.Team, error)

	// Switch 切换当前团队
	// @route PUT /teams/{teamID}/switch
	Switch(*gin.Context, *protobase.TeamIdOption) (*ginrpc.Empty, error)

	// GetInvitationToken 获取邀请token
	// @route GET /teams/{teamID}/invitation_tokens
	GetInvitationToken(*gin.Context, *protobase.TeamIdOption) (*protobase.InvitationTokenOption, error)

	// ResetInvitationToken 重置团队邀请token
	// @route PUT /teams/{teamID}/invitation_tokens
	ResetInvitationToken(*gin.Context, *protobase.TeamIdOption) (*protobase.InvitationTokenOption, error)

	// CheckInvitationToken 检查团队邀请token
	// @route GET /teams/invite/check
	CheckInvitationToken(*gin.Context, *protobase.InvitationTokenOption) (*teamresponse.TeamInviteContent, error)

	// Join 加入团队 需要邀请码
	// @route POST /team/join
	Join(*gin.Context, *protobase.InvitationTokenOption) (*ginrpc.Empty, error)

	// Setting 团队设置
	// @route PUT /teams/{teamID}/setting
	Setting(*gin.Context, *teamrequest.SettingOption) (*ginrpc.Empty, error)

	// Transfer 团队转让
	// @route PUT /team/{teamID}/transfer
	Transfer(*gin.Context, *teamrequest.GetTeamMemberOption) (*ginrpc.Empty, error)

	// Delete 删除团队
	// @route DELETE /team/{teamID}
	Delete(*gin.Context, *protobase.TeamIdOption) (*ginrpc.Empty, error)
}

type TeamMemberApi interface {
	// MemberList 团队成员
	// @route GET /teams/{teamID}/members
	MemberList(*gin.Context, *teamrequest.MembersOption) (*teamresponse.TeamMemberList, error)

	// UpdateMember 编辑团队成员 目前仅能编辑权限
	// @route PUT /teams/{teamID}/members/{memberID}
	UpdateMember(*gin.Context, *teamrequest.UpdateTeamMemberOption) (*teamresponse.TeamMember, error)

	// DeleteMember 删除团队成员
	// @route DELETE /teams/{teamID}/members/{memberID}
	DeleteMember(*gin.Context, *teamrequest.GetTeamMemberOption) (*ginrpc.Empty, error)

	// Quit 退出团队
	// @route DELETE /teams/{teamID}/members
	Quit(*gin.Context, *protobase.TeamIdOption) (*ginrpc.Empty, error)
}
