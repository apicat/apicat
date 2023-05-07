package api

import (
	"encoding/json"
	"net/http"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"

	"github.com/gin-gonic/gin"
)

type ResponseID struct {
	ResponseID uint `uri:"response-id" binding:"required,gte=0"`
}

func (r *ResponseID) CheckResponse(ctx *gin.Context) (*models.Commons, error) {
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&r)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Response.NotFound"}),
		})
		return nil, err
	}

	response, err := models.NewCommons(r.ResponseID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Response.NotFound"}),
		})
		return nil, err
	}

	return response, nil
}

type Responses struct {
	Code        int                             `json:"code" binding:"required"`
	Description string                          `json:"description" binding:"required,lte=255"`
	Header      []ResponsesHeader               `json:"header,omitempty" binding:"omitempty,dive"`
	Content     map[string]ResponsesContentType `json:"content" binding:"required"`
}

type ResponsesHeader struct {
	Name     string                `json:"name" binding:"required,lte=255"`
	Required bool                  `json:"required,omitempty"`
	Schema   ResponsesHeaderSchema `json:"schema" binding:"required"`
}

type ResponsesHeaderSchema struct {
	Type        string `json:"type" binding:"required,oneof=integer number string array"`
	Default     string `json:"default,omitempty" binding:"lte=255"`
	Example     string `json:"example,omitempty" binding:"lte=255"`
	Description string `json:"description,omitempty" binding:"lte=255"`
}

type ResponsesContentType struct {
	Ref    string                 `json:"$ref,omitempty" binding:"lte=255"`
	Schema ResponsesContentSchema `json:"schema"`
}

type ResponsesContentSchema struct {
	Type          string                  `json:"type" binding:"required,lte=255"`
	Required      *[]string               `json:"required,omitempty" binding:"omitempty,gte=0,dive,gte=0,lte=255"`
	Format        *string                 `json:"format,omitempty" binding:"omitempty,lte=255"`
	Properties    *map[string]interface{} `json:"properties,omitempty" binding:"omitempty,lte=255"`
	Items         *map[string]interface{} `json:"items,omitempty" binding:"omitempty,lte=255"`
	Description   *string                 `json:"description,omitempty" binding:"omitempty,lte=255"`
	Example       *string                 `json:"example,omitempty" binding:"omitempty,lte=255"`
	Ref           *string                 `json:"$ref,omitempty" binding:"omitempty,lte=255"`
	XApicatOrders *[]string               `json:"x-apicat-orders,omitempty" binding:"omitempty,gte=0,dive,gte=0,lte=255"`
	XApicatMock   *string                 `json:"x-apicat-mock,omitempty" binding:"omitempty,lte=255"`
}

func PublicResponsesList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	response, _ := models.NewCommons()
	response.ProjectId = project.ID
	response.Type = "response"
	if responseList, err := response.List(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		result := []map[string]interface{}{}
		for _, v := range responseList {
			data := Responses{}
			if err := json.Unmarshal([]byte(v.Content), &data); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}

			result = append(result, map[string]interface{}{
				"id":          v.ID,
				"description": data.Description,
				"code":        data.Code,
			})
		}

		ctx.JSON(http.StatusOK, result)
	}
}

func PublicResponsesDetails(ctx *gin.Context) {
	responseID := ResponseID{}
	if response, err := responseID.CheckResponse(ctx); err != nil {
		return
	} else {
		data := Responses{}
		if err := json.Unmarshal([]byte(response.Content), &data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, data)
	}
}

func PublicResponsesAdd(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	response := Responses{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&response)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, v := range response.Content {
		if v.Schema.Type == "object" {
			if v.Schema.XApicatOrders == nil || len(*v.Schema.Required) == 0 {
				v.Schema.XApicatOrders = &[]string{}
			}
			if v.Schema.Required == nil || len(*v.Schema.Required) == 0 {
				v.Schema.Required = &[]string{}
			}
			if v.Schema.Properties == nil || len(*v.Schema.Properties) == 0 {
				v.Schema.Properties = &map[string]interface{}{}
			}
		} else if v.Schema.Type == "array" {
			if v.Schema.Items == nil || len(*v.Schema.Items) == 0 {
				v.Schema.Items = &map[string]interface{}{}
			}
		}
	}

	content, err := json.Marshal(response)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	common, _ := models.NewCommons()
	common.ProjectId = project.ID
	common.Type = "response"
	common.Content = string(content)

	if err := common.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Responses.CreateFail"}),
		})
		return
	}

	result := map[string]interface{}{
		"id":          common.ID,
		"description": response.Description,
		"code":        response.Code,
		"content":     response.Content,
	}

	if len(response.Header) > 0 {
		result["header"] = response.Header
	}

	ctx.JSON(http.StatusCreated, result)
}

func PublicResponsesEdit(ctx *gin.Context) {
	responseID := ResponseID{}
	if response, err := responseID.CheckResponse(ctx); err != nil {
		return
	} else {
		responseData := Responses{}
		if err := json.Unmarshal([]byte(response.Content), &responseData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		editData := Responses{}
		if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&editData)); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		for _, v := range editData.Content {
			if v.Schema.Type == "object" {
				if v.Schema.XApicatOrders == nil || len(*v.Schema.Required) == 0 {
					v.Schema.XApicatOrders = &[]string{}
				}
				if v.Schema.Required == nil || len(*v.Schema.Required) == 0 {
					v.Schema.Required = &[]string{}
				}
				if v.Schema.Properties == nil || len(*v.Schema.Properties) == 0 {
					v.Schema.Properties = &map[string]interface{}{}
				}
			} else if v.Schema.Type == "array" {
				if v.Schema.Items == nil || len(*v.Schema.Items) == 0 {
					v.Schema.Items = &map[string]interface{}{}
				}
			}
		}

		responseData.Code = editData.Code
		responseData.Description = editData.Description
		responseData.Header = editData.Header
		responseData.Content = editData.Content

		content, err := json.Marshal(responseData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		response.Content = string(content)
		if err := response.Update(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Responses.UpdateFail"}),
			})
			return
		}

		ctx.Status(http.StatusCreated)
	}
}

func PublicResponsesDelete(ctx *gin.Context) {
	responseID := ResponseID{}
	if response, err := responseID.CheckResponse(ctx); err != nil {
		return
	} else {
		if err := response.Delete(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Responses.DeleteFail"}),
			})
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}
