package request

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	projectbase "github.com/apicat/apicat/v2/backend/route/proto/project/base"
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
