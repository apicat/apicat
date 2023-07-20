package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type CollectionHistoryListData struct {
	ID       uint                        `json:"id"`
	Title    string                      `json:"title"`
	Type     string                      `json:"type"`
	SubNodes []CollectionHistoryListData `json:"sub_nodes,omitempty"`
}

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

	r1 := map[string][]CollectionHistoryListData{}
	for _, v := range histories {
		month := v.CreatedAt.Format("2006-01")

		date := v.CreatedAt.Format("01月02日 15:04")
		var username string
		if _, ok := userDict[v.CreatedBy]; ok {
			username = userDict[v.CreatedBy].Username
		}

		r1[month] = append(r1[month], CollectionHistoryListData{
			ID:    v.ID,
			Title: fmt.Sprintf("%s(%s)", date, username),
			Type:  v.Type,
		})
	}

	r2 := []CollectionHistoryListData{}
	for k, v := range r1 {
		r2 = append(r2, CollectionHistoryListData{
			ID:       0,
			Title:    fmt.Sprintf("%s月", strings.Replace(k, "-", "年", -1)),
			Type:     "category",
			SubNodes: v,
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
