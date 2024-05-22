package reference

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	referencerelation "github.com/apicat/apicat/v2/backend/model/reference_relation"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
)

// UpdateSchemaRef 更新公共模型引用关系
// oldScheamIDs 需要更新的shcema之前引用的公共模型id
func UpdateSchemaRef(ctx context.Context, s *definition.DefinitionSchema, oldScheamIDs []uint) error {
	newSchemaIDs, err := ParseRefSchemasFromSchema(s)
	if err != nil {
		return err
	}
	if err := updateRefSchemaToschemas(ctx, s.ID, oldScheamIDs, newSchemaIDs); err != nil {
		return err
	}

	// 预留schema还会引用其他公共组建的情况
	return nil
}

// updateRefSchemaToschema 更新schema引用schemas的引用关系
func updateRefSchemaToschemas(ctx context.Context, sID uint, oldSchemaIDs, newSchemaIDs []uint) error {
	oldRefs, err := referencerelation.GetRefSchemaSchemas(ctx, oldSchemaIDs...)
	if err != nil {
		return err
	}

	oldRefsMap := make(map[uint]*referencerelation.RefSchemaSchemas, 0)
	for _, v := range oldRefs {
		oldRefsMap[v.RefSchemaID] = v
	}

	wantPop := make([]uint, 0)
	wantPush := make([]*referencerelation.RefSchemaSchemas, 0)

	// 删除老引用关系中存在但当前引用中不存在的引用
	for key, value := range oldRefsMap {
		if !arrutil.InArray(key, newSchemaIDs) {
			wantPop = append(wantPop, value.ID)
		}
	}

	// 添加老引用关系中不存在但当前引用中存在的引用
	for _, v := range newSchemaIDs {
		if _, ok := oldRefsMap[v]; !ok {
			wantPush = append(wantPush, &referencerelation.RefSchemaSchemas{
				SchemaID:    sID,
				RefSchemaID: v,
			})
		}
	}

	if err := referencerelation.BatchCreateRefSchemaSchemas(ctx, wantPush); err != nil {
		return err
	}
	return referencerelation.BatchDelRefSchemaSchemas(ctx, wantPop...)
}

func DerefSchema(ctx context.Context, s *definition.DefinitionSchema, deref bool) error {
	refSchemaIDs, err := ParseRefSchemasFromSchema(s)
	if err != nil {
		return err
	}

	rc := referencerelation.RefSchemaCollections{RefSchemaID: s.ID}
	cIDs, err := rc.GetCollectionIDs(ctx)
	if err != nil {
		return err
	}

	// 在collection中解引用schema
	if err := derefSchemaFromCollections(ctx, s, cIDs, deref); err != nil {
		return err
	}

	rr := referencerelation.RefSchemaResponses{RefSchemaID: s.ID}
	rIDs, err := rr.GetResponseIDs(ctx)
	if err != nil {
		return err
	}

	// 在response中解引用schema
	if err := derefSchemaFromResponses(ctx, s, rIDs, deref); err != nil {
		return err
	}

	rs := referencerelation.RefSchemaSchemas{RefSchemaID: s.ID}
	schemaIDs, err := rs.GetSchemaIDs(ctx)
	if err != nil {
		return err
	}

	// 在schema中解引用schema
	if err := derefSchemaFromSchemas(ctx, s, schemaIDs, deref); err != nil {
		return err
	}

	if deref {
		// 建立引用自身的(collections -> self 中的collections)与自身引用的(self -> schemas 中的schemas)之间的引用关系
		if err := linkSchemaRefContextParentCWithChildS(ctx, cIDs, refSchemaIDs); err != nil {
			return err
		}

		// 建立引用自身的(responses -> self 中的responses)与自身引用的(self -> schemas 中的schemas)之间的引用关系
		if err := linkSchemaRefContextParentRWithChildS(ctx, rIDs, refSchemaIDs); err != nil {
			return err
		}

		// 建立引用自身的(schemas -> self 中的schemas)与自身引用的(self -> schemas 中的schemas)之间的引用关系
		if err := linkSchemaRefContextParentSWithChildS(ctx, schemaIDs, refSchemaIDs); err != nil {
			return err
		}
	}

	// 清除引用关系(collections -> self)
	if err := clearRefCollectionsToSchema(ctx, s.ID); err != nil {
		return err
	}

	// 清除引用关系(responses -> self)
	if err := clearRefResponsesToSchema(ctx, s.ID); err != nil {
		return err
	}

	// 清除引用关系(schemas -> self)
	if err := clearRefSchemasToSchema(ctx, s.ID); err != nil {
		return err
	}

	// 清除引用关系(self -> scheams)
	return clearRefSchemaToSchemas(ctx, refSchemaIDs, s.ID)
}

