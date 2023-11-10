package ai

import (
	"encoding/json"
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/definition"
	"github.com/apicat/apicat/backend/model/iteration"
	"github.com/apicat/apicat/backend/model/project"
	"net/http"
	"strconv"
	"time"

	"github.com/apicat/apicat/backend/app/util"
	"github.com/apicat/apicat/backend/common/openai"
	"github.com/apicat/apicat/backend/common/spec/plugin/openapi"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/config"
	"github.com/apicat/apicat/backend/enum"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type AICreateCollectionStructure struct {
	ParentID    uint   `json:"parent_id" binding:"gte=0"`              // 父级id
	Title       string `json:"title" binding:"required,lte=255"`       // 名称
	SchemaID    uint   `json:"schema_id" binding:"gte=0"`              // 模型id
	Path        string `json:"path" binding:"lte=255"`                 // 请求路径
	Method      string `json:"method" binding:"lte=255"`               // 请求方法
	IterationID string `json:"iteration_id" binding:"omitempty,gte=0"` // 迭代id
}

type AICreateSchemaStructure struct {
	ParentID uint   `json:"parent_id" binding:"gte=0"`       // 父级id
	Name     string `json:"name" binding:"required,lte=255"` // 名称
}

type AICreateApiNameStructure struct {
	SchemaID uint `form:"schema_id" binding:"gt=0"` // 模型id
}

func AICreateCollection(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		openapiContent string
		schema         *definition.DefinitionSchemas
		err            error
	)

	data := &AICreateCollectionStructure{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	lang := util.GetUserLanguage(ctx)

	if data.IterationID != "" {
		_, err := iteration.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.CollectionCreateFail"}),
			})
			return
		}
	}

	if data.SchemaID > 0 {
		schema, err = definition.NewDefinitionSchemas(data.SchemaID)
		if err != nil {
			slog.DebugCtx(ctx, "DefinitionSchemas get failed", slog.String("err", err.Error()), slog.String("SchemaID", strconv.Itoa(int(data.SchemaID))))
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
			})
			return
		}

		o := openai.NewOpenAI(config.GetSysConfig().OpenAI, lang)
		o.SetMaxTokens(3000)
		openapiContent, err = o.CreateApiBySchema(data.Title, data.Path, data.Method, schema.Schema)
		if err != nil || openapiContent == "" {
			slog.DebugCtx(ctx, "CreateApiBySchema Failed", slog.String("err", err.Error()), slog.String("openapiContent", openapiContent))
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.CollectionCreateFail"}),
			})
			return
		}
	} else {
		o := openai.NewOpenAI(config.GetSysConfig().OpenAI, lang)
		o.SetMaxTokens(2000)
		openapiContent, err = o.CreateApi(data.Title)
		if err != nil || openapiContent == "" {
			slog.DebugCtx(ctx, "CreateApi Failed", slog.String("err", err.Error()), slog.String("openapiContent", openapiContent))
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.CollectionCreateFail"}),
			})
			return
		}
	}

	content, err := openapi.Decode([]byte(openapiContent))
	if err != nil {
		slog.DebugCtx(ctx, "JSON Unmarshal Failed", slog.String("err", err.Error()), slog.String("openapiContent", openapiContent))
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.CollectionCreateFail"}),
		})
		return
	}

	if len(content.Collections) == 0 {
		slog.DebugCtx(ctx, "No collection item")
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.CollectionCreateFail"}),
		})
		return
	}

	currentProject, _ := ctx.Get("CurrentProject")
	refContentVirtualIDToId := &model.RefContentVirtualIDToId{
		DefinitionSchemas:    definition.DefinitionSchemasImport(currentProject.(*project.Projects).ID, content.Definitions.Schemas),
		DefinitionResponses:  definition.DefinitionResponsesImport(currentProject.(*project.Projects).ID, content.Definitions.Responses),
		DefinitionParameters: definition.DefinitionParametersImport(currentProject.(*project.Projects).ID, content.Definitions.Parameters),
	}
	records := collection.CollectionsImport(currentProject.(*project.Projects).ID, data.ParentID, content.Collections, refContentVirtualIDToId)

	if len(records) == 0 {
		slog.DebugCtx(ctx, "CollectionsImport Failed")
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.CollectionCreateFail"}),
		})
		return
	}

	if data.IterationID != "" {
		for _, v := range records {
			i, err := iteration.NewIterations(data.IterationID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.CollectionCreateFail"}),
				})
				return
			}

			ia, _ := iteration.NewIterationApis()
			ia.IterationID = i.ID
			ia.CollectionID = v.ID
			ia.CollectionType = v.Type
			if err := ia.Create(); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.CollectionCreateFail"}),
				})
				return
			}
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         records[0].ID,
		"parent_id":  records[0].ParentId,
		"title":      records[0].Title,
		"type":       records[0].Type,
		"content":    records[0].Content,
		"created_at": records[0].CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by": records[0].Creator(),
		"updated_at": records[0].UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by": records[0].Updater(),
	})
}

