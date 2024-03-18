package definitionrelations

import (
	"apicat-cloud/backend/model/collection"
	"apicat-cloud/backend/model/definition"
	"apicat-cloud/backend/model/project"
	referencerelationship "apicat-cloud/backend/model/reference_relationship"
	"apicat-cloud/backend/model/team"
	"context"
	"encoding/json"
	"log/slog"
	"regexp"
	"strconv"

	"apicat-cloud/backend/module/array_operation"
	"apicat-cloud/backend/module/spec"
)

// ReadDefinitionResponseReference 读取collection中引用的response
func ReadDefinitionResponseReference(ctx context.Context, content string) []uint {
	// 定义正则表达式
	re := regexp.MustCompile(`"\$ref":"#/definitions/responses/(\d+)"`)

	// 在字符串中查找匹配项 matches: [["$ref":"#/definitions/responses/2050" 2050] ["$ref":"#/definitions/responses/2051" 2051]]
	matches := re.FindAllStringSubmatch(content, -1)

	// 遍历匹配项
	list := make([]uint, 0)
	for _, match := range matches {
		if len(match) >= 2 {
			// 第一个匹配项是整个匹配，从第二个匹配项开始是捕获组
			refID, err := strconv.Atoi(match[1])
			if err == nil {
				list = append(list, uint(refID))
			}
		}
	}
	return list
}

func UpdateResponseReference(ctx context.Context, r *definition.DefinitionResponse) error {
	oldResponseRef, err := referencerelationship.GetResponseReferencesByResponse(ctx, r.ProjectID, r.ID)
	if err != nil {
		return err
	}

	oldResponseRefSchemaDict := make(map[uint]*referencerelationship.ResponseReference, 0)
	for _, v := range oldResponseRef {
		oldResponseRefSchemaDict[v.RefSchemaID] = v
	}

	responseRefIDs := ReadDefinitionSchemaReference(ctx, r.Content)

	wantPop := make([]uint, 0)
	wantPush := make([]*referencerelationship.ResponseReference, 0)

	// 删除老引用关系中存在但当前引用中不存在的引用
	for key, value := range oldResponseRefSchemaDict {
		if !array_operation.InArray[uint](key, responseRefIDs) {
			wantPop = append(wantPop, value.ID)
		}
	}
	// 添加老引用关系中不存在但当前引用中存在的引用
	for _, v := range responseRefIDs {
		if _, ok := oldResponseRefSchemaDict[v]; !ok {
			wantPush = append(wantPush, &referencerelationship.ResponseReference{
				ProjectID:   r.ProjectID,
				ResponseID:  r.ID,
				RefSchemaID: v,
			})
		}
	}

	if err := referencerelationship.BatchCreateResponseReference(ctx, wantPush); err != nil {
		return err
	}
	return referencerelationship.BatchDeleteResponseReference(ctx, wantPop...)
}

// dereferenceDefinitionResponseInCollection 处理collection中引用的response
// deref: true: 展开引用，false: 删除引用
func dereferenceDefinitionResponseInCollection(ctx context.Context, c *collection.Collection, dr *definition.DefinitionResponse, deref bool) error {
	if c == nil || dr == nil {
		return nil
	}
	if c.Type == collection.CategoryType || dr.Type == definition.ResponseCategory {
		return nil
	}

	collection := &spec.Collection{
		ID:    c.ID,
		Title: c.Title,
		Type:  spec.CollectionType(c.Type),
	}
	if err := json.Unmarshal([]byte(c.Content), &collection.Content); err != nil {
		return err
	}

	if deref {
		response := &spec.HTTPResponseDefine{
			ID:          int64(dr.ID),
			Name:        dr.Name,
			Type:        dr.Type,
			ParentId:    uint64(dr.ParentID),
			Description: dr.Description,
		}
		if err := json.Unmarshal([]byte(dr.Header), &response.Header); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(dr.Content), &response.Content); err != nil {
			return err
		}

		if err := collection.DerefResponse(response); err != nil {
			return err
		}
	} else {
		if err := collection.DelResponseByRefId(int64(dr.ID)); err != nil {
			return err
		}
	}

	content, err := json.Marshal(collection.Content)
	if err != nil {
		return err
	}
	c.Content = string(content)
	return nil
}

