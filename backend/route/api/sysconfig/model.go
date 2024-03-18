package sysconfig

import (
	"apicat-cloud/backend/config"
	"apicat-cloud/backend/i18n"
	"apicat-cloud/backend/model/sysconfig"
	"apicat-cloud/backend/module/llm"
	"apicat-cloud/backend/module/llm/openai"
	protosysconfig "apicat-cloud/backend/route/proto/sysconfig"
	sysconfigbase "apicat-cloud/backend/route/proto/sysconfig/base"
	sysconfigrequest "apicat-cloud/backend/route/proto/sysconfig/request"
	"encoding/json"
	"log/slog"
	"net/http"

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

	mc := modelConfig.ToMapInterface()
	if ai, err := openai.NewOpenAI(mc["OpenAI"].(map[string]interface{})); err != nil {
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

	mc := modelConfig.ToMapInterface()
	if ai, err := openai.NewOpenAI(mc["AzureOpenAI"].(map[string]interface{})); err != nil {
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
