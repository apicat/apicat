package request

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
)

type GetProjectMemberListOption struct {
	protobase.PaginationOption
	protobase.ProjectIdOption
}

type CreateProjectMemberOption struct {
	protobase.ProjectIdOption
	protobase.ProjectMemberPermission
	MemberIDs []uint `json:"memberIDs" binding:"omitempty,dive,gt=0"`
}

type UpdateProjectMemberOption struct {
	ProjectMemberIDOption
	protobase.ProjectMemberPermission
}
