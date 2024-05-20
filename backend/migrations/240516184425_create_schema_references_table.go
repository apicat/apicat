package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184425",
		Migrate: func(tx *gorm.DB) error {

			type SchemaReference struct {
				ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
				SchemaID    uint `gorm:"type:bigint;index;not null;comment:公共模型id"`
				RefSchemaID uint `gorm:"type:bigint;index;not null;comment:引用的公共模型id"`
				CreatedAt   time.Time
				UpdatedAt   time.Time
			}

			if tx.Migrator().HasTable(&SchemaReference{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&SchemaReference{})
		},
	}

	MigrationHelper.Register(m)
}
