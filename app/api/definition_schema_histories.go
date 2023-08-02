package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/enum"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type SchemaUriData struct {
	ProjectID string `uri:"project-id" binding:"required,gt=0"`
	SchemaID  uint   `uri:"schemas-id" binding:"required,gt=0"`
}

type SchemaHistoryUriData struct {
	ProjectID string `uri:"project-id" binding:"required,gt=0"`
	SchemaID  uint   `uri:"schemas-id" binding:"required,gt=0"`
	HistoryID uint   `uri:"history-id" binding:"required,gt=0"`
}

type SchemaHistoryListData struct {
	ID       uint                    `json:"id"`
	Name     string                  `json:"name"`
	Type     string                  `json:"type"`
	SubNodes []SchemaHistoryListData `json:"sub_nodes,omitempty"`
}

type SchemaHistoryDetailsData struct {
	ID            uint   `json:"id"`
	SchemaID      uint   `json:"schema_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Schema        string `json:"schema"`
	CreatedTime   string `json:"created_time"`
	LastUpdatedBy string `json:"last_updated_by"`
}

type SchemaHistoryDiffData struct {
	HistoryID1 uint `form:"history_id1"`
	HistoryID2 uint `form:"history_id2"`
}

func DefinitionSchemaHistoryList(ctx *gin.Context) {
	uriData := SchemaUriData{}
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

	dsh, _ := models.NewDefinitionSchemaHistories()
	histories, err := dsh.List(uriData.SchemaID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.QueryFailed"}),
		})
		return
	}

	r1 := map[string][]SchemaHistoryListData{}
	for _, v := range histories {
		month := v.CreatedAt.Format("2006-01")

		date := v.CreatedAt.Format("01月02日 15:04")
		var username string
		if _, ok := userDict[v.CreatedBy]; ok {
			username = userDict[v.CreatedBy].Username
		}

		r1[month] = append(r1[month], SchemaHistoryListData{
			ID:   v.ID,
			Name: fmt.Sprintf("%s(%s)", date, username),
			Type: v.Type,
		})
	}

	r2 := []SchemaHistoryListData{}
	for k, v := range r1 {
		r2 = append(r2, SchemaHistoryListData{
			ID:       0,
			Name:     fmt.Sprintf("%s月", strings.Replace(k, "-", "年", -1)),
			Type:     "category",
			SubNodes: v,
		})
	}

	ctx.JSON(http.StatusOK, r2)
}

func DefinitionSchemaHistoryDetails(ctx *gin.Context) {
	currentDefinitionSchema, _ := ctx.Get("CurrentDefinitionSchema")

	uriData := SchemaHistoryUriData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	dsh, err := models.NewDefinitionSchemaHistories(uriData.HistoryID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if currentDefinitionSchema.(*models.DefinitionSchemas).ID != dsh.SchemaID {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	u, err := models.NewUsers(dsh.CreatedBy)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, SchemaHistoryDetailsData{
		ID:            dsh.ID,
		SchemaID:      dsh.SchemaID,
		Name:          dsh.Name,
		Description:   dsh.Description,
		Schema:        dsh.Schema,
		CreatedTime:   dsh.CreatedAt.Format("2006-01-02 15:04"),
		LastUpdatedBy: u.Username,
	})
}

func DefinitionSchemaHistoryDiff(ctx *gin.Context) {
	currentDefinitionSchema, _ := ctx.Get("CurrentDefinitionSchema")

	var data CollectionHistoryDiffData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	dsh1, err := models.NewDefinitionSchemaHistories(data.HistoryID1)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}
	dsh2, err := models.NewDefinitionSchemaHistories(data.HistoryID2)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if dsh1.SchemaID != currentDefinitionSchema.(*models.DefinitionSchemas).ID || dsh2.SchemaID != currentDefinitionSchema.(*models.DefinitionSchemas).ID {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	u1, err := models.NewUsers(dsh1.CreatedBy)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
		})
		return
	}
	u2, err := models.NewUsers(dsh2.CreatedBy)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
		})
		return
	}

	res := map[string]SchemaHistoryDetailsData{}
	res["schema1"] = SchemaHistoryDetailsData{
		ID:            dsh1.ID,
		SchemaID:      dsh1.SchemaID,
		Name:          dsh1.Name,
		Description:   dsh1.Description,
		Schema:        dsh1.Schema,
		CreatedTime:   dsh1.CreatedAt.Format("2006-01-02 15:04"),
		LastUpdatedBy: u1.Username,
	}

	res["schema2"] = SchemaHistoryDetailsData{
		ID:            dsh2.ID,
		SchemaID:      dsh2.SchemaID,
		Name:          dsh2.Name,
		Description:   dsh2.Description,
		Schema:        dsh2.Schema,
		CreatedTime:   dsh2.CreatedAt.Format("2006-01-02 15:04"),
		LastUpdatedBy: u2.Username,
	}

	ctx.JSON(http.StatusOK, res)
}

func DefinitionSchemaHistoryRestore(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")
	currentDefinitionSchema, _ := ctx.Get("CurrentDefinitionSchema")

	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	uriData := SchemaHistoryUriData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	dsh, err := models.NewDefinitionSchemaHistories(uriData.HistoryID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if currentDefinitionSchema.(*models.DefinitionSchemas).ID != dsh.SchemaID {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if err := dsh.Restore(currentDefinitionSchema.(*models.DefinitionSchemas), currentUser.(*models.Users).ID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.RestoreFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
