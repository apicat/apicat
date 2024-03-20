package project

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/definition"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/route/middleware/access"
	"github.com/apicat/apicat/backend/route/middleware/jwt"
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	protoproject "github.com/apicat/apicat/backend/route/proto/project"
	projectrequest "github.com/apicat/apicat/backend/route/proto/project/request"
	projectresponse "github.com/apicat/apicat/backend/route/proto/project/response"
	"github.com/apicat/apicat/backend/service/ai"
	definitionrelations "github.com/apicat/apicat/backend/service/definition_relations"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type definitionSchemaApiImpl struct{}

func NewDefinitionSchemaApi() protoproject.DefinitionSchemaApi {
	return &definitionSchemaApiImpl{}
}

func (dsai *definitionSchemaApiImpl) Create(ctx *gin.Context, opt *projectrequest.CreateDefinitionSchemaOption) (*projectresponse.DefinitionSchema, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if opt.ParentID != 0 {
		parentDS := &definition.DefinitionSchema{ID: opt.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := parentDS.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "parentDS.Get", "err", err)
			if opt.Type == definition.SchemaCategory {
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.CreationFailed"))
			}
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.CreationFailed"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
	}

	ds := &definition.DefinitionSchema{
		ProjectID:   selfPM.ProjectID,
		ParentID:    opt.ParentID,
		Name:        opt.Name,
		Description: opt.Description,
		Type:        opt.Type,
		Schema:      opt.Schema,
	}

	if err := ds.Create(ctx, selfTM); err != nil {
		slog.ErrorContext(ctx, "ds.Create", "err", err)
		if ds.Type == definition.SchemaCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.CreationFailed"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.CreationFailed"))
	}

	userInfo := jwt.GetUser(ctx)
	return convertModelDefinitionSchema(ds, userInfo, userInfo), nil
}

func (dsai *definitionSchemaApiImpl) List(ctx *gin.Context, opt *protobase.ProjectIdOption) (*projectresponse.DefinitionSchemaTree, error) {
	selfP := access.GetSelfProject(ctx)
	list, err := definition.GetDefinitionSchemas(ctx, selfP)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetDefinitionSchemas", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.FailedToGetList"))
	}

	tree := buildDefinitionSchemaTree(0, list)
	return &tree, nil
}

func (dsai *definitionSchemaApiImpl) Get(ctx *gin.Context, opt *projectrequest.GetDefinitionSchemaOption) (*projectresponse.DefinitionSchema, error) {
	selfP := access.GetSelfProject(ctx)
	ds := &definition.DefinitionSchema{ID: opt.SchemaID, ProjectID: selfP.ID}
	exist, err := ds.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ds.Get", "err", err)
		if ds.Type == definition.SchemaCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToGet"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.FailedToGet"))
	}
	if !exist {
		if ds.Type == definition.SchemaCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.DoesNotExist"))
		}
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchema.DoesNotExist"))
	}

	cUserInfo, err := definition.UserInfo(ctx, ds.CreatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "ds.CreatedMember.UserInfo", "err", err)
		if ds.Type == definition.SchemaCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToGet"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.FailedToGet"))
	}
	uUserInfo, err := definition.UserInfo(ctx, ds.UpdatedBy, true)
	if err != nil {
		slog.ErrorContext(ctx, "ds.UpdatedMember.UserInfo", "err", err)
		if ds.Type == definition.SchemaCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToGet"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.FailedToGet"))
	}

	return convertModelDefinitionSchema(ds, cUserInfo, uUserInfo), nil
}

func (dsai *definitionSchemaApiImpl) Update(ctx *gin.Context, opt *projectrequest.UpdateDefinitionSchemaOption) (*ginrpc.Empty, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	ds := &definition.DefinitionSchema{ID: opt.SchemaID, ProjectID: selfPM.ProjectID}
	exist, err := ds.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ds.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchema.DoesNotExist"))
	}

	if err := ds.Update(ctx, opt.Name, opt.Description, opt.Schema, selfTM.ID); err != nil {
		slog.ErrorContext(ctx, "ds.Update", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	// 编辑模型后更新模型的引用关系
	if ds.Type != definition.SchemaCategory {
		// 更新模型引用关系
		if err := definitionrelations.UpdateSchemaReference(ctx, ds); err != nil {
			slog.ErrorContext(ctx, "definitionrelations.UpdateSchemaReference", "err", err)
		}
	}

	return &ginrpc.Empty{}, nil
}

func (dsai *definitionSchemaApiImpl) Delete(ctx *gin.Context, opt *projectrequest.DeleteDefinitionSchemaOption) (*ginrpc.Empty, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	ds := &definition.DefinitionSchema{ID: opt.SchemaID, ProjectID: selfPM.ProjectID}
	exist, err := ds.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ds.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.FailedToDelete"))
	}
	if !exist {
		if ds.Type == definition.SchemaCategory {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchema.DoesNotExist"))
	}

	if ds.Type == definition.SchemaCategory {
		hasChildren, err := ds.HasChildren(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "ds.HasChildren", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToDelete"))
		}
		if hasChildren {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("category.DeleteNonEmpty"))
		}
	}

	if err := ds.Delete(ctx, selfTM); err != nil {
		slog.ErrorContext(ctx, "ds.Delete", "err", err)
		if ds.Type == definition.SchemaCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.FailedToDelete"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.FailedToDelete"))
	}

	if ds.Type != definition.SchemaCategory {
		if opt.Deref {
			if err := definitionrelations.UnpackDefinitionSchemaReferences(ctx, ds); err != nil {
				slog.ErrorContext(ctx, "definitionrelations.UnpackDefinitionSchemaReferences", "err", err)
			}
		} else {
			if err := definitionrelations.RemoveDefinitionSchemaReferences(ctx, ds); err != nil {
				slog.ErrorContext(ctx, "definitionrelations.RemoveDefinitionSchemaReferences", "err", err)
			}
		}
	}

	return &ginrpc.Empty{}, nil
}

