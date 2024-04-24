package sysconfig

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/sysconfig"
	"github.com/apicat/apicat/v2/backend/module/llm"
	protosysconfig "github.com/apicat/apicat/v2/backend/route/proto/sysconfig"
	sysconfigbase "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/base"
	sysconfigrequest "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/request"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type modelApiImpl struct{}

func NewModelApi() protosysconfig.ModelApi {
	return &modelApiImpl{}
}

func (s *modelApiImpl) Get(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigbase.ConfigList, error) {
	list, err := sysconfig.GetList(ctx, "model")
	if err != nil {
		slog.ErrorContext(ctx, "sysconfig.GetList", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.FailedToGetModelList"))
	}
	slist := make(sysconfigbase.ConfigList, 0, len(list))
	for _, v := range list {
		slist = append(slist, &sysconfigbase.ConfigDetail{
			Driver: v.Driver,
			Use:    v.BeingUsed,
			Config: cfgFormat(v),
		})
	}
	return &slist, nil
}

func (s *modelApiImpl) UpdateOpenAI(ctx *gin.Context, opt *sysconfigrequest.OpenAIOption) (*ginrpc.Empty, error) {
	modelConfig := &config.LLM{
		Driver: llm.OPENAI,
		OpenAI: &config.OpenAI{
			ApiKey:         opt.ApiKey,
			OrganizationID: opt.OrganizationID,
			ApiBase:        opt.ApiBase,
			LLMName:        opt.LLMName,
		},
	}

	if ai, err := llm.NewLLM(modelConfig.ToCfg()); err != nil {
		slog.ErrorContext(ctx, "openai.NewOpenAI", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, err)
	} else {
		if err := ai.Check(); err != nil {
			slog.ErrorContext(ctx, "ai.Check", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, err)
		}
	}

	jsonData, err := json.Marshal(opt)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
	}

	storage := &sysconfig.Sysconfig{
		Type:      "model",
		Driver:    llm.OPENAI,
		BeingUsed: true,
		Config:    string(jsonData),
	}

	if err := sysconfig.UpdateOrCreate(ctx, storage); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
	}
	config.SetLLM(modelConfig)
	return nil, nil
}

func (s *modelApiImpl) UpdateAzureOpenAI(ctx *gin.Context, opt *sysconfigrequest.AzureOpenAIOption) (*ginrpc.Empty, error) {
	modelConfig := &config.LLM{
		Driver: llm.AZUREOPENAI,
		AzureOpenAI: &config.AzureOpenAI{
			ApiKey:   opt.ApiKey,
			Endpoint: opt.Endpoint,
			LLMName:  opt.LLMName,
		},
	}

	if ai, err := llm.NewLLM(modelConfig.ToCfg()); err != nil {
		slog.ErrorContext(ctx, "openai.NewOpenAI", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, err)
	} else {
		if err := ai.Check(); err != nil {
			slog.ErrorContext(ctx, "ai.Check", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, err)
		}
	}

	jsonData, err := json.Marshal(opt)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
	}

	storage := &sysconfig.Sysconfig{
		Type:      "model",
		Driver:    llm.AZUREOPENAI,
		BeingUsed: true,
		Config:    string(jsonData),
	}

	if err := sysconfig.UpdateOrCreate(ctx, storage); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
	}
	config.SetLLM(modelConfig)
	return nil, nil
}
