package access

import (
	"github.com/apicat/apicat/v2/backend/model/project"

	"github.com/gin-gonic/gin"
)

const (
	ctxProjectKey       = "selfproject"
	ctxProjectMemberKey = "selfprojectmember"
)

func setSelfProject(ctx *gin.Context, p *project.Project) {
	ctx.Set(ctxProjectKey, p)
}

func GetSelfProject(ctx *gin.Context) *project.Project {
	v, ok := ctx.Get(ctxProjectKey)
	if ok && v != nil {
		return v.(*project.Project)
	}
	return nil
}

func setSelfProjectMember(ctx *gin.Context, pm *project.ProjectMember) {
	ctx.Set(ctxProjectMemberKey, pm)
}

func GetSelfProjectMember(ctx *gin.Context) *project.ProjectMember {
	v, ok := ctx.Get(ctxProjectMemberKey)
	if ok && v != nil {
		return v.(*project.ProjectMember)
	}
	return nil
}
