package collection

import (
	"context"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
)

type Tag struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID    string `gorm:"type:varchar(24);index;not null;comment:项目id"`
	Name         string `gorm:"type:varchar(255);not null;comment:名称"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func init() {
	model.RegMigrate(&Tag{})
}

func (t *Tag) Create(ctx context.Context) error {
	return model.DB(ctx).Create(t).Error
}
