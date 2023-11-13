package proto

type AICreateCollectionStructure struct {
	ParentID    uint   `json:"parent_id" binding:"gte=0"`              // 父级id
	Title       string `json:"title" binding:"required,lte=255"`       // 名称
	SchemaID    uint   `json:"schema_id" binding:"gte=0"`              // 模型id
	Path        string `json:"path" binding:"lte=255"`                 // 请求路径
	Method      string `json:"method" binding:"lte=255"`               // 请求方法
	IterationID string `json:"iteration_id" binding:"omitempty,gte=0"` // 迭代id
}

type AICreateSchemaStructure struct {
	ParentID uint   `json:"parent_id" binding:"gte=0"`       // 父级id
	Name     string `json:"name" binding:"required,lte=255"` // 名称
}

type AICreateApiNameStructure struct {
	SchemaID uint `form:"schema_id" binding:"gt=0"` // 模型id
}
