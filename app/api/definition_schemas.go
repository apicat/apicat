package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"

	"github.com/gin-gonic/gin"
)

type DefinitionSchemaCreate struct {
	ParentId    uint                   `json:"parent_id" binding:"gte=0"`
	Name        string                 `json:"name" binding:"required,lte=255"`
	Description string                 `json:"description" binding:"lte=255"`
	Type        string                 `json:"type" binding:"required,oneof=category schema"`
	Schema      map[string]interface{} `json:"schema"`
}

type DefinitionSchemaUpdate struct {
	Name        string                 `json:"name" binding:"required,lte=255"`
	Description string                 `json:"description" binding:"lte=255"`
	Schema      map[string]interface{} `json:"schema"`
}

type DefinitionSchemaSearch struct {
	ParentId uint   `form:"parent_id" binding:"gte=0"`
	Name     string `form:"name" binding:"lte=255"`
	Type     string `form:"type" binding:"omitempty,oneof=category schema"`
}

type DefinitionSchemaID struct {
	ID uint `uri:"schemas-id" binding:"required,gte=0"`
}

type DefinitionSchemaMove struct {
	Target OrderContent `json:"target" binding:"required"`
	Origin OrderContent `json:"origin" binding:"required"`
}

type OrderContent struct {
	Pid uint   `json:"pid" binding:"gte=0"`
	Ids []uint `json:"ids" binding:"required,dive,gte=0"`
}

func DefinitionSchemasList(ctx *gin.Context) {
	var data DefinitionSchemaSearch

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	project, _ := ctx.Get("CurrentProject")

	definition, _ := models.NewDefinitionSchemas()
	definition.ProjectId = project.(*models.Projects).ID
	definition.ParentId = data.ParentId
	definition.Name = data.Name
	definition.Type = data.Type

	definitions, err := definition.List()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.QueryFailed"}),
		})
		return
	}

	result := make([]gin.H, 0)
	for _, d := range definitions {
		schema := make(map[string]interface{})
		if err := json.Unmarshal([]byte(d.Schema), &schema); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
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

func DefinitionSchemasCreate(ctx *gin.Context) {
	currentMember, _ := ctx.Get("CurrentMember")
	if !currentMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var data DefinitionSchemaCreate

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	schemaJson, err := json.Marshal(data.Schema)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}

	project, _ := ctx.Get("CurrentProject")
	definition, _ := models.NewDefinitionSchemas()
	definition.ProjectId = project.(*models.Projects).ID
	definition.Name = data.Name
	definitions, err := definition.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.QueryFailed"}),
		})
		return
	}
	if len(definitions) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.NameExists"}),
		})
		return
	}

	definition.Description = data.Description
	definition.Type = data.Type
	definition.Schema = string(schemaJson)
	if err := definition.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.CreateFail"}),
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

func DefinitionSchemasUpdate(ctx *gin.Context) {
	currentMember, _ := ctx.Get("CurrentMember")
	if !currentMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		uriData DefinitionSchemaID
		data    DefinitionSchemaUpdate
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

	definition, err := models.NewDefinitionSchemas(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
		})
		return
	}

	schemaJson, err := json.Marshal(data.Schema)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}

	definition.Name = data.Name
	definition.Description = data.Description
	definitions, err := definition.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.QueryFailed"}),
		})
		return
	}

	if len(definitions) > 0 && definitions[0].ID != definition.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.NameExists"}),
		})
		return
	}

	definition.Schema = string(schemaJson)
	if err := definition.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.UpdateFail"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func DefinitionSchemasDelete(ctx *gin.Context) {
	currentMember, _ := ctx.Get("CurrentMember")
	if !currentMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var data DefinitionSchemaID

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	definition, err := models.NewDefinitionSchemas(data.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
		})
		return
	}

	// 模型解引用
	isUnRefData := IsUnRefData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&isUnRefData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := models.DefinitionsSchemaUnRefByCollections(definition, isUnRefData.IsUnRef); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := models.DefinitionsSchemaUnRefByDefinitionsResponse(definition, isUnRefData.IsUnRef); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := models.DefinitionsSchemaUnRefByDefinitionsSchema(definition, isUnRefData.IsUnRef); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := definition.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.DeleteFail"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func DefinitionSchemasGet(ctx *gin.Context) {
	var data DefinitionSchemaID

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	definition, err := models.NewDefinitionSchemas(data.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
		})
		return
	}

	schema := make(map[string]interface{})
	if err := json.Unmarshal([]byte(definition.Schema), &schema); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
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

func DefinitionSchemasCopy(ctx *gin.Context) {
	currentMember, _ := ctx.Get("CurrentMember")
	if !currentMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var data DefinitionSchemaID

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	oldDefinition, err := models.NewDefinitionSchemas(data.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
		})
		return
	}

	schema := map[string]interface{}{}
	if err := json.Unmarshal([]byte(oldDefinition.Schema), &schema); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.CopyFail"}),
		})
		return
	}

	newDefinition, _ := models.NewDefinitionSchemas()
	newDefinition.ProjectId = oldDefinition.ProjectId
	newDefinition.Name = oldDefinition.Name
	newDefinition.Description = oldDefinition.Description
	newDefinition.Type = oldDefinition.Type
	newDefinition.Schema = oldDefinition.Schema
	if err := newDefinition.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.CopyFail"}),
		})
		return
	}

	newDefinition.Name = fmt.Sprintf("%s_%s", newDefinition.Name, strconv.Itoa(int(newDefinition.ID)))
	if err := newDefinition.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.CopyFail"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          newDefinition.ID,
		"parent_id":   newDefinition.ParentId,
		"name":        newDefinition.Name,
		"description": newDefinition.Description,
		"type":        newDefinition.Type,
		"schema":      schema,
		"created_at":  newDefinition.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by":  newDefinition.Creator(),
		"updated_at":  newDefinition.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by":  newDefinition.Updater(),
	})
}

func DefinitionSchemasMove(ctx *gin.Context) {
	currentMember, _ := ctx.Get("CurrentMember")
	if !currentMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var data DefinitionSchemaMove

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for i, id := range data.Target.Ids {
		if definition, err := models.NewDefinitionSchemas(id); err == nil {
			definition.ParentId = data.Target.Pid
			definition.DisplayOrder = i
			definition.Save()
		}
	}

	if data.Target.Pid != data.Origin.Pid {
		for i, id := range data.Origin.Ids {
			if definition, err := models.NewDefinitionSchemas(id); err == nil {
				definition.ParentId = data.Origin.Pid
				definition.DisplayOrder = i
				definition.Save()
			}
		}
	}

	ctx.Status(http.StatusCreated)
}
