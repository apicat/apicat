package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/apicat/apicat/common/encrypt"
	"github.com/apicat/apicat/common/random"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
)

type ProjectSharingSwitchData struct {
	Share string `json:"share" binding:"required,oneof=open close"`
}

type ProjectShareSecretkeyCheckData struct {
	SecretKey string `json:"secret_key" binding:"required,lte=255"`
}

type DocShareStatusData struct {
	PublicCollectionID string `uri:"public_collection_id" binding:"required,lte=255"`
}

type DocShareSecretkeyCheckUriData struct {
	ProjectID    string `uri:"project-id" binding:"required,lte=255"`
	CollectionID uint   `uri:"collection-id" binding:"required,gte=0"`
}

type ShareTokenContentData struct {
	SecretKey  string `json:"secret_key"`
	Expiration int64  `json:"expiration"`
}

func ProjectSharingSwitch(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	var (
		project   *models.Projects
		data      ProjectSharingSwitchData
		secretKey string
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

		project.SharePassword = secretKey
		if err := project.Save(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.ModifySharingStatusFail"}),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"project_public_id": project.PublicId,
			"secret_key":        secretKey,
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
	if project.Visibility != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.PublicProject"}),
		})
		return
	}

	if err = translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
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

func DocShareStatus(ctx *gin.Context) {
	var (
		data DocShareStatusData
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	collection, _ := models.NewCollections()
	collection.PublicId = data.PublicCollectionID
	if err := collection.GetByPublicId(); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.NotFound"}),
		})
		return
	}

	project, err := models.NewProjects(collection.ProjectId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"project_id":    project.PublicId,
		"collection_id": collection.ID,
	})
}

func DocShareSecretkeyCheck(ctx *gin.Context) {
	var (
		uriData DocShareSecretkeyCheckUriData
		data    ProjectShareSecretkeyCheckData
		err     error
	)

	if err = translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err = translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	collection, err := models.NewCollections(uriData.CollectionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.NotFound"}),
		})
		return
	}

	if data.SecretKey != collection.SharePassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.AccessPasswordError"}),
		})
		return
	}

	token := "d" + encrypt.GetMD5Encode(data.SecretKey+fmt.Sprint(time.Now().UnixNano()))

	stt := models.NewShareTmpTokens()
	stt.ShareToken = encrypt.GetMD5Encode(token)
	stt.Expiration = time.Now().Add(time.Hour * 24)
	stt.ProjectID = collection.ProjectId
	stt.CollectionID = collection.ID
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

func DocSharingSwitch(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	var (
		project   *models.Projects
		uriData   DocShareSecretkeyCheckUriData
		data      ProjectSharingSwitchData
		secretKey string
	)

	project = currentProject.(*models.Projects)
	if project.Visibility != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.PublicProject"}),
		})
		return
	}

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	collection, err := models.NewCollections(uriData.CollectionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.NotFound"}),
		})
		return
	}

	if data.Share == "open" {
		if collection.PublicId == "" {
			collection.PublicId = shortuuid.New()
		}

		secretKey = random.GenerateRandomString(4)
		collection.SharePassword = secretKey
		if err := collection.Update(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.ModifySharingStatusFail"}),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"collection_public_id": collection.PublicId,
			"secret_key":           secretKey,
		})
	} else {
		collection.SharePassword = ""
		if err := collection.Update(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.ModifySharingStatusFail"}),
			})
			return
		}

		ctx.Status(http.StatusCreated)
	}
}

func DocShareResetSecretKey(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	var (
		project   *models.Projects
		uriData   DocShareSecretkeyCheckUriData
		secretKey string
	)

	project = currentProject.(*models.Projects)
	if project.Visibility != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectShare.PublicProject"}),
		})
		return
	}

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	collection, err := models.NewCollections(uriData.CollectionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.NotFound"}),
		})
		return
	}

	secretKey = random.GenerateRandomString(4)
	collection.SharePassword = secretKey

	if err := collection.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.ResetKeyFail"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"secret_key": secretKey,
	})
}
