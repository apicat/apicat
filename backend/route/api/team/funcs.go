package team

import (
	"apicat-cloud/backend/i18n"
	"apicat-cloud/backend/model/team"

	"apicat-cloud/backend/route/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func switchTeam(ctx *gin.Context, teamID string) error {
	u := jwt.GetUser(ctx)
	t := &team.Team{ID: teamID}
	exist, err := t.Get(ctx)
	if err != nil {
		return err
	}
	if !exist {
		return i18n.NewErr("team.DoesNotExist")
	}

	tm := &team.TeamMember{UserID: u.ID}
	teamMemberExist, err := t.HasMember(ctx, tm)
	if err != nil {
		return err
	}
	if !teamMemberExist {
		return i18n.NewErr("teamMember.NotTeamMember")
	}

	if err := tm.UpdateActiveAt(ctx); err != nil {
		return err
	}

	return nil
}
