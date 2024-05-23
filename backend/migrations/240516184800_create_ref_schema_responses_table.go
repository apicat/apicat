package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	type RefSchemaResponses struct {
		ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		RefSchemaID uint `gorm:"type:bigint;index;not null;comment:referenced definition schema id"`
		ResponseID  uint `gorm:"type:bigint;not null;comment:response id"`
	}

	type ResponseReference struct {
		ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		ResponseID  uint `gorm:"type:bigint;index;not null;comment:definition response id"`
		RefSchemaID uint `gorm:"type:bigint;index;not null;comment:ref schema id"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	m := &gormigrate.Migration{
		ID: "240516184801",
		Migrate: func(tx *gorm.DB) error {
			if !tx.Migrator().HasTable(&RefSchemaResponses{}) {
				if err := tx.Migrator().CreateTable(&RefSchemaResponses{}); err != nil {
					return err
				}
			}

			var list []*ResponseReference
			if err := tx.Find(&list).Error; err != nil {
				return err
			}

			var newList []*RefSchemaResponses
			for _, item := range list {
				newList = append(newList, &RefSchemaResponses{
					RefSchemaID: item.RefSchemaID,
					ResponseID:  item.ResponseID,
				})
			}

			if len(newList) == 0 {
				return nil
			}
			return tx.Create(&newList).Error
		},
	}
	MigrationHelper.Register(m)
}
