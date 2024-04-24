package config

import "time"

type Database struct {
	Debug    bool   `yaml:"Debug"`
	Host     string `yaml:"Host"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Database string `yaml:"Database"`

	MaxOpenConns    int           `yaml:"MaxOpenConns"`
	MaxIdleConns    int           `yaml:"MaxIdleConns"`
	ConnMaxIdleTime time.Duration `yaml:"ConnMaxIdleTime"`
}
