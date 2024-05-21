package reference

import (
	"context"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	referencerelation "github.com/apicat/apicat/v2/backend/model/reference_relation"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
)

// UpdateResponseRef 更新公共响应引用关系
func UpdateResponseRef(ctx context.Context, r *definition.DefinitionResponse, oldScheamIDs []uint) error {
	newSchemaIDs, err := ParseRefSchemasFromResponse(r)
	if err != nil {
		return err
	}
	if err := updateRefResponseToScheams(ctx, r.ID, oldScheamIDs, newSchemaIDs); err != nil {
		return err
	}

	// 预留response还会引用其他公共组建的情况
	return nil
}

// updateRefResponseToScheam 更新公共响应引用公共模型的引用关系
func updateRefResponseToScheams(ctx context.Context, rID uint, oldScheamIDs, newSchemaIDs []uint) error {
	oldRefs, err := referencerelation.GetRefSchemaResponses(ctx, oldScheamIDs...)
	if err != nil {
		return err
	}

	oldRefsMap := make(map[uint]*referencerelation.RefSchemaResponses, 0)
	for _, v := range oldRefs {
		oldRefsMap[v.RefSchemaID] = v
	}

	wantPop := make([]uint, 0)
	wantPush := make([]*referencerelation.RefSchemaResponses, 0)

	// 删除老引用关系中存在但当前引用中不存在的引用
	for key, value := range oldRefsMap {
		if !arrutil.InArray(key, newSchemaIDs) {
			wantPop = append(wantPop, value.ID)
		}
	}

	// 添加老引用关系中不存在但当前引用中存在的引用
	for _, v := range newSchemaIDs {
		if _, ok := oldRefsMap[v]; !ok {
			wantPush = append(wantPush, &referencerelation.RefSchemaResponses{
				ResponseID:  rID,
				RefSchemaID: v,
			})
		}
	}

	if err := referencerelation.BatchCreateRefSchemaResponses(ctx, wantPush); err != nil {
		return err
	}
	return referencerelation.BatchDelRefSchemaResponses(ctx, wantPop...)
}

// DerefResponse 公共响应解引用
func DerefResponse(ctx context.Context, r *definition.DefinitionResponse, deref bool) error {
	nowSchemaIDs, err := ParseRefSchemasFromResponse(r)
	if err != nil {
		return err
	}

	rc := referencerelation.RefResponseCollections{RefResponserID: r.ID}
	cIDs, err := rc.GetCollectionIDs(ctx)
	if err != nil {
		return err
	}

	// 在collection中解引用response
	if err := derefResponseFromCollections(ctx, r, cIDs, deref); err != nil {
		return err
	}

	if deref {
		// 建立引用自身的(collections -> self 中的collections)与自身引用的(self -> schemas 中的schemas)之间的引用关系
		if err := linkResponseRefContextParentCWithChildS(ctx, cIDs, nowSchemaIDs); err != nil {
			return err
		}
	}

	// 清除引用关系(collections -> self)
	if err := clearRefCollectionsToResponse(ctx, r.ID); err != nil {
		return err
	}

	// 清除引用关系(self -> scheams)
	return clearRefResponseToSchemas(ctx, r.ID, nowSchemaIDs)
}

// derefResponseFromCollections 解开collection中引用response的地方
func derefResponseFromCollections(ctx context.Context, r *definition.DefinitionResponse, collectionIDs []uint, deref bool) error {
	collections, err := collection.GetCollections(ctx, r.ProjectID, collectionIDs...)
	if err != nil {
		return err
	}

	for _, c := range collections {
		if err := c.DelRefResponse(ctx, r, deref); err != nil {
			slog.ErrorContext(ctx, "c.DelRefResponse", "err", err)
		}
	}

	return nil
}

// clearRefCollectionsToResponse 清除collections引用response的引用关系
func clearRefCollectionsToResponse(ctx context.Context, rID uint) error {
	return referencerelation.DelRefResponseCollections(ctx, rID)
}

// clearRefResponseToSchemas 清除response引用schemas的引用关系
func clearRefResponseToSchemas(ctx context.Context, rID uint, schemaIDs []uint) error {
	return referencerelation.DelRefSchemaResponse(ctx, rID, schemaIDs...)
}

// linkResponseRefContextParentCWithChildS 建立respnose引用父级collections与引用子集的schemas之间的引用关系
func linkResponseRefContextParentCWithChildS(ctx context.Context, collectionIDs, schemaIDs []uint) error {
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
