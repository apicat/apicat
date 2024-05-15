package project

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/model/user"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	projectbase "github.com/apicat/apicat/v2/backend/route/proto/project/base"
	projectresponse "github.com/apicat/apicat/v2/backend/route/proto/project/response"
	prototeambase "github.com/apicat/apicat/v2/backend/route/proto/team/base"
	protouserbase "github.com/apicat/apicat/v2/backend/route/proto/user/base"
	protouserresponse "github.com/apicat/apicat/v2/backend/route/proto/user/response"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/plugin/openapi"
	"github.com/apicat/apicat/v2/backend/module/spec/plugin/postman"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

func convertModelGlobalparameter(gp *global.GlobalParameter) *projectresponse.GlobalParameter {
	return &projectresponse.GlobalParameter{
		OnlyIdInfo: protobase.OnlyIdInfo{
			ID: gp.ID,
		},
		GlobalParameterDataOption: projectbase.GlobalParameterDataOption{
			In:       gp.In,
			Name:     gp.Name,
			Required: gp.Required,
			Schema:   gp.Schema,
		},
	}
}

func convertModelProjectServer(s *project.Server) *projectresponse.ProjectServer {
	return &projectresponse.ProjectServer{
		OnlyIdInfo: protobase.OnlyIdInfo{
			ID: s.ID,
		},
		ProjectServerDataOption: projectbase.ProjectServerDataOption{
			URL:         s.URL,
			Description: s.Description,
		},
	}
}

func convertModelProjectMember(pm *project.ProjectMember, memberInfo *team.TeamMember, userInfo *user.User) *projectresponse.ProjectMember {
	return &projectresponse.ProjectMember{
		IdCreateTimeInfo: protobase.IdCreateTimeInfo{
			ID:        pm.ID,
			CreatedAt: pm.CreatedAt.Unix(),
		},
		TeamMemberStatusOption: prototeambase.TeamMemberStatusOption{
			Status: memberInfo.Status,
		},
		ProjectMemberPermission: protobase.ProjectMemberPermission{
			Permission: pm.Permission,
		},
		User: protouserresponse.UserData{
			EmailOption: protouserbase.EmailOption{
				Email: userInfo.Email,
			},
			NameOption: protouserbase.NameOption{
				Name: userInfo.Name,
			},
			AvatarOption: protouserbase.AvatarOption{
				Avatar: userInfo.Avatar,
			},
		},
	}
}

func convertModelProjectGroup(pg *project.ProjectGroup) *projectresponse.ProjectGroup {
	return &projectresponse.ProjectGroup{
		OnlyIdInfo: protobase.OnlyIdInfo{
			ID: pg.ID,
		},
		ProjectGroupNameOption: projectbase.ProjectGroupNameOption{
			Name: pg.Name,
		},
	}
}

func convertModelDefinitionSchema(ds *definition.DefinitionSchema, cUserInfo, uUserInfo *user.User) *projectresponse.DefinitionSchema {
	return &projectresponse.DefinitionSchema{
		EmbedInfo: protobase.EmbedInfo{
			ID:        ds.ID,
			CreatedAt: ds.CreatedAt.Unix(),
			UpdatedAt: ds.UpdatedAt.Unix(),
		},
		DefinitionSchemaDataOption: projectbase.DefinitionSchemaDataOption{
			Name:        ds.Name,
			Schema:      ds.Schema,
			Description: ds.Description,
		},
		DefinitionSchemaParentIDOption: projectbase.DefinitionSchemaParentIDOption{
			ParentID: ds.ParentID,
		},
		DefinitionSchemaTypeOption: projectbase.DefinitionSchemaTypeOption{
			Type: ds.Type,
		},
		OperatorID: projectbase.OperatorID{
			CreatedBy: cUserInfo.Name,
			UpdatedBy: uUserInfo.Name,
		},
	}
}

func convertModelDefinitionResponse(dr *definition.DefinitionResponse, cUserInfo, uUserInfo *user.User) *projectresponse.DefinitionResponse {
	return &projectresponse.DefinitionResponse{
		EmbedInfo: protobase.EmbedInfo{
			ID:        dr.ID,
			CreatedAt: dr.CreatedAt.Unix(),
			UpdatedAt: dr.UpdatedAt.Unix(),
		},
		DefinitionResponseDataOption: projectbase.DefinitionResponseDataOption{
			Name:        dr.Name,
			Description: dr.Description,
			Header:      dr.Header,
			Content:     dr.Content,
		},
		DefinitionResponseParentIDOption: projectbase.DefinitionResponseParentIDOption{
			ParentID: dr.ParentID,
		},
		DefinitionResponseTypeOption: projectbase.DefinitionResponseTypeOption{
			Type: dr.Type,
		},
		OperatorID: projectbase.OperatorID{
			CreatedBy: cUserInfo.Name,
			UpdatedBy: uUserInfo.Name,
		},
	}
}

