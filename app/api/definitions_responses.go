package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apicat/apicat/commom/apicat_struct"
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

func DefinitionsResponsesCreate(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	data := apicat_struct.ResponseObject{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	definitionsResponses, _ := models.NewDefinitionsResponses()
	definitionsResponses.ProjectID = project.ID
	definitionsResponses.Name = data.Name

	count, err := definitionsResponses.GetCountByName()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionsResponses.NameExists"}),
		})
		return
	}

	definitionsResponses.Code = data.Code
	definitionsResponses.Description = data.Description

	header, err := json.Marshal(data.Header)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	definitionsResponses.Header = string(header)

	content, err := json.Marshal(data.Content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	definitionsResponses.Content = string(content)

	if err := definitionsResponses.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          definitionsResponses.ID,
		"name":        definitionsResponses.Name,
		"code":        definitionsResponses.Code,
		"description": definitionsResponses.Description,
		"header":      data.Header,
		"content":     data.Content,
	})
}

func DefinitionsResponsesUpdate(ctx *gin.Context) {
	data := apicat_struct.ResponseObject{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	dr := DefinitionsResponsesID{}
	definitionsResponses, err := dr.CheckDefinitionsResponses(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	definitionsResponses.Name = data.Name
	count, err := definitionsResponses.GetCountExcludeTheID()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionsResponses.NameExists"}),
		})
		return
	}

	definitionsResponses.Code = data.Code
	definitionsResponses.Description = data.Description

	header, err := json.Marshal(data.Header)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	definitionsResponses.Header = string(header)

	content, err := json.Marshal(data.Content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	definitionsResponses.Content = string(content)

	if err := definitionsResponses.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func DefinitionsResponsesDelete(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	dr := DefinitionsResponsesID{}
	definitionsResponses, err := dr.CheckDefinitionsResponses(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	header := []*apicat_struct.Header{}
	if err := json.Unmarshal([]byte(definitionsResponses.Header), &header); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	content := apicat_struct.BodyObject{}
	if err := json.Unmarshal([]byte(definitionsResponses.Content), &content); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	responseDetail := apicat_struct.ResponseObject{
		Name:        definitionsResponses.Name,
		Code:        definitionsResponses.Code,
		Description: definitionsResponses.Description,
		Header:      header,
		Content:     content,
	}
	fmt.Printf("responseDetail: %+v\n", responseDetail)
	collections, _ := models.NewCollections()
	collections.ProjectId = project.ID
	collectionList, err := collections.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, collection := range collectionList {
		if collection.Type == "http" {
			fmt.Printf("collection.Content: %+v\n", collection.Content)
			docContent := []map[string]interface{}{}
			if err := json.Unmarshal([]byte(collection.Content), &docContent); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}

			var response []byte
			for _, v := range docContent {
				if v["type"] == "apicat-http-response" {
					response, err = json.Marshal(v["attrs"])
					if err != nil {
						ctx.JSON(http.StatusBadRequest, gin.H{
							"message": err.Error(),
						})
						return
					}
				}
			}

			apicatResponseList := apicat_struct.ResponseObjectList{}
			if err := json.Unmarshal(response, &apicatResponseList); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}

			apicatResponseList.Dereference(&responseDetail)
		}
	}

	if err := definitionsResponses.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
