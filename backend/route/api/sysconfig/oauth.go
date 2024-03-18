package sysconfig

import (
	"apicat-cloud/backend/config"
	"apicat-cloud/backend/i18n"
	"apicat-cloud/backend/model/sysconfig"
	"apicat-cloud/backend/module/oauth2"
	protosysconfig "apicat-cloud/backend/route/proto/sysconfig"
	sysconfigbase "apicat-cloud/backend/route/proto/sysconfig/base"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type oauthApiImpl struct{}

func NewOauthApi() protosysconfig.OauthApi {
	return &oauthApiImpl{}
}

func (o *oauthApiImpl) GetGithubClientID(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigbase.GitHubClientID, error) {
	oauth := config.Get().Oauth2
	if v, ok := oauth["github"]; !ok {
		return &sysconfigbase.GitHubClientID{
			ClientID: "",
		}, nil
	} else {
		return &sysconfigbase.GitHubClientID{
			ClientID: v.ClientID,
		}, nil
	}
}

func (o *oauthApiImpl) Get(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigbase.GitHubOption, error) {
	oauth := config.Get().Oauth2
	if v, ok := oauth["github"]; !ok {
		return &sysconfigbase.GitHubOption{
			GitHubClientID: sysconfigbase.GitHubClientID{
				ClientID: "",
			},
			ClientSecret: "",
		}, nil
	} else {
		return &sysconfigbase.GitHubOption{
			GitHubClientID: sysconfigbase.GitHubClientID{
				ClientID: v.ClientID,
			},
			ClientSecret: v.ClientSecret,
		}, nil
	}
}

func (o *oauthApiImpl) Update(ctx *gin.Context, opt *sysconfigbase.GitHubOption) (*ginrpc.Empty, error) {
	jsonData, err := json.Marshal(opt)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.OauthUpdateFailed"))
	}

	oauth := &sysconfig.Sysconfig{
		Type:      "oauth",
		Driver:    "github",
		BeingUsed: true,
		Config:    string(jsonData),
	}

	if err := sysconfig.UpdateOrCreate(ctx, oauth); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.OauthUpdateFailed"))
	}

	syscfg := config.Get()
	syscfg.Oauth2 = map[string]oauth2.Config{
		"github": {
			ClientID:     opt.ClientID,
			ClientSecret: opt.ClientSecret,
		},
	}

	return &ginrpc.Empty{}, nil
}
