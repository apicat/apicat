package config

import (
	"os"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type App struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Log struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

type DB struct {
	Driver   string `yaml:"driver"`
	Path     string `yaml:"path"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

type OpenAI struct {
	Token string `yaml:"token"`
}

type Sysconfig struct {
	App    App    `yaml:"application"`
	Log    Log    `yaml:"log"`
	DB     DB     `yaml:"database"`
	OpenAI OpenAI `yaml:"openai"`
}

func createDefault() *Sysconfig {
	return &Sysconfig{
		App: App{
			Name: "apicat",
			Host: "0.0.0.0",
			Port: 8000,
		},
		Log: Log{
			Path:  "logs/",
			Level: "debug",
		},
		DB: DB{
			Driver:  "sqlite",
			Path:    "data/",
			Host:    "127.0.0.1",
			Port:    3306,
			User:    "root",
			Dbname:  "apicat",
			Charset: "utf8mb4",
		},
	}
}

var SysConfig *Sysconfig

func InitConfig(configFile string) {
	cfg := createDefault()
	if configFile == "" {
		SysConfig = cfg
		return
	}
	file, err := os.ReadFile(configFile)
	if err != nil {
		panic(err.Error())
	}
	var userCfg Sysconfig
	if err := yaml.Unmarshal(file, &userCfg); err != nil {
		panic(err.Error())
	}
	mergo.Merge(&userCfg, cfg) //nolint:errcheck
	SysConfig = &userCfg
}
