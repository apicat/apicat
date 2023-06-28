package middleware

import (
	"net/http"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/enum"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

func CheckProjectMemberHalfLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		project, exists := ctx.Get("CurrentProject")
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
			})
			ctx.Abort()
		}

		user, exists := ctx.Get("CurrentUser")
		if exists {
			member, _ := models.NewProjectMembers()
			member.UserID = user.(*models.Users).ID
			member.ProjectID = project.(*models.Projects).ID

			if err := member.GetByUserIDAndProjectID(); err == nil {
				ctx.Set("CurrentProjectMember", member)
				return
			}
		}

		// 判断是否为公开项目
		if project.(*models.Projects).Visibility == 1 {
			return
		}

		// 判断项目是否被分享
		if project.(*models.Projects).SharePassword != "" {
			return
		}

		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		ctx.Abort()
	}
}
