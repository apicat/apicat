package request

import (
	protobase "apicat-cloud/backend/route/proto/base"
	"apicat-cloud/backend/route/proto/collection/base"
)

type IterationIDNotRequiredOption struct {
	IterationID string `uri:"iterationID" json:"iterationID" query:"iterationID" binding:"omitempty,len=24"`
}

type OrderNode struct {
	ParentID uint   `json:"parentID" binding:"omitempty,numeric,gte=0"`
	IDs      []uint `json:"ids" binding:"required,dive,gte=0"`
}

type CreateCollectionOption struct {
	protobase.ProjectIdOption
	IterationIDNotRequiredOption
	base.CollectionData
	base.CollectionParentIDOption
	base.CollectionTypeOption
}

type GetCollectionListOption struct {
	protobase.ProjectIdOption
	IterationIDNotRequiredOption
}

type UpdateCollectionOption struct {
	base.ProjectCollectionIDOption
	base.CollectionData
}

type DeleteCollectionOption struct {
	base.ProjectCollectionIDOption
	IterationIDNotRequiredOption
}

type MoveCollectionOption struct {
	protobase.ProjectIdOption
	Target OrderNode `json:"target" binding:"required"`
	Origin OrderNode `json:"origin" binding:"required"`
}

type CopyCollectionOption struct {
	base.ProjectCollectionIDOption
	IterationIDNotRequiredOption
}

type RestoreOption struct {
	protobase.ProjectIdOption
	base.CollectionIDsOption
}

type GetExportPathOption struct {
	base.ProjectCollectionIDOption
	Type     string `query:"type" binding:"required,oneof=apicat swagger openapi3.0.0 openapi3.0.1 openapi3.0.2 openapi3.1.0 HTML md"`
	Download bool   `query:"download"`
}

type ExportCodeOption struct {
	base.ProjectCollectionIDOption
	Code string `uri:"code" binding:"required,len=32"`
}

type AIGenerateCollectionOption struct {
	protobase.ProjectIdOption
	base.CollectionParentIDOption
	Prompt string `json:"prompt" binding:"required"`
	IterationIDNotRequiredOption
}
