package request

import (
	protobase "apicat-cloud/backend/route/proto/base"
	projectbase "apicat-cloud/backend/route/proto/project/base"
)

type CreateProjectGroupOption struct {
	protobase.TeamIdOption
	projectbase.ProjectGroupNameOption
}

type RenameProjectGroupOption struct {
	ProjectGroupIdOption
	projectbase.ProjectGroupNameOption
}

type SortProjectGroupOption struct {
	protobase.TeamIdOption
	GroupIDs []uint `json:"groupIDs" binding:"omitempty,dive,gte=0"`
}
