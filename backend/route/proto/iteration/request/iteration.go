package request

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	collectionbase "github.com/apicat/apicat/backend/route/proto/collection/base"
	"github.com/apicat/apicat/backend/route/proto/iteration/base"
)

type CreateIterationOption struct {
	protobase.TeamIdOption
	protobase.ProjectIdOption
	base.IterationData
	collectionbase.CollectionIDsOption
}

type GetIterationListOption struct {
	protobase.TeamIdOption
	protobase.PaginationOption
	ProjectID string `uri:"projectID" json:"projectID" query:"projectID"`
}

type UpdateIterationOption struct {
	base.IterationIDOption
	base.IterationData
	collectionbase.CollectionIDsOption
}
