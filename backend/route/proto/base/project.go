package base

import (
	"apicat-cloud/backend/model/project"
)

type ProjectIdOption struct {
	ProjectID string `uri:"projectID" json:"projectID" query:"projectID" binding:"required,len=24"`
}

type ProjectVisibilityOption struct {
	Visibility string `json:"visibility" binding:"required,oneof=public private"`
}

type ProjectMemberPermission struct {
	Permission project.Permission `json:"permission" binding:"required,oneof=manage write read"`
}
