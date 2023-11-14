package definition

import (
	"encoding/json"
	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/definition"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/module/spec"
	"github.com/apicat/apicat/backend/route/api/global"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefinitionResponsesID struct {
	DefinitionResponsesID uint `uri:"response-id" binding:"required,gt=0"`
}

func (cr *DefinitionResponsesID) CheckDefinitionResponses(ctx *gin.Context) (*definition.DefinitionResponses, error) {
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindUri(&cr)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "GlobalParameters.NotFound"}),
		})
		return nil, err
	}

	definitionResponses, err := definition.NewDefinitionResponses(cr.DefinitionResponsesID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "GlobalParameters.NotFound"}),
		})
		return nil, err
	}

	return definitionResponses, nil
}

func DefinitionResponsesList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	p, _ := currentProject.(*project.Projects)

	definitionResponses, _ := definition.NewDefinitionResponses()
	definitionResponses.ProjectID = p.ID
	definitionResponsesList, err := definitionResponses.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "DefinitionResponses.QueryFailed"}),
		})
		return
	}

	result := []map[string]interface{}{}
	for _, v := range definitionResponsesList {
		header := []*spec.Schema{}
		if v.Header != "" {
			if err := json.Unmarshal([]byte(v.Header), &header); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.ContentParsingFailed"}),
				})
				return
			}
		}

		content := map[string]spec.Schema{}
		if v.Content != "" {
			if err := json.Unmarshal([]byte(v.Content), &content); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.ContentParsingFailed"}),
				})
				return
			}
		}

		result = append(result, map[string]interface{}{
			"id":          v.ID,
			"name":        v.Name,
			"description": v.Description,
			"type":        v.Type,
			"header":      header,
			"content":     content,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func DefinitionResponsesDetail(ctx *gin.Context) {
	cr := DefinitionResponsesID{}
	definitionResponses, err := cr.CheckDefinitionResponses(ctx)
	if err != nil {
		return
	}

	header := []*spec.Schema{}
	if definitionResponses.Header != "" {
		if err := json.Unmarshal([]byte(definitionResponses.Header), &header); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.ContentParsingFailed"}),
			})
			return
		}
	}
	content := map[string]spec.Schema{}
	if definitionResponses.Content != "" {
		if err := json.Unmarshal([]byte(definitionResponses.Content), &content); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.ContentParsingFailed"}),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":          definitionResponses.ID,
		"name":        definitionResponses.Name,
		"description": definitionResponses.Description,
		"type":        definitionResponses.Type,
		"header":      header,
		"content":     content,
	})
}

func DefinitionResponsesCreate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	currentProject, _ := ctx.Get("CurrentProject")
	p, _ := currentProject.(*project.Projects)

	data := proto.ResponseDetailData{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	definitionResponses, _ := definition.NewDefinitionResponses()
	definitionResponses.ProjectID = p.ID
	definitionResponses.Name = data.Name
	definitionResponses.Type = data.Type

	count, err := definitionResponses.GetCountByName()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "DefinitionResponses.QueryFailed"}),
		})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.NameExists"}),
		})
		return
	}

	definitionResponses.Description = data.Description

	responseHeader := make([]*spec.Schema, 0)
	if len(data.Header) > 0 {
		responseHeader = data.Header
	}

	header, err := json.Marshal(responseHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}
	definitionResponses.Header = string(header)

	content, err := json.Marshal(data.Content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}
	definitionResponses.Content = string(content)

	if err := definitionResponses.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "DefinitionResponses.CreateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          definitionResponses.ID,
		"name":        definitionResponses.Name,
		"description": definitionResponses.Description,
		"type":        definitionResponses.Type,
		"header":      data.Header,
		"content":     data.Content,
	})
}

func DefinitionResponsesUpdate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := proto.ResponseDetailData{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	cr := DefinitionResponsesID{}
	definitionResponses, err := cr.CheckDefinitionResponses(ctx)
	if err != nil {
		return
	}

	definitionResponses.Name = data.Name
	count, err := definitionResponses.GetCountExcludeTheID()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "DefinitionResponses.QueryFailed"}),
		})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.NameExists"}),
		})
		return
	}

	definitionResponses.Description = data.Description

	header, err := json.Marshal(data.Header)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}
	definitionResponses.Header = string(header)

	content, err := json.Marshal(data.Content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}
	definitionResponses.Content = string(content)

	if err := definitionResponses.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "DefinitionResponses.UpdateFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func DefinitionResponsesDelete(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	cr := DefinitionResponsesID{}
	definitionResponses, err := cr.CheckDefinitionResponses(ctx)
	if err != nil {
		return
	}

	data := global.IsUnRefData{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.IsUnRef == 1 {
		err = collection.DefinitionsResponseUnRef(definitionResponses)
	} else {
		err = collection.DefinitionsResponseDelRef(definitionResponses)
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "DefinitionResponses.DeleteFailed"}),
		})
		return
	}

	if err := definitionResponses.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "DefinitionResponses.DeleteFailed"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
