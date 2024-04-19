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
		if r.Driver == storage.LOCAL || r.Driver == storage.CLOUDFLARE || r.Driver == storage.QINIU {
			if err := json.Unmarshal([]byte(r.Config), &cfg); err == nil {
				config.SetStorage(&cfg)
			}
		}
	}
	defaultCfg := config.GetStorageDefault()
	defaultCfg.LocalDisk.Url = config.Get().App.AppUrl + "/uploads"
	config.SetStorage(defaultCfg)

	configJson, _ := json.Marshal(defaultCfg)
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
		if r.Driver == mailmodule.SMTP || r.Driver == mailmodule.SENDCLOUD {
			if err := json.Unmarshal([]byte(r.Config), &cfg); err == nil {
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
		if r.Driver == llm.OPENAI || r.Driver == llm.AZUREOPENAI {
			config.SetLLM(&cfg)
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
