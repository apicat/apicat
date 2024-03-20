package access

import (
	"github.com/apicat/apicat/v2/backend/model/team"

	"github.com/gin-gonic/gin"
)

const (
	ctxTeamKey       = "selfteam"
	ctxTeamMemberKey = "selfteammember"
)

func setSelfTeam(ctx *gin.Context, t *team.Team) {
	ctx.Set(ctxTeamKey, t)
}

func GetSelfTeam(ctx *gin.Context) *team.Team {
	v, ok := ctx.Get(ctxTeamKey)
	if ok && v != nil {
		return v.(*team.Team)
	}
	return nil
}

func setSelfTeamMember(ctx *gin.Context, tm *team.TeamMember) {
	ctx.Set(ctxTeamMemberKey, tm)
}

func GetSelfTeamMember(ctx *gin.Context) *team.TeamMember {
	v, ok := ctx.Get(ctxTeamMemberKey)
	if ok && v != nil {
		return v.(*team.TeamMember)
	}
	return nil
}
