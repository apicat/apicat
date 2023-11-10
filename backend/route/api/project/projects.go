package project

import (
	"encoding/base64"
	"encoding/json"
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/definition"
	"github.com/apicat/apicat/backend/model/global"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/server"
	"github.com/apicat/apicat/backend/model/share"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/module/spec"
	"github.com/apicat/apicat/backend/module/spec/plugin/export"
	"github.com/apicat/apicat/backend/module/spec/plugin/openapi"
	"github.com/apicat/apicat/backend/module/spec/plugin/postman"
	"github.com/apicat/apicat/backend/module/translator"
	"net/http"
	"strings"
	"time"

	"github.com/apicat/apicat/backend/app/util"
	"github.com/apicat/apicat/backend/enum"
	"golang.org/x/exp/slog"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
)

type CreateProject struct {
	Title      string `json:"title" binding:"required,lte=255"`
	Data       string `json:"data"`
	Cover      string `json:"cover" binding:"lte=255"`
	Visibility string `json:"visibility" binding:"required,oneof=private public"`
	DataType   string `json:"data_type" binding:"omitempty,oneof=apicat swagger openapi postman"`
	GroupID    uint   `json:"group_id" binding:"omitempty"`
}

type UpdateProject struct {
	Title       string `json:"title" binding:"required,lte=255"`
	Description string `json:"description" binding:"lte=255"`
	Cover       string `json:"cover" binding:"lte=255"`
	Visibility  string `json:"visibility" binding:"required,oneof=private public"`
}

type ProjectID struct {
	ID string `uri:"project-id" binding:"required"`
}

type ExportProject struct {
	Type     string `form:"type" binding:"required,oneof=apicat swagger openapi3.0.0 openapi3.0.1 openapi3.0.2 openapi3.1.0 HTML md"`
	Download string `form:"download" binding:"omitempty,oneof=true false"`
}

type TranslateProject struct {
	MemberID uint `json:"member_id" binding:"required,lte=255"`
}

type ProjectsListData struct {
	Auth       []string `form:"auth" binding:"omitempty,dive,oneof=manage write read"`
	GroupID    uint     `form:"group_id"`
	IsFollowed bool     `form:"is_followed"`
}

type ProjectChangeGroupData struct {
	TargetGroupID uint `json:"target_group_id" binding:"lte=255"`
}

