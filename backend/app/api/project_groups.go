package api

import (
	"net/http"

	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/apicat/apicat/backend/models"
	"github.com/gin-gonic/gin"
)

type ProjectGroupCreateData struct {
	Name string `json:"name" binding:"required,lte=255"`
}

type ProjectGroupRenameUriData struct {
	ID uint `uri:"group_id" binding:"required,gt=0"`
}

type ProjectGroupRenameData struct {
	Name string `json:"name" binding:"required,lte=255"`
}

type ProjectGroupDeleteUriData struct {
	ID uint `uri:"group_id" binding:"required,gt=0"`
}

type ProjectGroupOrderData struct {
	IDs []uint `json:"ids" binding:"required,dive,gte=0"`
}

func ProjectGroupList(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	pg, _ := models.NewProjectGroups()
	pg.UserID = currentUser.(*models.Users).ID
	projectGroups, err := pg.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.QueryFailed"}),
		})
		return
	}

	res := []gin.H{}
	if len(projectGroups) == 0 {
		ctx.JSON(http.StatusOK, res)
		return
	}

	for _, v := range projectGroups {
		res = append(res, gin.H{
			"id":   v.ID,
			"name": v.Name,
		})
	}

	ctx.JSON(http.StatusOK, res)
}

func ProjectGroupCreate(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var data ProjectGroupCreateData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	count, _ := models.GetProjectGroupCountByName(currentUser.(*models.Users).ID, data.Name)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.NameExists"}),
		})
		return
	}

	displayOrder, _ := models.GetProjectGroupDisplayOrder(currentUser.(*models.Users).ID)

	pg, _ := models.NewProjectGroups()
	pg.UserID = currentUser.(*models.Users).ID
	pg.Name = data.Name
	pg.DisplayOrder = displayOrder + 1
	if err := pg.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.CreateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":   pg.ID,
		"name": pg.Name,
	})
}

func ProjectGroupRename(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		data    ProjectGroupRenameData
		uriData ProjectGroupRenameUriData
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pg, err := models.NewProjectGroups(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.NotFound"}),
		})
		return
	}

	if pg.UserID != currentUser.(*models.Users).ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.UpdateFailed"}),
		})
		return
	}

	count, _ := models.GetProjectGroupCountExcludeTheID(currentUser.(*models.Users).ID, data.Name, pg.ID)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.NameExists"}),
		})
		return
	}

	pg.Name = data.Name
	if err := pg.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.UpdateFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func ProjectGroupDelete(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var uriData ProjectGroupRenameUriData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pg, err := models.NewProjectGroups(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.NotFound"}),
		})
		return
	}

	if pg.UserID != currentUser.(*models.Users).ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.DeleteFailed"}),
		})
		return
	}

	if err := pg.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.DeleteFailed"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func ProjectGroupOrder(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		data         ProjectGroupOrderData
		displayOrder int = 1
	)
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, v := range data.IDs {
		pg, err := models.NewProjectGroups(v)
		if err != nil {
			continue
		}

		if pg.UserID != currentUser.(*models.Users).ID {
			continue
		}

		pg.DisplayOrder = displayOrder
		if err := pg.Update(); err != nil {
			continue
		}

		displayOrder++
	}

	ctx.Status(http.StatusCreated)
}
