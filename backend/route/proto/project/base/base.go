package base

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
)

type OperatorID struct {
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
}

// 移动到proto.base
type DerefOption struct {
	Deref bool `query:"deref" json:"deref" binding:"omitempty,boolean"`
}

type ProjectGroupNameOption struct {
	Name string `json:"name" binding:"required"`
}

// 移动到proto.base
type ProjectDataOption struct {
	Title string `json:"title" binding:"required"`
	protobase.ProjectVisibilityOption
	Cover       string `json:"cover" binding:"required"`
	Description string `json:"description"`
}

// 移动到proto.base
type ShareCode struct {
	ShareCode  string `json:"shareCode"`
	Expiration int64  `json:"expiration"`
}

type DefinitionResponseParentIDOption struct {
	ParentID uint `json:"parentID" binding:"gte=0"`
}

type DefinitionResponseTypeOption struct {
	Type string `json:"type" binding:"required,oneof=category response"`
}

type DefinitionResponseDataOption struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Header      string `json:"header"`
	Content     string `json:"content"`
}

type DefinitionSchemaParentIDOption struct {
	ParentID uint `json:"parentID" binding:"numeric,gte=0"`
}

type DefinitionSchemaTypeOption struct {
	Type string `json:"type" binding:"required,oneof=category schema"`
}

type DefinitionSchemaDataOption struct {
	Name        string `json:"name" binding:"required,gte=1"`
	Schema      string `json:"schema"`
	Description string `json:"description"`
}

type GlobalParameterDataOption struct {
	In       string `json:"in" binding:"required,oneof=header cookie query path"`
	Name     string `json:"name" binding:"required,gte=1"`
	Required bool   `json:"required" binding:"boolean"`
	Schema   string `json:"schema" binding:"required"`
}

type ProjectServerDataOption struct {
	URL         string `json:"url" binding:"required,startswith=http://|startswith=https://"`
	Description string `json:"description"`
}
