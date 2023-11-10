package app

import (
	"github.com/apicat/apicat/backend/config"
	"github.com/apicat/apicat/backend/route"
	"github.com/apicat/apicat/frontend"
	"github.com/gin-gonic/gin"
	"html/template"
)

func Run() {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.ContextWithFallback = true

	t, _ := template.ParseFS(frontend.FrontDist, "dist/templates/*.tmpl")
	r.SetHTMLTemplate(t)

	route.InitApiRouter(r)
	r.Run(config.GetSysConfig().App.Host.Value + ":" + config.GetSysConfig().App.Port.Value)
}
