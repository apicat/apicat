package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/apicat/apicat/backend/app/util"
	"github.com/apicat/apicat/backend/common/spec"
	"github.com/apicat/apicat/backend/common/spec/plugin/export"
	"github.com/apicat/apicat/backend/common/spec/plugin/openapi"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/apicat/apicat/backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type CollectionDataGetData struct {
	ProjectID    string `uri:"project-id" binding:"required,gt=0"`
	CollectionID uint   `uri:"collection-id" binding:"required,gt=0"`
}

type ExportCollection struct {
	Type     string `form:"type" binding:"required,oneof=apicat swagger openapi3.0.0 openapi3.0.1 openapi3.0.2 openapi3.1.0 HTML md"`
	Download string `form:"download" binding:"omitempty,oneof=true false"`
}

type CollectionList struct {
	ID       uint              `json:"id"`
	ParentID uint              `json:"parent_id"`
	Title    string            `json:"title"`
	Type     string            `json:"type"`
	Selected *bool             `json:"selected,omitempty"`
	Items    []*CollectionList `json:"items"`
}

type CollectionCreate struct {
	ParentID    uint   `json:"parent_id" binding:"gte=0"`                       // 父级id
	Title       string `json:"title" binding:"required,lte=255"`                // 名称
	Type        string `json:"type" binding:"required,oneof=category doc http"` // 类型: category,doc,http
	Content     string `json:"content"`                                         // 内容
	IterationID string `json:"iteration_id" binding:"omitempty,gte=0"`          // 迭代id
}

type CollectionUpdate struct {
	Title   string `json:"title" binding:"required,lte=255"`
	Content string `json:"content"`
}

type CollectionCopyData struct {
	IterationID string `json:"iteration_id" binding:"omitempty,gte=0"`
}

type CollectionMovement struct {
	Target CollectionOrderContent `json:"target" binding:"required"`
	Origin CollectionOrderContent `json:"origin" binding:"required"`
}

type CollectionOrderContent struct {
	Pid uint   `json:"pid" binding:"gte=0"`
	Ids []uint `json:"ids" binding:"required,dive,gte=0"`
}

type CollectionDeleteData struct {
	IterationID string `form:"iteration_id" binding:"omitempty,gte=0"`
}

type CollectionsListData struct {
	IterationID string `form:"iteration_id" binding:"omitempty,gte=0"`
}

func CollectionsList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	var data CollectionsListData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	collection, _ := models.NewCollections()
	collection.ProjectId = project.ID
	collections, err := collection.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.QueryFailed"}),
		})
	}

	if data.IterationID == "" {
		ctx.JSON(http.StatusOK, buildProjectTree(0, collections))
	} else {
		iteration, err := models.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.QueryFailed"}),
			})
			return
		}

		ia, _ := models.NewIterationApis()
		cIDs, err := ia.GetCollectionIDByIterationID(iteration.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.QueryFailed"}),
			})
			return
		}

		ctx.JSON(http.StatusOK, buildIterationTree(0, collections, cIDs))
	}
}

func buildProjectTree(parentID uint, collections []*models.Collections) []*CollectionList {
	return buildTree(parentID, collections, false)
}

func buildIterationTree(parentID uint, collections []*models.Collections, selectCIDs []uint) []*CollectionList {
	return buildTree(parentID, collections, true, selectCIDs...)
}

func buildTree(parentID uint, collections []*models.Collections, isIteration bool, selectCIDs ...uint) []*CollectionList {
	result := make([]*CollectionList, 0)

	for _, c := range collections {
		if c.ParentId == parentID {
			children := buildTree(c.ID, collections, isIteration, selectCIDs...)

			c := CollectionList{
				ID:       c.ID,
				ParentID: c.ParentId,
				Title:    c.Title,
				Type:     c.Type,
				Items:    children,
			}

			isSelected := false
			if isIteration {
				for _, cid := range selectCIDs {
					if cid == c.ID {
						isSelected = true
						break
					}
					if !isSelected {
						for _, v := range c.Items {
							if *v.Selected {
								isSelected = true
								break
							}
						}
					}
				}
				c.Selected = &isSelected
			}

			result = append(result, &c)
		}
	}

	return result
}

