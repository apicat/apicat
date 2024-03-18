package request

import (
	"apicat-cloud/backend/model/team"
	protobase "apicat-cloud/backend/route/proto/base"
)

type RolesOption struct {
	Roles []team.Role `query:"roles" json:"roles" binding:"omitempty,dive,oneof=owner admin member"`
}

type GetTeamMemberOption struct {
	protobase.TeamIdOption
	MemberID uint `uri:"memberID" json:"memberID" query:"memberID" binding:"required"`
}
