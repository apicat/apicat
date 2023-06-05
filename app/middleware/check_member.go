package middleware

import (
	"net/http"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

func CheckMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		project, _ := ctx.Get("CurrentProject")
		user, _ := ctx.Get("CurrentUser")

		member, _ := models.NewProjectMembers()
		member.UserID = user.(*models.Users).ID
		member.ProjectID = project.(*models.Projects).ID

		if err := member.Get(); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
			})
			ctx.Abort()
		}

		ctx.Set("CurrentMember", member)
	}
}
