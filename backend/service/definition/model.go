package definition

import (
	"context"
	"errors"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/core/content_suggestion"
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/definition"
	referencerelation "github.com/apicat/apicat/v2/backend/model/reference_relation"
	"github.com/apicat/apicat/v2/backend/module/vector"
	"github.com/apicat/apicat/v2/backend/service/collection"
)

type DefinitionModelService struct {
	ctx context.Context
}

func NewDefinitionModelService(ctx context.Context) *DefinitionModelService {
	return &DefinitionModelService{ctx: ctx}
}

// 获取所有直接或间接引用了模型 id 为入参的模型 ID
func (dms *DefinitionModelService) GetRefModelIDs(id uint) ([]uint, error) {
	if id < 1 {
		return nil, errors.New("invalid model id")
	}

	var ids []uint
	if err := model.DB(dms.ctx).Model(&referencerelation.RefSchemaSchemas{}).Where("ref_schema_id = ?", id).Select("schema_id").Scan(&ids).Error; err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	result := make([]uint, 0)
	for _, v := range ids {
		result = append(result, v)
		if tmp, err := dms.GetRefModelIDs(v); err != nil {
			return nil, err
		} else {
			if len(tmp) > 0 {
				result = append(result, tmp...)
			}
		}
	}
	return result, nil
}

// 获取所有引用了模型 id 为入参的 Collection ID
func (dms *DefinitionModelService) GetRefCollectionIDs(id uint) ([]uint, error) {
	if id < 1 {
		return nil, errors.New("invalid model id")
	}

	var ids []uint
	if err := model.DB(dms.ctx).Model(&referencerelation.RefSchemaCollections{}).Where("ref_schema_id = ?", id).Select("collection_id").Scan(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (dms *DefinitionModelService) UpdateVector(dm *definition.DefinitionSchema) {
	if dm == nil {
		slog.ErrorContext(dms.ctx, "UpdateVector", "err", errors.New("definition schema is nil"))
		return
	}
	if dm.ProjectID == "" {
		slog.ErrorContext(dms.ctx, "UpdateVector", "err", errors.New("project id is empty"))
		return
	}

	modelVector, err := content_suggestion.NewDefinitionModelVector(dm.ProjectID)
	if err != nil {
		slog.ErrorContext(dms.ctx, "content_suggestion.NewDefinitionModelVector", "err", err)
	}
	modelVector.CreateLater(dm.ID)

	// 更新所有引用了该模型的集合的向量
	collectionService := collection.NewCollectionService(dms.ctx)
	if refCollectionIDs, err := dms.GetRefCollectionIDs(dm.ID); err == nil && len(refCollectionIDs) > 0 {
		collectionService.CreateVector(dm.ProjectID, refCollectionIDs...)
	}

	// 更新所有引用了该模型的模型的向量
	if refModelIDs, err := dms.GetRefModelIDs(dm.ID); err == nil && len(refModelIDs) > 0 {
		for _, v := range refModelIDs {
			modelVector.CreateLater(v)
			if refCollectionIDs, err := dms.GetRefCollectionIDs(v); err == nil && len(refCollectionIDs) > 0 {
				collectionService.CreateVector(dm.ProjectID, refCollectionIDs...)
			}
		}
	}
}

func (dms *DefinitionModelService) DelVector(dm *definition.DefinitionSchema) {
	if dm == nil {
		slog.ErrorContext(dms.ctx, "DelVector", "err", errors.New("definition schema is nil"))
		return
	}
	if dm.ProjectID == "" {
		slog.ErrorContext(dms.ctx, "DelVector", "err", errors.New("project id is empty"))
		return
	}
	if dm.VectorID == "" {
		return
	}

	if vectorDB, err := vector.NewVector(config.GetVector().ToModuleStruct()); err != nil {
		slog.ErrorContext(dms.ctx, "vector.NewVector", "err", err)
	} else {
		if err := vectorDB.DeleteObject(dm.ProjectID, dm.VectorID); err != nil {
			slog.ErrorContext(dms.ctx, "vectorDB.DeleteObject", "err", err)
		}
	}

	// 更新所有引用了该模型的集合的向量
	collectionService := collection.NewCollectionService(dms.ctx)
	if refCollectionIDs, err := dms.GetRefCollectionIDs(dm.ID); err == nil && len(refCollectionIDs) > 0 {
		collectionService.CreateVector(dm.ProjectID, refCollectionIDs...)
	}

	// 更新所有引用了该模型的模型的向量
	if refModelIDs, err := dms.GetRefModelIDs(dm.ID); err == nil && len(refModelIDs) > 0 {
		modelVector, err := content_suggestion.NewDefinitionModelVector(dm.ProjectID)
		if err != nil {
			slog.ErrorContext(dms.ctx, "content_suggestion.NewDefinitionModelVector", "err", err)
			return
		}

		for _, v := range refModelIDs {
			modelVector.CreateLater(v)
			if refCollectionIDs, err := dms.GetRefCollectionIDs(v); err == nil && len(refCollectionIDs) > 0 {
				collectionService.CreateVector(dm.ProjectID, refCollectionIDs...)
			}
		}
	}
}
