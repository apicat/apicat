package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/apicat/apicat/commom/translator"
	"github.com/apicat/apicat/models"

	"github.com/gin-gonic/gin"
)

type DefinitionCreate struct {
	ParentId    uint                   `json:"parent_id" binding:"gte=0"`
	Name        string                 `json:"name" binding:"required,lte=255"`
	Description string                 `json:"description" binding:"lte=255"`
	Type        string                 `json:"type" binding:"required,oneof=category schema"`
	Schema      map[string]interface{} `json:"schema"`
}

type DefinitionUpdate struct {
	Name        string                 `json:"name" binding:"required,lte=255"`
	Description string                 `json:"description" binding:"lte=255"`
	Schema      map[string]interface{} `json:"schema"`
}

type DefinitionSearch struct {
	ParentId uint   `form:"parent_id" binding:"gte=0"`
	Name     string `form:"name" binding:"lte=255"`
	Type     string `form:"type" binding:"omitempty,oneof=category schema"`
}

type DefinitionID struct {
	ID uint `uri:"definition-id" binding:"required,gte=0"`
}

type DefinitionMove struct {
	Target OrderContent `json:"target" binding:"required"`
	Origin OrderContent `json:"origin" binding:"required"`
}

type OrderContent struct {
	Pid uint   `json:"pid" binding:"gte=0"`
	Ids []uint `json:"ids" binding:"required,dive,gte=0"`
}

func DefinitionsList(ctx *gin.Context) {
	var data DefinitionSearch

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	project, _ := ctx.Get("CurrentProject")

	definition, _ := models.NewDefinitions()
	definition.ProjectId = project.(*models.Projects).ID
	definition.ParentId = data.ParentId
	definition.Name = data.Name
	definition.Type = data.Type

	definitions, err := definition.List()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.NotFound"}),
		})
		return
	}

	result := make([]gin.H, 0)
	for _, d := range definitions {
		schema := make(map[string]interface{})
		if err := json.Unmarshal([]byte(d.Schema), &schema); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		result = append(result, gin.H{
			"id":          d.ID,
			"parent_id":   d.ParentId,
			"name":        d.Name,
			"description": d.Description,
			"type":        d.Type,
			"schema":      schema,
		})
	}
	ctx.JSON(http.StatusOK, result)
}

func DefinitionsCreate(ctx *gin.Context) {
	var data DefinitionCreate

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	schemaJson, err := json.Marshal(data.Schema)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	project, _ := ctx.Get("CurrentProject")
	definition, _ := models.NewDefinitions()
	definition.ProjectId = project.(*models.Projects).ID
	definition.Name = data.Name
	definitions, err := definition.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if len(definitions) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.NameExists"}),
		})
		return
	}

	definition.Description = data.Description
	definition.Type = data.Type
	definition.Schema = string(schemaJson)
	if err := definition.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.CreateFail"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          definition.ID,
		"parent_id":   definition.ParentId,
		"name":        definition.Name,
		"description": definition.Description,
		"type":        definition.Type,
		"schema":      data.Schema,
		"created_at":  definition.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by":  definition.Creator(),
		"updated_at":  definition.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by":  definition.Updater(),
	})
}

func DefinitionsUpdate(ctx *gin.Context) {
	var (
		uriData DefinitionID
		data    DefinitionUpdate
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
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

	definition, err := models.NewDefinitions(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.NotFound"}),
		})
		return
	}

	schemaJson, err := json.Marshal(data.Schema)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	definition.Name = data.Name
	definition.Description = data.Description
	definitions, err := definition.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if len(definitions) > 0 && definitions[0].ID != definition.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.NameExists"}),
		})
		return
	}

	definition.Schema = string(schemaJson)
	if err := definition.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.UpdateFail"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func DefinitionsDelete(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	var data DefinitionID

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	definition, err := models.NewDefinitions(data.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.NotFound"}),
		})
		return
	}

	// 判断该模型是否被使用
	ref := "{\"$ref\":\"#/definitions/schemas/" + strconv.FormatUint(uint64(definition.ID), 10) + "\"}"

	collections, _ := models.NewCollections()
	collections.ProjectId = project.ID
	collectionList, err := collections.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	for _, v := range collectionList {
		if strings.Contains(v.Content, ref) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.InUse"}),
			})
			return
		}
	}

	definitions, _ := models.NewDefinitions()
	definitions.ProjectId = project.ID
	definitionList, err := definitions.List()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	for _, v := range definitionList {
		if strings.Contains(v.Schema, ref) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.InUse"}),
			})
			return
		}
	}

	commonResponses, _ := models.NewCommonResponses()
	commonResponses.ProjectID = project.ID
	commonResponsesList, err := commonResponses.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	for _, v := range commonResponsesList {
		if strings.Contains(v.Content, ref) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.InUse"}),
			})
			return
		}
	}

	if err := definition.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.DeleteFail"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func DefinitionsGet(ctx *gin.Context) {
	var data DefinitionID

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	definition, err := models.NewDefinitions(data.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.NotFound"}),
		})
		return
	}

	schema := make(map[string]interface{})
	if err := json.Unmarshal([]byte(definition.Schema), &schema); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          definition.ID,
		"parent_id":   definition.ParentId,
		"name":        definition.Name,
		"description": definition.Description,
		"type":        definition.Type,
		"schema":      schema,
		"created_at":  definition.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by":  definition.Creator(),
		"updated_at":  definition.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by":  definition.Updater(),
	})
}

func DefinitionsCopy(ctx *gin.Context) {
	var data DefinitionID

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	oldDefinition, err := models.NewDefinitions(data.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.NotFound"}),
		})
		return
	}

	newDefinition, _ := models.NewDefinitions()
	newDefinition.ProjectId = oldDefinition.ProjectId
	newDefinition.Name = fmt.Sprintf("%s (copy)", oldDefinition.Name)
	newDefinition.Description = oldDefinition.Description
	newDefinition.Type = oldDefinition.Type
	newDefinition.Schema = oldDefinition.Schema
	if err := newDefinition.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Definitions.CopyFail"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          newDefinition.ID,
		"parent_id":   newDefinition.ParentId,
		"name":        newDefinition.Name,
		"description": newDefinition.Description,
		"type":        newDefinition.Type,
		"schema":      newDefinition.Schema,
		"created_at":  newDefinition.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by":  newDefinition.Creator(),
		"updated_at":  newDefinition.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by":  newDefinition.Updater(),
	})
}

func DefinitionsMove(ctx *gin.Context) {
	var data DefinitionMove

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for i, id := range data.Target.Ids {
		if definition, err := models.NewDefinitions(id); err == nil {
			definition.ParentId = data.Target.Pid
			definition.DisplayOrder = i
			definition.Save()
		}
	}

	if data.Target.Pid != data.Origin.Pid {
		for i, id := range data.Origin.Ids {
			if definition, err := models.NewDefinitions(id); err == nil {
				definition.ParentId = data.Origin.Pid
				definition.DisplayOrder = i
				definition.Save()
			}
		}
	}

	ctx.Status(http.StatusCreated)
}
