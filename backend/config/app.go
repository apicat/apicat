package config

import (
	"encoding/json"
)

type App struct {
	Debug          bool   `yaml:"Debug"`
	AppName        string `yaml:"AppName"`
	AppUrl         string `yaml:"AppUrl"`
	AppServerBind  string `yaml:"AppServerBind"`
	MockUrl        string `yaml:"MockUrl"`
	MockServerBind string `yaml:"MockServerBind"`
}

func GetAppDefault() *App {
	return &App{
		AppName:        "ApiCat",
		AppUrl:         "http://localhost:8000",
		AppServerBind:  "0.0.0.0:8000",
		MockUrl:        "http://localhost:8001",
		MockServerBind: "0.0.0.0:8001",
	}
}

func SetApp(appConfig *App) {
	globalConf.App.AppName = appConfig.AppName
	globalConf.App.AppUrl = appConfig.AppUrl
	globalConf.App.MockUrl = appConfig.MockUrl
}

func (a *App) ToMapInterface() map[string]interface{} {
	var (
		res      map[string]interface{}
		jsonByte []byte
	)
	jsonByte, _ = json.Marshal(a)
	_ = json.Unmarshal(jsonByte, &res)
	return res
}
