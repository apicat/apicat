package base

import (
	protobase "apicat-cloud/backend/route/proto/base"
)

type CollectionIDOption struct {
	CollectionID uint `uri:"collectionID" json:"collectionID" query:"collectionID" binding:"required"`
}

// 移动到proto.base
type ShareCode struct {
	ShareCode  string `json:"shareCode"`
	Expiration int64  `json:"expiration"`
}

type ProjectCollectionIDOption struct {
	protobase.ProjectIdOption
	CollectionIDOption
}

type CollectionPublicIDOption struct {
	CollectionPublicID string `uri:"collectionPublicID" json:"collectionPublicID" query:"collectionPublicID" binding:"required"`
}

type CollectionData struct {
	Title   string `json:"title" binding:"required,gte=1,lte=255"`
	Content string `json:"content"`
}

type CollectionTypeOption struct {
	Type string `json:"type" binding:"required,oneof=category doc http"`
}

type CollectionParentIDOption struct {
	ParentID uint `json:"parentID" binding:"gte=0"`
}

type CollectionIDsOption struct {
	CollectionIDs []uint `json:"collectionIDs" binding:"omitempty,dive,gte=0"`
}
