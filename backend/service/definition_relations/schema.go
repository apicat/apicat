package definitionrelations

import (
	"context"
	"encoding/json"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/project"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/module/spec"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
)

// ReadDefinitionSchemaReference 读取collection中引用的schema
func ReadDefinitionSchemaReference(ctx context.Context, content string) []uint {
	// 定义正则表达式
	re := regexp.MustCompile(`"\$ref":"#/definitions/schemas/(\d+)"`)

	// 在字符串中查找匹配项 matches: [["$ref":"#/definitions/schemas/2050" 2050] ["$ref":"#/definitions/schemas/2051" 2051]]
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

func UpdateSchemaReference(ctx context.Context, s *definition.DefinitionSchema) error {
	oldSchemaRef, err := referencerelationship.GetSchemaReferencesBySchema(ctx, s.ProjectID, s.ID)
	if err != nil {
		return err
	}

	oldSchemaRefSchemaDict := make(map[uint]*referencerelationship.SchemaReference, 0)
	for _, v := range oldSchemaRef {
		oldSchemaRefSchemaDict[v.RefSchemaID] = v
	}

	schemaRefIDs := ReadDefinitionSchemaReference(ctx, s.Schema)

	wantPop := make([]uint, 0)
	wantPush := make([]*referencerelationship.SchemaReference, 0)

	// 删除老引用关系中存在但当前引用中不存在的引用
	for key, value := range oldSchemaRefSchemaDict {
		if !arrutil.InArray[uint](key, schemaRefIDs) {
			wantPop = append(wantPop, value.ID)
		}
	}
	// 添加老引用关系中不存在但当前引用中存在的引用
	for _, v := range schemaRefIDs {
		if _, ok := oldSchemaRefSchemaDict[v]; !ok {
			wantPush = append(wantPush, &referencerelationship.SchemaReference{
				ProjectID:   s.ProjectID,
				SchemaID:    s.ID,
				RefSchemaID: v,
			})
		}
	}

	if err := referencerelationship.BatchCreateSchemaReference(ctx, wantPush); err != nil {
		return err
	}
	return referencerelationship.BatchDeleteSchemaReference(ctx, wantPop...)
}

// dereferenceDefinitionSchemaInCollection 处理collection中引用的schema
// deref: true: 展开引用，false: 删除引用
func dereferenceDefinitionSchemaInCollection(ctx context.Context, c *collection.Collection, ds *definition.DefinitionSchema, deref bool) error {
	if c == nil || ds == nil {
		return nil
	}
	if c.Type == collection.CategoryType || ds.Type == definition.SchemaCategory {
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

	schema, err := ds.ToSpec()
	if err != nil {
		return err
	}

	if deref {
		if err := collection.DerefSchema(schema); err != nil {
			return err
		}
	} else {
		if err := collection.DelRefSchema(schema); err != nil {
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

// dereferenceDefinitionSchemaInCollections 处理引用了schema的collections，包含 collections.content 和 collections与schema的引用关系
// deref: true: 展开引用，false: 删除引用
func dereferenceDefinitionSchemaInCollections(ctx context.Context, ds *definition.DefinitionSchema, deref bool) ([]*referencerelationship.CollectionReference, error) {
	// 查引用了schema的collections
	ctosRecords, err := referencerelationship.GetCollectionReferencesByRef(ctx, ds.ProjectID, ds.ID, referencerelationship.ReferenceSchema)
	if err != nil {
		return nil, err
	}
	if len(ctosRecords) == 0 {
		return nil, nil
	}
	var (
		collectionIDs    []uint
		collectionRefIDs []uint
	)
	for _, ref := range ctosRecords {
		collectionIDs = append(collectionIDs, ref.CollectionID)
		collectionRefIDs = append(collectionRefIDs, ref.ID)
	}

	// 处理collections.content
	collections, err := collection.GetCollections(ctx, &project.Project{ID: ds.ProjectID}, collectionIDs...)
	if err != nil {
		return nil, err
	}
	for _, c := range collections {
		if err := dereferenceDefinitionSchemaInCollection(ctx, c, ds, deref); err != nil {
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

	return ctosRecords, nil
}

// dereferenceDefinitionSchemaInSchema 处理schema中引用的schema
// deref: true: 展开引用，false: 删除引用
func dereferenceDefinitionSchemaInSchema(ctx context.Context, ds *definition.DefinitionSchema, refDs *definition.DefinitionSchema, deref bool) error {
	if ds == nil || refDs == nil {
		return nil
	}
	if ds.Type == definition.SchemaCategory || refDs.Type == definition.SchemaCategory {
		return nil
	}

	schema, err := ds.ToSpec()
	if err != nil {
		return err
	}
	refSchema, err := refDs.ToSpec()
	if err != nil {
		return err
	}

	if deref {
		schema.Deref(refSchema)
	} else {
		if err := schema.DelRef(refSchema); err != nil {
			return err
		}
	}

	content, err := json.Marshal(schema.Schema)
	if err != nil {
		return err
	}
	ds.Schema = string(content)
	return nil
}

// dereferenceDefinitionSchemaInSchemas 处理引用了schema的schemas，包含 schema.schema 和 schemas与schema的引用关系
// deref: true: 展开引用，false: 删除引用
func dereferenceDefinitionSchemaInSchemas(ctx context.Context, ds *definition.DefinitionSchema, deref bool) ([]*referencerelationship.SchemaReference, error) {
	// 查引用了要删除的schema的schemas
	schemaRecords, err := referencerelationship.GetSchemaReferencesByRefSchema(ctx, ds.ProjectID, ds.ID)
	if err != nil {
		return nil, err
	}
	if len(schemaRecords) == 0 {
		return nil, nil
	}
	var (
		schemasIDs   []uint
		schemaRefIDs []uint
	)
	for _, ref := range schemaRecords {
		schemasIDs = append(schemasIDs, ref.SchemaID)
		schemaRefIDs = append(schemaRefIDs, ref.ID)
	}

	// 处理schemas.content
	schemas, err := definition.GetDefinitionSchemas(ctx, &project.Project{ID: ds.ProjectID}, schemasIDs...)
	if err != nil {
		return nil, err
	}
	for _, s := range schemas {
		if err := dereferenceDefinitionSchemaInSchema(ctx, s, ds, deref); err != nil {
			return nil, err
		}
		if err := s.Update(ctx, s.Name, s.Description, s.Schema, s.UpdatedBy); err != nil {
			return nil, err
		}
	}

	// 删除schemas引用记录
	if err := referencerelationship.BatchDeleteSchemaReference(ctx, schemaRefIDs...); err != nil {
		return nil, err
	}

	return schemaRecords, nil
}

// dereferenceDefinitionSchemaInResponse 处理response中引用的schema
// deref: true: 展开引用，false: 删除引用
func dereferenceDefinitionSchemaInResponse(ctx context.Context, dr *definition.DefinitionResponse, ds *definition.DefinitionSchema, deref bool) error {
	if dr == nil || ds == nil {
		return nil
	}
	if dr.Type == definition.ResponseCategory || ds.Type == definition.SchemaCategory {
		return nil
	}

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

	schema, err := ds.ToSpec()
	if err != nil {
		return err
	}

	if deref {
		response.DerefSchema(schema)
	} else {
		if err := response.DelRefSchema(schema); err != nil {
			return err
		}
	}

	header, err := json.Marshal(response.Header)
	if err != nil {
		return err
	}
	content, err := json.Marshal(response.Content)
	if err != nil {
		return err
	}
	dr.Header = string(header)
	dr.Content = string(content)
	return nil
}

// dereferenceDefinitionSchemaInResponses 处理引用了schema的responses，包含 responses.content 和 responses与schema的引用关系
// deref: true: 展开引用，false: 删除引用
func dereferenceDefinitionSchemaInResponses(ctx context.Context, ds *definition.DefinitionSchema, deref bool) ([]*referencerelationship.ResponseReference, error) {
	// 查引用了要删除的schema的responses
	responseRecords, err := referencerelationship.GetResponseReferencesByRefSchema(ctx, ds.ProjectID, ds.ID)
	if err != nil {
		return nil, err
	}
	if len(responseRecords) == 0 {
		return nil, nil
	}

	var (
		responseIDs    []uint
		responseRefIDs []uint
	)
	for _, ref := range responseRecords {
		responseIDs = append(responseIDs, ref.ResponseID)
		responseRefIDs = append(responseRefIDs, ref.ID)
	}

	// 处理responses.content
	responses, err := definition.GetDefinitionResponses(ctx, &project.Project{ID: ds.ProjectID}, responseIDs...)
	if err != nil {
		return nil, err
	}
	for _, r := range responses {
		if err := dereferenceDefinitionSchemaInResponse(ctx, r, ds, deref); err != nil {
			return nil, err
		}
		if err := r.Update(ctx, r.UpdatedBy); err != nil {
			return nil, err
		}
	}

	// 删除response引用记录
	if err := referencerelationship.BatchDeleteResponseReference(ctx, responseRefIDs...); err != nil {
		return nil, err
	}

	return responseRecords, nil
}

// RemoveDefinitionSchemaReferences 清除schema引用，将修改数据库中collection.content,definition_response.content,definition_schema.schema
func RemoveDefinitionSchemaReferences(ctx context.Context, ds *definition.DefinitionSchema) error {
	_, err := dereferenceDefinitionSchemaInCollections(ctx, ds, false)
	if err != nil {
		return err
	}
	_, err = dereferenceDefinitionSchemaInResponses(ctx, ds, false)
	if err != nil {
		return err
	}
	_, err = dereferenceDefinitionSchemaInSchemas(ctx, ds, false)
	if err != nil {
		return err
	}

	// 查询要删除的schema引用的schema
	schemaRefRecord, err := referencerelationship.GetSchemaReferencesBySchema(ctx, ds.ProjectID, ds.ID)
	if err != nil {
		return err
	}
	var schemaRefIDs []uint
	for _, ref := range schemaRefRecord {
		schemaRefIDs = append(schemaRefIDs, ref.ID)
	}

	return referencerelationship.BatchDeleteSchemaReference(ctx, schemaRefIDs...)
}

// UnpackDefinitionSchemaReferences 展开schema引用，将修改数据库中collection.content,definition_response.content,definition_schema.schema
func UnpackDefinitionSchemaReferences(ctx context.Context, ds *definition.DefinitionSchema) error {
	collectionRefs, err := dereferenceDefinitionSchemaInCollections(ctx, ds, true)
	if err != nil {
		return err
	}
	responseRefs, err := dereferenceDefinitionSchemaInResponses(ctx, ds, true)
	if err != nil {
		return err
	}
	schemaRefs, err := dereferenceDefinitionSchemaInSchemas(ctx, ds, true)
	if err != nil {
		return err
	}

	// 查询要删除的schema引用的schema
	schemaRefRecord, err := referencerelationship.GetSchemaReferencesBySchema(ctx, ds.ProjectID, ds.ID)
	if err != nil {
		return err
	}
	schemaRefIDs := make([]uint, 0)
	collectionWantPush := make([]*referencerelationship.CollectionReference, 0)
	schemaWantPush := make([]*referencerelationship.SchemaReference, 0)
	responseWantPush := make([]*referencerelationship.ResponseReference, 0)
	for _, v := range schemaRefRecord {
		schemaRefIDs = append(schemaRefIDs, v.ID)
		for _, c := range collectionRefs {
			collectionWantPush = append(collectionWantPush, &referencerelationship.CollectionReference{
				ProjectID:    c.ProjectID,
				CollectionID: c.CollectionID,
				RefID:        v.RefSchemaID,
				RefType:      referencerelationship.ReferenceSchema,
			})
		}
		for _, r := range responseRefs {
			responseWantPush = append(responseWantPush, &referencerelationship.ResponseReference{
				ProjectID:   r.ProjectID,
				ResponseID:  r.ResponseID,
				RefSchemaID: v.RefSchemaID,
			})
		}
		for _, s := range schemaRefs {
			schemaWantPush = append(schemaWantPush, &referencerelationship.SchemaReference{
				ProjectID:   s.ProjectID,
				SchemaID:    s.SchemaID,
				RefSchemaID: v.RefSchemaID,
			})
		}
	}

	// 将要删除的schema引用的其他schema引用关系添加到collections
	if err := referencerelationship.BatchCreateCollectionReference(ctx, collectionWantPush); err != nil {
		return err
	}
	// 将要删除的schema引用的其他schema引用关系添加到responsers
	if err := referencerelationship.BatchCreateResponseReference(ctx, responseWantPush); err != nil {
		return err
	}
	// 将要删除的schema引用的其他schema引用关系添加到引用了要删除的schema的schemas
	if err := referencerelationship.BatchCreateSchemaReference(ctx, schemaWantPush); err != nil {
		return err
	}

	return referencerelationship.BatchDeleteSchemaReference(ctx, schemaRefIDs...)
}

func ImportDefinitionSchemas(ctx context.Context, projectID string, schemas spec.Schemas, tm *team.TeamMember, parentID uint) collection.VirtualIDToIDMap {
	res := collection.VirtualIDToIDMap{}
	if len(schemas) == 0 {
		return res
	}

	for i, schema := range schemas {
		record := &definition.DefinitionSchema{
			ProjectID:    projectID,
			ParentID:     parentID,
			Name:         schema.Name,
			DisplayOrder: uint(i),
		}

		if len(schema.Items) > 0 || schema.Type == definition.SchemaCategory {
			record.Type = definition.SchemaCategory
			if err := record.Create(ctx, tm); err == nil {
				res[schema.ID] = record.ID
				res.Merge(ImportDefinitionSchemas(ctx, projectID, schema.Items, tm, record.ID))
			}
		} else {
			if schemaStr, err := json.Marshal(schema.Schema); err == nil {
				record.Type = definition.SchemaSchema
				record.Description = schema.Description
				record.Schema = string(schemaStr)
				if err := record.Create(ctx, tm); err == nil {
					res[schema.ID] = record.ID
				}
			}
		}
	}

	list, err := definition.GetDefinitionSchemas(ctx, &project.Project{ID: projectID})
	if err != nil {
		slog.ErrorContext(ctx, "GetDefinitionSchemas", "err", err)
		return res
	}

	for _, v := range list {
		schema := collection.ReplaceVirtualIDToID(v.Schema, res, "#/definitions/schemas/")
		v.Update(ctx, v.Name, v.Description, schema, tm.ID)

		if err := UpdateSchemaReference(ctx, v); err != nil {
			slog.ErrorContext(ctx, "ImportDefinitionSchemas.UpdateSchemaReference", "err", err)
		}
	}

	return res
}