func ProjectsList(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var data ProjectsListData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var (
		projectMembers []project.ProjectMembers
		err            error
	)
	if data.GroupID > 0 {
		projectMembers, err = project.GetProjectGroupedByUser(currentUser.(*user.Users).ID, data.GroupID)
	} else if data.IsFollowed {
		projectMembers, err = project.GetProjectFollowedByUser(currentUser.(*user.Users).ID)
	} else if len(data.Auth) > 0 {
		projectMembers, err = project.GetUserInvolvedProject(currentUser.(*user.Users).ID, data.Auth...)
	} else {
		projectMembers, err = project.GetUserInvolvedProject(currentUser.(*user.Users).ID)
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.QueryFailed"}),
		})
		return
	}
	if len(projectMembers) == 0 {
		ctx.JSON(http.StatusOK, []gin.H{})
		return
	}

	projectIDs := []uint{}
	for _, v := range projectMembers {
		projectIDs = append(projectIDs, v.ProjectID)
	}

	p, _ := project.NewProjects()
	projects, err := p.List(projectIDs...)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	followProjects := projectMembers
	if !data.IsFollowed {
		followProjects, err = project.GetProjectFollowedByUser(currentUser.(*user.Users).ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.QueryFailed"}),
			})
			return
		}
	}

	projectsList := []gin.H{}
	for _, v := range projects {
		isFollow := false
		for _, followProject := range followProjects {
			if v.ID == followProject.ProjectID {
				isFollow = true
				break
			}
		}
		groupID := 0
		for _, pm := range projectMembers {
			if v.ID == pm.ProjectID {
				groupID = int(pm.GroupID)
				break
			}
		}

		projectsList = append(projectsList, gin.H{
			"id":          v.PublicId,
			"title":       v.Title,
			"description": v.Description,
			"cover":       v.Cover,
			"is_followed": isFollow,
			"group_id":    groupID,
			"created_at":  v.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  v.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, projectsList)
}

func ProjectsGet(ctx *gin.Context) {
	currentProjectMember, currentProjectMemberExists := ctx.Get("CurrentProjectMember")
	currentProject, _ := ctx.Get("CurrentProject")
	p := currentProject.(*project.Projects)

	var (
		data       ProjectID
		authority  string
		visibility string
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if currentProjectMemberExists {
		authority = currentProjectMember.(*project.ProjectMembers).Authority
	} else {
		authority = "none"
	}

	if p.Visibility == 0 {
		visibility = "private"
	} else {
		visibility = "public"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          p.PublicId,
		"title":       p.Title,
		"description": p.Description,
		"cover":       p.Cover,
		"authority":   authority,
		"visibility":  visibility,
		"secret_key":  p.SharePassword,
		"created_at":  p.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":  p.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func ProjectsCreate(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	u, _ := CurrentUser.(*user.Users)
	if u.Role == "user" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.MemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		data    CreateProject
		content *spec.Spec
		err     error
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.GroupID > 0 {
		pg, err := project.NewProjectGroups(data.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.NotFound"}),
			})
			return
		}
		if pg.UserID != u.ID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.NotFound"}),
			})
			return
		}
	}

	p, _ := project.NewProjects()

	if data.DataType != "" {
		switch data.DataType {
		case "apicat":
			content, err = apicatFileParse(data.Data)
		case "swagger":
			content, err = openapiAndSwaggerFileParse(data.Data)
		case "openapi":
			content, err = openapiAndSwaggerFileParse(data.Data)
		case "postman":
			content, err = postmanFileParse(data.Data)
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.ImportFail"}),
			})
			return
		}

		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.ImportFail"}),
			})
			return
		}

		p.Description = content.Info.Description
	}

	if data.Visibility == "private" {
		p.Visibility = 0
	} else {
		p.Visibility = 1
	}

	p.Title = data.Title
	p.PublicId = shortuuid.New()
	p.Cover = data.Cover
	if err := p.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.CreateFail"}),
		})
		return
	}

	pm, _ := project.NewProjectMembers()
	pm.ProjectID = p.ID
	pm.UserID = u.ID
	pm.Authority = project.ProjectMembersManage
	pm.GroupID = data.GroupID
	if err := pm.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.CreateFail"}),
		})
		return
	}

	// 进行数据导入工作
	if data.Data != "" {
		server.ServersImport(p.ID, content.Servers)

		refContentVirtualIDToId := &model.RefContentVirtualIDToId{
			DefinitionSchemas:    definition.DefinitionSchemasImport(p.ID, content.Definitions.Schemas),
			DefinitionResponses:  definition.DefinitionResponsesImport(p.ID, content.Definitions.Responses),
			DefinitionParameters: definition.DefinitionParametersImport(p.ID, content.Definitions.Parameters),
			GolbalParameters:     global.GlobalParametersImport(p.ID, &content.Globals.Parameters),
		}

		collection.CollectionsImport(p.ID, 0, content.Collections, refContentVirtualIDToId)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          p.PublicId,
		"title":       p.Title,
		"description": p.Description,
		"cover":       p.Cover,
		"created_at":  p.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":  p.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// openapi & swagger 导出文件解析
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
		return nil, err
	}

	return openapi.Decode(rawContent)
}

// apicat 导出文件解析
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
		return nil, err
	}

	return spec.ParseJSON(rawContent)
}

// postman 导出文件解析
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
		return nil, err
	}

	return postman.Import(rawContent)
}

