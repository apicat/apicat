package apicat

import (
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/module/logger"
	"github.com/apicat/apicat/backend/module/translator"
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
		translator.Init,
		logger.Init,
		model.Init,
		route.Init,
	}
	for _, v := range inits {
		v()
	}
	return nil
}
