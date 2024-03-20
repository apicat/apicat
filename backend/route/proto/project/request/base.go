package request

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
)

type OrderNode struct {
	ParentID uint   `json:"parentID" binding:"omitempty,numeric,gte=0"`
	IDs      []uint `json:"ids" binding:"required,dive,gte=0"`
}

type ProjectMemberIDOption struct {
	protobase.ProjectIdOption
	MemberID uint `uri:"memberID" json:"memberID" query:"memberID" binding:"required"`
}

type ProjectGroupIdOption struct {
	GroupID uint `uri:"groupID" json:"groupID" query:"groupID" binding:"required"`
}

type SortOption struct {
	protobase.ProjectIdOption
	Target OrderNode `json:"target" binding:"required"`
	Origin OrderNode `json:"origin" binding:"required"`
}
