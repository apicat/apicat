package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/apicat/apicat/commom/translator"
	"github.com/apicat/apicat/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Content struct {
	Name     string `json:"name" binding:"required,lte=255"`
	Required bool   `json:"required,omitempty"`
	Schema   Schema `json:"schema" binding:"required"`
}

type Schema struct {
	Type        string `json:"type" binding:"required,oneof=integer number string array"`
	Default     string `json:"default,omitempty" binding:"lte=255"`
	Example     string `json:"example,omitempty" binding:"lte=255"`
	Description string `json:"description,omitempty" binding:"lte=255"`
}

type SetParameter struct {
	In      string    `json:"in" binding:"required,oneof=header cookie query"`
	Content []Content `json:"content" binding:"required,dive"`
}

func PublicParametersList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	parameter, _ := models.NewCommons()
	parameter.ProjectId = project.ID
	parameter.Type = "parameter"
	if err := parameter.Get(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{
				"header": []Content{},
				"cookie": []Content{},
				"query":  []Content{},
			})
			return
		}

		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := map[string][]Content{}
	if len(parameter.Content) == 0 {
		parameter.Content = "{}"
	}
	if err := json.Unmarshal([]byte(parameter.Content), &result); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if _, ok := result["header"]; !ok {
		result["header"] = []Content{}
	}
	if _, ok := result["cookie"]; !ok {
		result["cookie"] = []Content{}
	}
	if _, ok := result["query"]; !ok {
		result["query"] = []Content{}
	}

	ctx.JSON(http.StatusOK, result)
}

func PublicParametersSettings(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	data := SetParameter{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	result := map[string][]Content{}

	parameter, _ := models.NewCommons()
	parameter.ProjectId = project.ID
	parameter.Type = "parameter"
	if !parameter.IsExist() {
		result[data.In] = data.Content
		jsonString, err := json.Marshal(result)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		parameter.Content = string(jsonString)
		if err := parameter.Create(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.Status(http.StatusCreated)
		return
	}

	if err := parameter.Get(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if len(parameter.Content) > 0 {
		if err := json.Unmarshal([]byte(parameter.Content), &result); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	result[data.In] = data.Content

	jsonString, err := json.Marshal(result)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	parameter.Content = string(jsonString)
	if err := parameter.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
