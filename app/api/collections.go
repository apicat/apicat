package api

import (
	"fmt"
	"net/http"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type CollectionCheck struct {
	CollectionID uint `uri:"collection-id" binding:"required,gt=0"`
}

func (cc *CollectionCheck) CheckCollection(ctx *gin.Context) (*models.Collections, error) {
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&cc)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collection.NotFound"}),
		})
		return nil, err
	}

	collection, err := models.NewCollections(cc.CollectionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collection.NotFound"}),
		})
		return nil, err
	}

	return collection, nil
}

type CollectionList struct {
	ID       uint              `json:"id"`
	ParentID uint              `json:"parent_id"`
	Title    string            `json:"title"`
	Type     string            `json:"type"`
	Items    []*CollectionList `json:"items"`
}

type CollectionCreate struct {
	ParentID uint   `json:"parent_id" binding:"gte=0"`                       // 父级id
	Title    string `json:"title" binding:"required,lte=255"`                // 名称
	Type     string `json:"type" binding:"required,oneof=category doc http"` // 类型: category,doc,http
	Content  string `json:"content"`                                         // 内容
}

type CollectionUpdate struct {
	Title   string `json:"title" binding:"required,lte=255"`
	Content string `json:"content"`
}

type CollectionMovement struct {
	Target CollectionOrderContent `json:"target" binding:"required"`
	Origin CollectionOrderContent `json:"origin" binding:"required"`
}

type CollectionOrderContent struct {
	Pid uint   `json:"pid" binding:"gte=0"`
	Ids []uint `json:"ids" binding:"required,dive,gte=0"`
}

func CollectionsList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	collection, _ := models.NewCollections()
	collection.ProjectId = project.ID
	collections, err := collection.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.QueryFailed"}),
		})
	}

	ctx.JSON(http.StatusOK, buildTree(0, collections))
}

func buildTree(parentID uint, collections []*models.Collections) []*CollectionList {
	result := make([]*CollectionList, 0)

	for _, c := range collections {
		if c.ParentId == parentID {
			children := buildTree(c.ID, collections)
			result = append(result, &CollectionList{
				ID:       c.ID,
				ParentID: c.ParentId,
				Title:    c.Title,
				Type:     c.Type,
				Items:    children,
			})
		}
	}

	return result
}

func CollectionsGet(ctx *gin.Context) {
	cc := CollectionCheck{}
	collection, err := cc.CheckCollection(ctx)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         collection.ID,
		"parent_id":  collection.ParentId,
		"title":      collection.Title,
		"type":       collection.Type,
		"content":    collection.Content,
		"created_at": collection.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by": collection.Creator(),
		"updated_at": collection.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by": collection.Updater(),
	})
}

func CollectionsCreate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("currentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := CollectionCreate{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	collection, _ := models.NewCollections()
	collection.ProjectId = project.ID
	collection.ParentId = data.ParentID
	collection.Title = data.Title
	collection.Type = data.Type
	collection.Content = data.Content
	collection.CreatedBy = currentProjectMember.(*models.ProjectMembers).UserID
	if err := collection.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.CreateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         collection.ID,
		"parent_id":  collection.ParentId,
		"title":      collection.Title,
		"type":       collection.Type,
		"content":    collection.Content,
		"created_at": collection.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by": collection.Creator(),
		"updated_at": collection.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by": collection.Updater(),
	})
}

func CollectionsUpdate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("currentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	cc := CollectionCheck{}
	collection, err := cc.CheckCollection(ctx)
	if err != nil {
		return
	}

	data := CollectionUpdate{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	collection.Title = data.Title
	collection.Content = data.Content
	collection.UpdatedBy = currentProjectMember.(*models.ProjectMembers).UserID
	if err := collection.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.UpdateFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func CollectionsCopy(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("currentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	cc := CollectionCheck{}
	collection, err := cc.CheckCollection(ctx)
	if err != nil {
		return
	}

	newCollection := models.Collections{
		ProjectId:    collection.ProjectId,
		ParentId:     collection.ParentId,
		Title:        fmt.Sprintf("%s (copy)", collection.Title),
		Type:         collection.Type,
		Content:      collection.Content,
		DisplayOrder: collection.DisplayOrder,
		CreatedBy:    currentProjectMember.(*models.ProjectMembers).UserID,
	}

	if err := newCollection.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.CreateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         newCollection.ID,
		"parent_id":  newCollection.ParentId,
		"title":      newCollection.Title,
		"type":       newCollection.Type,
		"content":    newCollection.Content,
		"created_at": newCollection.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by": newCollection.Creator(),
		"updated_at": newCollection.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by": newCollection.Updater(),
	})
}

func CollectionsMovement(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("currentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := CollectionMovement{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for i, id := range data.Target.Ids {
		if collection, err := models.NewCollections(id); err == nil {
			collection.ParentId = data.Target.Pid
			collection.DisplayOrder = i
			collection.Update()
		}
	}

	if data.Target.Pid != data.Origin.Pid {
		for i, id := range data.Origin.Ids {
			if collection, err := models.NewCollections(id); err == nil {
				collection.ParentId = data.Origin.Pid
				collection.DisplayOrder = i
				collection.Update()
			}
		}
	}

	ctx.Status(http.StatusCreated)
}

func CollectionsDelete(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("currentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	cc := CollectionCheck{}
	collection, err := cc.CheckCollection(ctx)
	if err != nil {
		return
	}

	if err := models.Deletes(collection.ID, models.Conn, currentProjectMember.(*models.ProjectMembers).UserID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.DeleteFailed"}),
		})
		return
	}
	ctx.Status(http.StatusNoContent)
}
