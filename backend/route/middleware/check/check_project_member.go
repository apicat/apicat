package check

import (
	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckProjectMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p, exists := ctx.Get("CurrentProject")
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Projects.NotFound"}),
			})
			ctx.Abort()
			return
		}

		u, exists := ctx.Get("CurrentUser")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    proto.InvalidOrIncorrectLoginToken,
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Auth.TokenParsingFailed"}),
			})
			ctx.Abort()
			return
		}

		member, _ := project.NewProjectMembers()
		member.UserID = u.(*user.Users).ID
		member.ProjectID = p.(*project.Projects).ID

		if err := member.GetByUserIDAndProjectID(); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    proto.ProjectMemberInsufficientPermissionsCode,
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
			})
			ctx.Abort()
			return
		}

		ctx.Set("CurrentProjectMember", member)
	}
}