func (dsai *definitionSchemaApiImpl) Move(ctx *gin.Context, opt *projectrequest.SortDefinitionSchemaOption) (*ginrpc.Empty, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if opt.Target.ParentID != 0 {
		targetParentDS := &definition.DefinitionSchema{ID: opt.Target.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := targetParentDS.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "targetParentDS.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.FailedToMove"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
		if targetParentDS.Type != definition.SchemaCategory {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("category.IsNotCategory"))
		}
	}

	if opt.Origin.ParentID != opt.Target.ParentID && opt.Origin.ParentID != 0 {
		originParentDS := &definition.DefinitionSchema{ID: opt.Origin.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := originParentDS.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "originParentDS.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.FailedToMove"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
		if originParentDS.Type != definition.SchemaCategory {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("category.IsNotCategory"))
		}
	}

	for i, dsID := range opt.Target.IDs {
		ds := &definition.DefinitionSchema{ID: dsID, ProjectID: selfPM.ProjectID}
		exist, err := ds.Get(ctx)
		if err != nil || !exist {
			slog.ErrorContext(ctx, "ds.Get", "err", err)
			continue
		}

		if err := ds.Sort(ctx, opt.Target.ParentID, uint(i+1)); err != nil {
			slog.ErrorContext(ctx, "ds.Sort", "err", err)
			continue
		}
	}

	if opt.Target.ParentID != opt.Origin.ParentID {
		for i, dsID := range opt.Origin.IDs {
			ds := &definition.DefinitionSchema{ID: dsID, ProjectID: selfPM.ProjectID}
			exist, err := ds.Get(ctx)
			if err != nil || !exist {
				slog.ErrorContext(ctx, "ds.Get", "err", err)
				continue
			}

			if err := ds.Sort(ctx, opt.Origin.ParentID, uint(i+1)); err != nil {
				slog.ErrorContext(ctx, "ds.Sort", "err", err)
				continue
			}
		}
	}

	return &ginrpc.Empty{}, nil
}

func (dsai *definitionSchemaApiImpl) Copy(ctx *gin.Context, opt *projectrequest.GetDefinitionSchemaOption) (*projectresponse.DefinitionSchema, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	ds := &definition.DefinitionSchema{ID: opt.SchemaID, ProjectID: selfPM.ProjectID}
	exist, err := ds.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ds.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.CopyFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("definitionSchema.DoesNotExist"))
	}

	if ds.Type == definition.SchemaCategory {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.CannotBeCopied"))
	}

	newDS := &definition.DefinitionSchema{
		ProjectID:    selfPM.ProjectID,
		ParentID:     ds.ParentID,
		Name:         fmt.Sprintf("%s (copy)", ds.Name),
		Description:  ds.Description,
		Type:         ds.Type,
		Schema:       ds.Schema,
		DisplayOrder: ds.DisplayOrder,
	}

	if err = newDS.Create(ctx, selfTM); err != nil {
		slog.ErrorContext(ctx, "newDS.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.CopyFailed"))
	}

	userInfo := jwt.GetUser(ctx)
	return convertModelDefinitionSchema(newDS, userInfo, userInfo), nil
}

func (dsai *definitionSchemaApiImpl) AIGenerate(ctx *gin.Context, opt *projectrequest.AIGenerateSchemaOption) (*projectresponse.DefinitionSchema, error) {
	selfTM := access.GetSelfTeamMember(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if opt.ParentID != 0 {
		parentDS := &definition.DefinitionSchema{ID: opt.ParentID, ProjectID: selfPM.ProjectID}
		exist, err := parentDS.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "parentDS.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.GenerationFailed"))
		}
		if !exist {
			return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("category.DoesNotExist"))
		}
	}

	ds, err := ai.SchemaGenerate(ctx, opt.Prompt)
	if err != nil {
		slog.ErrorContext(ctx, "ai.SchemaGenerate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.GenerationFailed"))
	}

	ds.ProjectID = selfPM.ProjectID
	ds.ParentID = opt.ParentID
	if err := ds.Create(ctx, selfTM); err != nil {
		slog.ErrorContext(ctx, "ds.Create", "err", err)
		if ds.Type == definition.SchemaCategory {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("category.CreationFailed"))
		}
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("definitionSchema.GenerationFailed"))
	}

	userInfo := jwt.GetUser(ctx)
	return convertModelDefinitionSchema(ds, userInfo, userInfo), nil
}