// dereferenceDefinitionResponseInCollections 处理引用了response的collections，包含 collectio.content 和 collections与response的引用关系
// deref: true: 展开引用，false: 删除引用
func dereferenceDefinitionResponseInCollections(ctx context.Context, dr *definition.DefinitionResponse, deref bool) ([]*referencerelationship.CollectionReference, error) {
	// 查引用了response的collection
	collectionRecords, err := referencerelationship.GetCollectionReferencesByRef(ctx, dr.ProjectID, dr.ID, referencerelationship.ReferenceResponse)
	if err != nil {
		return nil, err
	}
	if len(collectionRecords) == 0 {
		return nil, nil
	}

	var (
		collectionIDs    []uint
		collectionRefIDs []uint
	)
	for _, ref := range collectionRecords {
		collectionIDs = append(collectionIDs, ref.CollectionID)
		collectionRefIDs = append(collectionRefIDs, ref.ID)
	}

	// 处理collection.content
	collections, err := collection.GetCollections(ctx, &project.Project{ID: dr.ProjectID}, collectionIDs...)
	if err != nil {
		return nil, err
	}
	for _, c := range collections {
		if err := dereferenceDefinitionResponseInCollection(ctx, c, dr, deref); err != nil {
			return nil, err
		}
		if err := c.Update(ctx, c.Title, c.Content, c.UpdatedBy); err != nil {
			return nil, err
		}
	}

	// 删除collection引用记录
	if err := referencerelationship.BatchDeleteCollectionReference(ctx, collectionRefIDs...); err != nil {
		return nil, err
	}

	return collectionRecords, nil
}

// RemoveResponseReferences 清除response引用，将修改数据库中collection.content
func RemoveResponseReferences(ctx context.Context, dr *definition.DefinitionResponse) error {
	_, err := dereferenceDefinitionResponseInCollections(ctx, dr, false)
	if err != nil {
		return err
	}

	// 查response引用的schema
	responseRef, err := referencerelationship.GetResponseReferencesByResponse(ctx, dr.ProjectID, dr.ID)
	if err != nil {
		return err
	}

	// 删除response引用记录
	var responseRefIDs []uint
	for _, ref := range responseRef {
		responseRefIDs = append(responseRefIDs, ref.ID)
	}
	return referencerelationship.BatchDeleteResponseReference(ctx, responseRefIDs...)
}

// UnpackResponseReferences 展开response引用，将修改数据库中collection.content
func UnpackResponseReferences(ctx context.Context, dr *definition.DefinitionResponse) error {
	collectionRef, err := dereferenceDefinitionResponseInCollections(ctx, dr, true)
	if err != nil {
		return err
	}

	// 查response引用的schema
	responseRef, err := referencerelationship.GetResponseReferencesByResponse(ctx, dr.ProjectID, dr.ID)
	if err != nil {
		return err
	}

	// 将response引用的schema添加到collection
	responseRefIDs := make([]uint, 0)
	wantPush := make([]*referencerelationship.CollectionReference, 0)
	for _, r := range responseRef {
		responseRefIDs = append(responseRefIDs, r.ID)
		for _, c := range collectionRef {
			wantPush = append(wantPush, &referencerelationship.CollectionReference{
				ProjectID:    c.ProjectID,
				CollectionID: c.CollectionID,
				RefID:        r.RefSchemaID,
				RefType:      referencerelationship.ReferenceSchema,
			})
		}
	}
	if err := referencerelationship.BatchCreateCollectionReference(ctx, wantPush); err != nil {
		return err
	}

	// 删除response引用记录
	return referencerelationship.BatchDeleteResponseReference(ctx, responseRefIDs...)
}

func ImportDefinitionResponses(ctx context.Context, projectID string, responses spec.HTTPResponseDefines, tm *team.TeamMember, schemaVirtualIDToID collection.VirtualIDToIDMap, parentID uint) collection.VirtualIDToIDMap {
	res := collection.VirtualIDToIDMap{}
	if len(responses) == 0 {
		return res
	}

	for i, response := range responses {
		record := &definition.DefinitionResponse{
			ProjectID:    projectID,
			ParentID:     parentID,
			Name:         response.Name,
			DisplayOrder: uint(i),
		}

		if len(response.Items) > 0 || response.Type == definition.ResponseCategory {
			record.Type = definition.ResponseCategory
			if err := record.Create(ctx, tm); err == nil {
				res[response.ID] = record.ID
				res.Merge(ImportDefinitionResponses(ctx, projectID, response.Items, tm, schemaVirtualIDToID, record.ID))
			}
		} else {
			var (
				header  string
				content string
			)
			if headerStr, err := json.Marshal(response.Header); err == nil {
				header = string(headerStr)
			}
			if contentStr, err := json.Marshal(response.Content); err == nil {
				content = collection.ReplaceVirtualIDToID(string(contentStr), schemaVirtualIDToID, "#/definitions/schemas/")
			}

			record.Type = definition.ResponseResponse
			record.Description = response.Description
			record.Header = header
			record.Content = content
			if err := record.Create(ctx, tm); err == nil {
				res[response.ID] = record.ID
			}

			if err := UpdateResponseReference(ctx, record); err != nil {
				slog.ErrorContext(ctx, "ImportDefinitionResponses.UpdateResponseReference", "err", err)
			}
		}
	}

	return res
}
