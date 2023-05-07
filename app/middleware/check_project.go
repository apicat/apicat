package middleware

import (
	"net/http"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"

	"github.com/gin-gonic/gin"
)

type ProjectID struct {
	ID string `uri:"id" binding:"required,lte=255"`
}

func CheckProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data ProjectID

		if err := ctx.ShouldBindUri(&data); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
			})
			ctx.Abort()
		} else {
			project, err := models.NewProjects(data.ID)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
				})
				ctx.Abort()
			}
			ctx.Set("CurrentProject", project)
		}
	}
}
