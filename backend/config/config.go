package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/oauth2"

	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App      *App                     `yaml:"App"`
	Log      *lumberjack.Logger       `yaml:"Log"`
	Database *Database                `yaml:"Database"`
	Cache    *Cache                   `yaml:"Cache"`
	Storage  *Storage                 `yaml:"Storage"`
	Email    *Email                   `yaml:"Email"`
	Oauth2   map[string]oauth2.Config `yaml:"Oauth2"`
	LLM      *LLM                     `yaml:"LLM"`
}

var globalConf = getDefault()

func Load(path string) error {
	b, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(b, globalConf); err != nil {
			return err
		}
	}
	return nil
}

func LoadFromEnv() {
	if v, exists := os.LookupEnv("APICAT_APP_SERVER_BIND"); exists {
		globalConf.App.AppServerBind = v
	}
	if v, exists := os.LookupEnv("APICAT_MOCK_SERVER_BIND"); exists {
		globalConf.App.MockServerBind = v
	}
	if v, exists := os.LookupEnv("APICAT_DEBUG"); exists {
		globalConf.Database.Debug = strings.ToLower(v) == "true"
	}
	if v, exists := os.LookupEnv("APICAT_DB_HOST"); exists {
		globalConf.Database.Host = v
	}
	if v, exists := os.LookupEnv("APICAT_DB_USERNAME"); exists {
		globalConf.Database.Username = v
	}
	if v, exists := os.LookupEnv("APICAT_DB_PASSWORD"); exists {
		globalConf.Database.Password = v
	}
	if v, exists := os.LookupEnv("APICAT_DB_DATABASE"); exists {
		globalConf.Database.Database = v
	}
	if v, exists := os.LookupEnv("APICAT_CACHE_HOST"); exists {
		globalConf.Cache.Host = v
	}
	if v, exists := os.LookupEnv("APICAT_CACHE_PASSWORD"); exists {
		globalConf.Cache.Password = v
	}
	if v, exists := os.LookupEnv("APICAT_CACHE_DB"); exists {
		if i, err := strconv.Atoi(v); err == nil {
			globalConf.Cache.DB = i
		} else {
			globalConf.Cache.DB = 0
		}
	}
}

func Get() *Config {
	return globalConf
}

func Check() error {
	if globalConf.App == nil {
		return errors.New("app config is nil")
	}
	if globalConf.App.AppServerBind == "" {
		return errors.New("app server bind is empty")
	}
	if globalConf.App.MockServerBind == "" {
		return errors.New("mock server bind is empty")
	}
	if globalConf.Database == nil {
		return errors.New("database config is nil")
	}
	if globalConf.Database.Host == "" {
		return errors.New("database host is empty")
	}
	if globalConf.Database.Username == "" {
		return errors.New("database username is empty")
	}
	if globalConf.Database.Database == "" {
		return errors.New("database name is empty")
	}
	if globalConf.Cache == nil {
		return errors.New("cache config is nil")
	}
	if globalConf.Cache.Host == "" {
		return errors.New("cache host is empty")
	}
	if globalConf.Cache.DB < 0 || globalConf.Cache.DB > 15 {
		return errors.New("cache db is invalid")
	}
	return nil
}

func getDefault() *Config {
	return &Config{
		App:      GetAppDefault(),
		Database: &Database{},
		Cache:    &Cache{},
		Storage:  GetStorageDefault(),
	}
}
