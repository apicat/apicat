package api

import (
	"net/http"

	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/apicat/apicat/backend/models"

	"github.com/gin-gonic/gin"
)

type CreateServer struct {
	Description string `json:"description" binding:"lte=255"`
	Url         string `json:"url" binding:"required,lte=255"`
}

func UrlList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	server := models.NewServers()
	servers, err := server.GetByProjectId(project.ID)
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
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	data := []CreateServer{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	resule := []*models.Servers{}

	for k, v := range data {
		server := models.NewServers()
		server.ProjectId = project.ID
		server.Description = v.Description
		server.Url = v.Url
		server.DisplayOrder = k
		resule = append(resule, server)
	}

	server := models.NewServers()
	if err := server.DeleteAndCreateServers(project.ID, resule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Servers.SetFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
