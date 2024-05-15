package project

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	protoproject "github.com/apicat/apicat/v2/backend/route/proto/project"
	projectrequest "github.com/apicat/apicat/v2/backend/route/proto/project/request"
	projectresponse "github.com/apicat/apicat/v2/backend/route/proto/project/response"
	"github.com/apicat/apicat/v2/backend/service/reference"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type definitionResponseApiImpl struct{}

func NewDefinitionResponseApi() protoproject.DefinitionResponseApi {
	return &definitionResponseApiImpl{}
}

func (drai *definitionResponseApiImpl) Create(ctx *gin.Context, opt *projectrequest.CreateDefinitionResponseOption) (*projectresponse.DefinitionResponse, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if opt.ParentID != 0 {
		parentDR := &definition.DefinitionResponse{ID: opt.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := parentDR.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "parentDR.Get", "err", err)
			if opt.Type == definition.ResponseCategory {
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.CreationFailed"))
			}
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.CreationFailed"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
	}

	dr := &definition.DefinitionResponse{
		ProjectID:   selfPM.ProjectID,
		ParentID:    opt.ParentID,
		Name:        opt.Name,
		Description: opt.Description,
		Type:        opt.Type,
		Header:      opt.Header,
		Content:     opt.Content,
	}

	if err := dr.Create(ctx, selfTM); err != nil {
		slog.ErrorContext(ctx, "dr.Create", "err", err)
		if dr.Type == definition.ResponseCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.CreationFailed"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.CreationFailed"))
	}

	userInfo := jwt.GetUser(ctx)
	return convertModelDefinitionResponse(dr, userInfo, userInfo), nil
}

func (drai *definitionResponseApiImpl) List(ctx *gin.Context, opt *protobase.ProjectIdOption) (*projectresponse.DefinitionResponseTree, error) {
	selfP := access.GetSelfProject(ctx)
	list, err := definition.GetDefinitionResponses(ctx, selfP.ID)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetDefinitionResponses", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.FailedToGetList"))
	}

	tree := buildDefinitionResponseTree(0, list)
	return &tree, nil
}

func (drai *definitionResponseApiImpl) Get(ctx *gin.Context, opt *projectrequest.GetDefinitionResponseOption) (*projectresponse.DefinitionResponse, error) {
	selfP := access.GetSelfProject(ctx)
	dr := &definition.DefinitionResponse{ID: opt.ResponseID, ProjectID: selfP.ID}
	exist, err := dr.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "dr.Get", "err", err)
		if dr.Type == definition.ResponseCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToGet"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.FailedToGet"))
	}
	if !exist {
		if dr.Type == definition.ResponseCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.DoesNotExist"))
		}
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionResponse.DoesNotExist"))
	}

	cUserInfo, err := definition.UserInfo(ctx, dr.CreatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "dr.CreatedMember.UserInfo", "err", err)
		if dr.Type == definition.ResponseCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToGet"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.FailedToGet"))
	}
	uUserInfo, err := definition.UserInfo(ctx, dr.UpdatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "dr.UpdatedMember.UserInfo", "err", err)
		if dr.Type == definition.ResponseCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToGet"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.FailedToGet"))
	}

	return convertModelDefinitionResponse(dr, cUserInfo, uUserInfo), nil
}

