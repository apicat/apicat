package request

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	projectbase "github.com/apicat/apicat/v2/backend/route/proto/project/base"
)

type GetDefinitionSchemaOption struct {
	protobase.ProjectIdOption
	SchemaID uint `uri:"schemaID" json:"schemaID" query:"schemaID" binding:"required,numeric,gt=0"`
}

type CreateDefinitionSchemaOption struct {
	protobase.ProjectIdOption
	projectbase.DefinitionSchemaDataOption
	projectbase.DefinitionSchemaParentIDOption
	projectbase.DefinitionSchemaTypeOption
}

type UpdateDefinitionSchemaOption struct {
	GetDefinitionSchemaOption
	projectbase.DefinitionSchemaDataOption
}

type DeleteDefinitionSchemaOption struct {
	GetDefinitionSchemaOption
	projectbase.DerefOption
}

type SortDefinitionSchemaOption struct {
	SortOption
}

type AIGenerateSchemaOption struct {
	protobase.ProjectIdOption
	projectbase.DefinitionSchemaParentIDOption
	Prompt string `json:"prompt" binding:"required"`
}
