package collection

import (
	"context"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
)

type TagToCollection struct {
	ID           uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	TagID        uint `gorm:"type:bigint;index;not null;comment:tag id"`
	CollectionID uint `gorm:"type:bigint;not null;comment:collection id"`
	DisplayOrder int  `gorm:"type:int(11);not null;default:0;comment:display order"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (ttc *TagToCollection) Create(ctx context.Context) error {
	return model.DB(ctx).Create(ttc).Error
}
