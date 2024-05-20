package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184424",
		Migrate: func(tx *gorm.DB) error {

			type ResponseReference struct {
				ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
				ResponseID  uint `gorm:"type:bigint;index;not null;comment:公共响应id"`
				RefSchemaID uint `gorm:"type:bigint;index;not null;comment:引用的公共模型id"`
				CreatedAt   time.Time
				UpdatedAt   time.Time
			}

			if tx.Migrator().HasTable(&ResponseReference{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&ResponseReference{})
		},
	}

	MigrationHelper.Register(m)
}
