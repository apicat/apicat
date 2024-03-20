package iteration

import (
	"context"

	"github.com/apicat/apicat/backend/model"
)

type IterationApi struct {
	ID             uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	IterationID    string `gorm:"type:varchar(24);index;not null;comment:迭代id"`
	CollectionID   uint   `gorm:"type:bigint;not null;comment:集合id"`
	CollectionType string `gorm:"type:varchar(255);not null;comment:集合类型:category,doc,http"`
	model.TimeModel
}

func init() {
	model.RegMigrate(&IterationApi{})
}

func (ia *IterationApi) Create(ctx context.Context) error {
	return model.DB(ctx).Create(ia).Error
}

func (ia *IterationApi) Delete(ctx context.Context) error {
	return model.DB(ctx).Delete(ia).Error
}
