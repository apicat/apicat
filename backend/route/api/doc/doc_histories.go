package doc

import (
	"fmt"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/user"
	collection2 "github.com/apicat/apicat/backend/route/api/collection"
	"net/http"
	"strings"

	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/gin-gonic/gin"
)

type CollectionHistoryListData struct {
	ID       uint                        `json:"id"`
	Title    string                      `json:"title"`
	Type     string                      `json:"type"`
	SubNodes []CollectionHistoryListData `json:"sub_nodes,omitempty"`
}

type CollectionHistoryUriData struct {
	ProjectID    string `uri:"project-id" binding:"required,gt=0"`
	CollectionID uint   `uri:"collection-id" binding:"required,gt=0"`
	HistoryID    uint   `uri:"history-id" binding:"required,gt=0"`
}

type CollectionHistoryDiffData struct {
	HistoryID1 uint `form:"history_id1"`
	HistoryID2 uint `form:"history_id2"`
}

type CollectionHistoryDetailsData struct {
	ID            uint   `json:"id"`
	CollectionID  uint   `json:"collection_id"`
	Content       string `json:"content"`
	CreatedTime   string `json:"created_time"`
	LastUpdatedBy string `json:"last_updated_by"`
	Title         string `json:"title"`
}

func CollectionHistoryList(ctx *gin.Context) {
	uriData := collection2.CollectionDataGetData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	u, _ := user.NewUsers()
	users, err := u.List(0, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.QueryFailed"}),
		})
		return
	}

	userDict := map[uint]user.Users{}
	for _, u := range users {
		userDict[u.ID] = u
	}

	ch, _ := collection.NewCollectionHistories()
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
	currentCollection, _ := ctx.Get("CurrentCollection")

	uriData := CollectionHistoryUriData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ch, err := collection.NewCollectionHistories(uriData.HistoryID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Redirect404Page,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if currentCollection.(*collection.Collections).ID != ch.CollectionId {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Redirect404Page,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	u, err := user.NewUsers(ch.CreatedBy)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, CollectionHistoryDetailsData{
		ID:            ch.ID,
		CollectionID:  ch.CollectionId,
		Title:         ch.Title,
		Content:       ch.Content,
		CreatedTime:   ch.CreatedAt.Format("2006-01-02 15:04"),
		LastUpdatedBy: u.Username,
	})
}

func CollectionHistoryDiff(ctx *gin.Context) {
	currentCollection, _ := ctx.Get("CurrentCollection")

	var data CollectionHistoryDiffData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	res := map[string]CollectionHistoryDetailsData{}

	ch1, err := collection.NewCollectionHistories(data.HistoryID1)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}
	if ch1.CollectionId != currentCollection.(*collection.Collections).ID {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}
	u1, err := user.NewUsers(ch1.CreatedBy)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
		})
		return
	}

	res["doc1"] = CollectionHistoryDetailsData{
		ID:            ch1.ID,
		CollectionID:  ch1.CollectionId,
		Title:         ch1.Title,
		Content:       ch1.Content,
		CreatedTime:   ch1.CreatedAt.Format("2006-01-02 15:04"),
		LastUpdatedBy: u1.Username,
	}

	if data.HistoryID2 == 0 {
		u2, err := user.NewUsers(currentCollection.(*collection.Collections).UpdatedBy)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    enum.Display404ErrorMessage,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
			})
			return
		}

		res["doc2"] = CollectionHistoryDetailsData{
			ID:            0,
			CollectionID:  currentCollection.(*collection.Collections).ID,
			Title:         currentCollection.(*collection.Collections).Title,
			Content:       currentCollection.(*collection.Collections).Content,
			CreatedTime:   currentCollection.(*collection.Collections).CreatedAt.Format("2006-01-02 15:04"),
			LastUpdatedBy: u2.Username,
		}

		ctx.JSON(http.StatusOK, res)
		return
	}

	ch2, err := collection.NewCollectionHistories(data.HistoryID2)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}
	if ch2.CollectionId != currentCollection.(*collection.Collections).ID {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}
	u2, err := user.NewUsers(ch2.CreatedBy)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
		})
		return
	}

	res["doc2"] = CollectionHistoryDetailsData{
		ID:            ch2.ID,
		CollectionID:  ch2.CollectionId,
		Title:         ch2.Title,
		Content:       ch2.Content,
		CreatedTime:   ch2.CreatedAt.Format("2006-01-02 15:04"),
		LastUpdatedBy: u2.Username,
	}

	ctx.JSON(http.StatusOK, res)
}

func CollectionHistoryRestore(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")
	currentCollection, _ := ctx.Get("CurrentCollection")

	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	uriData := CollectionHistoryUriData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ch, err := collection.NewCollectionHistories(uriData.HistoryID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if currentCollection.(*collection.Collections).ID != ch.CollectionId {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if err := ch.Restore(currentCollection.(*collection.Collections), currentUser.(*user.Users).ID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.RestoreFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
