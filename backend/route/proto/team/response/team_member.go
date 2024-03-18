package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
	teambase "apicat-cloud/backend/route/proto/team/base"
	userresponse "apicat-cloud/backend/route/proto/user/response"
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
