package db

import (
	"github.com/apicat/apicat/backend/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type skipPath struct {
	URL    string
	Method string
}

var skipPaths = []skipPath{
	{"/config/db", "GET"},
	{"/api/config/db", "PUT"},
}

func CheckDBConnStatus(skip ...string) gin.HandlerFunc {
	skipPrefix := make(map[string]bool)
	for i := range skip {
		skipPrefix[skip[i]] = true
	}
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		for k := range skipPrefix {
			if strings.HasPrefix(path, k) {
				return
			}
		}

		for _, skipPath := range skipPaths {
			if skipPath.URL == ctx.Request.URL.Path && skipPath.Method == ctx.Request.Method {
				return
			}
		}

		connStatus, _ := model.DBConnStatus()
		if connStatus != 1 {
			ctx.Redirect(http.StatusFound, "/config/db")
			ctx.Abort()
			return
		}
	}
}
