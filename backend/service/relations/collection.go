package relations

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/model/iteration"
	referencerelation "github.com/apicat/apicat/v2/backend/model/reference_relation"
	"github.com/apicat/apicat/v2/backend/model/share"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/service/except"
	"github.com/apicat/apicat/v2/backend/service/reference"
)

// DeleteCollections 删除集合并清理关联数据
func DeleteCollections(ctx context.Context, pID string, c *collection.Collection, tm *team.TeamMember) error {
	var collections []*collection.Collection
	if err := model.DB(ctx).Where("parent_id = ?", c.ID).Find(&collections).Error; err != nil {
		return err
	}

	var (
		ids []uint
		cs  []*collection.Collection
	)
	for _, subNode := range collections {
		ids = append(ids, subNode.ID)
		cs = append(cs, subNode)
	}

	for _, subNode := range collections {
		if err := DeleteCollections(ctx, pID, subNode, tm); err != nil {
			return err
		}
	}

	ids = append(ids, c.ID)
	cs = append(cs, c)

	// 集合解引用
	for _, c := range cs {
		specCollection, err := CollectionDerefWithSpec(ctx, c)
		if err != nil {
			slog.ErrorContext(ctx, "collection_relations.DeleteCollections.CollectionDerefWithSpec", "err", err)
			continue
		}

		contentByte, err := json.Marshal(specCollection.Content)
		if err != nil {
			slog.ErrorContext(ctx, "collection_relations.DeleteCollections.json.Marshal", "err", err)
			continue
		}

		c.Update(ctx, c.Title, string(contentByte), tm.ID)
	}

	// 删除集合在迭代中的该集合
	if err := iteration.BatchDeleteIterationApi(ctx, ids...); err != nil {
		slog.ErrorContext(ctx, "collection.Deletes.BatchDeleteIterationApi", "err", err)
		return err
	}

	// 删除该集合的分享令牌
	if err := share.DeleteCollectionShareTmpTokens(ctx, ids...); err != nil {
		slog.ErrorContext(ctx, "collection.Deletes.DeleteCollectionShareTmpTokens", "err", err)
		return err
	}

	// 删除该集合的引用关系
	responseIDs, err := reference.ParseRefResponsesFromCollection(c)
	if err != nil {
		slog.ErrorContext(ctx, "collection.Deletes.ParseRefResponsesFromCollection", "err", err)
		return err
	}
	if err := referencerelation.DelRefResponseCollection(ctx, c.ID, responseIDs...); err != nil {
		slog.ErrorContext(ctx, "collection.Deletes.DelRefResponseCollection", "err", err)
		return err
	}

	schemaIDs, err := reference.ParseRefSchemasFromCollection(c)
	if err != nil {
		slog.ErrorContext(ctx, "collection.Deletes.ParseRefSchemasFromCollection", "err", err)
		return err
	}
	if err := referencerelation.DelRefSchemaCollection(ctx, c.ID, schemaIDs...); err != nil {
		slog.ErrorContext(ctx, "collection.Deletes.DelRefSchemaCollection", "err", err)
		return err
	}

	// 删除该集合的被排除关系
	paramIDs, err := except.ParseExceptParamsFromCollection(c)
	if err != nil {
		slog.ErrorContext(ctx, "collection.Deletes.ParseExceptParamsFromCollection", "err", err)
		return err
	}
	if err := referencerelation.DelExceptParamCollection(ctx, c.ID, paramIDs...); err != nil {
		slog.ErrorContext(ctx, "collection.Deletes.DelExceptParamCollection", "err", err)
		return err
	}

	return collection.BatchDeleteCollections(ctx, tm.ID, ids...)
}

// CollectionDerefWithSpec 将集合解引用并转为spec.collection结构
func CollectionDerefWithSpec(ctx context.Context, c *collection.Collection) (*spec.Collection, error) {
	collectionSpec, err := c.ToSpec()
	if err != nil {
		return nil, err
	}

	specDefinitions := &spec.Definitions{}
	specDefinitions.Schemas, err = definition.GetDefinitionSchemasWithSpec(ctx, c.ProjectID)
	if err != nil {
		return nil, err
	}
	specDefinitions.Responses, err = definition.GetDefinitionResponsesWithSpec(ctx, c.ProjectID)
	if err != nil {
		return nil, err
	}

	specGlobalParameters, err := global.GetGlobalParametersWithSpec(ctx, c.ProjectID)
	if err != nil {
		return nil, err
	}

	if err := collectionSpec.Content.DeepDerefAll(specGlobalParameters, specDefinitions); err != nil {
		return nil, err
	} else {
		return collectionSpec, nil
	}
}

