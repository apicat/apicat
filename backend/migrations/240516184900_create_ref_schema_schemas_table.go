package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	type RefSchemaSchemas struct {
		ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		RefSchemaID uint `gorm:"type:bigint;index;not null;comment:被引用的公共模型id"`
		SchemaID    uint `gorm:"type:bigint;not null;comment:引用ref_schema_id的公共模型id"`
	}

	type SchemaReference struct {
		ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		SchemaID    uint `gorm:"type:bigint;index;not null;comment:公共模型id"`
		RefSchemaID uint `gorm:"type:bigint;index;not null;comment:引用的公共模型id"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	m := &gormigrate.Migration{
		ID: "240516184901",
		Migrate: func(tx *gorm.DB) error {
			if tx.Migrator().HasTable(&RefSchemaSchemas{}) {
				if err := tx.Migrator().CreateTable(&RefSchemaSchemas{}); err != nil {
					return err
				}
			}

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
		},
	}
	MigrationHelper.Register(m)
}
