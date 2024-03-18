package request

import (
	protobase "apicat-cloud/backend/route/proto/base"
	projectbase "apicat-cloud/backend/route/proto/project/base"
)

type GetProjectServerOption struct {
	protobase.ProjectIdOption
	ServerID uint `uri:"serverID" json:"serverID" query:"serverID" binding:"required,numeric,gt=0"`
}

type CreateProjectServerOption struct {
	protobase.ProjectIdOption
	projectbase.ProjectServerDataOption
}

type UpdateProjectServerOption struct {
	GetProjectServerOption
	projectbase.ProjectServerDataOption
}

type SortProjectServerOpt struct {
	protobase.ProjectIdOption
	ServerIDs []uint `json:"serverIDs" binding:"omitempty,dive,gt=0"`
}
