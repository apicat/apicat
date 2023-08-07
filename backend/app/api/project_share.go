package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/apicat/apicat/backend/common/encrypt"
	"github.com/apicat/apicat/backend/common/random"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/apicat/apicat/backend/models"
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

	if currentProject.(*models.Projects).Visibility == 0 {
		visibility = "private"
	} else {
		visibility = "public"
	}

	if currentProject.(*models.Projects).SharePassword == "" {
		hasShared = false
	} else {
		hasShared = true
	}

	if currentUserExists {
		member, _ := models.NewProjectMembers()
		member.UserID = currentUser.(*models.Users).ID
		member.ProjectID = currentProject.(*models.Projects).ID

		if err := member.GetByUserIDAndProjectID(); err == nil {
			authority = member.Authority
		}
	} else {
		authority = "none"
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

	if currentProject.(*models.Projects).Visibility == 0 {
		if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
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

	if currentProject.(*models.Projects).Visibility == 0 {
		visibility = "private"
	} else {
		visibility = "public"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"authority":  currentProjectMember.(*models.ProjectMembers).Authority,
		"visibility": visibility,
		"secret_key": currentProject.(*models.Projects).SharePassword,
	})
}

func ProjectSharingSwitch(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		project *models.Projects
		data    ProjectSharingSwitchData
	)

	project = currentProject.(*models.Projects)
	if project.Visibility != 0 {
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
		if project.SharePassword == "" {
			project.SharePassword = random.GenerateRandomString(4)

			if err := project.Save(); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ModifySharingStatusFail"}),
				})
				return
			}
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"project_public_id": project.PublicId,
			"secret_key":        project.SharePassword,
		})
	} else {
		stt := models.NewShareTmpTokens()
		stt.ProjectID = project.ID
		if err := stt.DeleteByProjectIDAndCollectionID(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ModifySharingStatusFail"}),
			})
			return
		}

		project.SharePassword = ""
		if err := project.Save(); err != nil {
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
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var (
		project   *models.Projects
		secretKey string
	)

	project = currentProject.(*models.Projects)
	if project.Visibility != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.PublicProject"}),
		})
		return
	}

	stt := models.NewShareTmpTokens()
	stt.ProjectID = project.ID
	if err := stt.DeleteByProjectIDAndCollectionID(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ResetKeyFail"}),
		})
		return
	}

	secretKey = random.GenerateRandomString(4)

	project.SharePassword = secretKey
	if err := project.Save(); err != nil {
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
		project *models.Projects
		data    ProjectShareSecretkeyCheckData
		err     error
	)

	project = currentProject.(*models.Projects)
	if err = translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.SecretKey != project.SharePassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.AccessPasswordError"}),
		})
		return
	}

	token := "p" + encrypt.GetMD5Encode(data.SecretKey+fmt.Sprint(time.Now().UnixNano()))

	stt := models.NewShareTmpTokens()
	stt.ShareToken = encrypt.GetMD5Encode(token)
	stt.Expiration = time.Now().Add(time.Hour * 24)
	stt.ProjectID = project.ID
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
