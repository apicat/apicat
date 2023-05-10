package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/apicat/apicat/common/spec"
	"github.com/apicat/apicat/common/spec/plugin/openapi"
	"github.com/apicat/apicat/common/translator"
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
	ID string `uri:"id" binding:"required"`
}

type ExportProject struct {
	Type     string `form:"type" binding:"required,oneof=apicat swagger openapi3.0.0 openapi3.0.1 openapi3.0.2 openapi3.1.0"`
	Download string `form:"download" binding:"omitempty,oneof=true false"`
}

func ProjectsList(ctx *gin.Context) {
	project, _ := models.NewProjects()
	projects, err := project.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	result := make([]gin.H, 0)
	for _, p := range projects {
		result = append(result, gin.H{
			"id":          p.PublicId,
			"title":       p.Title,
			"description": p.Description,
			"cover":       p.Cover,
			"created_at":  p.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  p.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func ProjectsCreate(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	user, _ := CurrentUser.(*models.Users)
	if user.Role == "user" {
		ctx.Status(http.StatusForbidden)
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

	// 进行数据导入工作
	if data.Data != "" {
		models.ServersImport(project.ID, content.Servers)
		definitionSchemas := models.DefinitionSchemasImport(project.ID, content.Definitions.Schemas)
		models.CollectionsImport(project.ID, 0, content.Collections, definitionSchemas)
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
	CurrentUser, _ := ctx.Get("CurrentUser")
	user, _ := CurrentUser.(*models.Users)
	if user.Role == "user" {
		ctx.Status(http.StatusForbidden)
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
	CurrentUser, _ := ctx.Get("CurrentUser")
	user, _ := CurrentUser.(*models.Users)
	if user.Role == "user" {
		ctx.Status(http.StatusForbidden)
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

func ProjectsGet(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, gin.H{
		"id":          project.PublicId,
		"title":       project.Title,
		"description": project.Description,
		"cover":       project.Cover,
		"created_at":  project.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":  project.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func ProjectDataGet(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	user, _ := CurrentUser.(*models.Users)
	if user.Role == "user" {
		ctx.Status(http.StatusForbidden)
		return
	}

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
	apicatData.Common.Responses = models.CommonResponsesExport(project.ID)
	apicatData.Definitions.Schemas = models.DefinitionSchemasExport(project.ID)
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
