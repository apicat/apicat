package referencerelationship

import (
	"context"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
)

type ParameterExcept struct {
	ID                 uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID          string `gorm:"type:varchar(24);index;not null;comment:项目id"`
	ParameterID        uint   `gorm:"type:bigint;index;not null;comment:全局参数id"`
	ExceptCollectionID uint   `gorm:"type:bigint;index;not null;comment:排除集合id"`
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
