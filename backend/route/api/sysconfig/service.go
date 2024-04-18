package sysconfig

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/sysconfig"
	protosysconfig "github.com/apicat/apicat/v2/backend/route/proto/sysconfig"
	sysconfigbase "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/base"

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
		AppName: app.AppName,
		AppUrl:  app.AppUrl,
		MockUrl: app.MockUrl,
	}, nil
}

func (s *serviceApiImpl) Update(ctx *gin.Context, opt *sysconfigbase.ServiceOption) (*ginrpc.Empty, error) {
	appConfig := &config.App{
		AppName: opt.AppName,
		AppUrl:  opt.AppUrl,
		MockUrl: opt.MockUrl,
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
