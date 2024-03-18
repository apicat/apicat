package request

import (
	protobase "apicat-cloud/backend/route/proto/base"
	collectionbase "apicat-cloud/backend/route/proto/collection/base"
	"apicat-cloud/backend/route/proto/iteration/base"
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
