package iteration

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
)

type IterationApi struct {
	ID             uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	IterationID    string `gorm:"type:varchar(24);index;not null;comment:iteration id"`
	CollectionID   uint   `gorm:"type:bigint;not null;comment:collection id"`
	CollectionType string `gorm:"type:varchar(255);not null;comment:collection type:category,doc,http"`
	model.TimeModel
}

func (ia *IterationApi) Create(ctx context.Context) error {
	return model.DB(ctx).Create(ia).Error
}

func (ia *IterationApi) Delete(ctx context.Context) error {
	return model.DB(ctx).Delete(ia).Error
}
