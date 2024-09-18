package project

import (
	"strings"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/core/content_suggestion"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	protoproject "github.com/apicat/apicat/v2/backend/route/proto/project"
	projectbase "github.com/apicat/apicat/v2/backend/route/proto/project/base"
	projectrequest "github.com/apicat/apicat/v2/backend/route/proto/project/request"
	projectresponse "github.com/apicat/apicat/v2/backend/route/proto/project/response"
	"github.com/apicat/apicat/v2/backend/service/relations"
	"github.com/apicat/apicat/v2/backend/utils/onetime_token"

	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/plugin/export"
	"github.com/apicat/apicat/v2/backend/module/spec/plugin/openapi"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type projectApiImpl struct{}

func NewProjectApi() protoproject.ProjectApi {
	return &projectApiImpl{}
}

// Create 创建项目
func (pai *projectApiImpl) Create(ctx *gin.Context, opt *projectrequest.CreateProjectOption) (*projectresponse.ProjectListItem, error) {
	selfMember := access.GetSelfTeamMember(ctx)
	if selfMember.Role.Lower(team.RoleAdmin) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	if opt.GroupID != 0 {
		pg := &project.ProjectGroup{ID: opt.GroupID}
		if exsit, err := pg.Get(ctx); err != nil || !exsit {
			return nil, ginrpc.NewError(
				http.StatusNotFound,
				i18n.NewErr("projectGroup.DoesNotExist"),
			)
		}
		if pg.MemberID != selfMember.ID {
			return nil, ginrpc.NewError(
				http.StatusNotFound,
				i18n.NewErr("projectGroup.DoesNotExist"),
			)
		}
	}

	pjt := &project.Project{
		TeamID:      selfMember.TeamID,
		Title:       opt.Title,
		Visibility:  opt.Visibility,
		Cover:       opt.Cover,
		Description: opt.Description,
	}

	var (
		content *spec.Spec
		err     error
	)
	if opt.Type != "" {
		switch opt.Type {
		case "apicat":
			content, err = apicatFileParse(opt.Data)
		case "swagger", "openapi":
			content, err = openapiAndSwaggerFileParse(opt.Data)
		case "postman":
			content, err = postmanFileParse(opt.Data)
		default:
			return nil, ginrpc.NewError(
				http.StatusBadRequest,
				i18n.NewErr("project.NotSupportFileType", opt.Type),
			)
		}

		if err != nil {
			slog.ErrorContext(ctx, "fileParse", "err", err)
			return nil, ginrpc.NewError(http.StatusUnprocessableEntity, i18n.NewErr("project.FileParseFailed"))
		}

		if pjt.Description == "" && content.Info.Description != "" {
			pjt.Description = content.Info.Description
		}

		contentStr, err := json.Marshal(content)
		if err != nil {
			slog.ErrorContext(ctx, "json.Marshal", "err", err)
			return nil, ginrpc.NewError(http.StatusUnprocessableEntity, i18n.NewErr("project.FileParseFailed"))
		}
		slog.InfoContext(ctx, "import", opt.Type, contentStr)
	}

	p, err := pjt.Create(ctx, selfMember, opt.GroupID)
	if err != nil {
		slog.ErrorContext(ctx, "pjt.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("project.CreationFailed"))
	}

	pm := &project.ProjectMember{
		ProjectID: p.ID,
		MemberID:  selfMember.ID,
	}
	if _, err := pm.Get(ctx); err != nil {
		slog.ErrorContext(ctx, "pjt.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("project.CreationFailed"))
	}

	// 进行数据导入工作
	if opt.Type != "" {
		project.ServersImport(ctx, p.ID, content.Servers)

		refContentVirtualIDToID := &collection.RefContentVirtualIDToId{}
		refContentVirtualIDToID.GlobalParameters = relations.ImportGlobalParameters(ctx, p.ID, content.Globals.Parameters)
		// refContentVirtualIDToID.DefinitionParameters = relations.ImportDefinitionParameters(ctx, p.ID, content.Definitions.Parameters)
		refContentVirtualIDToID.DefinitionSchemas = relations.ImportDefinitionSchemas(ctx, p.ID, content.Definitions.Schemas, selfMember, 0)
		refContentVirtualIDToID.DefinitionResponses = relations.ImportDefinitionResponses(ctx, p.ID, content.Definitions.Responses, selfMember, refContentVirtualIDToID.DefinitionSchemas, 0)
		relations.CollectionImport(ctx, selfMember, p.ID, 0, content.Collections, refContentVirtualIDToID)
	}

	return &projectresponse.ProjectListItem{
		OnlyIdInfo: protobase.OnlyIdInfo{
			ID: p.ID,
		},
		ProjectDataOption: projectbase.ProjectDataOption{
			Title: p.Title,
			ProjectVisibilityOption: protobase.ProjectVisibilityOption{
				Visibility: p.Visibility,
			},
			Cover:       p.Cover,
			Description: p.Description,
		},
		SelfMember: projectresponse.ProjectSelfMemberInfo{
			GroupID:    pm.GroupID,
			IsFollowed: pm.FollowedAt != nil,
			ProjectMemberPermission: protobase.ProjectMemberPermission{
				Permission: pm.Permission,
			},
		},
	}, nil
}

// Get 获取项目
func (pai *projectApiImpl) Get(ctx *gin.Context, opt *projectrequest.GetProjectDetailOption) (*projectresponse.ProjectDetail, error) {
	// 中间件已经对 Projct 进行了检查，可以直接取
	p := access.GetSelfProject(ctx)
	pm := access.GetSelfProjectMember(ctx)

	if jwt.GetUser(ctx) == nil || pm == nil {
		pm = &project.ProjectMember{ProjectID: p.ID, Permission: project.ProjectMemberNone}
	}

	if init, err := content_suggestion.NewVectorInitializer(ctx, p.ID); err != nil {
		slog.ErrorContext(ctx, "content_suggestion.NewVectorInitializer", "err", err)
	} else {
		init.Run()
	}

	cfg := config.GetApp()
	return &projectresponse.ProjectDetail{
		ProjectListItem: projectresponse.ProjectListItem{
			OnlyIdInfo: protobase.OnlyIdInfo{
				ID: p.ID,
			},
			ProjectDataOption: projectbase.ProjectDataOption{
				Title: p.Title,
				ProjectVisibilityOption: protobase.ProjectVisibilityOption{
					Visibility: p.Visibility,
				},
				Cover:       p.Cover,
				Description: p.Description,
			},
			SelfMember: projectresponse.ProjectSelfMemberInfo{
				GroupID:    pm.GroupID,
				IsFollowed: pm.FollowedAt != nil,
				ProjectMemberPermission: protobase.ProjectMemberPermission{
					Permission: pm.Permission,
				},
			},
		},
		MockURL: fmt.Sprintf("%s/mock/%s", strings.TrimSuffix(cfg.MockUrl, "/"), p.ID),
	}, nil
}

// List 获取项目列表
func (pai *projectApiImpl) List(ctx *gin.Context, opt *projectrequest.GetProjectListOption) (*projectresponse.GetProjectsResponse, error) {
	selfMember := access.GetSelfTeamMember(ctx)

	var (
		projects        []*project.Project
		records         []*project.ProjectMember
		projectToMember = make(map[string]map[string]interface{})
		projectIds      []string
		err             error
	)

	if opt.GroupID != 0 {
		pg := &project.ProjectGroup{ID: opt.GroupID}
		if exsit, err := pg.Get(ctx); err != nil || !exsit {
			return nil, ginrpc.NewError(
				http.StatusNotFound,
				i18n.NewErr("projectGroup.DoesNotExist"),
			)
		}
		if pg.MemberID != selfMember.ID {
			return nil, ginrpc.NewError(
				http.StatusNotFound,
				i18n.NewErr("projectGroup.DoesNotExist"),
			)
		}

		records, err = project.GetProjectMemberRecordsByGroupID(ctx, selfMember.ID, opt.GroupID)
	} else if opt.IsFollowed {
		records, err = project.GetProjectMemberRecordsWithFollowed(ctx, selfMember.ID)
	} else if len(opt.Permissions) > 0 {
		records, err = project.GetProjectMemberRecordsByPermissions(ctx, selfMember.ID, opt.Permissions...)
	} else {
		records, err = project.GetProjectMemberRecordsByMemberID(ctx, selfMember.ID)
	}
	if err != nil {
		slog.ErrorContext(ctx, "project.GetProjectMemberRecords*", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("project.FailedToGetList"),
		)
	}

	resp := make(projectresponse.GetProjectsResponse, 0)
	if len(records) == 0 {
		return &resp, nil
	}

	for _, r := range records {
		projectIds = append(projectIds, r.ProjectID)
		projectToMember[r.ProjectID] = map[string]interface{}{
			"is_followed": r.FollowedAt != nil,
			"group_id":    r.GroupID,
			"permission":  r.Permission,
		}
	}

	projects, err = project.GetProjectsByIds(ctx, projectIds)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetProjectsByIds", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("project.FailedToGetList"),
		)
	}

	for _, v := range projects {
		resp = append(resp, &projectresponse.ProjectListItem{
			OnlyIdInfo: protobase.OnlyIdInfo{
				ID: v.ID,
			},
			ProjectDataOption: projectbase.ProjectDataOption{
				Title: v.Title,
				ProjectVisibilityOption: protobase.ProjectVisibilityOption{
					Visibility: v.Visibility,
				},
				Cover:       v.Cover,
				Description: v.Description,
			},
			SelfMember: projectresponse.ProjectSelfMemberInfo{
				GroupID:    projectToMember[v.ID]["group_id"].(uint),
				IsFollowed: projectToMember[v.ID]["is_followed"].(bool),
				ProjectMemberPermission: protobase.ProjectMemberPermission{
					Permission: projectToMember[v.ID]["permission"].(project.Permission),
				},
			},
		})
	}
	return &resp, nil
}

// ChangeGroup 切换项目分组
func (pai *projectApiImpl) ChangeGroup(ctx *gin.Context, opt *projectrequest.SwitchProjectGroupOption) (*ginrpc.Empty, error) {
	selfMember := access.GetSelfTeamMember(ctx)
	// 项目分组设为0，表示取消分组
	if opt.GroupID != 0 {
		pg := &project.ProjectGroup{ID: opt.GroupID}
		if exsit, err := pg.Get(ctx); err != nil || !exsit {
			return nil, ginrpc.NewError(
				http.StatusBadRequest,
				i18n.NewErr("project.FailedToGetList"),
			)
		}
		if pg.MemberID != selfMember.ID {
			return nil, ginrpc.NewError(
				http.StatusBadRequest,
				i18n.NewErr("project.FailedToGetList"),
			)
		}
	}
	pm := access.GetSelfProjectMember(ctx)
	pm.GroupID = opt.GroupID
	if err := pm.Update(ctx); err != nil {
		slog.ErrorContext(ctx, "pm.Update", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("projectGroup.GroupingFailed"),
		)
	}

	return &ginrpc.Empty{}, nil
}

// Follow 关注项目
func (pai *projectApiImpl) Follow(ctx *gin.Context, opt *protobase.ProjectIdOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	t := time.Now()
	if err := pm.UpdateFollow(ctx, &t); err != nil {
		slog.ErrorContext(ctx, "pm.UpdateFollow", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("projectGroup.FailedToFollowProject"),
		)
	}
	return &ginrpc.Empty{}, nil
}

// UnFollow 取消关注项目
func (pai *projectApiImpl) UnFollow(ctx *gin.Context, opt *protobase.ProjectIdOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if err := pm.UpdateFollow(ctx, nil); err != nil {
		slog.ErrorContext(ctx, "pm.UpdateFollow", "err", err)
		return nil, ginrpc.NewError(
			http.StatusInternalServerError,
			i18n.NewErr("project.FailedToUnfollowProject"),
		)
	}

	return &ginrpc.Empty{}, nil
}

// Setting 项目设置
func (pai *projectApiImpl) Setting(ctx *gin.Context, opt *projectrequest.UpdateProjectOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberManage) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	p := access.GetSelfProject(ctx)
	p.Title = opt.Title
	p.Visibility = opt.Visibility
	p.Cover = opt.Cover
	p.Description = opt.Description
	if err := p.Update(ctx); err != nil {
		slog.ErrorContext(ctx, "p.Update", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	return &ginrpc.Empty{}, nil
}

// Delete 删除项目
func (pai *projectApiImpl) Delete(ctx *gin.Context, opt *protobase.ProjectIdOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberManage) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	p := access.GetSelfProject(ctx)
	if err := p.Delete(ctx); err != nil {
		slog.ErrorContext(ctx, "p.Delete", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("project.FailedToDelete"))
	}

	return &ginrpc.Empty{}, nil
}

// Transfer 移交项目
func (pai *projectApiImpl) Transfer(ctx *gin.Context, opt *projectrequest.ProjectMemberIDOption) (*ginrpc.Empty, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberManage) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("common.PermissionDenied"),
		)
	}

	targetMember := &project.ProjectMember{ID: opt.MemberID}
	exist, err := targetMember.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "targetMember.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("project.TransferFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("projectMember.NotInTheProject"))
	}
	targetMemberInfo, err := targetMember.MemberInfo(ctx, false)
	if err != nil {
		slog.ErrorContext(ctx, "targetMember.MemberInfo", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("project.TransferFailed"))
	}

	if targetMember.ProjectID != selfPM.ProjectID {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("projectMember.NotInTheProject"))
	}
	if targetMember.Permission != project.ProjectMemberWrite {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("project.TransferToErrMember"))
	}
	if targetMemberInfo.Status == team.MemberStatusDeactive {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("project.TransferToDisabledMember"))
	}

	if err := project.TransferProject(ctx, selfP, selfPM, targetMember); err != nil {
		slog.ErrorContext(ctx, "project.TransferProject", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("project.TransferFailed"))
	}

	return &ginrpc.Empty{}, nil
}

