package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184430",
		Migrate: func(tx *gorm.DB) error {

			type RefSchemaResponses struct {
				ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
				RefSchemaID uint `gorm:"type:bigint;index;not null;comment:被引用的公共模型id"`
				ResponseID  uint `gorm:"type:bigint;not null;comment:引用ref_schema_id的公共响应id"`
			}

			if tx.Migrator().HasTable(&RefSchemaResponses{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&RefSchemaResponses{})
		},
	}

	MigrationHelper.Register(m)
}