// derefSchemaFromSchemas 从公共模型中解引用公共模型
func derefSchemaFromSchemas(ctx context.Context, s *definition.DefinitionSchema, schemaIDs []uint, deref bool) error {
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

// derefSchemaFromResponses 从公共响应中解引用公共模型
func derefSchemaFromResponses(ctx context.Context, s *definition.DefinitionSchema, responseIDs []uint, deref bool) error {
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

// derefSchemaFromCollections 从集合中解引用公共模型
func derefSchemaFromCollections(ctx context.Context, s *definition.DefinitionSchema, collectionIDs []uint, deref bool) error {
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

// clearRefCollectionsToSchema 清除collections引用schema的引用关系
func clearRefCollectionsToSchema(ctx context.Context, sID uint) error {
	return referencerelation.DelRefSchemaCollections(ctx, sID)
}

// clearRefResponsesToSchema 清除responses引用schema的引用关系
func clearRefResponsesToSchema(ctx context.Context, sID uint) error {
	return referencerelation.DelRefSchemaResponses(ctx, sID)
}

// clearRefSchemasToSchema 清除schemas引用schema的引用关系
func clearRefSchemasToSchema(ctx context.Context, sID uint) error {
	return referencerelation.DelRefSchemaSchemas(ctx, sID)
}

// clearRefSchemaToSchemas 清除schema引用schemas的引用关系
func clearRefSchemaToSchemas(ctx context.Context, schemaIDs []uint, sID uint) error {
	return referencerelation.DelRefSchemaSchema(ctx, sID, schemaIDs...)
}

// linkSchemaRefContextParentCWithChildS 建立schema引用父级collections与引用子集的schemas之间的引用关系
func linkSchemaRefContextParentCWithChildS(ctx context.Context, collectionIDs, schemaIDs []uint) error {
	wantPush := make([]*referencerelation.RefSchemaCollections, 0)
	for _, v := range schemaIDs {
		for _, c := range collectionIDs {
			wantPush = append(wantPush, &referencerelation.RefSchemaCollections{
				RefSchemaID:  v,
				CollectionID: c,
			})
		}
	}
	return referencerelation.BatchCreateRefSchemaCollections(ctx, wantPush)
}

// linkSchemaRefContextParentRWithChildS 建立schema引用父级responses与引用子集的schemas之间的引用关系
func linkSchemaRefContextParentRWithChildS(ctx context.Context, responseIDs, schemaIDs []uint) error {
	wantPush := make([]*referencerelation.RefSchemaResponses, 0)
	for _, v := range schemaIDs {
		for _, r := range responseIDs {
			wantPush = append(wantPush, &referencerelation.RefSchemaResponses{
				RefSchemaID: v,
				ResponseID:  r,
			})
		}
	}
	return referencerelation.BatchCreateRefSchemaResponses(ctx, wantPush)
}

// linkSchemaRefContextParentSWithChildS 建立schema引用父级schemas与引用子集的schemas之间的引用关系
func linkSchemaRefContextParentSWithChildS(ctx context.Context, schemaIDs, childSchemaIDs []uint) error {
	wantPush := make([]*referencerelation.RefSchemaSchemas, 0)
	for _, v := range childSchemaIDs {
		for _, s := range schemaIDs {
			wantPush = append(wantPush, &referencerelation.RefSchemaSchemas{
				RefSchemaID: v,
				SchemaID:    s,
			})
		}
	}
	return referencerelation.BatchCreateRefSchemaSchemas(ctx, wantPush)
}
