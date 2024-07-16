package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Database struct {
	Host     string
	Username string
	Password string
	Database string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
}

func LoadDatabaseConfig() {
	globalConf.Database = &Database{}
	if v, exists := os.LookupEnv("DB_HOST"); exists {
		globalConf.Database.Host = v
	}
	if v, exists := os.LookupEnv("DB_USERNAME"); exists {
		globalConf.Database.Username = v
	}
	if v, exists := os.LookupEnv("DB_PASSWORD"); exists {
		globalConf.Database.Password = v
	}
	if v, exists := os.LookupEnv("DB_DATABASE"); exists {
		globalConf.Database.Database = v
	}
	if v, exists := os.LookupEnv("DB_MAX_OPEN_CONNS"); exists && v != "" {
		if num, err := strconv.Atoi(v); err == nil {
			globalConf.Database.MaxOpenConns = num
		}
	}
	if v, exists := os.LookupEnv("DB_MAX_IDLE_CONNS"); exists && v != "" {
		if num, err := strconv.Atoi(v); err == nil {
			globalConf.Database.MaxIdleConns = num
		}
	}
	if v, exists := os.LookupEnv("DB_CONN_MAX_IDLE_TIME"); exists && v != "" {
		if second, err := strconv.Atoi(v); err == nil {
			globalConf.Database.ConnMaxIdleTime = time.Duration(second)
		}
	}
}

func CheckDatabaseConfig() error {
	if globalConf.Database.Host == "" {
		return errors.New("database host is empty")
	}
	if globalConf.Database.Username == "" {
		return errors.New("database username is empty")
	}
	if globalConf.Database.Database == "" {
		return errors.New("database name is empty")
	}
	return nil
}
