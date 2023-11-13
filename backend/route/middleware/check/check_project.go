package check

import (
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/module/translator"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProjectID struct {
	ID string `uri:"project-id" binding:"required,lte=255"`
}

func CheckProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data ProjectID

		responseCode := proto.Display404ErrorMessage
		if ctx.Request.Method == "GET" {
			responseCode = proto.Redirect404Page
		}

		if err := ctx.ShouldBindUri(&data); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    responseCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
			})
			ctx.Abort()
			return
		}

		p, err := project.NewProjects(data.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    responseCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
			})
			ctx.Abort()
			return
		}

		ctx.Set("CurrentProject", p)
	}
}
