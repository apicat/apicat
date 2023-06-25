package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ExportResponse(exportType, download, filename string, content []byte, ctx *gin.Context) {
	switch download {
	case "true":
		switch exportType {
		case "HTML":
			ctx.Header("Content-Disposition", "attachment; filename="+filename+".html")
		case "md":
			ctx.Header("Content-Disposition", "attachment; filename="+filename+".md")
		default:
			ctx.Header("Content-Disposition", "attachment; filename="+filename+".json")
		}
		ctx.Data(http.StatusOK, "application/octet-stream", content)
	default:
		switch exportType {
		case "HTML":
			ctx.Data(http.StatusOK, "text/html; charset=utf-8", content)
		case "md":
			ctx.Data(http.StatusOK, "text/markdown; charset=utf-8", content)
		default:
			ctx.Data(http.StatusOK, "application/json", content)
		}
	}
}
