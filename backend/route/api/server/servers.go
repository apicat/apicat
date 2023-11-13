package server

import (
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/server"
	"github.com/apicat/apicat/backend/module/translator"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UrlList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	p, _ := currentProject.(*project.Projects)

	s := server.NewServers()
	servers, err := s.GetByProjectId(p.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Servers.GetFailed"}),
		})
		return
	}

	result := []gin.H{}
	for _, s := range servers {
		result = append(result, gin.H{
			"description": s.Description,
			"url":         s.Url,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func UrlSettings(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	currentProject, _ := ctx.Get("CurrentProject")
	p, _ := currentProject.(*project.Projects)

	data := []proto.CreateServer{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	resule := []*server.Servers{}

	for k, v := range data {
		s := server.NewServers()
		s.ProjectId = p.ID
		s.Description = v.Description
		s.Url = v.Url
		s.DisplayOrder = k
		resule = append(resule, s)
	}

	s := server.NewServers()
	if err := s.DeleteAndCreateServers(p.ID, resule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Servers.SetFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
