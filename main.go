package main

import (
	"flag"

	"github.com/apicat/apicat/app"
	"github.com/apicat/apicat/common/log"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/config"
	"github.com/apicat/apicat/models"
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
