package api

import (
	"encoding/json"
	"net/http"

	"github.com/apicat/apicat/backend/common/array_operations"
	"github.com/apicat/apicat/backend/common/spec"
	"github.com/apicat/apicat/backend/common/spec/jsonschema"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/apicat/apicat/backend/models"
	"github.com/gin-gonic/gin"
)

type GlobalParameterDetails struct {
	ID       uint            `json:"id" binding:"required"`
	In       string          `json:"in" binding:"required,oneof=header query path cookie"`
	Name     string          `json:"name" binding:"required,lte=255"`
	Required bool            `json:"required"`
	Schema   ParameterSchema `json:"schema" binding:"required"`
}

type ParameterSchema struct {
	Type        string `json:"type" binding:"required,oneof=string number integer array"`
	Default     string `json:"default" binding:"omitempty,lte=255"`
	Example     string `json:"example" binding:"omitempty,lte=255"`
	Description string `json:"description" binding:"omitempty,lte=255"`
}

type GlobalParametersData struct {
	In       string          `json:"in" binding:"required,oneof=header query path cookie"`
	Name     string          `json:"name" binding:"required,lte=255"`
	Required bool            `json:"required"`
	Schema   ParameterSchema `json:"schema" binding:"required"`
}

type GlobalParametersID struct {
	ParameterID uint `uri:"parameter-id" binding:"required,gt=0"`
}

func (gp *GlobalParametersID) CheckGlobalParameters(ctx *gin.Context) (*models.GlobalParameters, error) {
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&gp)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.NotFound"}),
		})
		return nil, err
	}

	globalParameters, err := models.NewGlobalParameters(gp.ParameterID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.NotFound"}),
		})
		return nil, err
	}

	return globalParameters, nil
}

func GlobalParametersList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	globalParameters := &models.GlobalParameters{
		ProjectID: project.ID,
	}
	globalParametersList, err := globalParameters.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.QueryFailed"}),
		})
		return
	}

	result := map[string][]GlobalParameterDetails{}
	result["header"] = []GlobalParameterDetails{}
	result["cookie"] = []GlobalParameterDetails{}
	result["path"] = []GlobalParameterDetails{}
	result["query"] = []GlobalParameterDetails{}
	for _, v := range globalParametersList {
		var schema ParameterSchema
		if err := json.Unmarshal([]byte(v.Schema), &schema); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
			})
			return
		}

		result[v.In] = append(result[v.In], GlobalParameterDetails{
			ID:       v.ID,
			In:       v.In,
			Name:     v.Name,
			Required: v.Required == 1,
			Schema:   schema,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func GlobalParametersCreate(ctx *gin.Context) {
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

	var data GlobalParametersData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	globalParameters, _ := models.NewGlobalParameters()
	globalParameters.ProjectID = project.ID
	globalParameters.Name = data.Name
	globalParameters.In = data.In
	count, err := globalParameters.GetCountByName()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.QueryFailed"}),
		})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.NameExists"}),
		})
		return
	}

	required := 0
	if data.Required {
		required = 1
	}

	jsonSchema, err := json.Marshal(data.Schema)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}

	globalParameters.Required = required
	globalParameters.Schema = string(jsonSchema)
	if err := globalParameters.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.CreateFailed"}),
		})
		return
	}

	globalParameterDetails := &GlobalParameterDetails{
		ID:       globalParameters.ID,
		In:       globalParameters.In,
		Name:     globalParameters.Name,
		Required: data.Required,
		Schema:   data.Schema,
	}

	ctx.JSON(http.StatusCreated, globalParameterDetails)
}

func GlobalParametersUpdate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var data GlobalParametersData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	gp := GlobalParametersID{}
	globalParameters, err := gp.CheckGlobalParameters(ctx)
	if err != nil {
		return
	}

	globalParameters.Name = data.Name
	count, err := globalParameters.GetCountExcludeTheID()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.QueryFailed"}),
		})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.NameExists"}),
		})
		return
	}

	required := 0
	if data.Required {
		required = 1
	}
	jsonSchema, err := json.Marshal(data.Schema)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}
	globalParameters.Required = required
	globalParameters.Schema = string(jsonSchema)

	if err := globalParameters.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.UpdateFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

type IsUnRefData struct {
	IsUnRef int `form:"is_unref"`
}

func GlobalParametersDelete(ctx *gin.Context) {
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

	data := IsUnRefData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	gp := GlobalParametersID{}
	globalParameter, err := gp.CheckGlobalParameters(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var paramScheam *jsonschema.Schema
	if err := json.Unmarshal([]byte(globalParameter.Schema), &paramScheam); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	param := &spec.Schema{
		Name:     globalParameter.Name,
		Required: globalParameter.Required == 1,
		Schema:   paramScheam,
	}

	collections, _ := models.NewCollections()
	collections.ProjectId = project.ID
	collectionList, err := collections.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.QueryFailed"})})
		return
	}

	for _, collection := range collectionList {
		if collection.Type == "http" {
			var content []*spec.NodeProxy
			if err := json.Unmarshal([]byte(collection.Content), &content); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			var request *spec.HTTPNode[spec.HTTPRequestNode]
			for _, i := range content {
				switch nx := i.Node.(type) {
				case *spec.HTTPNode[spec.HTTPRequestNode]:
					request = nx
				}
			}

			request.Attrs.GlobalExcepts[globalParameter.In] = array_operations.Remove[int64](int64(globalParameter.ID), request.Attrs.GlobalExcepts[globalParameter.In])
			if data.IsUnRef == 1 {
				if !array_operations.InArray[int64](int64(globalParameter.ID), request.Attrs.GlobalExcepts[globalParameter.In]) {
					request.Attrs.Parameters.Add(globalParameter.In, param)
				}
			}

			newContent, err := json.Marshal(content)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if err := collection.UpdateContent(false, collection.Title, string(newContent), currentProjectMember.(*models.ProjectMembers).UserID); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
		}
	}

	if err := globalParameter.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.DeleteFailed"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
