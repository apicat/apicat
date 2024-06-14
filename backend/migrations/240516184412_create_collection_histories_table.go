package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184412",
		Migrate: func(tx *gorm.DB) error {

			type CollectionHistory struct {
				ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				CollectionID uint   `gorm:"type:bigint;index;not null;comment:collection id"`
				Title        string `gorm:"type:varchar(255);not null;comment:collection title"`
				Content      string `gorm:"type:mediumtext;comment:doc content"`
				CreatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:created by member id"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&CollectionHistory{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&CollectionHistory{})
		},
	}

	MigrationHelper.Register(m)
}
