package collection

import (
	"apicat-cloud/backend/model"
	"context"
	"time"
)

type TagToCollection struct {
	ID           uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	TagID        uint `gorm:"type:bigint;index;not null;comment:标签id"`
	CollectionID uint `gorm:"type:bigint;not null;comment:集合id"`
	DisplayOrder int  `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func init() {
	model.RegMigrate(&TagToCollection{})
}

func (ttc *TagToCollection) Create(ctx context.Context) error {
	return model.DB(ctx).Create(ttc).Error
}
