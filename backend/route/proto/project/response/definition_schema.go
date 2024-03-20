package response

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	projectbase "github.com/apicat/apicat/backend/route/proto/project/base"
)

type DefinitionSchema struct {
	protobase.EmbedInfo
	projectbase.DefinitionSchemaDataOption
	projectbase.DefinitionSchemaParentIDOption
	projectbase.DefinitionSchemaTypeOption
	projectbase.OperatorID
}

type DefinitionSchemaTree []*DefinitionSchemaNode

type DefinitionSchemaNode struct {
	protobase.EmbedInfo
	projectbase.DefinitionSchemaDataOption
	projectbase.DefinitionSchemaParentIDOption
	projectbase.DefinitionSchemaTypeOption
	Items []*DefinitionSchemaNode `json:"items"`
}
