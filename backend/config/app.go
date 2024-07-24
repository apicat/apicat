package config

import (
	"errors"
	"os"
	"strings"
)

type App struct {
	Debug          bool
	AppUrl         string
	AppServerBind  string
	MockUrl        string
	MockServerBind string
}

func LoadAppConfig() {
	globalConf.App = &App{
		AppUrl:         "http://localhost:8000",
		AppServerBind:  "0.0.0.0:8000",
		MockUrl:        "http://localhost:8001",
		MockServerBind: "0.0.0.0:8001",
	}

	if v, exists := os.LookupEnv("APP_DEBUG"); exists {
		globalConf.App.Debug = strings.ToLower(v) == "true"
	}
	if v, exists := os.LookupEnv("APP_URL"); exists {
		globalConf.App.AppUrl = v
	}
	if v, exists := os.LookupEnv("APP_SERVER_BIND"); exists {
		globalConf.App.AppServerBind = v
	}
	if v, exists := os.LookupEnv("MOCK_URL"); exists {
		globalConf.App.MockUrl = v
	}
	if v, exists := os.LookupEnv("MOCK_SERVER_BIND"); exists {
		globalConf.App.MockServerBind = v
	}
}

func CheckAppConfig() error {
	if globalConf.App.AppServerBind == "" {
		return errors.New("app server bind is empty")
	}
	if globalConf.App.MockServerBind == "" {
		return errors.New("mock server bind is empty")
	}
	return nil
}

func GetApp() *App {
	return globalConf.App
}

func SetApp(appConfig *App) {
	globalConf.App.AppUrl = appConfig.AppUrl
	globalConf.App.MockUrl = appConfig.MockUrl
}
