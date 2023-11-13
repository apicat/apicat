package proto

import "github.com/apicat/apicat/backend/module/spec"

type ResponseDetailData struct {
	Name        string                 `json:"name" binding:"required,lte=255"`
	Description string                 `json:"description" binding:"lte=255"`
	Type        string                 `json:"type" binding:"required,oneof=category response"`
	Header      []*spec.Schema         `json:"header,omitempty" binding:"omitempty,dive"`
	Content     map[string]spec.Schema `json:"content,omitempty" binding:"required"`
	Ref         string                 `json:"$ref,omitempty" binding:"omitempty,lte=255"`
}

type SchemaUriData struct {
	ProjectID string `uri:"project-id" binding:"required,gt=0"`
	SchemaID  uint   `uri:"schemas-id" binding:"required,gt=0"`
}

type SchemaHistoryUriData struct {
	ProjectID string `uri:"project-id" binding:"required,gt=0"`
	SchemaID  uint   `uri:"schemas-id" binding:"required,gt=0"`
	HistoryID uint   `uri:"history-id" binding:"required,gt=0"`
}

type SchemaHistoryListData struct {
	ID       uint                    `json:"id"`
	Name     string                  `json:"name"`
	Type     string                  `json:"type"`
	SubNodes []SchemaHistoryListData `json:"sub_nodes,omitempty"`
}

type SchemaHistoryDetailsData struct {
	ID            uint           `json:"id"`
	SchemaID      uint           `json:"schema_id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Schema        map[string]any `json:"schema"`
	CreatedTime   string         `json:"created_time"`
	LastUpdatedBy string         `json:"last_updated_by"`
}

type SchemaHistoryDiffData struct {
	HistoryID1 uint `form:"history_id1"`
	HistoryID2 uint `form:"history_id2"`
}

type DefinitionSchemaCreate struct {
	ParentId    uint                   `json:"parent_id" binding:"gte=0"`
	Name        string                 `json:"name" binding:"required,lte=255"`
	Description string                 `json:"description" binding:"lte=255"`
	Type        string                 `json:"type" binding:"required,oneof=category schema"`
	Schema      map[string]interface{} `json:"schema"`
}

type DefinitionSchemaUpdate struct {
	Name        string                 `json:"name" binding:"required,lte=255"`
	Description string                 `json:"description" binding:"lte=255"`
	Schema      map[string]interface{} `json:"schema"`
}

type DefinitionSchemaSearch struct {
	ParentId uint   `form:"parent_id" binding:"gte=0"`
	Name     string `form:"name" binding:"lte=255"`
	Type     string `form:"type" binding:"omitempty,oneof=category schema"`
}

type DefinitionSchemaID struct {
	ID uint `uri:"schemas-id" binding:"required,gte=0"`
}

type DefinitionSchemaMove struct {
	Target OrderContent `json:"target" binding:"required"`
	Origin OrderContent `json:"origin" binding:"required"`
}

type OrderContent struct {
	Pid uint   `json:"pid" binding:"gte=0"`
	Ids []uint `json:"ids" binding:"required,dive,gte=0"`
}
