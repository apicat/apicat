package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184429",
		Migrate: func(tx *gorm.DB) error {

			type RefSchemaCollections struct {
				ID           uint `gorm:"type:bigint;primaryKey;autoIncrement"`
				RefSchemaID  uint `gorm:"type:bigint;index;not null;comment:被引用的公共模型id"`
				CollectionID uint `gorm:"type:bigint;not null;comment:引用ref_schema_id的文档id"`
			}

			if tx.Migrator().HasTable(&RefSchemaCollections{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&RefSchemaCollections{})
		},
	}

	MigrationHelper.Register(m)
}
