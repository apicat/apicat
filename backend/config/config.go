package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/apicat/apicat/backend/module/oauth2"

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
		var conf Config
		if err := yaml.Unmarshal(b, &conf); err != nil {
			return err
		}
		globalConf = &conf
	}
	return nil
}

func LoadFromEnv() {
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
	if v, exists := os.LookupEnv("APICAT_CACHE_DRIVER"); exists {
		globalConf.Cache.Driver = v
	}
	if v, exists := os.LookupEnv("APICAT_CACHE_HOST"); exists {
		globalConf.Cache.Redis.Host = v
	}
	if v, exists := os.LookupEnv("APICAT_CACHE_PASSWORD"); exists {
		globalConf.Cache.Redis.Password = v
	}
	if v, exists := os.LookupEnv("APICAT_CACHE_DB"); exists {
		if i, err := strconv.Atoi(v); err == nil {
			globalConf.Cache.Redis.DB = i
		} else {
			globalConf.Cache.Redis.DB = 0
		}
	}
}

func Get() *Config {
	return globalConf
}

func getDefault() *Config {
	return &Config{
		App:      GetAppDefault(),
		Database: GetDatabaseDefault(),
		Cache:    GetCacheDefault(),
		Storage:  GetStorageDefault(),
	}
}
