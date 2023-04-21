package api

import (
	"encoding/json"
	"net/http"

	"github.com/apicat/apicat/commom/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type GlobalParameterDetails struct {
	ID       uint                  `uri:"id" binding:"required"`
	In       string                `json:"in" binding:"required,oneof=header query path cookie"`
	Name     string                `json:"name" binding:"required,lte=255"`
	Required bool                  `json:"required" binding:"required"`
	Schema   GlobalParameterSchema `json:"schema" binding:"required"`
}

type GlobalParameterSchema struct {
	Type        string `json:"type" binding:"required,oneof=string number integer array"`
	Default     string `json:"default" binding:"omitempty,lte=255"`
	Example     string `json:"example" binding:"omitempty,lte=255"`
	Description string `json:"description" binding:"omitempty,lte=255"`
}

type GlobalParametersCreateData struct {
	In       string                `json:"in" binding:"required,oneof=header query path cookie"`
	Name     string                `json:"name" binding:"required,lte=255"`
	Required bool                  `json:"required" binding:"required"`
	Schema   GlobalParameterSchema `json:"schema" binding:"required"`
}

func GlobalParametersCreate(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	var data GlobalParametersCreateData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
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
			"message": err.Error(),
		})
		return
	}

	globalParameters := &models.GlobalParameters{
		ProjectID: project.ID,
		In:        data.In,
		Name:      data.Name,
		Required:  required,
		Schema:    string(jsonSchema),
	}

	if err := globalParameters.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
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
