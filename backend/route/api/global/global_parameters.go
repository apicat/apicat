package global

import (
	"encoding/json"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/global"
	"github.com/apicat/apicat/backend/model/project"
	"net/http"

	"github.com/apicat/apicat/backend/common/apicat_struct"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
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

func (gp *GlobalParametersID) CheckGlobalParameters(ctx *gin.Context) (*global.GlobalParameters, error) {
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&gp)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.NotFound"}),
		})
		return nil, err
	}

	globalParameters, err := global.NewGlobalParameters(gp.ParameterID)
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
	p, _ := currentProject.(*project.Projects)

	globalParameters := &global.GlobalParameters{
		ProjectID: p.ID,
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
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	currentProject, _ := ctx.Get("CurrentProject")
	p, _ := currentProject.(*project.Projects)

	var data GlobalParametersData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	globalParameters, _ := global.NewGlobalParameters()
	globalParameters.ProjectID = p.ID
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
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
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
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	currentProject, _ := ctx.Get("CurrentProject")
	p, _ := currentProject.(*project.Projects)

	gp := GlobalParametersID{}
	globalParameters, err := gp.CheckGlobalParameters(ctx)
	if err != nil {
		return
	}
	data := IsUnRefData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	globalParameterSchema := ParameterSchema{}
	if err := json.Unmarshal([]byte(globalParameters.Schema), &globalParameterSchema); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}

	parametersSchema := apicat_struct.SchemaObject{
		Name:     globalParameters.Name,
		Required: globalParameters.Required == 1,
		Schema: apicat_struct.Schema{
			Type:        globalParameterSchema.Type,
			Example:     globalParameterSchema.Example,
			Description: globalParameterSchema.Description,
		},
	}

	collections, _ := collection.NewCollections()
	collections.ProjectId = p.ID
	collectionList, err := collections.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.QueryFailed"}),
		})
		return
	}
	for _, c := range collectionList {
		if c.Type == "http" {
			// 解析文档内容
			docContent := []map[string]interface{}{}
			if err := json.Unmarshal([]byte(c.Content), &docContent); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
				})
				return
			}

			var request []byte
			for _, v := range docContent {
				if v["type"] == "apicat-http-request" {
					request, err = json.Marshal(v["attrs"])
					if err != nil {
						ctx.JSON(http.StatusBadRequest, gin.H{
							"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
						})
						return
					}
				}

			}

			apicatRequest := apicat_struct.RequestObject{}
			if err := json.Unmarshal(request, &apicatRequest); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
				})
				return
			}

			// 删除GlobalExcepts中的全局参数，检查id是否在GlobalExcepts中存在，如果存在则删除GlobalExcepts中这个id。如果不存在则将解引用的参数补充在parameters中的第一位
			switch globalParameters.In {
			case "query":
				if !apicatRequest.GlobalExcepts.CheckQueryRef(int(globalParameters.ID)) {
					if data.IsUnRef == 1 {
						apicatRequest.Parameters.CheckQueryRef(parametersSchema)
					}
				}
			case "header":
				if !apicatRequest.GlobalExcepts.CheckHeaderRef(int(globalParameters.ID)) {
					if data.IsUnRef == 1 {
						apicatRequest.Parameters.CheckHeaderRef(parametersSchema)
					}
				}
			case "path":
				if !apicatRequest.GlobalExcepts.CheckPathRef(int(globalParameters.ID)) {
					if data.IsUnRef == 1 {
						apicatRequest.Parameters.CheckPathRef(parametersSchema)
					}
				}
			case "cookie":
				if !apicatRequest.GlobalExcepts.CheckCookieRef(int(globalParameters.ID)) {
					if data.IsUnRef == 1 {
						apicatRequest.Parameters.CheckCookieRef(parametersSchema)
					}
				}
			default:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.TypeDoesNotExist"}),
				})
				return
			}

			// 将修改后的参数重新写入文档内容
			for i, v := range docContent {
				if v["type"] == "apicat-http-request" {
					docContent[i]["attrs"] = apicatRequest
				}
			}

			if newContent, err := json.Marshal(docContent); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
				})
				return
			} else {
				c.Content = string(newContent)
				if err := c.Update(); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.UpdateFailed"}),
					})
					return
				}
			}
		}

		if err := globalParameters.Delete(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "GlobalParameters.DeleteFailed"}),
			})
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}
