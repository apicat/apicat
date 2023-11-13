package proto

type CollectionDataGetData struct {
	ProjectID    string `uri:"project-id" binding:"required,gt=0"`
	CollectionID uint   `uri:"collection-id" binding:"required,gt=0"`
}

type ExportCollection struct {
	Type     string `form:"type" binding:"required,oneof=apicat swagger openapi3.0.0 openapi3.0.1 openapi3.0.2 openapi3.1.0 HTML md"`
	Download string `form:"download" binding:"omitempty,oneof=true false"`
}

type CollectionList struct {
	ID       uint              `json:"id"`
	ParentID uint              `json:"parent_id"`
	Title    string            `json:"title"`
	Type     string            `json:"type"`
	Selected *bool             `json:"selected,omitempty"`
	Items    []*CollectionList `json:"items"`
}

type CollectionCreate struct {
	ParentID    uint   `json:"parent_id" binding:"gte=0"`                       // 父级id
	Title       string `json:"title" binding:"required,lte=255"`                // 名称
	Type        string `json:"type" binding:"required,oneof=category doc http"` // 类型: category,doc,http
	Content     string `json:"content"`                                         // 内容
	IterationID string `json:"iteration_id" binding:"omitempty,gte=0"`          // 迭代id
}

type CollectionUpdate struct {
	Title   string `json:"title" binding:"required,lte=255"`
	Content string `json:"content"`
}

type CollectionCopyData struct {
	IterationID string `json:"iteration_id" binding:"omitempty,gte=0"`
}

type CollectionMovement struct {
	Target CollectionOrderContent `json:"target" binding:"required"`
	Origin CollectionOrderContent `json:"origin" binding:"required"`
}

type CollectionOrderContent struct {
	Pid uint   `json:"pid" binding:"gte=0"`
	Ids []uint `json:"ids" binding:"required,dive,gte=0"`
}

type CollectionDeleteData struct {
	IterationID string `form:"iteration_id" binding:"omitempty,gte=0"`
}

type CollectionsListData struct {
	IterationID string `form:"iteration_id" binding:"omitempty,gte=0"`
}
