package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184427",
		Migrate: func(tx *gorm.DB) error {

			type ExceptParamCollection struct {
				ID            uint `gorm:"type:bigint;primaryKey;autoIncrement"`
				ExceptParamID uint `gorm:"type:bigint;index;not null;comment:被排除的全局参数id"`
				CollectionID  uint `gorm:"type:bigint;not null;comment:排除except_param_id的文档id"`
			}

			if tx.Migrator().HasTable(&ExceptParamCollection{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&ExceptParamCollection{})
		},
	}

	MigrationHelper.Register(m)
}
