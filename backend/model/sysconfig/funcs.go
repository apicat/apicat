package sysconfig

import (
	"context"
	"encoding/json"
	"net/mail"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/module/cache"
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
	initCacheConfig()
	initStorageConfig()
	initEmailConfig()
	initModelConfig()
	initOauthConfig()
}

func initAppConfig() {
	var appConfig *config.App
	r := &Sysconfig{
		Type:   "service",
		Driver: "default",
	}
	exist, _ := r.Get(context.Background())
	if exist {
		if err := json.Unmarshal([]byte(r.Config), &appConfig); err == nil {
			config.SetApp(appConfig)
			return
		}
	}
	config.SetApp(config.GetAppDefault())
}

func initCacheConfig() {
	var cfg map[string]interface{}
	r := &Sysconfig{
		Type: "cache",
	}
	exist, _ := r.GetByUse(context.Background())
	if exist {
		if err := json.Unmarshal([]byte(r.Config), &cfg); err == nil {
			switch r.Driver {
			case cache.LOCAL:
				config.SetCache(&config.Cache{
					Driver: r.Driver,
				})
				return
			case cache.REDIS:
				config.SetCache(&config.Cache{
					Driver: r.Driver,
					Redis: &config.Redis{
						Host:     cfg["host"].(string),
						Password: cfg["password"].(string),
						DB:       int(cfg["database"].(float64)),
					},
				})
				return
			}
		}
	}
	defaultCfg := config.GetCacheDefault()
	config.SetCache(defaultCfg)

	configMap := map[string]interface{}{
		"host":     defaultCfg.Redis.Host,
		"password": defaultCfg.Redis.Password,
		"database": defaultCfg.Redis.DB,
	}
	configJson, _ := json.Marshal(configMap)
	r.Driver = defaultCfg.Driver
	r.BeingUsed = true
	r.Config = string(configJson)
	UpdateOrCreate(context.Background(), r)
}

func initStorageConfig() {
	var cfg map[string]interface{}
	r := &Sysconfig{
		Type: "storage",
	}
	exist, _ := r.GetByUse(context.Background())
	if exist {
		if err := json.Unmarshal([]byte(r.Config), &cfg); err == nil {
			switch r.Driver {
			case storage.LOCAL:
				config.SetStorage(&config.Storage{
					Driver: r.Driver,
					LocalDisk: &config.LocalDisk{
						Path: cfg["path"].(string),
						Url:  config.Get().App.AppUrl + "/uploads",
					},
				})
				return
			case storage.CLOUDFLARE:
				config.SetStorage(&config.Storage{
					Driver: r.Driver,
					Cloudflare: &config.Cloudflare{
						AccountID:       cfg["accountID"].(string),
						AccessKeyID:     cfg["accessKeyID"].(string),
						AccessKeySecret: cfg["accessKeySecret"].(string),
						BucketName:      cfg["bucketName"].(string),
						Url:             cfg["bucketUrl"].(string),
					},
				})
				return
			case storage.QINIU:
				config.SetStorage(&config.Storage{
					Driver: r.Driver,
					Qiniu: &config.Qiniu{
						AccessKeyID:     cfg["accessKey"].(string),
						AccessKeySecret: cfg["secretKey"].(string),
						BucketName:      cfg["bucketName"].(string),
						Url:             cfg["bucketUrl"].(string),
					},
				})
				return
			}
		}
	}
	defaultCfg := config.GetStorageDefault()
	defaultCfg.LocalDisk.Url = config.Get().App.AppUrl + "/uploads"
	config.SetStorage(defaultCfg)

	configMap := map[string]string{
		"path": defaultCfg.LocalDisk.Path,
	}
	configJson, _ := json.Marshal(configMap)
	r.Driver = defaultCfg.Driver
	r.BeingUsed = true
	r.Config = string(configJson)
	UpdateOrCreate(context.Background(), r)
}

func initEmailConfig() {
	var cfg map[string]interface{}
	r := &Sysconfig{
		Type: "email",
	}
	exist, _ := r.GetByUse(context.Background())
	if exist {
		if err := json.Unmarshal([]byte(r.Config), &cfg); err == nil {
			switch r.Driver {
			case mailmodule.SMTP:
				config.SetEmail(&config.Email{
					Driver: r.Driver,
					Smtp: &config.EmailSmtp{
						Host: cfg["host"].(string),
						From: mail.Address{
							Name:    cfg["user"].(string),
							Address: cfg["address"].(string),
						},
						Password: cfg["password"].(string),
					},
				})
			case mailmodule.SENDCLOUD:
				config.SetEmail(&config.Email{
					Driver: r.Driver,
					SendCloud: &config.EmailSendCloud{
						ApiUser:  cfg["apiUser"].(string),
						ApiKey:   cfg["apiKey"].(string),
						From:     cfg["fromEmail"].(string),
						FromName: cfg["fromName"].(string),
					},
				})
			}
		}
	}
}

func initModelConfig() {
	var cfg map[string]interface{}
	r := &Sysconfig{
		Type: "model",
	}
	exist, _ := r.GetByUse(context.Background())
	if exist {
		if err := json.Unmarshal([]byte(r.Config), &cfg); err == nil {
			switch r.Driver {
			case llm.OPENAI:
				config.SetLLM(&config.LLM{
					Driver: r.Driver,
					OpenAI: &config.OpenAI{
						ApiKey:         cfg["apiKey"].(string),
						OrganizationID: cfg["organizationID"].(string),
						ApiBase:        cfg["apiBase"].(string),
						LLMName:        cfg["llmName"].(string),
						EmbeddingName:  cfg["embeddingName"].(string),
					},
				})
			case llm.AZUREOPENAI:
				config.SetLLM(&config.LLM{
					Driver: r.Driver,
					AzureOpenAI: &config.AzureOpenAI{
						ApiKey:        cfg["apiKey"].(string),
						Endpoint:      cfg["endpoint"].(string),
						LLMName:       cfg["llmName"].(string),
						EmbeddingName: cfg["embeddingName"].(string),
					},
				})
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
