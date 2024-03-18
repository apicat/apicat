package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
	projectbase "apicat-cloud/backend/route/proto/project/base"
)

type ProjectGroup struct {
	protobase.OnlyIdInfo
	projectbase.ProjectGroupNameOption
}

type ProjectGroups []*ProjectGroup
