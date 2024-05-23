package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184414",
		Migrate: func(tx *gorm.DB) error {

			type TagToCollection struct {
				ID           uint `gorm:"type:bigint;primaryKey;autoIncrement"`
				TagID        uint `gorm:"type:bigint;index;not null;comment:tag id"`
				CollectionID uint `gorm:"type:bigint;not null;comment:collection id"`
				DisplayOrder int  `gorm:"type:int(11);not null;default:0;comment:display order"`
				CreatedAt    time.Time
				UpdatedAt    time.Time
			}

			if tx.Migrator().HasTable(&TagToCollection{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&TagToCollection{})
		},
	}

	MigrationHelper.Register(m)
}
