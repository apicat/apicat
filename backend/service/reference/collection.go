package reference

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model/collection"
	referencerelation "github.com/apicat/apicat/v2/backend/model/reference_relation"
	"github.com/apicat/apicat/v2/backend/service/except"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
)

// UpdateCollectionRef 更新集合引用关系（包含集合引用、排除集合）
func UpdateCollectionRef(ctx context.Context, c *collection.Collection, oldSchemaIDs, oldResponseIDs, oldExceptParamIDs []uint) error {
	newResponseIDs, err := ParseRefResponsesFromCollection(c)
	if err != nil {
		return err
	}
	if err := updateRefCollectionToResponses(ctx, c.ID, oldResponseIDs, newResponseIDs); err != nil {
		return err
	}

	newSchemaIDs, err := ParseRefSchemasFromCollection(c)
	if err != nil {
		return err
	}
	if err := updateRefCollectionToSchemas(ctx, c.ID, oldSchemaIDs, newSchemaIDs); err != nil {
		return err
	}

	newParamIDs, err := except.ParseExceptParamsFromCollection(c)
	if err != nil {
		return err
	}
	if err := updateExceptParamsToCollection(ctx, c.ID, oldExceptParamIDs, newParamIDs); err != nil {
		return err
	}

	return nil
}

// updateRefCollectionToResponses 更新集合引用响应的引用关系
func updateRefCollectionToResponses(ctx context.Context, cID uint, oldResponseIDs, newResponseIDs []uint) error {
	oldRefs, err := referencerelation.GetRefResponseCollections(ctx, oldResponseIDs...)
	if err != nil {
		return err
	}

	oldRefsMap := make(map[uint]*referencerelation.RefResponseCollections, 0)
	for _, v := range oldRefs {
		oldRefsMap[v.RefResponserID] = v
	}

	wantPop := make([]uint, 0)
	wantPush := make([]*referencerelation.RefResponseCollections, 0)

	for key, value := range oldRefsMap {
		if !arrutil.InArray(key, newResponseIDs) {
			wantPop = append(wantPop, value.ID)
		}
	}

	for _, v := range newResponseIDs {
		if _, ok := oldRefsMap[v]; !ok {
			wantPush = append(wantPush, &referencerelation.RefResponseCollections{
				CollectionID:   cID,
				RefResponserID: v,
			})
		}
	}

	if err := referencerelation.BatchCreateRefResponseCollections(ctx, wantPush); err != nil {
		return err
	}
	return referencerelation.BatchDelRefResponseCollections(ctx, wantPop...)
}

// updateRefCollectionToSchemas 更新集合引用公共模型的引用关系
func updateRefCollectionToSchemas(ctx context.Context, cID uint, oldScheamIDs, newSchemaIDs []uint) error {
	oldRefs, err := referencerelation.GetRefSchemaCollections(ctx, oldScheamIDs...)
	if err != nil {
		return err
	}

	oldRefsMap := make(map[uint]*referencerelation.RefSchemaCollections, 0)
	for _, v := range oldRefs {
		oldRefsMap[v.RefSchemaID] = v
	}

	wantPop := make([]uint, 0)
	wantPush := make([]*referencerelation.RefSchemaCollections, 0)

	for key, value := range oldRefsMap {
		if !arrutil.InArray(key, newSchemaIDs) {
			wantPop = append(wantPop, value.ID)
		}
	}

	for _, v := range newSchemaIDs {
		if _, ok := oldRefsMap[v]; !ok {
			wantPush = append(wantPush, &referencerelation.RefSchemaCollections{
				CollectionID: cID,
				RefSchemaID:  v,
			})
		}
	}

	if err := referencerelation.BatchCreateRefSchemaCollections(ctx, wantPush); err != nil {
		return err
	}
	return referencerelation.BatchDelRefSchemaCollections(ctx, wantPop...)
}

// updateExceptParamsToCollection 更新集合被全局参数排除关系
func updateExceptParamsToCollection(ctx context.Context, cID uint, oldParamIDs, newParamIDs []uint) error {
	oldRefs, err := referencerelation.GetExceptParamCollections(ctx, oldParamIDs...)
	if err != nil {
		return err
	}

	oldRefsMap := make(map[uint]*referencerelation.ExceptParamCollection, 0)
	for _, v := range oldRefs {
		oldRefsMap[v.ExceptParamID] = v
	}

	wantPop := make([]uint, 0)
	wantPush := make([]*referencerelation.ExceptParamCollection, 0)

	for key, value := range oldRefsMap {
		if !arrutil.InArray(key, newParamIDs) {
			wantPop = append(wantPop, value.ID)
		}
	}

	for _, v := range newParamIDs {
		if _, ok := oldRefsMap[v]; !ok {
			wantPush = append(wantPush, &referencerelation.ExceptParamCollection{
				CollectionID:  cID,
				ExceptParamID: v,
			})
		}
	}

	if err := referencerelation.BatchCreateExceptParamCollections(ctx, wantPush); err != nil {
		return err
	}
	return referencerelation.BatchDelExceptParamCollections(ctx, wantPop...)
}
