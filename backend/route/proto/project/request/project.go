package request

import (
	"apicat-cloud/backend/model/project"
	protobase "apicat-cloud/backend/route/proto/base"
	projectbase "apicat-cloud/backend/route/proto/project/base"
)

type ProjectImportDataOption struct {
	Data string `json:"data"`
	Type string `json:"type" binding:"omitempty,oneof=apicat openapi swagger postman"`
}

type GroupIdOption struct {
	GroupID uint `query:"groupID" json:"groupID"`
}

type CreateProjectOption struct {
	protobase.TeamIdOption
	projectbase.ProjectDataOption
	GroupIdOption
	ProjectImportDataOption
}

type GetProjectListOption struct {
	protobase.TeamIdOption
	GroupIdOption
	Permissions []project.Permission `query:"permissions" binding:"omitempty,dive,oneof=manage write read"`
	IsFollowed  bool                 `query:"isFollowed"`
}

type SwitchProjectGroupOption struct {
	protobase.ProjectIdOption
	GroupIdOption
}

type GetProjectDetailOption struct {
	protobase.ProjectIdOption
	ShareCode string `query:"shareCode"`
}

type UpdateProjectOption struct {
	protobase.ProjectIdOption
	projectbase.ProjectDataOption
}

type GetExportPathOption struct {
	protobase.ProjectIdOption
	Type     string `query:"type" binding:"required,oneof=apicat swagger openapi3.0.0 openapi3.0.1 openapi3.0.2 openapi3.1.0 HTML md"`
	Download bool   `query:"download"`
}

type ExportCodeOption struct {
	protobase.ProjectIdOption
	Code string `uri:"code" binding:"required,len=32"`
}
