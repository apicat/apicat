package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184422",
		Migrate: func(tx *gorm.DB) error {

			type IterationApi struct {
				ID             uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				IterationID    string `gorm:"type:varchar(24);index;not null;comment:iteration id"`
				CollectionID   uint   `gorm:"type:bigint;not null;comment:collection id"`
				CollectionType string `gorm:"type:varchar(255);not null;comment:collection type:category,doc,http"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&IterationApi{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&IterationApi{})
		},
	}

	MigrationHelper.Register(m)
}
