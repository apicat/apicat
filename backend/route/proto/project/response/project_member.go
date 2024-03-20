package response

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	teambase "github.com/apicat/apicat/backend/route/proto/team/base"
	userresponse "github.com/apicat/apicat/backend/route/proto/user/response"
)

type ProjectMember struct {
	protobase.IdCreateTimeInfo
	teambase.TeamMemberStatusOption
	protobase.ProjectMemberPermission
	User userresponse.UserData `json:"user"`
}

type ProjectMembers []*ProjectMember

type GetProjectMemberListResponse struct {
	protobase.PaginationInfo
	Items ProjectMembers `json:"items"`
}
