package app

import (
	"github.com/apicat/apicat/backend/app/router"
	"github.com/apicat/apicat/backend/config"
	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.ContextWithFallback = true

	router.InitApiRouter(r)
	r.Run(config.GetSysConfig().App.Host.Value + ":" + config.GetSysConfig().App.Port.Value)
}
