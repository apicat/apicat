package reference

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
)

// UpdateResponseRef 更新公共响应引用关系
func UpdateResponseRef(ctx context.Context, r *definition.DefinitionResponse) error {
	rr := &referencerelationship.ResponseReference{ResponseID: r.ID}
	lastRefs, err := rr.GetResponseRefs(ctx)
	if err != nil {
		return err
	}

	lastRefsMap := make(map[uint]*referencerelationship.ResponseReference, 0)
	for _, v := range lastRefs {
		lastRefsMap[v.RefSchemaID] = v
	}

	latestRefs := ParseRefSchemas(r.Content)

	wantPop := make([]uint, 0)
	wantPush := make([]*referencerelationship.ResponseReference, 0)

	// 删除老引用关系中存在但当前引用中不存在的引用
	for key, value := range lastRefsMap {
		if !arrutil.InArray[uint](key, latestRefs) {
			wantPop = append(wantPop, value.ID)
		}
	}

	// 添加老引用关系中不存在但当前引用中存在的引用
	for _, v := range latestRefs {
		if _, ok := lastRefsMap[v]; !ok {
			wantPush = append(wantPush, &referencerelationship.ResponseReference{
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

// DerefResponse 公共响应解引用
func DerefResponse(ctx context.Context, r *definition.DefinitionResponse, deref bool) error {
	if deref {
		return unpackResponseRef(ctx, r)
	} else {
		return clearResponseRef(ctx, r)
	}
}

// clearResponseRef 清除公共响应
func clearResponseRef(ctx context.Context, r *definition.DefinitionResponse) error {
	// 清除被引用关系(collections -> self)
	cr := &referencerelationship.CollectionReference{RefID: r.ID, RefType: referencerelationship.ReferenceResponse}
	collectionIDs, err := cr.GetCollectionIDsByRef(ctx)
	if err != nil {
		return err
	}

	if err := derefResponseFromCollections(ctx, r, collectionIDs, false); err != nil {
		return err
	}

	// 清除引用关系(self -> scheams)
	return deleteResponseReference(ctx, r)
}

// unpackResponseRef 展开公共响应
func unpackResponseRef(ctx context.Context, r *definition.DefinitionResponse) error {
	// 清除被引用关系(collections -> self)
	cr := &referencerelationship.CollectionReference{RefID: r.ID, RefType: referencerelationship.ReferenceResponse}
	collectionIDs, err := cr.GetCollectionIDsByRef(ctx)
	if err != nil {
		return err
	}

	if err := derefResponseFromCollections(ctx, r, collectionIDs, true); err != nil {
		return err
	}

	// 建立自身引用的(self -> schemas 中的schemas)与引用自身的(collections -> self 中的collections)之间的引用关系
	rr := referencerelationship.ResponseReference{ResponseID: r.ID}
	responseRefRecords, err := rr.GetResponseRefs(ctx)
	if err != nil {
		return err
	}

	collectionWantPush := make([]*referencerelationship.CollectionReference, 0)
	for _, v := range responseRefRecords {
		for _, c := range collectionIDs {
			collectionWantPush = append(collectionWantPush, &referencerelationship.CollectionReference{
				CollectionID: c,
				RefID:        v.RefSchemaID,
				RefType:      referencerelationship.ReferenceSchema,
			})
		}
	}
	if err := referencerelationship.BatchCreateCollectionReference(ctx, collectionWantPush); err != nil {
		return err
	}

	// 清除引用关系(self -> scheams)
	return deleteResponseReference(ctx, r)
}

// derefResponseFromCollections 从集合中解引用公共响应
func derefResponseFromCollections(ctx context.Context, r *definition.DefinitionResponse, collectionIDs []uint, deref bool) error {
	collections, err := collection.GetCollections(ctx, r.ProjectID, collectionIDs...)
	if err != nil {
		return err
	}

	for _, c := range collections {
		c.DelRefResponse(ctx, r, deref)
	}

	cr := &referencerelationship.CollectionReference{RefID: r.ID, RefType: referencerelationship.ReferenceResponse}
	return cr.DelByRef(ctx)
}

// deleteResponseReference 删除公共响应引用关系
func deleteResponseReference(ctx context.Context, r *definition.DefinitionResponse) error {
	rr := referencerelationship.ResponseReference{ResponseID: r.ID}
	refs, err := rr.GetResponseRefs(ctx)
	if err != nil {
		return err
	}

	var ids []uint
	for _, item := range refs {
		ids = append(ids, item.ID)
	}

	return referencerelationship.BatchDeleteResponseReference(ctx, ids...)
}
