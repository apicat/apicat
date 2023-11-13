package project

import (
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/module/translator"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProjectGroupList(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	pg, _ := project.NewProjectGroups()
	pg.UserID = currentUser.(*user.Users).ID
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

	var data proto.ProjectGroupCreateData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	count, _ := project.GetProjectGroupCountByName(currentUser.(*user.Users).ID, data.Name)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.NameExists"}),
		})
		return
	}

	displayOrder, _ := project.GetProjectGroupDisplayOrder(currentUser.(*user.Users).ID)

	pg, _ := project.NewProjectGroups()
	pg.UserID = currentUser.(*user.Users).ID
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
		data    proto.ProjectGroupRenameData
		uriData proto.ProjectGroupRenameUriData
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

	pg, err := project.NewProjectGroups(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.NotFound"}),
		})
		return
	}

	if pg.UserID != currentUser.(*user.Users).ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.UpdateFailed"}),
		})
		return
	}

	count, _ := project.GetProjectGroupCountExcludeTheID(currentUser.(*user.Users).ID, data.Name, pg.ID)
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

	var uriData proto.ProjectGroupRenameUriData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pg, err := project.NewProjectGroups(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectGroups.NotFound"}),
		})
		return
	}

	if pg.UserID != currentUser.(*user.Users).ID {
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
		data         proto.ProjectGroupOrderData
		displayOrder int = 1
	)
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, v := range data.IDs {
		pg, err := project.NewProjectGroups(v)
		if err != nil {
			continue
		}

		if pg.UserID != currentUser.(*user.Users).ID {
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
