package config

import (
	"fmt"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"regexp"
)

type AppFile struct {
	Name string `yaml:"name" env:"APICAT_APP_NAME"`
	Host string `yaml:"host" env:"APICAT_APP_HOST"`
	Port string `yaml:"port" env:"APICAT_APP_PORT"`
}

type LogFile struct {
	Path  string `yaml:"path" env:"APICAT_LOG_PATH"`
	Level string `yaml:"level" env:"APICAT_LOG_LEVEL"`
}

type DBFile struct {
	Driver   string `yaml:"driver" env:"APICAT_DB_DRIVER"`
	Path     string `yaml:"path" env:"APICAT_DB_PATH"`
	Host     string `yaml:"host" env:"APICAT_DB_HOST"`
	Port     string `yaml:"port" env:"APICAT_DB_PORT"`
	User     string `yaml:"user" env:"APICAT_DB_USER"`
	Password string `yaml:"password" env:"APICAT_DB_PASSWORD"`
	Dbname   string `yaml:"dbname" env:"APICAT_DB_NAME"`
}

type OpenAIFile struct {
	Source   string `yaml:"source" env:"APICAT_OPENAI_SOURCE"`
	Key      string `yaml:"key" env:"APICAT_OPENAI_KEY"`
	Endpoint string `yaml:"endpoint" env:"APICAT_OPENAI_ENDPOINT"`
}

type FileConfig struct {
	App    AppFile    `yaml:"application"`
	Log    LogFile    `yaml:"log"`
	DB     DBFile     `yaml:"database"`
	OpenAI OpenAIFile `yaml:"openai"`
}

type ConfigItem struct {
	Value      string
	DataSource string
	EnvName    string
}

type App struct {
	Name ConfigItem `env:"APICAT_APP_NAME"`
	Host ConfigItem `env:"APICAT_APP_HOST"`
	Port ConfigItem `env:"APICAT_APP_PORT"`
}

type Log struct {
	Path  ConfigItem `env:"APICAT_LOG_PATH"`
	Level ConfigItem `env:"APICAT_LOG_LEVEL"`
}

type DB struct {
	Driver   ConfigItem `env:"APICAT_DB_DRIVER"`
	Path     ConfigItem `env:"APICAT_DB_PATH"`
	Host     ConfigItem `env:"APICAT_DB_HOST"`
	Port     ConfigItem `env:"APICAT_DB_PORT"`
	User     ConfigItem `env:"APICAT_DB_USER"`
	Password ConfigItem `env:"APICAT_DB_PASSWORD"`
	Dbname   ConfigItem `env:"APICAT_DB_NAME"`
}

type OpenAI struct {
	Source   ConfigItem `env:"APICAT_OPENAI_SOURCE"`
	Key      ConfigItem `env:"APICAT_OPENAI_KEY"`
	Endpoint ConfigItem `env:"APICAT_OPENAI_ENDPOINT"`
}

type SysConfig struct {
	App    App
	Log    Log
	DB     DB
	OpenAI OpenAI
}

var (
	sysConfig *SysConfig
	// FilePath 启动项目时指定的配置文件
	FilePath string
	// DefaultConfigFilePath 默认配置文件路径
	DefaultConfigFilePath string = "./setting.default.yaml"
)

func createDefault() *SysConfig {
	return &SysConfig{
		App: App{
			Name: ConfigItem{
				Value:      "apicat",
				DataSource: "value",
			},
			Host: ConfigItem{
				Value:      "0.0.0.0",
				DataSource: "value",
			},
			Port: ConfigItem{
				Value:      "8000",
				DataSource: "value",
			},
		},
		Log: Log{
			Path: ConfigItem{
				Value:      "",
				DataSource: "value",
			},
			Level: ConfigItem{
				Value:      "debug",
				DataSource: "value",
			},
		},
		DB: DB{
			Driver: ConfigItem{
				Value:      "mysql",
				DataSource: "value",
			},
			Path: ConfigItem{
				Value:      "data/",
				DataSource: "value",
			},
			Host: ConfigItem{
				Value:      "127.0.0.1",
				DataSource: "value",
			},
			Port: ConfigItem{
				Value:      "3306",
				DataSource: "value",
			},
			User: ConfigItem{
				Value:      "root",
				DataSource: "value",
			},
			Dbname: ConfigItem{
				Value:      "apicat",
				DataSource: "value",
			},
		},
	}
}

func getEnvConfig() SysConfig {
	envConfig := SysConfig{}

	setEnvValues := func(structPtr interface{}, tag string) {
		valueOf := reflect.ValueOf(structPtr).Elem()
		valueType := reflect.TypeOf(structPtr).Elem()

		for i := 0; i < valueType.NumField(); i++ {
			envName := valueType.Field(i).Tag.Get(tag)
			if ev, exist := os.LookupEnv(envName); exist {
				field := valueOf.Field(i)
				field.FieldByName("Value").SetString(ev)
				field.FieldByName("DataSource").SetString("env")
				field.FieldByName("EnvName").SetString(envName)
			}
		}
	}

	setEnvValues(&envConfig.App, "env")
	setEnvValues(&envConfig.Log, "env")
	setEnvValues(&envConfig.DB, "env")
	setEnvValues(&envConfig.OpenAI, "env")

	return envConfig
}

