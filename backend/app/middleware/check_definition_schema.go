package middleware

import (
	"net/http"

	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/models"
	"github.com/gin-gonic/gin"
)

type SchemaUriData struct {
	ProjectID string `uri:"project-id" binding:"required,gt=0"`
	SchemaID  uint   `uri:"schemas-id" binding:"required,gt=0"`
}

// 需要先通过CheckProject中间件。检验模型是否存在，是否所属请求对应的项目
func CheckDefinitionSchema() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentProject, _ := ctx.Get("CurrentProject")

		var data SchemaUriData
		if err := ctx.ShouldBindUri(&data); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
			})
			ctx.Abort()
			return
		}

		ds, err := models.NewDefinitionSchemas(data.SchemaID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
			})
			ctx.Abort()
			return
		}

		if ds.ProjectId != currentProject.(*models.Projects).ID {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
			})
			ctx.Abort()
			return
		}

		ctx.Set("CurrentDefinitionSchema", ds)
	}
}
