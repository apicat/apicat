package response

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	teambase "github.com/apicat/apicat/v2/backend/route/proto/team/base"
	userresponse "github.com/apicat/apicat/v2/backend/route/proto/user/response"
)

type TeamMember struct {
	protobase.IdCreateTimeInfo
	teambase.TeamMemberData
	protobase.TeamIdOption
	User userresponse.UserData `json:"user"`
}

type TeamMembers []*TeamMember

type TeamMemberList struct {
	protobase.PaginationInfo
	Items TeamMembers `json:"items"`
}