func replaceEnvVars(fileConfig *FileConfig, sysConfig *SysConfig) {
	setEnvValues := func(fileStructPtr interface{}, sysStructPtr interface{}) {
		fileConfigValue := reflect.ValueOf(fileStructPtr).Elem()
		sysConfigValue := reflect.ValueOf(sysStructPtr).Elem() // 获取Sysconfig的可修改副本

		for i := 0; i < fileConfigValue.NumField(); i++ {
			fileField := fileConfigValue.Field(i)
			fileFieldValue := fileField.String()
			sysField := sysConfigValue.Field(i)

			envVarPattern := regexp.MustCompile(`\$\{(.+)\}`)
			if envVarPattern.MatchString(fileFieldValue) {
				matches := envVarPattern.FindAllStringSubmatch(fileConfigValue.Field(i).String(), -1)

				for _, match := range matches {
					if len(match) == 2 {
						envName := match[1]
						if ev, exist := os.LookupEnv(envName); exist {
							sysField.FieldByName("Value").SetString(ev)
							sysField.FieldByName("DataSource").SetString("env")
							sysField.FieldByName("EnvName").SetString(envName)
						}
					}
				}
			} else {
				sysField.FieldByName("Value").SetString(fileConfigValue.Field(i).String())
				sysField.FieldByName("DataSource").SetString("value")
			}
		}
	}

	setEnvValues(&fileConfig.App, &sysConfig.App)
	setEnvValues(&fileConfig.Log, &sysConfig.Log)
	setEnvValues(&fileConfig.DB, &sysConfig.DB)
	setEnvValues(&fileConfig.OpenAI, &sysConfig.OpenAI)
}

func loadConfig(filepath string) (*SysConfig, error) {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var fileConfig FileConfig
	if err := yaml.Unmarshal(fileContent, &fileConfig); err != nil {
		return nil, err
	}

	var sysConfig SysConfig
	replaceEnvVars(&fileConfig, &sysConfig)
	if err != nil {
		return nil, err
	}

	return &sysConfig, nil
}

func sysToFile(sysConfig *SysConfig) *FileConfig {
	fileConfig := &FileConfig{}

	setFileValues := func(sysStructPtr interface{}, fileStructPtr interface{}) {
		sysConfigValue := reflect.ValueOf(sysStructPtr).Elem()
		fileConfigValue := reflect.ValueOf(fileStructPtr).Elem()

		for i := 0; i < sysConfigValue.NumField(); i++ {
			sysField := sysConfigValue.Field(i)

			if sysField.FieldByName("DataSource").String() == "env" {
				fileConfigValue.Field(i).SetString(fmt.Sprintf("${%s}", sysField.FieldByName("EnvName").String()))
			} else {
				fileConfigValue.Field(i).SetString(sysField.FieldByName("Value").String())

			}
		}
	}

	setFileValues(&sysConfig.App, &fileConfig.App)
	setFileValues(&sysConfig.Log, &fileConfig.Log)
	setFileValues(&sysConfig.DB, &fileConfig.DB)
	setFileValues(&sysConfig.OpenAI, &fileConfig.OpenAI)

	return fileConfig
}

func SaveConfig(cfg *SysConfig) error {
	path := DefaultConfigFilePath
	if FilePath != "" {
		path = FilePath
	}

	fileConfig := sysToFile(cfg)

	data, err := yaml.Marshal(fileConfig)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(filePath string) (bool, error) {
	// 使用os.Stat()函数来获取文件信息
	_, err := os.Stat(filePath)

	// 如果文件存在
	if err == nil {
		// 文件存在
		return true, nil
	} else if os.IsNotExist(err) {
		// 文件不存在
		return false, nil
	} else {
		// 无法确定文件是否存在
		return false, err
	}
}

func GetSysConfig() *SysConfig {
	return sysConfig
}

func InitConfig() {
	cfg := createDefault()

	exist, err := fileExists(DefaultConfigFilePath)
	if err != nil {
		fmt.Println(err)
	}

	// 不存在默认配置文件并且未指定配置文件时，使用默认配置，创建默认配置文件
	if !exist && FilePath == "" {
		envCfg := getEnvConfig()
		mergo.Merge(&envCfg, cfg)
		sysConfig = &envCfg

		if err := SaveConfig(cfg); err != nil {
			fmt.Println(err)
		}

		return
	}

	// 指定配置文件时，使用指定配置文件，未指定时，使用默认配置文件
	filepath := DefaultConfigFilePath
	if FilePath != "" {
		filepath = FilePath
	}

	userCfg, err := loadConfig(filepath)
	if err != nil {
		panic(err.Error())
	}

	mergo.Merge(userCfg, cfg)
	sysConfig = userCfg
}
