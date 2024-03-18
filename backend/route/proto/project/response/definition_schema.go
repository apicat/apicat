package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
	projectbase "apicat-cloud/backend/route/proto/project/base"
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
