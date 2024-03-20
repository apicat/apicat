package request

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	teambase "github.com/apicat/apicat/backend/route/proto/team/base"
)

type ListOption struct {
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
}

type MembersOption struct {
	protobase.TeamIdOption
	RolesOption
	ListOption
	Status string `query:"status" validate:"omitempty,oneof=active deactive"`
}

type UpdateTeamMemberOption struct {
	GetTeamMemberOption
	teambase.TeamMemberData
}
