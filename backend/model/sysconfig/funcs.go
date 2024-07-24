package sysconfig

import (
	"context"
	"encoding/json"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/model"
	mailmodule "github.com/apicat/apicat/v2/backend/module/mail"
	aimodel "github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/oauth2"

	"gorm.io/gorm"
)

func GetList(ctx context.Context, t string) ([]*Sysconfig, error) {
	var list []*Sysconfig
	err := model.DB(ctx).Where("type = ?", t).Find(&list).Error
	return list, err
}

func GetUseingList(ctx context.Context, t string) ([]*Sysconfig, error) {
	var list []*Sysconfig
	err := model.DB(ctx).Where("type = ? and being_used = ?", t, 1).Find(&list).Error
	return list, err
}

func UpdateOrCreate(ctx context.Context, sc *Sysconfig) error {
	r := &Sysconfig{
		Type:   sc.Type,
		Driver: sc.Driver,
	}
	exist, err := r.Get(ctx)
	if err != nil {
		return err
	}

	return model.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Sysconfig{}).Where("type = ?", sc.Type).Update("being_used", 0).Error; err != nil {
			return err
		}

		if !exist {
			if err := tx.Model(&Sysconfig{}).Create(sc).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&Sysconfig{}).Where("id = ?", r.ID).Updates(map[string]interface{}{
				"being_used": 1,
				"config":     sc.Config,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func ClearModelStatus(ctx context.Context) error {
	return model.DB(ctx).Model(&Sysconfig{}).Where("type = ?", "model").Updates(map[string]interface{}{"extra": "", "being_used": 0}).Error
}

func Load() {
	loadAppConfig()
	loadEmailConfig()
	loadModelConfig()
	loadOauthConfig()
}

func loadAppConfig() {
	r := &Sysconfig{
		Type:   "service",
		Driver: "default",
	}
	exist, _ := r.Get(context.Background())
	if exist {
		var appConfig *config.App
		if err := json.Unmarshal([]byte(r.Config), &appConfig); err == nil {
			config.SetApp(appConfig)
			return
		}
	}
}

func loadEmailConfig() {
	var cfg config.Email
	r := &Sysconfig{
		Type: "email",
	}
	exist, _ := r.GetByUse(context.Background())
	if exist {
		switch r.Driver {
		case mailmodule.SMTP:
			if err := json.Unmarshal([]byte(r.Config), &cfg.Smtp); err == nil {
				cfg.Driver = mailmodule.SMTP
				config.SetEmail(&cfg)
			}
		case mailmodule.SENDCLOUD:
			if err := json.Unmarshal([]byte(r.Config), &cfg.SendCloud); err == nil {
				cfg.Driver = mailmodule.SENDCLOUD
				config.SetEmail(&cfg)
			}
		}
	}
}

func loadModelConfig() {
	if models, err := GetUseingList(context.Background(), "model"); err == nil {
		var cfg config.Model
		for _, m := range models {
			switch m.Driver {
			case aimodel.OPENAI:
				if err := json.Unmarshal([]byte(m.Config), &cfg.OpenAI); err == nil {
					if m.Extra == "llm" {
						cfg.LLMDriver = aimodel.OPENAI
					} else if m.Extra == "embedding" {
						cfg.EmbeddingDriver = aimodel.OPENAI
					} else if m.Extra == "llm,embedding" {
						cfg.LLMDriver = aimodel.OPENAI
						cfg.EmbeddingDriver = aimodel.OPENAI
					}
				}
			case aimodel.AZURE_OPENAI:
				if err := json.Unmarshal([]byte(m.Config), &cfg.AzureOpenAI); err == nil {
					if m.Extra == "llm" {
						cfg.LLMDriver = aimodel.AZURE_OPENAI
					} else if m.Extra == "embedding" {
						cfg.EmbeddingDriver = aimodel.AZURE_OPENAI
					} else if m.Extra == "llm,embedding" {
						cfg.LLMDriver = aimodel.AZURE_OPENAI
						cfg.EmbeddingDriver = aimodel.AZURE_OPENAI
					}
				}
			}
		}
		config.SetModel(&cfg)
	}
}

func loadOauthConfig() {
	var cfg map[string]interface{}
	r := &Sysconfig{
		Type: "oauth",
	}
	exist, _ := r.GetByUse(context.Background())
	if exist {
		if err := json.Unmarshal([]byte(r.Config), &cfg); err == nil {
			switch r.Driver {
			case "github":
				syscfg := config.Get()
				syscfg.Oauth2 = map[string]oauth2.Config{
					"github": {
						ClientID:     cfg["clientID"].(string),
						ClientSecret: cfg["clientSecret"].(string),
					},
				}
			}
		}
	}
}
