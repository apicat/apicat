package middleware

import (
	"fmt"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type skipPath struct {
	URL    string
	Method string
}

var skipPaths = []skipPath{
	{"/api/config/db", "GET"},
	{"/api/config/db", "PUT"},
}

func CheckDBConnStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, skipPath := range skipPaths {
			if skipPath.URL == ctx.Request.URL.Path && skipPath.Method == ctx.Request.Method {
				ctx.Next()
				return
			}
		}

		connStatus, err := models.DBConnStatus()

		var tm string
		if connStatus != 1 {
			switch connStatus {
			case 2:
				tm = translator.Trasnlate(ctx, &translator.TT{ID: "DB.ConnectFailed"})
			case 3:
				tm = translator.Trasnlate(ctx, &translator.TT{ID: "DB.NotFound"})
			default:
				tm = translator.Trasnlate(ctx, &translator.TT{ID: "DB.ConnectFailed"})
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf(tm, err.Error()),
			})
			ctx.Abort()
			return
		}
	}
}
