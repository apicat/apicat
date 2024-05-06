package reference

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
)

func UpdateSchemaRef(ctx context.Context, s *definition.DefinitionSchema) error {
	sr := &referencerelationship.SchemaReference{SchemaID: s.ID}
	lastRefs, err := sr.GetSchemaRefs(ctx)
	if err != nil {
		return err
	}

	lastRefsMap := make(map[uint]*referencerelationship.SchemaReference, 0)
	for _, v := range lastRefs {
		lastRefsMap[v.RefSchemaID] = v
	}

	latestRefs := ParseRefSchemas(s.Schema)

	wantPop := make([]uint, 0)
	wantPush := make([]*referencerelationship.SchemaReference, 0)

	// 删除老引用关系中存在但当前引用中不存在的引用
	for key, value := range lastRefsMap {
		if !arrutil.InArray[uint](key, latestRefs) {
			wantPop = append(wantPop, value.ID)
		}
	}

	// 添加老引用关系中不存在但当前引用中存在的引用
	for _, v := range latestRefs {
		if _, ok := lastRefsMap[v]; !ok {
			wantPush = append(wantPush, &referencerelationship.SchemaReference{
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

func DeleteSchemaReference(ctx context.Context, s *definition.DefinitionSchema) error {
	sr := &referencerelationship.SchemaReference{SchemaID: s.ID}
	refs, err := sr.GetSchemaRefs(ctx)
	if err != nil {
		return err
	}

	var ids []uint
	for _, item := range refs {
		ids = append(ids, item.ID)
	}

	return referencerelationship.BatchDeleteSchemaReference(ctx, ids...)
}

func DerefSchemaFromSchemas(ctx context.Context, s *definition.DefinitionSchema, schemaIDs []uint, deref bool) error {
	schemas, err := definition.GetDefinitionSchemas(ctx, s.ProjectID, schemaIDs...)
	if err != nil {
		return err
	}

	for _, schema := range schemas {
		schema.DelRef(ctx, s, deref)
	}

	sr := &referencerelationship.SchemaReference{RefSchemaID: s.ID}
	return sr.DelByRefSchemaID(ctx)
}

func DerefSchemaFromResponses(ctx context.Context, s *definition.DefinitionSchema, responseIDs []uint, deref bool) error {
	responses, err := definition.GetDefinitionResponses(ctx, s.ProjectID, responseIDs...)
	if err != nil {
		return err
	}

	for _, response := range responses {
		response.DelRef(ctx, s, deref)
	}

	rr := &referencerelationship.ResponseReference{RefSchemaID: s.ID}
	return rr.DelByRefSchemaID(ctx)
}

func DerefSchemaFromCollections(ctx context.Context, s *definition.DefinitionSchema, collectionIDs []uint, deref bool) error {
	collections, err := collection.GetCollections(ctx, s.ProjectID, collectionIDs...)
	if err != nil {
		return err
	}

	for _, c := range collections {
		c.DelRefSchema(ctx, s, deref)
	}

	cr := &referencerelationship.CollectionReference{RefID: s.ID, RefType: referencerelationship.ReferenceSchema}
	return cr.DelByRef(ctx)
}

func ClearSchemaRef(ctx context.Context, s *definition.DefinitionSchema) error {
	sr := &referencerelationship.SchemaReference{RefSchemaID: s.ID}
	schemaIDs, err := sr.GetSchemaIDsByRef(ctx)
	if err != nil {
		return err
	}
	if err := DerefSchemaFromSchemas(ctx, s, schemaIDs, false); err != nil {
		return err
	}

	rr := &referencerelationship.ResponseReference{RefSchemaID: s.ID}
	responseIDs, err := rr.GetResponseIDsByRef(ctx)
	if err != nil {
		return err
	}
	if err := DerefSchemaFromResponses(ctx, s, responseIDs, false); err != nil {
		return err
	}

	cr := &referencerelationship.CollectionReference{RefID: s.ID, RefType: referencerelationship.ReferenceSchema}
	collectionIDs, err := cr.GetCollectionIDsByRef(ctx)
	if err != nil {
		return err
	}
	if err := DerefSchemaFromCollections(ctx, s, collectionIDs, false); err != nil {
		return err
	}

	return DeleteSchemaReference(ctx, s)
}

func UnpackSchemaRef(ctx context.Context, s *definition.DefinitionSchema) error {
	sr := &referencerelationship.SchemaReference{RefSchemaID: s.ID}
	schemaIDs, err := sr.GetSchemaIDsByRef(ctx)
	if err != nil {
		return err
	}
	if err := DerefSchemaFromSchemas(ctx, s, schemaIDs, true); err != nil {
		return err
	}

	rr := &referencerelationship.ResponseReference{RefSchemaID: s.ID}
	responseIDs, err := rr.GetResponseIDsByRef(ctx)
	if err != nil {
		return err
	}
	if err := DerefSchemaFromResponses(ctx, s, responseIDs, true); err != nil {
		return err
	}

	cr := &referencerelationship.CollectionReference{RefID: s.ID, RefType: referencerelationship.ReferenceSchema}
	collectionIDs, err := cr.GetCollectionIDsByRef(ctx)
	if err != nil {
		return err
	}
	if err := DerefSchemaFromCollections(ctx, s, collectionIDs, true); err != nil {
		return err
	}

	sr = &referencerelationship.SchemaReference{SchemaID: s.ID}
	schemaRefRecords, err := sr.GetSchemaRefs(ctx)
	if err != nil {
		return err
	}

	// 引用目标schema的schemas、responses、collections 与 目标schema的所引用的schema 建立引用关系
	collectionWantPush := make([]*referencerelationship.CollectionReference, 0)
	schemaWantPush := make([]*referencerelationship.SchemaReference, 0)
	responseWantPush := make([]*referencerelationship.ResponseReference, 0)
	for _, v := range schemaRefRecords {
		for _, s := range schemaIDs {
			schemaWantPush = append(schemaWantPush, &referencerelationship.SchemaReference{
				SchemaID:    s,
				RefSchemaID: v.RefSchemaID,
			})
		}
		for _, r := range responseIDs {
			responseWantPush = append(responseWantPush, &referencerelationship.ResponseReference{
				ResponseID:  r,
				RefSchemaID: v.RefSchemaID,
			})
		}
		for _, c := range collectionIDs {
			collectionWantPush = append(collectionWantPush, &referencerelationship.CollectionReference{
				CollectionID: c,
				RefID:        v.RefSchemaID,
				RefType:      referencerelationship.ReferenceSchema,
			})
		}
	}
	if err := referencerelationship.BatchCreateSchemaReference(ctx, schemaWantPush); err != nil {
		return err
	}
	if err := referencerelationship.BatchCreateResponseReference(ctx, responseWantPush); err != nil {
		return err
	}

	if err := referencerelationship.BatchCreateCollectionReference(ctx, collectionWantPush); err != nil {
		return err
	}
	return DeleteSchemaReference(ctx, s)
}
