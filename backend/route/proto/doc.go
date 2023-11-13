package proto

type CollectionHistoryListData struct {
	ID       uint                        `json:"id"`
	Title    string                      `json:"title"`
	Type     string                      `json:"type"`
	SubNodes []CollectionHistoryListData `json:"sub_nodes,omitempty"`
}

type CollectionHistoryUriData struct {
	ProjectID    string `uri:"project-id" binding:"required,gt=0"`
	CollectionID uint   `uri:"collection-id" binding:"required,gt=0"`
	HistoryID    uint   `uri:"history-id" binding:"required,gt=0"`
}

type CollectionHistoryDiffData struct {
	HistoryID1 uint `form:"history_id1"`
	HistoryID2 uint `form:"history_id2"`
}

type CollectionHistoryDetailsData struct {
	ID            uint   `json:"id"`
	CollectionID  uint   `json:"collection_id"`
	Content       string `json:"content"`
	CreatedTime   string `json:"created_time"`
	LastUpdatedBy string `json:"last_updated_by"`
	Title         string `json:"title"`
}

type DocShareStatusData struct {
	PublicCollectionID string `uri:"public_collection_id" binding:"required,lte=255"`
}