// Exit 退出项目
func (pai *projectApiImpl) Exit(ctx *gin.Context, opt *protobase.ProjectIdOption) (*ginrpc.Empty, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Equal(project.ProjectMemberManage) {
		return nil, ginrpc.NewError(
			http.StatusForbidden,
			i18n.NewErr("project.CanNotQuitOwnProject"),
		)
	}

	if err := selfPM.Delete(ctx); err != nil {
		slog.ErrorContext(ctx, "selfPM.Delete", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("project.QuitFailed"))
	}

	return &ginrpc.Empty{}, nil
}

func (pai *projectApiImpl) GetExportPath(ctx *gin.Context, opt *projectrequest.GetExportPathOption) (*projectresponse.ExportProject, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	selfTM := access.GetSelfTeamMember(ctx)
	if selfPM.Permission.Equal(project.ProjectMemberRead) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	tokenKey := fmt.Sprintf(
		"ExportProject-%d-%d",
		selfTM.ID,
		time.Now().Unix(),
	)
	c, err := cache.NewCache(config.Get().Cache.ToModuleStruct())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("project.ExportFailed"))
	}
	token, err := onetime_token.NewTokenHelper(c).GenerateToken(tokenKey, opt, time.Minute)
	if err != nil {
		slog.ErrorContext(ctx, "onetime_token.GenerateToken", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("project.ExportFailed"))
	}

	return &projectresponse.ExportProject{
		Path: fmt.Sprintf("/api/projects/%s/export/%s", selfPM.ProjectID, token),
	}, nil
}

