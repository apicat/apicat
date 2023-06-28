package api

import (
	"net/http"

	"github.com/apicat/apicat/common/encrypt"
	"github.com/apicat/apicat/common/random"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type ProjectSharingSwitchData struct {
	Share string `json:"share" binding:"required,oneof=open close"`
}

type ProjectShareSecretkeyCheckData struct {
	SecretKey string `json:"secret_key" binding:"required,lte=255"`
}

func ProjectSharingSwitch(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	var (
		project   *models.Projects
		data      ProjectSharingSwitchData
		secretKey string
		link      string
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
		secretKey = random.GenerateRandomString(4)
		link = ctx.Request.Host + "/project/" + project.PublicId

		project.SharePassword = secretKey
		if err := project.Save(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ModifySharingStatusFail"}),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"link":       link,
			"secret_key": secretKey,
		})
	} else {
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

func ProjectShareResetSecretKey(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

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

	secretKey = random.GenerateRandomString(4)

	project.SharePassword = secretKey
	if err := project.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ResetKeyFail"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func ProjectShareSecretkeyCheck(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	var (
		project *models.Projects
		data    ProjectShareSecretkeyCheckData
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

	if data.SecretKey != project.SharePassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.AccessPasswordError"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"token": encrypt.Get16MD5Encode(data.SecretKey),
	})
}
