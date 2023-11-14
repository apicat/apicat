package apicat

import (
	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/module/logger"
	"github.com/apicat/apicat/backend/route"

	"github.com/apicat/apicat/backend/config"
)

type App struct{}

func NewApp(conf string) *App {
	config.FilePath = conf
	config.InitConfig()
	return &App{}
}

func (a *App) Run() error {
	inits := []func(){
		i18n.Init,
		logger.Init,
		model.Init,
		route.Init,
	}
	for _, v := range inits {
		v()
	}
	return nil
}
