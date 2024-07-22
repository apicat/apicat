package sysconfig

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/sysconfig"
	"github.com/apicat/apicat/v2/backend/module/model"
	protosysconfig "github.com/apicat/apicat/v2/backend/route/proto/sysconfig"
	sysconfigrequest "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/request"
	sysconfigresponse "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/response"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type modelApiImpl struct{}

var modelDrivers = []string{model.OPENAI, model.AZURE_OPENAI}
var openaiLLMSupport = []string{"gpt-4", "gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-3.5-turbo"}
var openaiEmbeddingSupport = []string{"text-embedding-3-small", "text-embedding-3-large", "text-embedding-ada-002"}

func NewModelApi() protosysconfig.ModelApi {
	return &modelApiImpl{}
}

func (s *modelApiImpl) Get(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigresponse.ModelConfigList, error) {
	configMap := make(map[string]*sysconfigresponse.ModelConfigDetail)

	modelCfg := config.GetModel()
	drivers := make([]string, 0)
	if modelCfg.LLMDriver != "" {
		drivers = append(drivers, modelCfg.LLMDriver)
	}
	if modelCfg.EmbeddingDriver != "" && modelCfg.EmbeddingDriver != modelCfg.LLMDriver {
		drivers = append(drivers, modelCfg.EmbeddingDriver)
	}

	for _, driver := range drivers {
		switch driver {
		case model.OPENAI:
			if js, err := json.Marshal(modelCfg.OpenAI); err == nil {
				configMap[model.OPENAI] = &sysconfigresponse.ModelConfigDetail{
					Driver: model.OPENAI,
					Config: cfgFormat(&sysconfig.Sysconfig{Config: string(js)}),
					Models: sysconfigresponse.ModelNameList{
						"llm":       openaiLLMSupport,
						"embedding": openaiEmbeddingSupport,
					},
				}
			}
		case model.AZURE_OPENAI:
			if js, err := json.Marshal(modelCfg.AzureOpenAI); err == nil {
				configMap[model.AZURE_OPENAI] = &sysconfigresponse.ModelConfigDetail{
					Driver: model.AZURE_OPENAI,
					Config: cfgFormat(&sysconfig.Sysconfig{Config: string(js)}),
				}
			}
		}
	}

	if list, err := sysconfig.GetList(ctx, "model"); err == nil {
		for _, v := range list {
			if _, ok := configMap[v.Driver]; !ok {
				switch v.Driver {
				case model.OPENAI:
					configMap[model.OPENAI] = &sysconfigresponse.ModelConfigDetail{
						Driver: v.Driver,
						Config: cfgFormat(v),
						Models: sysconfigresponse.ModelNameList{
							"llm":       openaiLLMSupport,
							"embedding": openaiEmbeddingSupport,
						},
					}
				case model.AZURE_OPENAI:
					configMap[model.AZURE_OPENAI] = &sysconfigresponse.ModelConfigDetail{
						Driver: v.Driver,
						Config: cfgFormat(v),
					}
				}
			}
		}
	}

	slist := make(sysconfigresponse.ModelConfigList, 0)
	for _, v := range modelDrivers {
		if _, ok := configMap[v]; ok {
			slist = append(slist, configMap[v])
		}
	}
	return &slist, nil
}

