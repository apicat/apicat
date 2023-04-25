package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/apicat/apicat/commom/spec"
	"github.com/apicat/apicat/commom/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type CommonResponsesID struct {
	CommonResponsesID uint `uri:"response-id" binding:"required,gt=0"`
}

func (cr *CommonResponsesID) CheckCommonResponses(ctx *gin.Context) (*models.CommonResponses, error) {
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&cr)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.NotFound"}),
		})
		return nil, err
	}

	commonResponses, err := models.NewCommonResponses(cr.CommonResponsesID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "response not found",
		})
		return nil, err
	}

	return commonResponses, nil
}

type ResponseDetailData struct {
	Name        string                 `json:"name" binding:"required,lte=255"`
	Code        int                    `json:"code" binding:"required"`
	Description string                 `json:"description" binding:"required,lte=255"`
	Header      []*spec.Schema         `json:"header,omitempty" binding:"omitempty,dive"`
	Content     map[string]spec.Schema `json:"content,omitempty" binding:"required"`
	Ref         string                 `json:"$ref,omitempty" binding:"omitempty,lte=255"`
}

func CommonResponsesList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	commonResponses, _ := models.NewCommonResponses()
	commonResponses.ProjectID = project.ID
	commonResponsesList, err := commonResponses.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := []map[string]interface{}{}
	for _, v := range commonResponsesList {
		header := []*spec.Schema{}
		if err := json.Unmarshal([]byte(v.Header), &header); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		content := map[string]spec.Schema{}
		if err := json.Unmarshal([]byte(v.Content), &content); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		result = append(result, map[string]interface{}{
			"id":          v.ID,
			"name":        v.Name,
			"code":        v.Code,
			"description": v.Description,
			"header":      header,
			"content":     content,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func CommonResponsesDetail(ctx *gin.Context) {
	cr := CommonResponsesID{}
	commonResponses, err := cr.CheckCommonResponses(ctx)
	if err != nil {
		return
	}

	header := []*spec.Schema{}
	if err := json.Unmarshal([]byte(commonResponses.Header), &header); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	content := map[string]spec.Schema{}
	if err := json.Unmarshal([]byte(commonResponses.Content), &content); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":          commonResponses.ID,
		"name":        commonResponses.Name,
		"code":        commonResponses.Code,
		"description": commonResponses.Description,
		"header":      header,
		"content":     content,
	})
}

func CommonResponsesCreate(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	data := ResponseDetailData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	commonResponses, _ := models.NewCommonResponses()
	commonResponses.ProjectID = project.ID
	commonResponses.Name = data.Name

	count, err := commonResponses.GetCountByName()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "CommonResponses.NameExists"}),
		})
		return
	}

	commonResponses.Code = data.Code
	commonResponses.Description = data.Description

	responseHeader := make([]*spec.Schema, 0)
	if len(data.Header) > 0 {
		responseHeader = data.Header
	}

	header, err := json.Marshal(responseHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	commonResponses.Header = string(header)

	content, err := json.Marshal(data.Content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	commonResponses.Content = string(content)

	if err := commonResponses.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          commonResponses.ID,
		"name":        commonResponses.Name,
		"code":        commonResponses.Code,
		"description": commonResponses.Description,
		"header":      data.Header,
		"content":     data.Content,
	})
}

func CommonResponsesUpdate(ctx *gin.Context) {
	data := ResponseDetailData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	cr := CommonResponsesID{}
	commonResponses, err := cr.CheckCommonResponses(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	commonResponses.Name = data.Name
	count, err := commonResponses.GetCountExcludeTheID()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "CommonResponses.NameExists"}),
		})
		return
	}

	commonResponses.Code = data.Code
	commonResponses.Description = data.Description

	header, err := json.Marshal(data.Header)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	commonResponses.Header = string(header)

	content, err := json.Marshal(data.Content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	commonResponses.Content = string(content)

	if err := commonResponses.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func CommonResponsesDelete(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	cr := CommonResponsesID{}
	commonResponses, err := cr.CheckCommonResponses(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	header := []*spec.Schema{}
	if err := json.Unmarshal([]byte(commonResponses.Header), &header); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	content := map[string]spec.Schema{}
	if err := json.Unmarshal([]byte(commonResponses.Content), &content); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	responseDetail := ResponseDetailData{
		Name:        commonResponses.Name,
		Code:        commonResponses.Code,
		Description: commonResponses.Description,
		Header:      header,
		Content:     content,
	}

	collections, _ := models.NewCollections()
	collections.ProjectId = project.ID
	collectionList, err := collections.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ref := "{\"$ref\":\"#/commons/responses/" + strconv.FormatUint(uint64(commonResponses.ID), 10) + "\"}"
	responseDetailJson, err := json.Marshal(responseDetail)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, collection := range collectionList {
		if collection.Type == "http" {
			if strings.Contains(collection.Content, ref) {
				newContent := strings.Replace(collection.Content, ref, string(responseDetailJson), -1)
				collection.Content = newContent
				if err := collection.Update(); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": err.Error(),
					})
					return
				}
			}
		}
	}

	if err := commonResponses.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
