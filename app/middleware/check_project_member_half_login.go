package middleware

import (
	"net/http"
	"strconv"

	"github.com/apicat/apicat/common/encrypt"
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

		token := ctx.Query("token")
		if token == "" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    enum.ProjectMemberInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
			})
			ctx.Abort()
		}

		if token[:1] == "p" {
			if token[:1] != encrypt.GetMD5Encode(project.(*models.Projects).SharePassword) {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.AccessPasswordError"}),
				})
			}

			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    enum.ProjectMemberInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
			})
			ctx.Abort()
		}

		if token[:1] == "d" {
			collectionIDStr := ctx.Param("collection-id")
			if collectionIDStr == "" {
				ctx.JSON(http.StatusForbidden, gin.H{
					"code":    enum.ProjectMemberInsufficientPermissionsCode,
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
				})
				ctx.Abort()
			}

			collectionID, err := strconv.Atoi(collectionIDStr)
			if err != nil {
				ctx.JSON(http.StatusForbidden, gin.H{
					"code":    enum.ProjectMemberInsufficientPermissionsCode,
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
				})
				ctx.Abort()
			}

			collection, err := models.NewCollections(uint(collectionID))
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.NotFound"}),
				})
				ctx.Abort()
			}

			if token[:1] != encrypt.GetMD5Encode(collection.SharePassword) {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.AccessPasswordError"}),
				})
				ctx.Abort()
			}

			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    enum.ProjectMemberInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
			})
			ctx.Abort()
		}

		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		ctx.Abort()
	}
}
