package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/apicat/apicat/common/spec"
	"github.com/apicat/apicat/common/spec/plugin/export"
	"github.com/apicat/apicat/common/spec/plugin/openapi"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/enum"
	"github.com/apicat/apicat/models"
	"golang.org/x/exp/slog"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
)

type CreateProject struct {
	Title string `json:"title" binding:"required,lte=255"`
	Data  string `json:"data"`
	Cover string `json:"cover" binding:"lte=255"`
}

type UpdateProject struct {
	Title       string `json:"title" binding:"required,lte=255"`
	Description string `json:"description" binding:"lte=255"`
	Cover       string `json:"cover" binding:"lte=255"`
}

type ProjectID struct {
	ID string `uri:"project-id" binding:"required"`
}

type ExportProject struct {
	Type     string `form:"type" binding:"required,oneof=apicat swagger openapi3.0.0 openapi3.0.1 openapi3.0.2 openapi3.1.0 HTML md"`
	Download string `form:"download" binding:"omitempty,oneof=true false"`
}

type TranslateProject struct {
	MemberID uint `json:"member_id" binding:"required,lte=255"`
}

func ProjectsList(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	project, _ := models.NewProjects()
	projects, err := project.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	projectMembers, err := models.GetUserInvolvedProject(currentUser.(*models.Users).ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	projectIDs := []uint{}
	for _, v := range projectMembers {
		projectIDs = append(projectIDs, v.ProjectID)
	}

	projectsList := []gin.H{}
	for _, i := range projectIDs {
		for _, v := range projects {
			if i == v.ID {
				projectsList = append(projectsList, gin.H{
					"id":          v.PublicId,
					"title":       v.Title,
					"description": v.Description,
					"cover":       v.Cover,
					"created_at":  v.CreatedAt.Format("2006-01-02 15:04:05"),
					"updated_at":  v.UpdatedAt.Format("2006-01-02 15:04:05"),
				})
			}
		}
	}

	ctx.JSON(http.StatusOK, projectsList)
}

func ProjectsGet(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var data ProjectID

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	project, err := models.NewProjects(data.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	projectMember, _ := models.NewProjectMembers()
	projectMember.ProjectID = project.ID
	projectMember.UserID = currentUser.(*models.Users).ID

	authority := "none"
	if err := projectMember.GetByUserIDAndProjectID(); err == nil {
		authority = projectMember.Authority
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          project.PublicId,
		"title":       project.Title,
		"description": project.Description,
		"cover":       project.Cover,
		"authority":   authority,
		"created_at":  project.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":  project.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func ProjectsCreate(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	user, _ := CurrentUser.(*models.Users)
	if user.Role == "user" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.MemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		data       CreateProject
		content    *spec.Spec
		rawContent []byte
		err        error
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	project, _ := models.NewProjects()
	if data.Data != "" {
		var base64Content string
		if strings.Contains(data.Data, "data:application/json;base64,") {
			base64Content = strings.Replace(data.Data, "data:application/json;base64,", "", 1)
		} else {
			base64Content = strings.Replace(data.Data, "data:application/x-yaml;base64,", "", 1)
		}
		rawContent, err = base64.StdEncoding.DecodeString(base64Content)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.ImportFail"}),
			})
			return
		}

		content, err = openapi.Decode(rawContent)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.ImportFail"}),
			})
			return
		}
		project.Description = content.Info.Description
	}
	project.Title = data.Title
	project.PublicId = shortuuid.New()
	project.Visibility = 0
	project.Cover = data.Cover
	if err := project.Create(); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.CreateFail"}),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.ProjectID = project.ID
	pm.UserID = user.ID
	pm.Authority = models.ProjectMembersManage
	if err := pm.Create(); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.CreateFail"}),
		})
		return
	}

	// 进行数据导入工作
	if data.Data != "" {
		models.ServersImport(project.ID, content.Servers)

		refContentVirtualIDToId := &models.RefContentVirtualIDToId{
			DefinitionSchemas:    models.DefinitionSchemasImport(project.ID, content.Definitions.Schemas),
			DefinitionResponses:  models.DefinitionResponsesImport(project.ID, content.Definitions.Responses),
			DefinitionParameters: models.DefinitionParametersImport(project.ID, content.Definitions.Parameters),
		}

		models.CollectionsImport(project.ID, 0, content.Collections, refContentVirtualIDToId)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          project.PublicId,
		"title":       project.Title,
		"description": project.Description,
		"cover":       project.Cover,
		"created_at":  project.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":  project.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func ProjectsUpdate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		uriData ProjectID
		data    UpdateProject
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	project, err := models.NewProjects(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	project.Title = data.Title
	project.Description = data.Description
	project.Cover = data.Cover
	if err := project.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.UpdateFail"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func ProjectsDelete(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var data ProjectID

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	project, err := models.NewProjects(data.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}
	if err := project.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.DeleteFail"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func ProjectDataGet(ctx *gin.Context) {
	var (
		uriData ProjectID
		data    ExportProject
		content []byte
		err     error
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	project, err := models.NewProjects(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	apicatData := &spec.Spec{}
	apicatData.ApiCat = "apicat"
	apicatData.Info = &spec.Info{
		ID:          project.PublicId,
		Title:       project.Title,
		Description: project.Description,
		Version:     "1.0.0",
	}

	apicatData.Servers = models.ServersExport(project.ID)
	apicatData.Globals.Parameters = models.GlobalParametersExport(project.ID)
	apicatData.Definitions.Schemas = models.DefinitionSchemasExport(project.ID)
	apicatData.Definitions.Parameters = models.DefinitionParametersExport(project.ID)
	apicatData.Definitions.Responses = models.DefinitionResponsesExport(project.ID)
	apicatData.Collections = models.CollectionsExport(project.ID)

	if apicatDataContent, err := json.Marshal(apicatData); err == nil {
		slog.InfoCtx(ctx, "Export", slog.String("apicat", string(apicatDataContent)))
	}

	switch data.Type {
	case "swagger":
		content, err = openapi.Encode(apicatData, "2.0")
	case "openapi3.0.0":
		content, err = openapi.Encode(apicatData, "3.0.0")
	case "openapi3.0.1":
		content, err = openapi.Encode(apicatData, "3.0.1")
	case "openapi3.0.2":
		content, err = openapi.Encode(apicatData, "3.0.2")
	case "openapi3.1.0":
		content, err = openapi.Encode(apicatData, "3.1.0")
	case "HTML":
		content, err = export.HTML(apicatData)
	case "md":
		content, err = export.Markdown(apicatData)
	default:
		content, err = apicatData.ToJSON(spec.JSONOption{Indent: "  "})
	}

	slog.InfoCtx(ctx, "Export", slog.String(data.Type, string(content)))

	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.ExportFail"}),
		})
		return
	}

	if data.Download == "true" {
		ctx.Header("Content-Disposition", "attachment; filename="+project.Title+"-"+data.Type+".json")
		ctx.Data(http.StatusOK, "application/octet-stream", content)
	} else {
		ctx.Data(http.StatusOK, "application/json", content)
	}
}

// ProjectExit handles the exit of a project member.
func ProjectExit(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if currentProjectMember.(*models.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	if err := currentProjectMember.(*models.ProjectMembers).Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.ExitFail"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func ProjectTransfer(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := TranslateProject{}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pm, err := models.NewProjectMembers(data.MemberID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.NotFound"}),
		})
		return
	}

	if pm.Authority != models.ProjectMembersWrite {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.TargetProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.TransferFail"}),
		})
		return
	}

	if pm.ProjectID != currentProjectMember.(*models.ProjectMembers).ProjectID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.TargetProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.TransferFail"}),
		})
		return
	}

	currentProjectMember.(*models.ProjectMembers).Authority = models.ProjectMembersWrite
	if err := currentProjectMember.(*models.ProjectMembers).Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.TransferFail"}),
		})
		return
	}

	pm.Authority = models.ProjectMembersManage
	if err := pm.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.TransferFail"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
