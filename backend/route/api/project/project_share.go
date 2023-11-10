package project

import (
	"fmt"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/share"
	"github.com/apicat/apicat/backend/model/user"
	"net/http"
	"time"

	"github.com/apicat/apicat/backend/common/encrypt"
	"github.com/apicat/apicat/backend/common/random"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/gin-gonic/gin"
)

type ProjectSharingSwitchData struct {
	Share string `json:"share" binding:"required,oneof=open close"`
}

type ProjectShareSecretkeyCheckData struct {
	SecretKey string `json:"secret_key" binding:"required,lte=255"`
}

func ProjectShareStatus(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	currentUser, currentUserExists := ctx.Get("CurrentUser")

	var (
		authority  string
		visibility string
		hasShared  bool
	)

	if currentProject.(*project.Projects).Visibility == 0 {
		visibility = "private"
	} else {
		visibility = "public"
	}

	if currentProject.(*project.Projects).SharePassword == "" {
		hasShared = false
	} else {
		hasShared = true
	}

	authority = "none"
	if currentUserExists {
		member, _ := project.NewProjectMembers()
		member.UserID = currentUser.(*user.Users).ID
		member.ProjectID = currentProject.(*project.Projects).ID

		if err := member.GetByUserIDAndProjectID(); err == nil {
			authority = member.Authority
		}
	}

	if authority == "none" && visibility == "private" && !hasShared {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.Redirect403Page,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"authority":  authority,
		"visibility": visibility,
		"has_shared": hasShared,
	})
}

func ProjectShareDetails(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")

	if currentProject.(*project.Projects).Visibility == 0 {
		if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    enum.ProjectMemberInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
			})
			return
		}
	}

	var (
		visibility string
	)

	if currentProject.(*project.Projects).Visibility == 0 {
		visibility = "private"
	} else {
		visibility = "public"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"authority":  currentProjectMember.(*project.ProjectMembers).Authority,
		"visibility": visibility,
		"secret_key": currentProject.(*project.Projects).SharePassword,
	})
}

func ProjectSharingSwitch(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		p    *project.Projects
		data ProjectSharingSwitchData
	)

	p = currentProject.(*project.Projects)
	if p.Visibility != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.PublicProject"}),
		})
		return
	}

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.Share == "open" {
		if p.SharePassword == "" {
			p.SharePassword = random.GenerateRandomString(4)

			if err := p.Save(); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ModifySharingStatusFail"}),
				})
				return
			}
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"project_public_id": p.PublicId,
			"secret_key":        p.SharePassword,
		})
	} else {
		stt := share.NewShareTmpTokens()
		stt.ProjectID = p.ID
		if err := stt.DeleteByProjectIDAndCollectionID(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ModifySharingStatusFail"}),
			})
			return
		}

		p.SharePassword = ""
		if err := p.Save(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ModifySharingStatusFail"}),
			})
			return
		}

		ctx.Status(http.StatusCreated)
	}
}

func ProjectShareReset(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		p         *project.Projects
		secretKey string
	)

	p = currentProject.(*project.Projects)
	if p.Visibility != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.PublicProject"}),
		})
		return
	}

	stt := share.NewShareTmpTokens()
	stt.ProjectID = p.ID
	if err := stt.DeleteByProjectIDAndCollectionID(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ResetKeyFail"}),
		})
		return
	}

	secretKey = random.GenerateRandomString(4)

	p.SharePassword = secretKey
	if err := p.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ResetKeyFail"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"secret_key": secretKey,
	})
}

func ProjectShareSecretkeyCheck(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	var (
		p    *project.Projects
		data ProjectShareSecretkeyCheckData
		err  error
	)

	p = currentProject.(*project.Projects)
	if err = translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.SecretKey != p.SharePassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.AccessPasswordError"}),
		})
		return
	}

	token := "p" + encrypt.GetMD5Encode(data.SecretKey+fmt.Sprint(time.Now().UnixNano()))

	stt := share.NewShareTmpTokens()
	stt.ShareToken = encrypt.GetMD5Encode(token)
	stt.Expiration = time.Now().Add(time.Hour * 24)
	stt.ProjectID = p.ID
	if err := stt.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.VerifyKeyFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"token":      token,
		"expiration": stt.Expiration.Format("2006-01-02 15:04:05"),
	})
}
