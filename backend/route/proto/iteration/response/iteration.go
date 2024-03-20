package response

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	"github.com/apicat/apicat/backend/route/proto/iteration/base"
	projectbase "github.com/apicat/apicat/backend/route/proto/project/base"
)

type Iteration struct {
	protobase.IdCreateTimeInfo
	base.IterationData
	Project *IterationProject `json:"project"`
}

type IterationProject struct {
	protobase.OnlyIdInfo
	projectbase.ProjectDataOption
	SelfMember protobase.ProjectMemberPermission `json:"selfMember"`
}

type IterationListItem struct {
	protobase.IdCreateTimeInfo
	base.IterationData
	ApisCount int64                 `json:"apisCount"`
	Project   *IterationListProject `json:"project"`
}

type IterationListProject struct {
	ID         string                            `json:"id"`
	Title      string                            `json:"title" binding:"required"`
	SelfMember protobase.ProjectMemberPermission `json:"selfMember"`
}

type IterationList struct {
	protobase.PaginationInfo
	Items []*IterationListItem `json:"items"`
}
