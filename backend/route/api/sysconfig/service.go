package sysconfig

import (
	"apicat-cloud/backend/config"
	"apicat-cloud/backend/i18n"
	"apicat-cloud/backend/model/sysconfig"
	protosysconfig "apicat-cloud/backend/route/proto/sysconfig"
	sysconfigbase "apicat-cloud/backend/route/proto/sysconfig/base"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type serviceApiImpl struct{}

func NewServiceApi() protosysconfig.ServiceApi {
	return &serviceApiImpl{}
}

func (s *serviceApiImpl) Get(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigbase.ServiceOption, error) {
	app := config.Get().App
	return &sysconfigbase.ServiceOption{
		AppName:        app.AppName,
		AppUrl:         app.AppUrl,
		AppServerBind:  app.AppServerBind,
		MockUrl:        app.MockUrl,
		MockServerBind: app.MockServerBind,
	}, nil
}

func (s *serviceApiImpl) Update(ctx *gin.Context, opt *sysconfigbase.ServiceOption) (*ginrpc.Empty, error) {
	if i := strings.Index(opt.AppServerBind, ":"); i < 7 {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.ServiceBindFailed"))
	}
	if i := strings.Index(opt.MockServerBind, ":"); i < 7 {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.ServiceBindFailed"))
	}

	appBind := strings.Split(opt.AppServerBind, ":")
	mockBind := strings.Split(opt.MockServerBind, ":")
	if len(appBind) != 2 || len(mockBind) != 2 {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.ServiceBindFailed"))
	}
	if appBind[1] == mockBind[1] {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.ServiceBindPortSame"))
	}

	appConfig := &config.App{
		AppName:        opt.AppName,
		AppUrl:         opt.AppUrl,
		AppServerBind:  opt.AppServerBind,
		MockUrl:        opt.MockUrl,
		MockServerBind: opt.MockServerBind,
	}

	jsonData, err := json.Marshal(appConfig)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ServiceUpdateFailed"))
	}

	app := &sysconfig.Sysconfig{
		Type:      "service",
		Driver:    "default",
		BeingUsed: true,
		Config:    string(jsonData),
	}

	if err := sysconfig.UpdateOrCreate(ctx, app); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ServiceUpdateFailed"))
	}
	config.SetApp(appConfig)
	config.SetLocalDiskUrl(opt.AppUrl)
	return &ginrpc.Empty{}, nil
}

func (s *serviceApiImpl) GetDB(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigbase.MySQLDetail, error) {
	db := config.Get().Database
	return &sysconfigbase.MySQLDetail{
		Host:     db.Host,
		Database: db.Database,
		Username: db.Username,
	}, nil
}
