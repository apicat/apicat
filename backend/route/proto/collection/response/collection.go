package response

import (
	"time"

	protobase "github.com/apicat/apicat/backend/route/proto/base"
	"github.com/apicat/apicat/backend/route/proto/collection/base"
	projectbase "github.com/apicat/apicat/backend/route/proto/project/base"
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
