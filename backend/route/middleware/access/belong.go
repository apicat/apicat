package access

import (
	"net/http"
	"strconv"

	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/iteration"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/team"
	"github.com/apicat/apicat/backend/route/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func BelongToTeam() func(*gin.Context) {
	return func(ctx *gin.Context) {
		var t *team.Team
		if ctx.Param("teamID") != "" {
			t = &team.Team{ID: ctx.Param("teamID")}
		} else if ctx.Param("projectID") != "" {
			p := &project.Project{ID: ctx.Param("projectID")}
			exist, err := p.Get(ctx)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
				return
			}
			if !exist {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": i18n.NewTran("project.DoesNotExist").Translate(ctx)})
				return
			}
			t = &team.Team{ID: p.TeamID}
		} else if ctx.Param("iterationID") != "" {
			i := &iteration.Iteration{ID: ctx.Param("iterationID")}
			exist, err := i.Get(ctx)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
				return
			}
			if !exist {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": i18n.NewTran("iteration.DoesNotExist").Translate(ctx)})
				return
			}
			t = &team.Team{ID: i.TeamID}
		} else if ctx.Param("groupID") != "" {
			groupID, err := strconv.ParseUint(ctx.Param("groupID"), 10, 64)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx)})
				return
			}
			pg := &project.ProjectGroup{ID: uint(groupID)}
			exist, err := pg.Get(ctx)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
				return
			}
			if !exist {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": i18n.NewTran("projectGroup.DoesNotExist").Translate(ctx)})
				return
			}
			tm, err := team.GetMember(ctx, pg.MemberID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
				return
			}
			t = &team.Team{ID: tm.TeamID}
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx)})
			return
		}

		exist, err := t.Get(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
			return
		}
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": i18n.NewTran("team.DoesNotExist").Translate(ctx)})
			return
		}

		u := jwt.GetUser(ctx)
		if u == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": i18n.NewTran("user.LoginStatusExpired").Translate(ctx),
				"action":  "login",
			})
			return
		}

		tm := &team.TeamMember{UserID: u.ID}
		exist, err = t.HasMember(ctx, tm)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
			return
		}
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": i18n.NewTran("common.PermissionDenied").Translate(ctx)})
			return
		}

		if tm.Status == team.MemberStatusDeactive {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": i18n.NewTran("common.PermissionDenied").Translate(ctx)})
			return
		}

		setSelfTeam(ctx, t)
		setSelfTeamMember(ctx, tm)
	}
}

func BelongToProject() func(*gin.Context) {
	return func(ctx *gin.Context) {
		var p *project.Project

		if ctx.Param("projectID") != "" {
			p = &project.Project{ID: ctx.Param("projectID")}
		} else if ctx.Param("iterationID") != "" {
			i := &iteration.Iteration{ID: ctx.Param("iterationID")}
			exist, err := i.Get(ctx)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
				return
			}
			if !exist {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": i18n.NewTran("iteration.DoesNotExist").Translate(ctx)})
				return
			}
			p = &project.Project{ID: i.ProjectID}
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx)})
			return
		}

		exist, err := p.Get(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
			return
		}
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": i18n.NewTran("project.DoesNotExist").Translate(ctx)})
			return
		}

		tm := GetSelfTeamMember(ctx)
		pm := &project.ProjectMember{ProjectID: p.ID, MemberID: tm.ID}
		exist, err = pm.Get(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": i18n.NewTran("common.GenericError").Translate(ctx)})
			return
		}
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": i18n.NewTran("common.PermissionDenied").Translate(ctx)})
			return
		}

		setSelfProject(ctx, p)
		setSelfProjectMember(ctx, pm)
	}
}
