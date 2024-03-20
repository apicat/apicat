package access

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/share"
	"github.com/apicat/apicat/backend/model/team"
	"github.com/apicat/apicat/backend/route/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func AllowGuest() func(*gin.Context) {
	return func(ctx *gin.Context) {
		projectID := ctx.Param("projectID")
		if projectID == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx)})
			return
		}
		p := &project.Project{ID: projectID}
		exist, err := p.Get(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
			return
		}
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": i18n.NewTran("project.DoesNotExist").Translate(ctx)})
			return
		}

		// 检查项目所属的团队还是否存在
		t := &team.Team{ID: p.TeamID}
		exist, err = t.Get(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
			return
		}
		if !exist {
			// 团队都不存在了，项目也就不存在了
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": i18n.NewTran("project.DoesNotExist").Translate(ctx)})
			return
		}

		setSelfProject(ctx, p)

		jumpTo := "login"
		if jwt.GetUser(ctx) != nil {
			// 用户登录了
			jumpTo = "home"
			tm := &team.TeamMember{UserID: jwt.GetUser(ctx).ID}
			// 查用户是否在这个项目的团队里
			exist, err = t.HasMember(ctx, tm)
			if err == nil && exist && tm.Status == team.MemberStatusActive {
				// 用户在团队，且是可用状态
				setSelfTeam(ctx, t)
				setSelfTeamMember(ctx, tm)

				// 查用户是否在这个项目里
				pm := &project.ProjectMember{ProjectID: p.ID, MemberID: tm.ID}
				exist, err = pm.Get(ctx)
				// 项目成员存在
				if err == nil && exist {
					setSelfProjectMember(ctx, pm)
					ctx.Next()
					return
				}
			}
		}

		if p.Visibility == project.VisibilityPublic {
			// 项目为公开项目
			ctx.Next()
			return
		}

		// 获取分享令牌
		shareCode := ctx.Query("shareCode")
		// token为32位md5加密后的字符串加前缀，前缀为p代表项目分享，d代表文档分享
		if len(shareCode) < 33 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx),
				"action":  jumpTo,
			})
			return
		}

		if shareCode[:1] != "p" && shareCode[:1] != "d" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx),
				"action":  jumpTo,
			})
			return
		}

		// 检查令牌合法性
		stt := &share.ShareTmpToken{ShareToken: fmt.Sprintf("%x", md5.Sum([]byte(shareCode)))}
		exist, err = stt.Get(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
			return
		}
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx),
				"action":  jumpTo,
			})
			return
		}

		if p.ID != stt.ProjectID {
			// 项目分享令牌里对应的项目id和请求参数项目id不一致
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx),
				"action":  jumpTo,
			})
			return
		}

		if shareCode[:1] == "d" {
			collectionIDStr := ctx.Param("collectionID")
			if collectionIDStr != "" && collectionIDStr != fmt.Sprintf("%d", stt.CollectionID) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx),
					"action":  jumpTo,
				})
				return
			}
		}
	}
}
