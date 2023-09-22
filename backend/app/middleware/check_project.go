package middleware

import (
	"net/http"

	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/apicat/apicat/backend/models"

	"github.com/gin-gonic/gin"
)

type ProjectID struct {
	ID string `uri:"project-id" binding:"required,lte=255"`
}

func CheckProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data ProjectID

		responseCode := enum.Display404ErrorMessage
		if ctx.Request.Method == "GET" {
			responseCode = enum.Redirect404Page
		}

		if err := ctx.ShouldBindUri(&data); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    responseCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
			})
			ctx.Abort()
			return
		}

		project, err := models.NewProjects(data.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    responseCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
			})
			ctx.Abort()
			return
		}

		ctx.Set("CurrentProject", project)
	}
}
