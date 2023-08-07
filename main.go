package main

import (
	"flag"

	"github.com/apicat/apicat/backend/app"
	"github.com/apicat/apicat/backend/common/log"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/config"
	"github.com/apicat/apicat/backend/models"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "c", "", "The config file path, if not set, it will start with the example config.")
	flag.Parse()

	config.InitConfig(configFilePath)
	translator.Init()
	log.Init()
	models.Init()
	app.Run()
}
