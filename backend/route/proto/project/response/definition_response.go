package response

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	projectbase "github.com/apicat/apicat/backend/route/proto/project/base"
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