func AICreateSchema(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		openapiContent string
		err            error
	)

	type jsonSchema struct {
		Title       string                 `json:"title"`
		Description string                 `json:"description"`
		Type        string                 `json:"type"`
		Required    []string               `json:"required"`
		Format      string                 `json:"format"`
		Properties  map[string]interface{} `json:"properties"`
		Items       interface{}            `json:"items"`
		Example     interface{}            `json:"example"`
	}

	data := &AICreateSchemaStructure{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	lang := util.GetUserLanguage(ctx)
	o := openai.NewOpenAI(config.GetSysConfig().OpenAI, lang)
	o.SetMaxTokens(2000)
	openapiContent, err = o.CreateSchema(data.Name)
	if err != nil || openapiContent == "" {
		slog.DebugCtx(ctx, "CreateSchema Failed", slog.String("err", err.Error()), slog.String("openapiContent", openapiContent))
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.SchemaCreateFail"}),
		})
		return
	}

	js := &jsonSchema{}
	if err := json.Unmarshal([]byte(openapiContent), js); err != nil {
		slog.DebugCtx(ctx, "JSON Unmarshal Failed", slog.String("err", err.Error()), slog.String("openapiContent", openapiContent))
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.SchemaCreateFail"}),
		})
		return
	}

	p, _ := ctx.Get("CurrentProject")
	d, _ := definition.NewDefinitionSchemas()
	d.ProjectId = p.(*project.Projects).ID
	d.Name = js.Title
	definitions, err := d.List()
	if err != nil {
		slog.DebugCtx(ctx, "definitions search Failed", slog.String("err", err.Error()), slog.String("ProjectId", strconv.FormatUint(uint64(d.ProjectId), 10)), slog.String("Name", d.Name))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.SchemaCreateFail"}),
		})
		return
	}
	if len(definitions) > 0 {
		d.Name = d.Name + time.Now().Format("20060102150405")
	}

	d.Description = js.Description
	d.Type = "schema"
	d.Schema = openapiContent
	if err := d.Create(); err != nil {
		slog.DebugCtx(ctx, "definition Create Failed", slog.String("err", err.Error()), slog.String("openapiContent", openapiContent))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.CreateFail"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          d.ID,
		"parent_id":   d.ParentId,
		"name":        d.Name,
		"description": d.Description,
		"type":        d.Type,
		"schema":      d.Schema,
		"created_at":  d.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by":  d.Creator(),
		"updated_at":  d.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by":  d.Updater(),
	})
}

func AICreateApiNames(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		openapiContent string
		err            error
	)

	data := &AICreateApiNameStructure{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	schema, err := definition.NewDefinitionSchemas(data.SchemaID)
	if err != nil {
		slog.DebugCtx(ctx, "DefinitionSchemas get failed", slog.String("err", err.Error()), slog.String("SchemaID", strconv.Itoa(int(data.SchemaID))))
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
		})
		return
	}

	lang := util.GetUserLanguage(ctx)
	o := openai.NewOpenAI(config.GetSysConfig().OpenAI, lang)
	openapiContent, err = o.ListApiBySchema(schema.Name)
	if err != nil || openapiContent == "" {
		slog.DebugCtx(ctx, "ListApiBySchema Failed", slog.String("err", err.Error()), slog.String("openapiContent", openapiContent))
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.CollectionCreateFail"}),
		})
		return
	}

	var arr []map[string]string
	if err := json.Unmarshal([]byte(openapiContent), &arr); err != nil {
		slog.DebugCtx(ctx, "JSON Unmarshal Failed", slog.String("err", err.Error()), slog.String("openapiContent", openapiContent))
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "AI.CollectionCreateFail"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, arr)
}
