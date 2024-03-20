package request

import (
	"github.com/apicat/apicat/v2/backend/route/proto/collection/base"
)

type GenerateTestCaseOption struct {
	base.ProjectCollectionIDOption
	Prompt     string `json:"prompt" binding:"omitempty"`
	Regenerate bool   `json:"regenerate" binding:"boolean,omitempty"`
}

type GetTestCaseOption struct {
	base.ProjectCollectionIDOption
	TestCaseID uint `uri:"testCaseID" json:"testCaseID" binding:"required,gt=0"`
}

type RegenerateTestCaseOption struct {
	base.ProjectCollectionIDOption
	TestCaseID uint   `uri:"testCaseID" json:"testCaseID" binding:"required,gt=0"`
	Prompt     string `json:"prompt" binding:"omitempty"`
}

type DeleteTestCaseOption struct {
	base.ProjectCollectionIDOption
	TestCaseID uint `uri:"testCaseID" json:"testCaseID" binding:"required,gt=0"`
}
