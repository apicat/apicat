package base

type EmbedInfo struct {
	ID        any   `json:"id"`
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

type OnlyIdInfo struct {
	ID any `json:"id"`
}

type IdCreateTimeInfo struct {
	ID        any   `json:"id"`
	CreatedAt int64 `json:"createdAt"`
}

type PaginationOption struct {
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
}

type PaginationInfo struct {
	Count       int `json:"count"`
	TotalPage   int `json:"totalPage"`
	CurrentPage int `json:"currentPage"`
}

type SecretKeyOption struct {
	SecretKey string `json:"secretKey"`
}

type TimeIntervalOption struct {
	StartTime int64 `query:"startTime" json:"startTime" binding:"required,numeric,gt=1584374400"`
	EndTime   int64 `query:"endTime" json:"endTime" binding:"required,numeric,gt=1584374400"`
}

type DiffOption struct {
	OriginalID uint `query:"originalID" json:"originalID" binding:"required,numeric,gt=0"`
	TargetID   uint `query:"targetID" json:"targetID" binding:"omitempty,numeric,gte=0"`
}
