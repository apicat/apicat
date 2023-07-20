package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

func CollectionHistoryList(ctx *gin.Context) {
	uriData := CollectionDataGetData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	u, _ := models.NewUsers()
	users, err := u.List(0, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.QueryFailed"}),
		})
		return
	}

	userDict := map[uint]*models.Users{}
	for _, user := range users {
		userDict[user.ID] = &user
	}

	ch, _ := models.NewCollectionHistories()
	histories, err := ch.List(uriData.CollectionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.QueryFailed"}),
		})
		return
	}

	r1 := map[string][]any{}
	for _, v := range histories {
		month := v.CreatedAt.Format("2006-01")

		date := v.CreatedAt.Format("01月02日 15:04")
		var username string
		if _, ok := userDict[v.CreatedBy]; ok {
			username = userDict[v.CreatedBy].Username
		}

		r1[month] = append(r1[month], map[string]any{
			"id":    v.ID,
			"title": fmt.Sprintf("%s(%s)", date, username),
		})
	}

	r2 := []any{}
	for k, v := range r1 {
		r2 = append(r2, map[string]any{
			"id":        0,
			"title":     fmt.Sprintf("%s月", strings.Replace(k, "-", "年", -1)),
			"type":      0,
			"sub_nodes": v,
		})
	}

	ctx.JSON(http.StatusOK, r2)
}

func CollectionHistoryDetails(ctx *gin.Context) {
}

func CollectionHistoryDiff(ctx *gin.Context) {
}

func CollectionHistoryRestore(ctx *gin.Context) {
}
