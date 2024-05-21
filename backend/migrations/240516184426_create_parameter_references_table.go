package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184426",
		Migrate: func(tx *gorm.DB) error {

			type ParameterExcept struct {
				ID                 uint `gorm:"type:bigint;primaryKey;autoIncrement"`
				ParameterID        uint `gorm:"type:bigint;index;not null;comment:全局参数id"`
				ExceptCollectionID uint `gorm:"type:bigint;index;not null;comment:排除集合id"`
				CreatedAt          time.Time
				UpdatedAt          time.Time
			}

			if tx.Migrator().HasTable(&ParameterExcept{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&ParameterExcept{})
		},
	}

	MigrationHelper.Register(m)
}