func (drai *definitionResponseApiImpl) Update(ctx *gin.Context, opt *projectrequest.UpdateDefinitionResponseOption) (*ginrpc.Empty, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	dr := &definition.DefinitionResponse{ID: opt.ResponseID, ProjectID: selfPM.ProjectID}
	exist, err := dr.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "dr.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionResponse.DoesNotExist"))
	}

	oldRefSchemaIDs, err := reference.ParseRefSchemasFromResponse(dr)
	if err != nil {
		slog.ErrorContext(ctx, "reference.ParseRefSchemasFromResponse", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	dr.Name = opt.Name
	dr.Description = opt.Description
	dr.Header = opt.Header
	dr.Content = opt.Content

	// 编辑响应时更新响应引用的模型
	if dr.Type != definition.ResponseCategory {
		if err := reference.UpdateResponseRef(ctx, dr, oldRefSchemaIDs); err != nil {
			slog.ErrorContext(ctx, "reference.UpdateResponseRef", "err", err)
		}
	}

	if err := dr.Update(ctx, selfTM.ID); err != nil {
		slog.ErrorContext(ctx, "dr.Update", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	return &ginrpc.Empty{}, nil
}

func (drai *definitionResponseApiImpl) Delete(ctx *gin.Context, opt *projectrequest.DeleteDefinitionResponseOption) (*ginrpc.Empty, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	dr := &definition.DefinitionResponse{ID: opt.ResponseID, ProjectID: selfPM.ProjectID}
	exist, err := dr.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "dr.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.FailedToDelete"))
	}
	if !exist {
		if dr.Type == definition.ResponseCategory {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionResponse.DoesNotExist"))
	}

	if dr.Type == definition.ResponseCategory {
		hasChildren, err := dr.HasChildren(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "dr.HasChildren", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToDelete"))
		}
		if hasChildren {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("category.DeleteNonEmpty"))
		}
	}

	if err := dr.Delete(ctx, selfTM); err != nil {
		slog.ErrorContext(ctx, "dr.Delete", "err", err)
		if dr.Type == definition.ResponseCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToDelete"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.FailedToDelete"))
	}

	if dr.Type != definition.ResponseCategory {
		if err := reference.DerefResponse(ctx, dr, opt.Deref); err != nil {
			slog.ErrorContext(ctx, "reference.DerefResponse", "err", err)
		}
	}

	return &ginrpc.Empty{}, nil
}

func (drai *definitionResponseApiImpl) Move(ctx *gin.Context, opt *projectrequest.SortDefinitionResponseOption) (*ginrpc.Empty, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if opt.Target.ParentID != 0 {
		targetParentDR := &definition.DefinitionResponse{ID: opt.Target.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := targetParentDR.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "targetParentDR.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.FailedToMove"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
		if targetParentDR.Type != definition.ResponseCategory {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("category.IsNotCategory"))
		}
	}

	if opt.Origin.ParentID != opt.Target.ParentID && opt.Origin.ParentID != 0 {
		originParentDR := &definition.DefinitionResponse{ID: opt.Origin.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := originParentDR.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "originParentDR.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.FailedToMove"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
		if originParentDR.Type != definition.ResponseCategory {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("category.IsNotCategory"))
		}
	}

	for i, dsID := range opt.Target.IDs {
		dr := &definition.DefinitionResponse{ID: dsID, ProjectID: selfPM.ProjectID}
		exist, err := dr.Get(ctx)
		if err != nil || !exist {
			slog.ErrorContext(ctx, "Target.dr.Get", "err", err)
			continue
		}

		if err := dr.Sort(ctx, opt.Target.ParentID, uint(i+1)); err != nil {
			slog.ErrorContext(ctx, "Target.dr.Sort", "err", err)
			continue
		}
	}

	if opt.Target.ParentID != opt.Origin.ParentID {
		for i, dsID := range opt.Origin.IDs {
			dr := &definition.DefinitionResponse{ID: dsID, ProjectID: selfPM.ProjectID}
			exist, err := dr.Get(ctx)
			if err != nil || !exist {
				slog.ErrorContext(ctx, "Origin.dr.Get", "err", err)
				continue
			}

			if err := dr.Sort(ctx, opt.Origin.ParentID, uint(i+1)); err != nil {
				slog.ErrorContext(ctx, "Origin.dr.Sort", "err", err)
				continue
			}
		}
	}

	return &ginrpc.Empty{}, nil
}

func (drai *definitionResponseApiImpl) Copy(ctx *gin.Context, opt *projectrequest.GetDefinitionResponseOption) (*projectresponse.DefinitionResponse, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	dr := &definition.DefinitionResponse{ID: opt.ResponseID, ProjectID: selfPM.ProjectID}
	exist, err := dr.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "dr.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.CopyFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionResponse.DoesNotExist"))
	}

	if dr.Type == definition.ResponseCategory {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.CannotBeCopied"))
	}

	newDR := &definition.DefinitionResponse{
		ProjectID:    selfPM.ProjectID,
		ParentID:     dr.ParentID,
		Name:         fmt.Sprintf("%s (copy)", dr.Name),
		Description:  dr.Description,
		Type:         dr.Type,
		Header:       dr.Header,
		Content:      dr.Content,
		DisplayOrder: dr.DisplayOrder,
	}

	if err = newDR.Create(ctx, selfTM); err != nil {
		slog.ErrorContext(ctx, "newDR.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionResponse.CopyFailed"))
	}

	userInfo := jwt.GetUser(ctx)
	return convertModelDefinitionResponse(newDR, userInfo, userInfo), nil
}
