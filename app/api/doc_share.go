package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/apicat/apicat/common/encrypt"
	"github.com/apicat/apicat/common/random"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/enum"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
)

type DocShareStatusData struct {
	PublicCollectionID string `uri:"public_collection_id" binding:"required,lte=255"`
}

type DocShareSecretkeyCheckUriData struct {
	ProjectID    string `uri:"project-id" binding:"required,lte=255"`
	CollectionID uint   `uri:"collection-id" binding:"required,gte=0"`
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

func DocShareDetails(ctx *gin.Context) {
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
		uriData            CollectionDataGetData
		visibility         string
		collectionPublicID string
		secretKey          string
	)

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

	if currentProject.(*models.Projects).Visibility == 0 {
		visibility = "private"
		collectionPublicID = collection.PublicId
		secretKey = collection.SharePassword
	} else {
		visibility = "public"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"visibility":           visibility,
		"collection_public_id": collectionPublicID,
		"secret_key":           secretKey,
	})
}

func DocShareSwitch(ctx *gin.Context) {
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
		uriData DocShareSecretkeyCheckUriData
		data    ProjectSharingSwitchData
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

		if collection.SharePassword == "" {
			collection.SharePassword = random.GenerateRandomString(4)
		}

		if err := collection.Update(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.ModifySharingStatusFail"}),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"collection_public_id": collection.PublicId,
			"secret_key":           collection.SharePassword,
		})
	} else {
		stt := models.NewShareTmpTokens()
		stt.CollectionID = collection.ID
		if err := stt.DeleteByCollectionID(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.ResetKeyFail"}),
			})
			return
		}

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

func DocShareReset(ctx *gin.Context) {
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

	stt := models.NewShareTmpTokens()
	stt.CollectionID = collection.ID
	if err := stt.DeleteByCollectionID(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.ResetKeyFail"}),
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

func DocShareCheck(ctx *gin.Context) {
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
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.AccessPasswordError"}),
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
