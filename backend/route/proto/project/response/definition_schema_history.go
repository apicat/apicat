package response

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	projectbase "github.com/apicat/apicat/backend/route/proto/project/base"
)

type DefinitionSchemaHistory struct {
	protobase.IdCreateTimeInfo
	DefinitionSchemaHistoryData
	CreatedBy string `json:"createdBy"`
}

type DefinitionSchemaHistoryData struct {
	projectbase.DefinitionSchemaDataOption
	SchemaID uint `json:"schemaID"`
}

type DefinitionSchemaHistoryList []*DefinitionSchemaHistoryItem

type DefinitionSchemaHistoryItem struct {
	protobase.IdCreateTimeInfo
	CreatedBy string `json:"createdBy"`
}

type DiffDefinitionSchemaHistories struct {
	Schema1 *DefinitionSchemaHistory `json:"schema1"`
	Schema2 *DefinitionSchemaHistory `json:"schema2"`
}
