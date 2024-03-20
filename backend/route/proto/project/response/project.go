package response

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	projectbase "github.com/apicat/apicat/backend/route/proto/project/base"
)

type ProjectSelfMemberInfo struct {
	GroupID    uint `json:"groupID"`
	IsFollowed bool `json:"isFollowed"`
	protobase.ProjectMemberPermission
}

type ProjectListItem struct {
	protobase.OnlyIdInfo
	projectbase.ProjectDataOption
	SelfMember ProjectSelfMemberInfo `json:"selfMember"`
}

type ProjectDetail struct {
	ProjectListItem
	MockURL string `json:"mockURL"`
}

type GetProjectsResponse []*ProjectListItem

type ExportProject struct {
	Path string `json:"path"`
}
