package request

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
)

type DefinitionSchemaHistoryIDOption struct {
	GetDefinitionSchemaOption
	HistoryID uint `uri:"historyID" json:"historyID" query:"historyID" binding:"required,numeric,gt=0"`
}

type GetDefinitionSchemaHistoryListOption struct {
	GetDefinitionSchemaOption
	protobase.TimeIntervalOption
}
type DiffDefinitionSchemaHistoriesOption struct {
	GetDefinitionSchemaOption
	protobase.DiffOption
}