// CollectionDerefWithApiCatSpec 将集合解引用并转为spec结构
func CollectionDerefWithApiCatSpec(ctx context.Context, c *collection.Collection) (*spec.Spec, error) {
	collectionSpec, err := CollectionDerefWithSpec(ctx, c)
	if err != nil {
		return nil, err
	}

	apicatStruct := spec.NewEmptySpec()
	apicatStruct.Collections = append(apicatStruct.Collections, collectionSpec)
	return apicatStruct, nil
}

// CollectionImport 导入集合
func CollectionImport(ctx context.Context, member *team.TeamMember, projectID string, parentID uint, collections []*spec.Collection, refContentNameToId *collection.RefContentVirtualIDToId) []*collection.Collection {
	collectionList := make([]*collection.Collection, 0)

	var emptySlice []uint
	for i, c := range collections {
		if len(c.Items) > 0 || c.Type == "category" {
			category := &collection.Collection{
				ProjectID: projectID,
				ParentID:  parentID,
				Title:     c.Title,
				Type:      collection.CategoryType,
			}
			if err := category.Create(ctx, member); err == nil {
				collectionList = append(collectionList, category)
				children := CollectionImport(ctx, member, projectID, category.ID, c.Items, refContentNameToId)
				collectionList = append(collectionList, children...)
			}
		} else {
			if collectionByte, err := json.Marshal(c.Content); err == nil {
				collectionStr := string(collectionByte)
				collectionStr = collection.ReplaceVirtualIDToID(collectionStr, refContentNameToId.DefinitionSchemas, "\"#/definitions/schemas/")
				collectionStr = collection.ReplaceVirtualIDToID(collectionStr, refContentNameToId.DefinitionResponses, "\"#/definitions/responses/")
				collectionStr = collection.ReplaceVirtualIDToID(collectionStr, refContentNameToId.DefinitionParameters, "\"#/definitions/parameters/")
				collectionStr = replaceGlobalParametersVirtualIDToID(collectionStr, refContentNameToId.GlobalParameters)

				record := &collection.Collection{
					ProjectID:    projectID,
					ParentID:     parentID,
					Title:        c.Title,
					Type:         collection.HttpType,
					Content:      collectionStr,
					DisplayOrder: i,
				}
				if err := record.Create(ctx, member); err == nil {
					collectionList = append(collectionList, record)
					collection.TagImport(ctx, projectID, record.ID, c.Tags)
				}

				if err := reference.UpdateCollectionRef(ctx, record, emptySlice, emptySlice, emptySlice); err != nil {
					slog.ErrorContext(ctx, "CollectionImport.UpdateCollectionRef", "err", err)
				}
			}
		}
	}

	return collectionList
}

// replaceGlobalParametersVirtualIDToID 将集合中的全局参数的虚拟ID替换为真实ID
func replaceGlobalParametersVirtualIDToID(content string, virtualIDToIDMap collection.VirtualIDToIDMap) string {
	specContent, err := spec.NewCollectionFromJson(content)
	if err != nil {
		return content
	}

	req := specContent.Content.GetRequest()
	if req != nil {
		for k, v := range req.Attrs.GlobalExcepts.Header {
			if id, ok := virtualIDToIDMap[v]; ok {
				req.Attrs.GlobalExcepts.Header[k] = int64(id)
			}
		}
		for k, v := range req.Attrs.GlobalExcepts.Query {
			if id, ok := virtualIDToIDMap[v]; ok {
				req.Attrs.GlobalExcepts.Query[k] = int64(id)
			}
		}
		for k, v := range req.Attrs.GlobalExcepts.Cookie {
			if id, ok := virtualIDToIDMap[v]; ok {
				req.Attrs.GlobalExcepts.Cookie[k] = int64(id)
			}
		}
	}

	newContent, err := specContent.ToJson()
	if err != nil {
		return content
	}
	return newContent
}
