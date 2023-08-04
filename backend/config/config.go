package config

import (
	"os"
	"reflect"
	"strconv"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type App struct {
	Name string `yaml:"name" env:"APICAT_APP_NAME"`
	Host string `yaml:"host" env:"APICAT_APP_HOST"`
	Port int    `yaml:"port" env:"APICAT_APP_PORT"`
}

type Log struct {
	Path  string `yaml:"path" env:"APICAT_LOG_PATH"`
	Level string `yaml:"level" env:"APICAT_LOG_LEVEL"`
}

type DB struct {
	Driver   string `yaml:"driver" env:"APICAT_DB_DRIVER"`
	Path     string `yaml:"path" env:"APICAT_DB_PATH"`
	Host     string `yaml:"host" env:"APICAT_DB_HOST"`
	Port     int    `yaml:"port" env:"APICAT_DB_PORT"`
	User     string `yaml:"user" env:"APICAT_DB_USER"`
	Password string `yaml:"password" env:"APICAT_DB_PASSWORD"`
	Dbname   string `yaml:"dbname" env:"APICAT_DB_NAME"`
}

type OpenAI struct {
	Key string `yaml:"key" env:"APICAT_OPENAI_KEY"`
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
			Path:  "",
			Level: "debug",
		},
		DB: DB{
			Driver: "mysql",
			Path:   "data/",
			Host:   "127.0.0.1",
			Port:   3306,
			User:   "root",
			Dbname: "apicat",
		},
	}
}

func getEnvConfig() Sysconfig {
	envConfig := Sysconfig{}

	v := reflect.ValueOf(&envConfig.App).Elem()
	t := reflect.TypeOf(envConfig.App)
	for i := 0; i < t.NumField(); i++ {
		envName := t.Field(i).Tag.Get("env")
		if ev, exist := os.LookupEnv(envName); exist {
			if t.Field(i).Name == "Port" {
				if p, err := strconv.Atoi(ev); err == nil {
					v.Field(i).SetInt(int64(p))
				}
			} else {
				v.Field(i).SetString(ev)
			}
		}
	}

	v = reflect.ValueOf(&envConfig.Log).Elem()
	t = reflect.TypeOf(envConfig.Log)
	for i := 0; i < t.NumField(); i++ {
		envName := t.Field(i).Tag.Get("env")
		if ev, exist := os.LookupEnv(envName); exist {
			v.Field(i).SetString(ev)
		}
	}

	v = reflect.ValueOf(&envConfig.DB).Elem()
	t = reflect.TypeOf(envConfig.DB)
	for i := 0; i < t.NumField(); i++ {
		envName := t.Field(i).Tag.Get("env")
		if ev, exist := os.LookupEnv(envName); exist {
			if t.Field(i).Name == "Port" {
				if p, err := strconv.Atoi(ev); err == nil {
					v.Field(i).SetInt(int64(p))
				}
			} else {
				v.Field(i).SetString(ev)
			}
		}
	}

	v = reflect.ValueOf(&envConfig.OpenAI).Elem()
	t = reflect.TypeOf(envConfig.OpenAI)
	for i := 0; i < t.NumField(); i++ {
		envName := t.Field(i).Tag.Get("env")
		if ev, exist := os.LookupEnv(envName); exist {
			v.Field(i).SetString(ev)
		}
	}

	return envConfig
}

var SysConfig *Sysconfig

func InitConfig(configFile string) {
	cfg := createDefault()
	if configFile == "" {
		envCfg := getEnvConfig()
		mergo.Merge(&envCfg, cfg)
		SysConfig = &envCfg
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
	mergo.Merge(&userCfg, cfg)

	SysConfig = &userCfg
}
