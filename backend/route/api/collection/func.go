package collection

import (
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/user"
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	collectionbase "github.com/apicat/apicat/backend/route/proto/collection/base"
	collectionresponse "github.com/apicat/apicat/backend/route/proto/collection/response"
	projectbase "github.com/apicat/apicat/backend/route/proto/project/base"

	"github.com/gin-gonic/gin"
)

func convertModelCollection(ctx *gin.Context, c *collection.Collection, cUserInfo, uUserInfo *user.User) *collectionresponse.Collection {
	return &collectionresponse.Collection{
		EmbedInfo: protobase.EmbedInfo{
			ID:        c.ID,
			CreatedAt: c.CreatedAt.Unix(),
			UpdatedAt: c.UpdatedAt.Unix(),
		},
		CollectionData: collectionbase.CollectionData{
			Title:   c.Title,
			Content: c.Content,
		},
		CollectionTypeOption: collectionbase.CollectionTypeOption{
			Type: c.Type,
		},
		CollectionParentIDOption: collectionbase.CollectionParentIDOption{
			ParentID: c.ParentID,
		},
		OperatorID: projectbase.OperatorID{
			CreatedBy: cUserInfo.Name,
			UpdatedBy: uUserInfo.Name,
		},
	}
}

func convertModelCollectionHistory(ctx *gin.Context, c *collection.CollectionHistory, userInfo *user.User) *collectionresponse.CollectionHistory {
	return &collectionresponse.CollectionHistory{
		IdCreateTimeInfo: protobase.IdCreateTimeInfo{
			ID:        c.ID,
			CreatedAt: c.CreatedAt.Unix(),
		},
		CollectionHistoryData: collectionresponse.CollectionHistoryData{
			CollectionIDOption: collectionbase.CollectionIDOption{
				CollectionID: c.CollectionID,
			},
			CollectionData: collectionbase.CollectionData{
				Title:   c.Title,
				Content: c.Content,
			},
		},
		CreatedBy: userInfo.Name,
	}
}

func buildTree(parentID uint, collections []*collection.Collection, selectCIDs []uint) collectionresponse.CollectionTree {
	result := make(collectionresponse.CollectionTree, 0)

	for _, c := range collections {
		if c.ParentID == parentID {
			children := buildTree(c.ID, collections, selectCIDs)

			cl := collectionresponse.CollectionNode{
				ID:    c.ID,
				Title: c.Title,
				Items: children,
				CollectionTypeOption: collectionbase.CollectionTypeOption{
					Type: c.Type,
				},
				CollectionParentIDOption: collectionbase.CollectionParentIDOption{
					ParentID: c.ParentID,
				},
			}

			if selectCIDs != nil {
				isSelected := false
				for _, cid := range selectCIDs {
					if cid == cl.ID {
						isSelected = true
						break
					}
					if !isSelected {
						for _, v := range cl.Items {
							if *v.Selected {
								isSelected = true
								break
							}
						}
					}
				}
				cl.Selected = &isSelected
			}

			result = append(result, &cl)
		}
	}

	return result
}
