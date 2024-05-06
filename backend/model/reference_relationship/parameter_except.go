package referencerelationship

import (
	"context"
	"errors"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
)

type ParameterExcept struct {
	ID                 uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	ParameterID        uint `gorm:"type:bigint;index;not null;comment:全局参数id"`
	ExceptCollectionID uint `gorm:"type:bigint;index;not null;comment:排除集合id"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func BatchCreateParameterExcept(ctx context.Context, list []*ParameterExcept) error {
	if len(list) == 0 {
		return nil
	}
	return model.DB(ctx).Create(list).Error
}

func BatchDeleteParameterExcept(ctx context.Context, ids ...uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Where("id in ?", ids).Delete(&ParameterExcept{}).Error
}

func GetParameterExceptsByParameter(ctx context.Context, projectID string, parameterID uint) ([]*ParameterExcept, error) {
	var list []*ParameterExcept
	tx := model.DB(ctx).Where("project_id = ? AND parameter_id = ?", projectID, parameterID).Find(&list)
	return list, tx.Error
}

func GetParameterExceptsByCollection(ctx context.Context, projectID string, exceptCollectionID uint) ([]*ParameterExcept, error) {
	var list []*ParameterExcept
	tx := model.DB(ctx).Where("project_id = ? AND except_collection_id = ?", projectID, exceptCollectionID).Find(&list)
	return list, tx.Error
}

func (p *ParameterExcept) GetParameterExcepts(ctx context.Context) ([]*ParameterExcept, error) {
	var list []*ParameterExcept

	tx := model.DB(ctx)
	if p.ParameterID != 0 {
		tx = tx.Where("parameter_id = ?", p.ParameterID)
	} else if p.ExceptCollectionID != 0 {
		tx = tx.Where("except_collection_id = ?", p.ExceptCollectionID)
	} else {
		return nil, errors.New("query condition error")
	}

	return list, tx.Find(&list).Error
}

func (p *ParameterExcept) GetParameterExceptCIDs(ctx context.Context) ([]uint, error) {
	var ids []uint
	tx := model.DB(ctx).Model(p).Where("parameter_id = ?", p.ParameterID).Select("except_collection_id").Find(&ids)
	return ids, tx.Error
}
