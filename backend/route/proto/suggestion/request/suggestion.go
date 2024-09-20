package request

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
)

type CollectionOption struct {
	protobase.ProjectIdOption
	RequestID string `json:"requestID" binding:"required,gt=1"`
	Title     string `json:"title" binding:"required"`
	Path      string `json:"path" binding:"omitempty,gt=0"`
}

type ModelOption struct {
	protobase.ProjectIdOption
	RequestID string `json:"requestID" binding:"required,gt=1"`
	Title     string `json:"title" binding:"required"`
}

type SchemaOption struct {
	protobase.ProjectIdOption
	RequestID string `json:"requestID" binding:"required,gt=1"`
	Title     string `json:"title" binding:"required"`
	Schema    string `json:"schema" binding:"required"`
}
