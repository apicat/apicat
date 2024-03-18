package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
	projectbase "apicat-cloud/backend/route/proto/project/base"
)

type DefinitionResponse struct {
	protobase.EmbedInfo
	projectbase.DefinitionResponseDataOption
	projectbase.DefinitionResponseParentIDOption
	projectbase.DefinitionResponseTypeOption
	projectbase.OperatorID
}

type DefinitionResponseTree []*DefinitionResponseNode

type DefinitionResponseNode struct {
	protobase.EmbedInfo
	projectbase.DefinitionResponseDataOption
	projectbase.DefinitionResponseParentIDOption
	projectbase.DefinitionResponseTypeOption
	Items []*DefinitionResponseNode `json:"items"`
}
