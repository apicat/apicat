package middleware

import (
	"net/http"

	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/apicat/apicat/backend/models"
	"github.com/gin-gonic/gin"
)

func CheckProjectMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		project, exists := ctx.Get("CurrentProject")
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
			})
			ctx.Abort()
			return
		}

		user, exists := ctx.Get("CurrentUser")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    enum.InvalidOrIncorrectLoginToken,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Auth.TokenParsingFailed"}),
			})
			ctx.Abort()
			return
		}

		member, _ := models.NewProjectMembers()
		member.UserID = user.(*models.Users).ID
		member.ProjectID = project.(*models.Projects).ID

		if err := member.GetByUserIDAndProjectID(); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    enum.ProjectMemberInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
			})
			ctx.Abort()
			return
		}

		ctx.Set("CurrentProjectMember", member)
	}
}