func convertModelDefinitionSchemaHistory(h *definition.DefinitionSchemaHistory, userInfo *user.User) *projectresponse.DefinitionSchemaHistory {
	return &projectresponse.DefinitionSchemaHistory{
		IdCreateTimeInfo: protobase.IdCreateTimeInfo{
			ID:        h.ID,
			CreatedAt: h.CreatedAt.Unix(),
		},
		DefinitionSchemaHistoryData: projectresponse.DefinitionSchemaHistoryData{
			DefinitionSchemaDataOption: projectbase.DefinitionSchemaDataOption{
				Name:        h.Name,
				Description: h.Description,
				Schema:      h.Schema,
			},
			SchemaID: h.SchemaID,
		},
		CreatedBy: userInfo.Name,
	}
}

func buildDefinitionSchemaTree(parentID uint, schemas []*definition.DefinitionSchema) projectresponse.DefinitionSchemaTree {
	result := make(projectresponse.DefinitionSchemaTree, 0)

	for _, s := range schemas {
		if s.ParentID == parentID {
			children := buildDefinitionSchemaTree(s.ID, schemas)

			cl := projectresponse.DefinitionSchemaNode{
				EmbedInfo: protobase.EmbedInfo{
					ID:        s.ID,
					CreatedAt: s.CreatedAt.Unix(),
					UpdatedAt: s.UpdatedAt.Unix(),
				},
				DefinitionSchemaDataOption: projectbase.DefinitionSchemaDataOption{
					Name:        s.Name,
					Schema:      s.Schema,
					Description: s.Description,
				},
				DefinitionSchemaParentIDOption: projectbase.DefinitionSchemaParentIDOption{
					ParentID: s.ParentID,
				},
				DefinitionSchemaTypeOption: projectbase.DefinitionSchemaTypeOption{
					Type: s.Type,
				},
				Items: children,
			}

			result = append(result, &cl)
		}
	}

	return result
}

func buildDefinitionResponseTree(parentID uint, responses []*definition.DefinitionResponse) projectresponse.DefinitionResponseTree {
	result := make(projectresponse.DefinitionResponseTree, 0)

	for _, r := range responses {
		if r.ParentID == parentID {
			children := buildDefinitionResponseTree(r.ID, responses)

			cl := projectresponse.DefinitionResponseNode{
				EmbedInfo: protobase.EmbedInfo{
					ID:        r.ID,
					CreatedAt: r.CreatedAt.Unix(),
					UpdatedAt: r.UpdatedAt.Unix(),
				},
				DefinitionResponseDataOption: projectbase.DefinitionResponseDataOption{
					Name:        r.Name,
					Description: r.Description,
					Header:      r.Header,
					Content:     r.Content,
				},
				DefinitionResponseParentIDOption: projectbase.DefinitionResponseParentIDOption{
					ParentID: r.ParentID,
				},
				DefinitionResponseTypeOption: projectbase.DefinitionResponseTypeOption{
					Type: r.Type,
				},
				Items: children,
			}

			result = append(result, &cl)
		}
	}

	return result
}

// openapi & swagger 文件解析
func openapiAndSwaggerFileParse(fileContent string) (*spec.Spec, error) {
	var (
		base64Content string
		rawContent    []byte
		err           error
	)

	if strings.Contains(fileContent, "data:application/json;base64,") {
		base64Content = strings.Replace(fileContent, "data:application/json;base64,", "", 1)
	} else {
		base64Content = strings.Replace(fileContent, "data:application/x-yaml;base64,", "", 1)
	}

	rawContent, err = base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return nil, ginrpc.NewError(http.StatusBadRequest, err)
	}

	return openapi.Parse(rawContent)
}

// apicat 文件解析
func apicatFileParse(fileContent string) (*spec.Spec, error) {
	var (
		base64Content string
		rawContent    []byte
		err           error
	)

	if strings.Contains(fileContent, "data:application/json;base64,") {
		base64Content = strings.Replace(fileContent, "data:application/json;base64,", "", 1)
	}

	rawContent, err = base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return nil, ginrpc.NewError(http.StatusBadRequest, err)
	}

	return spec.ParseJSON(rawContent)
}

// postman 文件解析
func postmanFileParse(fileContent string) (*spec.Spec, error) {
	var (
		base64Content string
		rawContent    []byte
		err           error
	)

	if strings.Contains(fileContent, "data:application/json;base64,") {
		base64Content = strings.Replace(fileContent, "data:application/json;base64,", "", 1)
	}

	rawContent, err = base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return nil, ginrpc.NewError(http.StatusBadRequest, err)
	}

	return postman.Import(rawContent)
}

func dsDerefWithSpec(ctx *gin.Context, ds *definition.DefinitionSchema) (*spec.DefinitionModel, error) {
	schemaSpec, err := ds.ToSpec()
	if err != nil {
		return nil, err
	}

	schemas, err := definition.GetDefinitionSchemasWithSpec(ctx, ds.ProjectID)
	if err != nil {
		return nil, err
	}

	schemaSpec.DeepDeref(schemas)

	return schemaSpec, nil
}
