package main

import (
	"flag"
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/module/logger"

	"github.com/apicat/apicat/backend/app"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/config"
)

func main() {
	flag.StringVar(&config.FilePath, "c", "", "The config file path, if not set, it will start with the example config.")
	flag.Parse()

	config.InitConfig()
	translator.Init()
	logger.Init()
	model.Init()
	app.Run()
}
