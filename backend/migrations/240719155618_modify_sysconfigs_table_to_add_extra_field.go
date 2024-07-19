package migrations

import (
	"encoding/json"

	"github.com/go-gormigrate/gormigrate/v2"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

func init() {
	type sysconfig struct {
		ID        uint   `gorm:"primarykey"`
		Type      string `gorm:"type:varchar(255);uniqueIndex:ukey;not null;comment:Configuration type"`
		Driver    string `gorm:"type:varchar(255);uniqueIndex:ukey;not null"`
		BeingUsed bool   `gorm:"type:tinyint;comment:is using"`
		Config    string `gorm:"type:varchar(512);"`
		Extra     string `gorm:"type:varchar(512);"`
	}

	type oldOpenAI struct {
		ApiKey         string `yaml:"apiKey"`
		OrganizationID string `yaml:"organizationID"`
		ApiBase        string `yaml:"apiBase"`
		LLMName        string `yaml:"llmName"`
		EmbeddingName  string `yaml:"embeddingName"`
	}
	type newOpenAI struct {
		ApiKey         string `json:"apiKey"`
		OrganizationID string `json:"organizationID"`
		ApiBase        string `json:"apiBase"`
		LLM            string `json:"llm"`
		Embedding      string `json:"embedding"`
	}

	type oldAzureOpenAI struct {
		ApiKey        string `yaml:"apiKey"`
		Endpoint      string `yaml:"endpoint"`
		LLMName       string `yaml:"llmName"`
		EmbeddingName string `yaml:"embeddingName"`
	}
	type newAzureOpenAI struct {
		ApiKey    string `json:"apiKey"`
		Endpoint  string `json:"endpoint"`
		LLM       string `json:"llm"`
		Embedding string `json:"embedding"`
	}

	m := &gormigrate.Migration{
		ID: "240719155618",
		Migrate: func(tx *gorm.DB) error {
			if tx.Migrator().HasTable(&sysconfig{}) {
				if !tx.Migrator().HasColumn(&sysconfig{}, "extra") {
					if err := tx.Migrator().AddColumn(&sysconfig{}, "extra"); err != nil {
						return err
					}
				}

				var list []*sysconfig
				if err := tx.Where("type = ?", "model").Find(&list).Error; err != nil {
					return err
				}

				for _, v := range list {
					if v.Driver == "openai" {
						var old oldOpenAI
						if err := yaml.Unmarshal([]byte(v.Config), &old); err != nil {
							return err
						}

						new := newOpenAI{
							ApiKey:         old.ApiKey,
							OrganizationID: old.OrganizationID,
							ApiBase:        old.ApiBase,
							LLM:            old.LLMName,
							Embedding:      old.EmbeddingName,
						}

						newConfig, err := json.Marshal(new)
						if err != nil {
							return err
						}

						data := map[string]interface{}{
							"config": string(newConfig),
						}
						if v.BeingUsed {
							data["extra"] = "llm"
						}

						if err := tx.Model(&sysconfig{}).Where("id = ?", v.ID).Updates(data).Error; err != nil {
							return err
						}
					} else if v.Driver == "azure-openai" {
						var old oldAzureOpenAI
						if err := yaml.Unmarshal([]byte(v.Config), &old); err != nil {
							return err
						}

						new := newAzureOpenAI{
							ApiKey:    old.ApiKey,
							Endpoint:  old.Endpoint,
							LLM:       old.LLMName,
							Embedding: old.EmbeddingName,
						}

						newConfig, err := json.Marshal(new)
						if err != nil {
							return err
						}

						data := map[string]interface{}{
							"config": string(newConfig),
						}
						if v.BeingUsed {
							data["extra"] = "llm"
						}

						if err := tx.Model(&sysconfig{}).Where("id = ?", v.ID).Updates(data).Error; err != nil {
							return err
						}
					}
				}
			}
			return nil
		},
	}
	MigrationHelper.Register(m)
}
