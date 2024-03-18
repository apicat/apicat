package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
	teambase "apicat-cloud/backend/route/proto/team/base"
	userresponse "apicat-cloud/backend/route/proto/user/response"
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
