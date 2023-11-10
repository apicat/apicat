package doc

import (
	"fmt"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/share"
	"github.com/apicat/apicat/backend/module/encrypt"
	"github.com/apicat/apicat/backend/module/random"
	"github.com/apicat/apicat/backend/module/translator"
	collection2 "github.com/apicat/apicat/backend/route/api/collection"
	project2 "github.com/apicat/apicat/backend/route/api/project"
	"net/http"
	"time"

	"github.com/apicat/apicat/backend/enum"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
)

type DocShareStatusData struct {
	PublicCollectionID string `uri:"public_collection_id" binding:"required,lte=255"`
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

	c, _ := collection.NewCollections()
	c.PublicId = data.PublicCollectionID
	if err := c.GetByPublicId(); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.NotFound"}),
		})
		return
	}

	p, err := project.NewProjects(c.ProjectId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"project_id":    p.PublicId,
		"collection_id": c.ID,
	})
}

func DocShareDetails(ctx *gin.Context) {
	currentCollection, _ := ctx.Get("CurrentCollection")
	c := currentCollection.(*collection.Collections)

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
		uriData            collection2.CollectionDataGetData
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

	if currentProject.(*project.Projects).Visibility == 0 {
		visibility = "private"
		collectionPublicID = c.PublicId
		secretKey = c.SharePassword
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
	currentCollection, _ := ctx.Get("CurrentCollection")
	c := currentCollection.(*collection.Collections)

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
		data project2.ProjectSharingSwitchData
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
		if c.PublicId == "" {
			c.PublicId = shortuuid.New()
		}

		if c.SharePassword == "" {
			c.SharePassword = random.GenerateRandomString(4)
		}

		if err := c.Update(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.ModifySharingStatusFail"}),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"collection_public_id": c.PublicId,
			"secret_key":           c.SharePassword,
		})
	} else {
		stt := share.NewShareTmpTokens()
		stt.CollectionID = c.ID
		if err := stt.DeleteByCollectionID(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.ResetKeyFail"}),
			})
			return
		}

		c.SharePassword = ""
		if err := c.Update(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.ModifySharingStatusFail"}),
			})
			return
		}

		ctx.Status(http.StatusCreated)
	}
}

func DocShareReset(ctx *gin.Context) {
	currentCollection, _ := ctx.Get("CurrentCollection")
	c := currentCollection.(*collection.Collections)

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
	stt.CollectionID = c.ID
	if err := stt.DeleteByCollectionID(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DocShare.ResetKeyFail"}),
		})
		return
	}

	secretKey = random.GenerateRandomString(4)
	c.SharePassword = secretKey

	if err := c.Update(); err != nil {
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
	currentCollection, _ := ctx.Get("CurrentCollection")
	c := currentCollection.(*collection.Collections)

	var (
		data project2.ProjectShareSecretkeyCheckData
		err  error
	)

	if err = translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.SecretKey != c.SharePassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Share.AccessPasswordError"}),
		})
		return
	}

	token := "d" + encrypt.GetMD5Encode(data.SecretKey+fmt.Sprint(time.Now().UnixNano()))

	stt := share.NewShareTmpTokens()
	stt.ShareToken = encrypt.GetMD5Encode(token)
	stt.Expiration = time.Now().Add(time.Hour * 24)
	stt.ProjectID = c.ProjectId
	stt.CollectionID = c.ID
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
