package relations

import (
	"log/slog"

	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/model/user"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	prototeambase "github.com/apicat/apicat/v2/backend/route/proto/team/base"
	prototeamresponse "github.com/apicat/apicat/v2/backend/route/proto/team/response"

	"github.com/gin-gonic/gin"
)

// JoinTeam 加入团队
func JoinTeam(ctx *gin.Context, token string, u *user.User) error {
	tm, err := team.GetMemberByToken(ctx, token)
	if err != nil {
		slog.ErrorContext(ctx, "team.GetMemberByToken", "err", err)
		return i18n.NewErr("team.FailedToJoinTeam")
	}

	t := &team.Team{ID: tm.TeamID}
	exist, err := t.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "t.Get", "err", err)
		return i18n.NewErr("team.FailedToJoinTeam")
	}
	if !exist {
		return i18n.NewErr("team.DoesNotExist")
	}

	targetTM := &team.TeamMember{UserID: u.ID}
	teamMemberExist, err := t.HasMember(ctx, targetTM)
	if err != nil {
		slog.ErrorContext(ctx, "t.HasMember", "err", err)
		return i18n.NewErr("team.FailedToJoinTeam")
	}
	if teamMemberExist {
		return i18n.NewErr("teamMember.JoinTeamRepeat")
	}

	if _, err = t.AddMember(ctx, tm.ID, u); err != nil {
		slog.ErrorContext(ctx, "t.AddMember", "err", err)
		return i18n.NewErr("team.FailedToJoinTeam")
	}

	return nil
}

// ConvertModelTeam 将model.team转换为proto.team
func ConvertModelTeam(ctx *gin.Context, t *team.Team) *prototeamresponse.Team {
	memberCount, _ := team.GetMembersCount(ctx, t.ID)
	return &prototeamresponse.Team{
		OnlyIdInfo: protobase.OnlyIdInfo{
			ID: t.ID,
		},
		TeamDataOption: prototeambase.TeamDataOption{
			Name:   t.Name,
			Avatar: t.Avatar,
		},
		MembersCount: int(memberCount),
	}
}

// ConvertModelTeamMember 将model.teamMember转换为proto.teamMember
func ConvertModelTeamMember(ctx *gin.Context, tm *team.TeamMember, userInfo *user.User) *prototeamresponse.TeamMember {
	return &prototeamresponse.TeamMember{
		IdCreateTimeInfo: protobase.IdCreateTimeInfo{
			ID:        tm.ID,
			CreatedAt: tm.CreatedAt.Unix(),
		},
		TeamMemberData: prototeambase.TeamMemberData{
			Role: tm.Role,
			TeamMemberStatusOption: prototeambase.TeamMemberStatusOption{
				Status: tm.Status,
			},
		},
		TeamIdOption: protobase.TeamIdOption{
			TeamID: tm.TeamID,
		},
		User: ConvertModelUser(ctx, userInfo).UserData,
	}
}
