package api

import (
	"net/http"

	"github.com/apicat/apicat/commom/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type TrashsRecoverQuery struct {
	CollectionID []uint `form:"collection-id" binding:"required,dive,gte=0"`
}

type TrashsRecoverBody struct {
	Category uint `json:"category" binding:"gte=0"`
}

func TrashsList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	collection, _ := models.NewCollections()
	collection.ProjectId = project.ID
	collections, err := collection.TrashList()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	result := []map[string]interface{}{}
	for _, v := range collections {
		result = append(result, map[string]interface{}{
			"id":         v.ID,
			"title":      v.Title,
			"type":       v.Type,
			"deleted_at": v.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func TrashsRecover(ctx *gin.Context) {
	trashsRecoverQuery := TrashsRecoverQuery{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&trashsRecoverQuery)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	trashsRecoverBody := TrashsRecoverBody{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&trashsRecoverBody)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	allOK := true

	for _, v := range trashsRecoverQuery.CollectionID {
		collection, _ := models.NewCollections()
		collection.ID = v
		collection.ProjectId = project.ID
		if err := collection.GetUnscopedCollections(); err != nil {
			allOK = false
			continue
		}
		collection.ParentId = trashsRecoverBody.Category
		collection.Restore()
	}

	if !allOK {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Trashs.PartialRecoveryFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
