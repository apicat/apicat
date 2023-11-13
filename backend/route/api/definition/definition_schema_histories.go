package definition

import (
	"encoding/json"
	"fmt"
	"github.com/apicat/apicat/backend/model/definition"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/module/translator"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DefinitionSchemaHistoryList(ctx *gin.Context) {
	uriData := proto.SchemaUriData{}
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

	dsh, _ := definition.NewDefinitionSchemaHistories()
	histories, err := dsh.List(uriData.SchemaID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.QueryFailed"}),
		})
		return
	}

	r1 := map[string][]proto.SchemaHistoryListData{}
	for _, v := range histories {
		month := v.CreatedAt.Format("2006-01")

		date := v.CreatedAt.Format("01月02日 15:04")
		var username string
		if _, ok := userDict[v.CreatedBy]; ok {
			username = userDict[v.CreatedBy].Username
		}

		r1[month] = append(r1[month], proto.SchemaHistoryListData{
			ID:   v.ID,
			Name: fmt.Sprintf("%s(%s)", date, username),
			Type: v.Type,
		})
	}

	r2 := []proto.SchemaHistoryListData{}
	for k, v := range r1 {
		r2 = append(r2, proto.SchemaHistoryListData{
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

	uriData := proto.SchemaHistoryUriData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	dsh, err := definition.NewDefinitionSchemaHistories(uriData.HistoryID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Redirect404Page,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if currentDefinitionSchema.(*definition.DefinitionSchemas).ID != dsh.SchemaID {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Redirect404Page,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	u, err := user.NewUsers(dsh.CreatedBy)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
		})
		return
	}

	schema := make(map[string]interface{})
	if err := json.Unmarshal([]byte(dsh.Schema), &schema); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, proto.SchemaHistoryDetailsData{
		ID:            dsh.ID,
		SchemaID:      dsh.SchemaID,
		Name:          dsh.Name,
		Description:   dsh.Description,
		Schema:        schema,
		CreatedTime:   dsh.CreatedAt.Format("2006-01-02 15:04"),
		LastUpdatedBy: u.Username,
	})
}

func DefinitionSchemaHistoryDiff(ctx *gin.Context) {
	currentDefinitionSchema, _ := ctx.Get("CurrentDefinitionSchema")

	var data proto.CollectionHistoryDiffData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	res := map[string]proto.SchemaHistoryDetailsData{}

	dsh1, err := definition.NewDefinitionSchemaHistories(data.HistoryID1)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}
	if dsh1.SchemaID != currentDefinitionSchema.(*definition.DefinitionSchemas).ID {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}
	u1, err := user.NewUsers(dsh1.CreatedBy)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
		})
		return
	}
	schema1 := make(map[string]interface{})
	if err := json.Unmarshal([]byte(dsh1.Schema), &schema1); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}
	res["schema1"] = proto.SchemaHistoryDetailsData{
		ID:            dsh1.ID,
		SchemaID:      dsh1.SchemaID,
		Name:          dsh1.Name,
		Description:   dsh1.Description,
		Schema:        schema1,
		CreatedTime:   dsh1.CreatedAt.Format("2006-01-02 15:04"),
		LastUpdatedBy: u1.Username,
	}

	if data.HistoryID2 == 0 {
		u2, err := user.NewUsers(currentDefinitionSchema.(*definition.DefinitionSchemas).UpdatedBy)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    proto.Display404ErrorMessage,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
			})
			return
		}
		schema2 := make(map[string]interface{})
		if err := json.Unmarshal([]byte(currentDefinitionSchema.(*definition.DefinitionSchemas).Schema), &schema2); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
			})
			return
		}
		res["schema2"] = proto.SchemaHistoryDetailsData{
			ID:            0,
			SchemaID:      currentDefinitionSchema.(*definition.DefinitionSchemas).ID,
			Name:          currentDefinitionSchema.(*definition.DefinitionSchemas).Name,
			Description:   currentDefinitionSchema.(*definition.DefinitionSchemas).Description,
			Schema:        schema2,
			CreatedTime:   currentDefinitionSchema.(*definition.DefinitionSchemas).CreatedAt.Format("2006-01-02 15:04"),
			LastUpdatedBy: u2.Username,
		}

		ctx.JSON(http.StatusOK, res)
		return
	}

	dsh2, err := definition.NewDefinitionSchemaHistories(data.HistoryID2)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if dsh2.SchemaID != currentDefinitionSchema.(*definition.DefinitionSchemas).ID {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	u2, err := user.NewUsers(dsh2.CreatedBy)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
		})
		return
	}
	schema2 := make(map[string]interface{})
	if err := json.Unmarshal([]byte(dsh2.Schema), &schema2); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}

	res["schema2"] = proto.SchemaHistoryDetailsData{
		ID:            dsh2.ID,
		SchemaID:      dsh2.SchemaID,
		Name:          dsh2.Name,
		Description:   dsh2.Description,
		Schema:        schema2,
		CreatedTime:   dsh2.CreatedAt.Format("2006-01-02 15:04"),
		LastUpdatedBy: u2.Username,
	}

	ctx.JSON(http.StatusOK, res)
}

func DefinitionSchemaHistoryRestore(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")
	currentDefinitionSchema, _ := ctx.Get("CurrentDefinitionSchema")

	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	uriData := proto.SchemaHistoryUriData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	dsh, err := definition.NewDefinitionSchemaHistories(uriData.HistoryID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if currentDefinitionSchema.(*definition.DefinitionSchemas).ID != dsh.SchemaID {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.NotFound"}),
		})
		return
	}

	if err := dsh.Restore(currentDefinitionSchema.(*definition.DefinitionSchemas), currentUser.(*user.Users).ID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "History.RestoreFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
