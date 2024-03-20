package project

import (
	"log/slog"
	"net/http"

	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/global"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/route/middleware/access"
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	protoproject "github.com/apicat/apicat/backend/route/proto/project"
	projectrequest "github.com/apicat/apicat/backend/route/proto/project/request"
	projectresponse "github.com/apicat/apicat/backend/route/proto/project/response"
	globalrelations "github.com/apicat/apicat/backend/service/global_relations"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type globalParameterApiImpl struct{}

func NewGlobalParameterAPI() protoproject.GlobalParameterApi {
	return &globalParameterApiImpl{}
}

func (gpai *globalParameterApiImpl) Create(ctx *gin.Context, data *projectrequest.CreateGlobalParameterOption) (*projectresponse.GlobalParameter, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	gp := global.GlobalParameter{
		ProjectID: pm.ProjectID,
		In:        data.In,
		Name:      data.Name,
	}
	exist, err := gp.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "gp.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("globalParameter.CreationFailed"))
	}
	if exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("globalParameter.HasBeenUsed"))
	}

	gp.Required = data.Required
	gp.Schema = data.Schema
	if err := gp.Create(ctx); err != nil {
		slog.ErrorContext(ctx, "gp.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("globalParameter.CreationFailed"))
	}

	return convertModelGlobalparameter(&gp), nil
}

func (gpai *globalParameterApiImpl) List(ctx *gin.Context, opt *protobase.ProjectIdOption) (*projectresponse.GlobalParameterList, error) {
	p := access.GetSelfProject(ctx)

	list, err := global.GetGlobalParameters(ctx, p.ID)
	if err != nil {
		slog.ErrorContext(ctx, "project.GetGlobalParameters", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("globalParameter.FailedToGetList"))
	}

	ret := &projectresponse.GlobalParameterList{
		Header: make([]*projectresponse.GlobalParameter, 0),
		Cookie: make([]*projectresponse.GlobalParameter, 0),
		Query:  make([]*projectresponse.GlobalParameter, 0),
		Path:   make([]*projectresponse.GlobalParameter, 0),
	}

	for _, gp := range list {
		switch gp.In {
		case global.ParameterInHeader:
			ret.Header = append(ret.Header, convertModelGlobalparameter(gp))
		case global.ParameterInCookie:
			ret.Cookie = append(ret.Cookie, convertModelGlobalparameter(gp))
		case global.ParameterInQuery:
			ret.Query = append(ret.Query, convertModelGlobalparameter(gp))
		case global.ParameterInPath:
			ret.Path = append(ret.Path, convertModelGlobalparameter(gp))
		}
	}

	return ret, nil
}

func (gpai *globalParameterApiImpl) Update(ctx *gin.Context, opt *projectrequest.UpdateGlobalParameterOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	gp := global.GlobalParameter{ID: opt.ParameterID}
	exist, err := gp.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "gp.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("globalParameter.DoesNotExist"),
		)
	}

	gp.In = opt.In
	gp.Name = opt.Name
	gp.Required = opt.Required
	gp.Schema = opt.Schema
	exist, err = gp.CheckRepeat(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "gp.CheckRepeat", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}
	if exist {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("globalParameter.HasBeenUsed"),
		)
	}

	if err := gp.Update(ctx); err != nil {
		slog.ErrorContext(ctx, "gp.Update", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	return &ginrpc.Empty{}, nil
}

func (gpai *globalParameterApiImpl) Delete(ctx *gin.Context, opt *projectrequest.DeleteGlobalParameterOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	gp := &global.GlobalParameter{ID: opt.ParameterID}
	exist, err := gp.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "gp.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(
			http.StatusBadRequest,
			i18n.NewErr("globalParameter.DoesNotExist"),
		)
	}

	if err := gp.Delete(ctx); err != nil {
		slog.ErrorContext(ctx, "gp.Delete", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("globalParameter.FailedToDelete"))
	}

	if opt.Deref {
		if err := globalrelations.UnpackExceptGlobalParameter(ctx, gp); err != nil {
			slog.ErrorContext(ctx, "globalrelations.UnpackGlobalParametersInCollection", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("globalParameter.FailedToDelete"))
		}
	} else {
		if err := globalrelations.RemoveExceptGlobalParameter(ctx, gp); err != nil {
			slog.ErrorContext(ctx, "globalrelations.RemoveExceptGlobalParameter", "err", err)
		}
	}

	return &ginrpc.Empty{}, nil
}

func (gpai *globalParameterApiImpl) Sort(ctx *gin.Context, opt *projectrequest.SortGlobalParameterOption) (*ginrpc.Empty, error) {
	pm := access.GetSelfProjectMember(ctx)
	if pm.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if err := global.SortGlobalParameters(ctx, pm.ProjectID, opt.In, opt.ParameterIDs); err != nil {
		slog.ErrorContext(ctx, "global.SortGlobalParameters", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("globalParameter.FailedToSort"))
	}

	return &ginrpc.Empty{}, nil
}
