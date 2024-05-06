package reference

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model/collection"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
)

func UpdateCollectionRef(ctx context.Context, c *collection.Collection) error {
	if err := updateParamExcept(ctx, c); err != nil {
		return err
	}

	return updateReference(ctx, c)
}

func updateReference(ctx context.Context, c *collection.Collection) error {
	cr := &referencerelationship.CollectionReference{CollectionID: c.ID}
	lastRefs, err := cr.GetCollectionRefs(ctx)
	if err != nil {
		return err
	}

	lastSchemaRefsMap := make(map[uint]*referencerelationship.CollectionReference, 0)
	lastResponseRefsMap := make(map[uint]*referencerelationship.CollectionReference, 0)
	for _, v := range lastRefs {
		switch v.RefType {
		case referencerelationship.ReferenceSchema:
			lastSchemaRefsMap[v.RefID] = v
		case referencerelationship.ReferenceResponse:
			lastResponseRefsMap[v.RefID] = v
		}
	}

	schemaWantPop, schemaWantPush := updateSchemaRef(c, lastSchemaRefsMap)
	responseWantPop, responseWantPush := updateResponseRef(c, lastResponseRefsMap)

	wantPop := append(schemaWantPop, responseWantPop...)
	wantPush := append(schemaWantPush, responseWantPush...)

	if err := referencerelationship.BatchCreateCollectionReference(ctx, wantPush); err != nil {
		return err
	}

	return referencerelationship.BatchDeleteCollectionReference(ctx, wantPop...)
}

func updateParamExcept(ctx context.Context, c *collection.Collection) error {
	pe := referencerelationship.ParameterExcept{ExceptCollectionID: c.ID}

	lastExcepts, err := pe.GetParameterExcepts(ctx)
	if err != nil {
		return err
	}

	lastExceptsMap := make(map[uint]*referencerelationship.ParameterExcept, 0)
	for _, v := range lastExcepts {
		lastExceptsMap[v.ParameterID] = v
	}

	latestExcepts := ParseExceptParameterFromCollection(c)

	wantPop := make([]uint, 0)
	wantPush := make([]*referencerelationship.ParameterExcept, 0)

	for key, value := range lastExceptsMap {
		if !arrutil.InArray[uint](key, latestExcepts) {
			wantPop = append(wantPop, value.ID)
		}
	}

	for _, v := range latestExcepts {
		if _, ok := lastExceptsMap[v]; !ok {
			wantPush = append(wantPush, &referencerelationship.ParameterExcept{
				ParameterID:        v,
				ExceptCollectionID: c.ID,
			})
		}
	}

	if err := referencerelationship.BatchCreateParameterExcept(ctx, wantPush); err != nil {
		return err
	}

	return referencerelationship.BatchDeleteParameterExcept(ctx, wantPop...)
}

func updateSchemaRef(c *collection.Collection, lastRefsMap map[uint]*referencerelationship.CollectionReference) (wantPop []uint, wantPush []*referencerelationship.CollectionReference) {
	latestSchemaRefs := ParseRefSchemas(c.Content)

	for key, value := range lastRefsMap {
		if !arrutil.InArray[uint](key, latestSchemaRefs) {
			wantPop = append(wantPop, value.ID)
		}
	}

	for _, v := range latestSchemaRefs {
		if _, ok := lastRefsMap[v]; !ok {
			wantPush = append(wantPush, &referencerelationship.CollectionReference{
				CollectionID: c.ID,
				RefID:        v,
				RefType:      referencerelationship.ReferenceSchema,
			})
		}
	}

	return wantPop, wantPush
}

func updateResponseRef(c *collection.Collection, lastRefsMap map[uint]*referencerelationship.CollectionReference) (wantPop []uint, wantPush []*referencerelationship.CollectionReference) {
	lastResponseRefs := ParseRefResponses(c.Content)

	for key, value := range lastRefsMap {
		if !arrutil.InArray[uint](key, lastResponseRefs) {
			wantPop = append(wantPop, value.ID)
		}
	}

	for _, v := range lastResponseRefs {
		if _, ok := lastRefsMap[v]; !ok {
			wantPush = append(wantPush, &referencerelationship.CollectionReference{
				CollectionID: c.ID,
				RefID:        v,
				RefType:      referencerelationship.ReferenceResponse,
			})
		}
	}

	return wantPop, wantPush
}
