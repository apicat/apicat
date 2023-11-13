package proto

type IterationSchemaData struct {
	ID           string `json:"id,omitempty" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	ProjectID    string `json:"project_id" binding:"required"`
	ProjectTitle string `json:"project_title" binding:"required"`
	ApiNum       int64  `json:"api_num"`
	Authority    string `json:"authority"`
	CreatedAt    string `json:"created_at" binding:"required"`
}

type IterationListData struct {
	ProjectID string `form:"project_id"`
	Page      int64  `form:"page"`
	PageSize  int64  `form:"page_size"`
}

type IterationListResData struct {
	CurrentPage int64                 `json:"current_page"`
	TotalPage   int64                 `json:"total_page"`
	Total       int64                 `json:"total"`
	Iterations  []IterationSchemaData `json:"iterations"`
}

type IterationCreateData struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	ProjectID     string `json:"project_id" binding:"required"`
	CollectionIDs []uint `json:"collection_ids"`
}

type IterationUriData struct {
	IterationID string `uri:"iteration-id" binding:"required"`
}

type IterationUpdateData struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	CollectionIDs []uint `json:"collection_ids"`
}
