package check

import (
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/share"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/module/encrypt"
	"github.com/apicat/apicat/backend/module/translator"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckProjectMemberHalfLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p, exists := ctx.Get("CurrentProject")
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    proto.Redirect404Page,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
			})
			ctx.Abort()
			return
		}

		u, exists := ctx.Get("CurrentUser")
		if exists {
			member, _ := project.NewProjectMembers()
			member.UserID = u.(*user.Users).ID
			member.ProjectID = p.(*project.Projects).ID

			if err := member.GetByUserIDAndProjectID(); err == nil {
				ctx.Set("CurrentProjectMember", member)
				return
			}
		}

		// 判断是否为公开项目
		if p.(*project.Projects).Visibility == 1 {
			return
		}

		token := ctx.Query("token")
		if token == "" || len(token) < 1 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    proto.InvalidOrIncorrectAccessToken,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
			})
			ctx.Abort()
			return
		}

		// 校验访问令牌是否存在
		stt := share.NewShareTmpTokens()
		stt.ShareToken = encrypt.GetMD5Encode(token)
		if err := stt.GetByShareToken(); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    proto.InvalidOrIncorrectAccessToken,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
			})
			ctx.Abort()
			return
		}

		// 校验访问令牌是否过期
		now := time.Now()
		if stt.Expiration.Before(now) {
			if err := stt.Delete(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.DeleteFailed"}),
				})
				ctx.Abort()
				return
			}

			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    proto.InvalidOrIncorrectAccessToken,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.tokenHasExpired"}),
			})
			ctx.Abort()
			return
		}

		// 分享项目的访问令牌
		if token[:1] == "p" {
			if p.(*project.Projects).ID != stt.ProjectID {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"code":    proto.InvalidOrIncorrectAccessToken,
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
				})
				ctx.Abort()
				return
			}
			return
		}

		// 分享文档的访问令牌
		if token[:1] == "d" {
			collectionIDStr := ctx.Param("collection-id")
			if collectionIDStr == "" {
				if p.(*project.Projects).ID != stt.ProjectID {
					ctx.JSON(http.StatusUnauthorized, gin.H{
						"code":    proto.InvalidOrIncorrectAccessToken,
						"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
					})
					ctx.Abort()
					return
				}
			} else {
				collectionID, err := strconv.Atoi(collectionIDStr)
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.NotFound"}),
					})
					ctx.Abort()
					return
				}

				c, err := collection.NewCollections(uint(collectionID))
				if err != nil {
					ctx.JSON(http.StatusNotFound, gin.H{
						"code":    proto.Redirect404Page,
						"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.NotFound"}),
					})
					ctx.Abort()
					return
				}

				if c.ID != stt.CollectionID {
					ctx.JSON(http.StatusUnauthorized, gin.H{
						"code":    proto.InvalidOrIncorrectAccessToken,
						"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.InvalidToken"}),
					})
					ctx.Abort()
					return
				}
			}
			return
		}

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    proto.InvalidOrIncorrectAccessToken,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		ctx.Abort()
		return
	}
}
