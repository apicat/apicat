package api

import (
	"net/http"

	"github.com/apicat/apicat/commom/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type DefinitionsResponsesID struct {
	DefinitionsResponsesID uint `uri:"response-id" binding:"required,gt=0"`
}

func (dr *DefinitionsResponsesID) CheckDefinitionsResponses(ctx *gin.Context) (*models.DefinitionsResponses, error) {
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&dr)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.NotFound"}),
		})
		return nil, err
	}

	definitionsResponses, err := models.NewDefinitionsResponses(dr.DefinitionsResponsesID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "response not found",
		})
		return nil, err
	}

	return definitionsResponses, nil
}

func DefinitionsResponsesList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	definitionsResponses, _ := models.NewDefinitionsResponses()
	definitionsResponses.ProjectID = project.ID
	definitionsResponsesList, err := definitionsResponses.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := []map[string]interface{}{}
	for _, v := range definitionsResponsesList {
		result = append(result, map[string]interface{}{
			"id":          v.ID,
			"code":        v.Code,
			"description": v.Description,
		})
	}

	ctx.JSON(http.StatusOK, result)
}