func CollectionsGet(ctx *gin.Context) {
	currentCollection, _ := ctx.Get("CurrentCollection")
	collection := currentCollection.(*models.Collections)

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
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
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

	if data.IterationID != "" {
		_, err := models.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}
	}

	collection, _ := models.NewCollections()
	collection.ProjectId = project.ID
	collection.ParentId = data.ParentID
	collection.Title = data.Title
	collection.Type = data.Type
	collection.Content = data.Content
	collection.CreatedBy = currentProjectMember.(*models.ProjectMembers).UserID
	collection.UpdatedBy = currentProjectMember.(*models.ProjectMembers).UserID

	var err error
	if collection.Type == "category" {
		err = collection.CreateCategory()
	} else {
		err = collection.CreateDoc()
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.CreateFailed"}),
		})
		return
	}

	if data.IterationID != "" {
		iteration, err := models.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}

		ia, _ := models.NewIterationApis()
		ia.IterationID = iteration.ID
		ia.CollectionID = collection.ID
		ia.CollectionType = collection.Type
		if err := ia.Create(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}
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
	currentCollection, _ := ctx.Get("CurrentCollection")
	collection := currentCollection.(*models.Collections)

	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := CollectionUpdate{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ch, _ := models.NewCollectionHistories()
	ch.CollectionId = collection.ID
	ch.Title = collection.Title
	ch.Type = collection.Type
	ch.Content = collection.Content
	ch.CreatedBy = currentProjectMember.(*models.ProjectMembers).UserID

	// 不是同一个人编辑的文档或5分钟后编辑文档内容，保存历史记录
	if collection.UpdatedBy != currentProjectMember.(*models.ProjectMembers).UserID || collection.UpdatedAt.Add(5*time.Minute).Before(time.Now()) {
		if err := ch.Create(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.UpdateFailed"}),
			})
			return
		}
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
	currentCollection, _ := ctx.Get("CurrentCollection")
	collection := currentCollection.(*models.Collections)

	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := CollectionCopyData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.IterationID != "" {
		_, err := models.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}
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

	if err := newCollection.CreateDoc(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.CreateFailed"}),
		})
		return
	}

	if data.IterationID != "" {
		iteration, err := models.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}

		ia, _ := models.NewIterationApis()
		ia.IterationID = iteration.ID
		ia.CollectionID = newCollection.ID
		ia.CollectionType = newCollection.Type
		if err := ia.Create(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}
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
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
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
	currentCollection, _ := ctx.Get("CurrentCollection")
	collection := currentCollection.(*models.Collections)

	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := CollectionDeleteData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.IterationID != "" {
		_, err := models.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.DeleteFailed"}),
			})
			return
		}

		collections, err := collection.GetSubCollectionsContainsSelf()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.DeleteFailed"}),
			})
			return
		}
		var cIDs []uint
		for _, v := range collections {
			cIDs = append(cIDs, v.ID)
		}

		if err := models.DeleteIterationApisByCollectionID(cIDs...); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.DeleteFailed"}),
			})
			return
		}
	}

	if err := models.Deletes(collection.ID, models.Conn, currentProjectMember.(*models.ProjectMembers).UserID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.DeleteFailed"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func CollectionDataGet(ctx *gin.Context) {
	uriData := CollectionDataGetData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	data := ExportCollection{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	project, err := models.NewProjects(uriData.ProjectID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	collection, err := models.NewCollections(uriData.CollectionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	apicatData := models.CollectionExport(project, collection)
	if apicatDataContent, err := json.Marshal(apicatData); err == nil {
		slog.InfoCtx(ctx, "Export", slog.String("apicat", string(apicatDataContent)))
	}

	var content []byte
	switch data.Type {
	case "swagger":
		content, err = openapi.Encode(apicatData, "2.0")
	case "openapi3.0.0":
		content, err = openapi.Encode(apicatData, "3.0.0")
	case "openapi3.0.1":
		content, err = openapi.Encode(apicatData, "3.0.1")
	case "openapi3.0.2":
		content, err = openapi.Encode(apicatData, "3.0.2")
	case "openapi3.1.0":
		content, err = openapi.Encode(apicatData, "3.1.0")
	case "HTML":
		content, err = export.HTML(apicatData)
	case "md":
		content, err = export.Markdown(apicatData)
	default:
		content, err = apicatData.ToJSON(spec.JSONOption{Indent: "  "})
	}

	slog.InfoCtx(ctx, "Export", slog.String(data.Type, string(content)))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Collections.ExportFail"}),
		})
		return
	}

	util.ExportResponse(data.Type, data.Download, project.Title+"-"+data.Type, content, ctx)
}
