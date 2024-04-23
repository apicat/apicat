package sysconfig

import (
	"context"
	"encoding/json"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/module/llm"
	mailmodule "github.com/apicat/apicat/v2/backend/module/mail"
	"github.com/apicat/apicat/v2/backend/module/oauth2"
	"github.com/apicat/apicat/v2/backend/module/storage"

	"gorm.io/gorm"
)

func GetList(ctx context.Context, t string) ([]*Sysconfig, error) {
	var list []*Sysconfig
	err := model.DB(ctx).Where("type = ?", t).Find(&list).Error
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

func Init() {
	initAppConfig()
	initStorageConfig()
	initEmailConfig()
	initModelConfig()
	initOauthConfig()
}

func initAppConfig() {
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

func initStorageConfig() {
	var cfg config.Storage
	r := &Sysconfig{
		Type: "storage",
	}
	exist, _ := r.GetByUse(context.Background())
	if exist {
		cfg.Driver = r.Driver

		switch r.Driver {
		case storage.LOCAL:
			if err := json.Unmarshal([]byte(r.Config), &cfg.LocalDisk); err == nil {
				config.SetStorage(&cfg)
				return
			}
		case storage.CLOUDFLARE:
			if err := json.Unmarshal([]byte(r.Config), &cfg.Cloudflare); err == nil {
				config.SetStorage(&cfg)
				return
			}
		case storage.QINIU:
			if err := json.Unmarshal([]byte(r.Config), &cfg.Qiniu); err == nil {
				config.SetStorage(&cfg)
				return
			}
		}
	}
	defaultCfg := config.GetStorageDefault()
	defaultCfg.LocalDisk.Url = config.Get().App.AppUrl + "/uploads"
	config.SetStorage(defaultCfg)

	configJson, _ := json.Marshal(defaultCfg.LocalDisk)
	r.Driver = defaultCfg.Driver
	r.BeingUsed = true
	r.Config = string(configJson)
	UpdateOrCreate(context.Background(), r)
}

func initEmailConfig() {
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

func initModelConfig() {
	var cfg config.LLM
	r := &Sysconfig{
		Type: "model",
	}
	exist, _ := r.GetByUse(context.Background())
	if exist {
		switch r.Driver {
		case llm.OPENAI:
			if err := json.Unmarshal([]byte(r.Config), &cfg.OpenAI); err == nil {
				cfg.Driver = llm.OPENAI
				config.SetLLM(&cfg)
			}
		case llm.AZUREOPENAI:
			if err := json.Unmarshal([]byte(r.Config), &cfg.AzureOpenAI); err == nil {
				cfg.Driver = llm.AZUREOPENAI
				config.SetLLM(&cfg)
			}
		}
	}
}

func initOauthConfig() {
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