func (s *modelApiImpl) UpdateOpenAI(ctx *gin.Context, opt *sysconfigrequest.OpenAIOption) (*ginrpc.Empty, error) {
	modelConfig := &config.Model{
		OpenAI: &config.OpenAI{
			ApiKey:         opt.ApiKey,
			OrganizationID: opt.OrganizationID,
			ApiBase:        opt.ApiBase,
			LLM:            opt.LLM,
			Embedding:      opt.Embedding,
		},
	}

	if opt.LLM != "" {
		modelConfig.LLMDriver = model.OPENAI
		modelConfig.OpenAI.LLM = opt.LLM

		if openAI, err := model.NewModel(modelConfig.ToModuleStruct("llm")); err != nil {
			slog.ErrorContext(ctx, "model.NewModel.OpenAI", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, err)
		} else {
			if err := openAI.Check("llm"); err != nil {
				slog.ErrorContext(ctx, "openAI.Check", "err", err)
				return nil, ginrpc.NewError(http.StatusBadRequest, err)
			}
		}
	}
	if opt.Embedding != "" {
		modelConfig.EmbeddingDriver = model.OPENAI
		modelConfig.OpenAI.Embedding = opt.Embedding

		if openAI, err := model.NewModel(modelConfig.ToModuleStruct("embedding")); err != nil {
			slog.ErrorContext(ctx, "model.NewModel.OpenAI", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, err)
		} else {
			if err := openAI.Check("embedding"); err != nil {
				slog.ErrorContext(ctx, "openAI.Check", "err", err)
				return nil, ginrpc.NewError(http.StatusBadRequest, err)
			}
		}
	}

	m := &sysconfig.Sysconfig{
		Type:   "model",
		Driver: model.OPENAI,
	}
	exist, err := m.Get(ctx)
	if err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
	}

	jsonData, err := json.Marshal(opt)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
	}

	m.Config = string(jsonData)
	if exist {
		if err := m.Update(ctx); err != nil {
			slog.ErrorContext(ctx, "m.Update", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
		}
	} else {
		if err := m.Create(ctx); err != nil {
			slog.ErrorContext(ctx, "m.Create", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
		}
	}
	return nil, nil
}

func (s *modelApiImpl) UpdateAzureOpenAI(ctx *gin.Context, opt *sysconfigrequest.AzureOpenAIOption) (*ginrpc.Empty, error) {
	modelConfig := &config.Model{
		AzureOpenAI: &config.AzureOpenAI{
			ApiKey:   opt.ApiKey,
			Endpoint: opt.Endpoint,
		},
	}

	if opt.LLM != "" {
		modelConfig.LLMDriver = model.AZURE_OPENAI
		modelConfig.AzureOpenAI.LLM = opt.LLM

		if azureOpenAI, err := model.NewModel(modelConfig.ToModuleStruct("llm")); err != nil {
			slog.ErrorContext(ctx, "model.NewModel.AzureOpenAI", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, err)
		} else {
			if err := azureOpenAI.Check("llm"); err != nil {
				slog.ErrorContext(ctx, "azureOpenAI.Check", "err", err)
				return nil, ginrpc.NewError(http.StatusBadRequest, err)
			}
		}
	}
	if opt.Embedding != "" {
		modelConfig.EmbeddingDriver = model.AZURE_OPENAI
		modelConfig.AzureOpenAI.Embedding = opt.Embedding

		if azureOpenAI, err := model.NewModel(modelConfig.ToModuleStruct("embedding")); err != nil {
			slog.ErrorContext(ctx, "model.NewModel.AzureOpenAI", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, err)
		} else {
			if err := azureOpenAI.Check("embedding"); err != nil {
				slog.ErrorContext(ctx, "azureOpenAI.Check", "err", err)
				return nil, ginrpc.NewError(http.StatusBadRequest, err)
			}
		}
	}

	m := &sysconfig.Sysconfig{
		Type:   "model",
		Driver: model.AZURE_OPENAI,
	}
	exist, err := m.Get(ctx)
	if err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
	}

	jsonData, err := json.Marshal(opt)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
	}

	m.Config = string(jsonData)
	if exist {
		if err := m.Update(ctx); err != nil {
			slog.ErrorContext(ctx, "m.Update", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
		}
	} else {
		if err := m.Create(ctx); err != nil {
			slog.ErrorContext(ctx, "m.Create", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
		}
	}
	return nil, nil
}

func (s *modelApiImpl) GetDefault(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigresponse.DefaultModelMap, error) {
	llmDefaultMap := make(map[string]*sysconfigresponse.DefaultModelDetail)
	embeddingDefaultMap := make(map[string]*sysconfigresponse.DefaultModelDetail)

	modelCfg := config.GetModel()
	var setDefault = func(typ, driver string) {
		switch driver {
		case model.OPENAI:
			if typ == "llm" {
				llmDefaultMap[model.OPENAI] = &sysconfigresponse.DefaultModelDetail{
					Driver:   model.OPENAI,
					Model:    modelCfg.OpenAI.LLM,
					Selected: true,
				}
			} else if typ == "embedding" {
				embeddingDefaultMap[model.OPENAI] = &sysconfigresponse.DefaultModelDetail{
					Driver:   model.OPENAI,
					Model:    modelCfg.OpenAI.Embedding,
					Selected: true,
				}
			}
		case model.AZURE_OPENAI:
			if typ == "llm" {
				llmDefaultMap[model.AZURE_OPENAI] = &sysconfigresponse.DefaultModelDetail{
					Driver:   model.AZURE_OPENAI,
					Model:    modelCfg.AzureOpenAI.LLM,
					Selected: true,
				}
			} else if typ == "embedding" {
				embeddingDefaultMap[model.AZURE_OPENAI] = &sysconfigresponse.DefaultModelDetail{
					Driver:   model.AZURE_OPENAI,
					Model:    modelCfg.AzureOpenAI.Embedding,
					Selected: true,
				}
			}
		}
	}

	if modelCfg.LLMDriver != "" {
		setDefault("llm", modelCfg.LLMDriver)
	}
	if modelCfg.EmbeddingDriver != "" {
		setDefault("embedding", modelCfg.EmbeddingDriver)
	}

	if list, err := sysconfig.GetList(ctx, "model"); err == nil {
		type modelName struct {
			LLM       string `json:"llm"`
			Embedding string `json:"embedding"`
		}
		for _, v := range list {
			var cfg modelName
			if err := json.Unmarshal([]byte(v.Config), &cfg); err == nil {
				if _, ok := llmDefaultMap[v.Driver]; !ok && cfg.LLM != "" {
					llmDefaultMap[v.Driver] = &sysconfigresponse.DefaultModelDetail{
						Driver:   v.Driver,
						Model:    cfg.LLM,
						Selected: v.BeingUsed && v.Extra == "llm",
					}
				}
				if _, ok := embeddingDefaultMap[v.Driver]; !ok && cfg.Embedding != "" {
					embeddingDefaultMap[v.Driver] = &sysconfigresponse.DefaultModelDetail{
						Driver:   v.Driver,
						Model:    cfg.Embedding,
						Selected: v.BeingUsed && v.Extra == "embedding",
					}
				}
			}
		}
	}

	defaultMap := sysconfigresponse.DefaultModelMap{
		"llm":       make(sysconfigresponse.DefaultModelList, 0),
		"embedding": make(sysconfigresponse.DefaultModelList, 0),
	}
	for _, v := range modelDrivers {
		if _, ok := llmDefaultMap[v]; ok {
			defaultMap["llm"] = append(defaultMap["llm"], llmDefaultMap[v])
		}
		if _, ok := embeddingDefaultMap[v]; ok {
			defaultMap["embedding"] = append(defaultMap["embedding"], embeddingDefaultMap[v])
		}
	}

	return &defaultMap, nil
}

func (s *modelApiImpl) UpdateDefault(ctx *gin.Context, opt *sysconfigrequest.DefaultModelMapOption) (*ginrpc.Empty, error) {
	type modelName struct {
		LLM       string `json:"llm"`
		Embedding string `json:"embedding"`
	}
	var cfg modelName

	llmModel := &sysconfig.Sysconfig{
		Type:   "model",
		Driver: opt.LLM.Driver,
	}
	exist, err := llmModel.Get(ctx)
	if err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.ReasoningModelNotExist"))
	}
	if err := json.Unmarshal([]byte(llmModel.Config), &cfg); err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.DefaultModelUpdateFailed"))
	}
	if cfg.LLM != opt.LLM.Model {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.ReasoningModelNotExist"))
	}

	embeddingModel := &sysconfig.Sysconfig{
		Type:   "model",
		Driver: opt.Embedding.Driver,
	}
	exist, err = embeddingModel.Get(ctx)
	if err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.ModelUpdateFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.EmbeddingModelNotExist"))
	}
	if err := json.Unmarshal([]byte(embeddingModel.Config), &cfg); err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.DefaultModelUpdateFailed"))
	}
	if cfg.Embedding != opt.Embedding.Model {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.EmbeddingModelNotExist"))
	}

	if err := sysconfig.ClearModelStatus(ctx); err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.DefaultModelUpdateFailed"))
	}

	llmModel.BeingUsed = true
	llmModel.Extra = "llm"
	if err := llmModel.Update(ctx); err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.DefaultModelUpdateFailed"))
	}
	embeddingModel.BeingUsed = true
	embeddingModel.Extra = "embedding"
	if err := embeddingModel.Update(ctx); err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.DefaultModelUpdateFailed"))
	}
	return nil, nil
}