// Export 导出项目
func Export(ctx *gin.Context) {
	// 解析和校验 URI 中的参数
	opt := &projectrequest.ExportCodeOption{}
	if err := ctx.ShouldBindUri(&opt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c, err := cache.NewCache(config.Get().Cache.ToModuleStruct())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.NewTran("project.ExportFailed").Translate(ctx),
		})
		return
	}
	tokenHelper := onetime_token.NewTokenHelper(c)

	t := projectrequest.GetExportPathOption{}
	if !tokenHelper.CheckToken(opt.Code, &t) {
		slog.ErrorContext(ctx, "onetime_token.CheckToken", "err", i18n.NewErr("project.ExportFailed"))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.NewTran("project.ExportFailed").Translate(ctx),
		})
		return
	}

	if t.ProjectID != opt.ProjectID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.NewTran("project.ExportFailed").Translate(ctx),
		})
		return
	}

	p := &project.Project{ID: t.ProjectID}
	exist, err := p.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "p.Get", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.NewTran("project.ExportFailed").Translate(ctx),
		})
		return
	}
	if !exist {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": i18n.NewTran("project.DoesNotExist").Translate(ctx),
		})
		return
	}

	apicatData := spec.NewEmptySpec()
	relations.SpecFillInfo(ctx, apicatData, p)
	relations.SpecFillServers(ctx, apicatData, p.ID)
	relations.SpecFillGlobals(ctx, apicatData, p.ID)
	relations.SpecFillDefinitions(ctx, apicatData, p.ID)
	relations.SpecFillCollections(ctx, apicatData, p.ID)
	if _, err := json.Marshal(apicatData); err != nil {
		slog.ErrorContext(ctx, "export", "marshalErr", err)
	}

	var content []byte
	switch t.Type {
	case "swagger":
		content, err = openapi.Generate(apicatData, "2.0", "json")
	case "openapi3.0.0":
		content, err = openapi.Generate(apicatData, "3.0.0", "json")
	case "openapi3.0.1":
		content, err = openapi.Generate(apicatData, "3.0.1", "json")
	case "openapi3.0.2":
		content, err = openapi.Generate(apicatData, "3.0.2", "json")
	case "openapi3.1.0":
		content, err = openapi.Generate(apicatData, "3.1.0", "json")
	case "HTML":
		content, err = export.HTML(apicatData)
	case "md":
		content, err = export.Markdown(apicatData)
	case "apicat":
		content, err = apicatData.ToJSON(spec.JSONOption{Indent: "  "})
	default:
		content, err = apicatData.ToJSON(spec.JSONOption{Indent: "  "})
	}
	if err != nil {
		slog.ErrorContext(ctx, "export", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.NewTran("project.ExportFailed").Translate(ctx),
		})
		return
	}

	slog.InfoContext(ctx, "export", t.Type, content)

	switch t.Download {
	case true:
		filename := fmt.Sprintf("%s-%s", p.Title, t.Type)
		switch t.Type {
		case "HTML":
			ctx.Header("Content-Disposition", "attachment; filename="+filename+".html")
		case "md":
			ctx.Header("Content-Disposition", "attachment; filename="+filename+".md")
		default:
			ctx.Header("Content-Disposition", "attachment; filename="+filename+".json")
		}
		ctx.Data(http.StatusOK, "application/octet-stream", content)
	default:
		switch t.Type {
		case "HTML":
			ctx.Data(http.StatusOK, "text/html; charset=utf-8", content)
		case "md":
			ctx.Data(http.StatusOK, "text/markdown; charset=utf-8", content)
		default:
			ctx.Data(http.StatusOK, "application/json", content)
		}
	}

	tokenHelper.DelToken(opt.Code)
}
