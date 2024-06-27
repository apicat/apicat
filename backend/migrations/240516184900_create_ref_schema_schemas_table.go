package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	type RefSchemaSchemas struct {
		ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		RefSchemaID uint `gorm:"type:bigint;index;not null;comment:referenced definition schema id"`
		SchemaID    uint `gorm:"type:bigint;not null;comment:schema id"`
	}

	type SchemaReference struct {
		ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		SchemaID    uint `gorm:"type:bigint;index;not null;comment:definition schema id"`
		RefSchemaID uint `gorm:"type:bigint;index;not null;comment:ref schema id"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	m := &gormigrate.Migration{
		ID: "240516184901",
		Migrate: func(tx *gorm.DB) error {
			if !tx.Migrator().HasTable(&RefSchemaSchemas{}) {
				if err := tx.Migrator().CreateTable(&RefSchemaSchemas{}); err != nil {
					return err
				}
			}

			if tx.Migrator().HasTable(&SchemaReference{}) {
				var list []*SchemaReference
				if err := tx.Find(&list).Error; err != nil {
					return err
				}

				var newList []*RefSchemaSchemas
				for _, item := range list {
					newList = append(newList, &RefSchemaSchemas{
						RefSchemaID: item.RefSchemaID,
						SchemaID:    item.SchemaID,
					})
				}

				if len(newList) == 0 {
					return nil
				}
				return tx.Create(&newList).Error
			}
			return nil
		},
	}
	MigrationHelper.Register(m)
}