func ProjectsUpdate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		uriData ProjectID
		data    UpdateProject
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	p, err := project.NewProjects(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	p.Title = data.Title
	p.Description = data.Description
	p.Cover = data.Cover
	if data.Visibility == "private" {
		p.Visibility = 0
	} else {
		p.Visibility = 1

		// 将项目分享密钥及项目下集合的分享密钥置为空
		p.SharePassword = ""
		c, _ := collection.NewCollections()
		c.SharePassword = ""
		if err := collection.BatchUpdateByProjectID(p.ID, map[string]any{"share_password": ""}); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.UpdateFail"}),
			})
		}

		stt := share.NewShareTmpTokens()
		stt.ProjectID = p.ID
		if err := stt.DeleteByProjectID(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.UpdateFail"}),
			})
		}
	}

	if err := p.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.UpdateFail"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func ProjectsDelete(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var data ProjectID

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	p, err := project.NewProjects(data.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}
	if err := p.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.DeleteFail"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func ProjectDataGet(ctx *gin.Context) {
	var (
		uriData ProjectID
		data    ExportProject
		content []byte
		err     error
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	p, err := project.NewProjects(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	apicatData := &spec.Spec{}
	apicatData.ApiCat = "apicat"
	apicatData.Info = &spec.Info{
		ID:          p.PublicId,
		Title:       p.Title,
		Description: p.Description,
		Version:     "1.0.0",
	}

	apicatData.Servers = server.ServersExport(p.ID)
	apicatData.Globals.Parameters = global.GlobalParametersExport(p.ID)
	apicatData.Definitions.Schemas = definition.DefinitionSchemasExport(p.ID)
	apicatData.Definitions.Parameters = definition.DefinitionParametersExport(p.ID)
	apicatData.Definitions.Responses = definition.DefinitionResponsesExport(p.ID)
	apicatData.Collections = collection.CollectionsExport(p.ID)

	if apicatDataContent, err := json.Marshal(apicatData); err == nil {
		slog.InfoCtx(ctx, "Export", slog.String("apicat", string(apicatDataContent)))
	}

	switch data.Type {
	case "swagger":
		content, err = openapi.Encode(apicatData, "2.0")
	case "openapi3.0.0":
		content, err = openapi.Encode(apicatData, "3.0.0")
	case "openapi3.0.1":
		content, err = openapi.Encode(apicatData, "3.0.1")
	case "openapi3.0.2":
		content, err = openapi.Encode(apicatData, "3.0.2")
	case "openapi3.1.0":
		content, err = openapi.Encode(apicatData, "3.1.0")
	case "HTML":
		content, err = export.HTML(apicatData)
	case "md":
		content, err = export.Markdown(apicatData)
	case "apicat":
		content, err = apicatData.ToJSON(spec.JSONOption{Indent: "  "})
	default:
		content, err = apicatData.ToJSON(spec.JSONOption{Indent: "  "})
	}

	slog.InfoCtx(ctx, "Export", slog.String(data.Type, string(content)))

	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.ExportFail"}),
		})
		return
	}

	util.ExportResponse(data.Type, data.Download, p.Title+"-"+data.Type, content, ctx)
}

// ProjectExit handles the exit of a project member.
func ProjectExit(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if currentProjectMember.(*project.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	if err := currentProjectMember.(*project.ProjectMembers).Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.ExitFail"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func ProjectTransfer(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := TranslateProject{}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pm, err := project.NewProjectMembers(data.MemberID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.NotFound"}),
		})
		return
	}

	if pm.Authority != project.ProjectMembersWrite {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.TargetProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.TransferFail"}),
		})
		return
	}

	if pm.ProjectID != currentProjectMember.(*project.ProjectMembers).ProjectID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.TargetProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.TransferFail"}),
		})
		return
	}

	currentProjectMember.(*project.ProjectMembers).Authority = project.ProjectMembersWrite
	if err := currentProjectMember.(*project.ProjectMembers).Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.TransferFail"}),
		})
		return
	}

	pm.Authority = project.ProjectMembersManage
	if err := pm.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.TransferFail"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func ProjectFollow(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")

	nowTime := time.Now()
	currentProjectMember.(*project.ProjectMembers).FollowedAt = &nowTime
	if err := currentProjectMember.(*project.ProjectMembers).Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectFollows.FollowFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func ProjectUnFollow(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")

	currentProjectMember.(*project.ProjectMembers).FollowedAt = nil
	if err := currentProjectMember.(*project.ProjectMembers).Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectFollows.UnfollowFailed"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func ProjectChangeGroup(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")

	var data ProjectChangeGroupData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	currentProjectMember.(*project.ProjectMembers).GroupID = data.TargetGroupID
	if err := currentProjectMember.(*project.ProjectMembers).Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.ChangeFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
