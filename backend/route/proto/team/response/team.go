package response

import (
	"github.com/apicat/apicat/v2/backend/model/team"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	teambase "github.com/apicat/apicat/v2/backend/route/proto/team/base"
)

type Team struct {
	protobase.OnlyIdInfo
	teambase.TeamDataOption
	MembersCount int `json:"membersCount"`
}

type TeamList struct {
	Items []*Team `json:"items"`
}

type CurrentTeamRes struct {
	Team
	Role team.Role `json:"role"`
}

type TeamInviteContent struct {
	Inviter string `json:"inviter"`
	Team    string `json:"team"`
}
