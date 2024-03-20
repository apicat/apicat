package response

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
)

type ProjectShareStatus struct {
	protobase.ProjectVisibilityOption
	protobase.ProjectMemberPermission
	HasShare bool `json:"hasShare"`
}

type ProjectShareDetail struct {
	protobase.ProjectMemberPermission
	protobase.ProjectVisibilityOption
	protobase.SecretKeyOption
}
