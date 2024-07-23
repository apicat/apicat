package config

import (
	"os"

	"github.com/apicat/apicat/v2/backend/module/oauth2"
	"github.com/joho/godotenv"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	App      *App
	Log      *lumberjack.Logger
	Database *Database
	Cache    *Cache
	Storage  *Storage
	Email    *Email
	Oauth2   map[string]oauth2.Config
	Model    *Model
	Vector   *Vector
}

var globalConf = &Config{}

func Load(path string) error {
	if _, err := os.Stat(path); err == nil {
		if err := godotenv.Load(path); err != nil {
			return err
		}
	}
	LoadAppConfig()
	LoadDatabaseConfig()
	LoadCacheConfig()
	LoadEmailConfig()
	LoadStorageConfig()
	LoadModelConfig()
	LoadVertorConfig()
	return nil
}

func Get() *Config {
	return globalConf
}

func Check() error {
	if err := CheckAppConfig(); err != nil {
		return err
	}
	if err := CheckDatabaseConfig(); err != nil {
		return err
	}
	if err := CheckCacheConfig(); err != nil {
		return err
	}
	if err := CheckStorageConfig(); err != nil {
		return err
	}
	if err := CheckEmailConfig(); err != nil {
		return err
	}
	if err := CheckModelConfig(); err != nil {
		return err
	}
	if err := CheckVectorConfig(); err != nil {
		return err
	}
	return nil
}
