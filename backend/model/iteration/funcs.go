package iteration

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"

	"gorm.io/gorm"
)

func GetIterations(ctx context.Context, teamID string, page, pageSize int, pIDs ...string) ([]*Iteration, error) {
	tx := model.DB(ctx).Where("team_id = ?", teamID)
	if len(pIDs) > 0 {
		tx = tx.Where("project_id in (?)", pIDs)
	}
	if page > 0 && pageSize > 0 {
		tx = tx.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	var iterations []*Iteration
	err := tx.Find(&iterations).Error
	return iterations, err
}

func GetIterationsCount(ctx context.Context, teamID string, pIDs ...string) (int64, error) {
	tx := model.DB(ctx)
	tx = tx.Model(&Iteration{}).Where("team_id = ?", teamID)
	if len(pIDs) > 0 {
		tx = tx.Where("project_id in (?)", pIDs)
	}
	var count int64
	err := tx.Count(&count).Error
	return count, err
}

func GetIterationApi(ctx context.Context, i *Iteration) ([]*IterationApi, error) {
	var iterationApi []*IterationApi
	return iterationApi, model.DB(ctx).Where("iteration_id = ?", i.ID).Find(&iterationApi).Error
}

func BatchDeleteIterationApi(ctx context.Context, cIDs ...uint) error {
	if len(cIDs) == 0 {
		return nil
	}
	return model.DB(ctx).Delete(&IterationApi{}, "collection_id IN (?)", cIDs).Error
}

func RestoreIterationApi(ctx context.Context, cIDs []uint) error {
	if len(cIDs) == 0 {
		return nil
	}

	return model.DB(ctx).Unscoped().Model(&IterationApi{}).Where("collection_id IN (?)", cIDs).Update("deleted_at", gorm.Expr("NULL")).Error
}
