package response

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	projectbase "github.com/apicat/apicat/v2/backend/route/proto/project/base"
)

type ProjectGroup struct {
	protobase.OnlyIdInfo
	projectbase.ProjectGroupNameOption
}

type ProjectGroups []*ProjectGroup
