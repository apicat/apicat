package trash

import (
	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TrashsList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	p, _ := currentProject.(*project.Projects)

	c, _ := collection.NewCollections()
	c.ProjectId = p.ID
	collections, err := c.TrashList()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Trashs.QueryFailed"}),
		})
	}

	result := []map[string]interface{}{}
	for _, v := range collections {
		result = append(result, map[string]interface{}{
			"id":         v.ID,
			"title":      v.Title,
			"type":       v.Type,
			"deleted_at": v.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			"deleted_by": v.Deleter(),
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func TrashsRecover(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	trashsRecoverQuery := proto.TrashsRecoverQuery{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&trashsRecoverQuery)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	trashsRecoverBody := proto.TrashsRecoverBody{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&trashsRecoverBody)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	currentProject, _ := ctx.Get("CurrentProject")
	p, _ := currentProject.(*project.Projects)

	allOK := true

	for _, v := range trashsRecoverQuery.CollectionID {
		c, _ := collection.NewCollections()
		c.ID = v
		c.ProjectId = p.ID
		if err := c.GetUnscopedCollections(); err != nil {
			allOK = false
			continue
		}
		c.ParentId = trashsRecoverBody.Category
		c.Restore()
	}

	if !allOK {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Trashs.PartialRecoveryFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
