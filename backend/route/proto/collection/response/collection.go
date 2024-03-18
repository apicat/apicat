package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
	"apicat-cloud/backend/route/proto/collection/base"
	projectbase "apicat-cloud/backend/route/proto/project/base"
	"time"
)

type Collection struct {
	protobase.EmbedInfo
	base.CollectionData
	base.CollectionTypeOption
	base.CollectionParentIDOption
	projectbase.OperatorID
}

type CollectionTree []*CollectionNode

type CollectionNode struct {
	ID       uint              `json:"id"`
	Title    string            `json:"title"`
	Selected *bool             `json:"selected,omitempty"`
	Items    []*CollectionNode `json:"items"`
	base.CollectionParentIDOption
	base.CollectionTypeOption
}

type DeleteInfo struct {
	DeletedAt time.Time `json:"deletedAt"`
	DeletedBy string    `json:"deletedBy"`
}

type Trash struct {
	base.CollectionIDOption
	CollectionTitle string `json:"collectionTitle"`
	DeleteInfo
}

type TrashList []*Trash

type RestoreNum struct {
	Num int `json:"num"`
}

type ExportCollection struct {
	Path string `json:"path"`
}
