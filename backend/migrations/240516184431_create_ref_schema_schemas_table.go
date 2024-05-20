package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184431",
		Migrate: func(tx *gorm.DB) error {

			type RefSchemaSchemas struct {
				ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
				RefSchemaID uint `gorm:"type:bigint;index;not null;comment:被引用的公共模型id"`
				SchemaID    uint `gorm:"type:bigint;not null;comment:引用ref_schema_id的公共模型id"`
			}

			if tx.Migrator().HasTable(&RefSchemaSchemas{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&RefSchemaSchemas{})
		},
	}

	MigrationHelper.Register(m)
}
