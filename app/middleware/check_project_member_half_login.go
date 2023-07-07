package middleware

import (
	"net/http"
	"strconv"
	"time"

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
			return
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
		if token == "" || len(token) < 1 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    enum.ShareTokenInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
			})
			ctx.Abort()
			return
		}

		// 校验访问令牌是否存在
		sr := models.NewShareRecords()
		sr.ShareToken = encrypt.GetMD5Encode(token)
		if err := sr.GetByShareToken(); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    enum.ShareTokenInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
			})
			ctx.Abort()
			return
		}

		// 校验访问令牌是否过期
		now := time.Now()
		if sr.Expiration.Before(now) {
			if err := sr.Delete(); err != nil {
				ctx.JSON(http.StatusForbidden, gin.H{
					"code":    enum.ShareTokenInsufficientPermissionsCode,
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
				})
				ctx.Abort()
				return
			}

			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    enum.ShareTokenInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.tokenHasExpired"}),
			})
			ctx.Abort()
			return
		}

		// 分享项目的访问令牌
		if token[:1] == "p" {
			if project.(*models.Projects).ID != sr.ProjectID {
				ctx.JSON(http.StatusForbidden, gin.H{
					"code":    enum.ShareTokenInsufficientPermissionsCode,
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
				})
				ctx.Abort()
				return
			}
		}

		// 分享文档的访问令牌
		if token[:1] == "d" {
			collectionIDStr := ctx.Param("collection-id")
			if collectionIDStr == "" {
				if project.(*models.Projects).ID != sr.ProjectID {
					ctx.JSON(http.StatusForbidden, gin.H{
						"code":    enum.ShareTokenInsufficientPermissionsCode,
						"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
					})
					ctx.Abort()
					return
				}
			} else {
				collectionID, err := strconv.Atoi(collectionIDStr)
				if err != nil {
					ctx.JSON(http.StatusNotFound, gin.H{
						"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.NotFound"}),
					})
					ctx.Abort()
					return
				}

				collection, err := models.NewCollections(uint(collectionID))
				if err != nil {
					ctx.JSON(http.StatusNotFound, gin.H{
						"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.NotFound"}),
					})
					ctx.Abort()
					return
				}

				if collection.ID != sr.CollectionID {
					ctx.JSON(http.StatusForbidden, gin.H{
						"code":    enum.ShareTokenInsufficientPermissionsCode,
						"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
					})
					ctx.Abort()
					return
				}
			}
			return
		}

		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		ctx.Abort()
	}
}
