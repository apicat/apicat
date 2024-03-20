package request

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	"github.com/apicat/apicat/backend/route/proto/collection/base"
)

type CollectionHistoryIDOption struct {
	base.ProjectCollectionIDOption
	HistoryID uint `uri:"historyID" json:"historyID" query:"historyID" binding:"required,numeric,gt=0"`
}

type GetCollectionHistoryListOption struct {
	base.ProjectCollectionIDOption
	protobase.TimeIntervalOption
}

type DiffCollectionHistoriesOption struct {
	base.ProjectCollectionIDOption
	protobase.DiffOption
}
