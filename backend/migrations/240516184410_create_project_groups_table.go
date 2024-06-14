package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184410",
		Migrate: func(tx *gorm.DB) error {

			type ProjectGroup struct {
				ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				MemberID     uint   `gorm:"type:bigint;uniqueIndex:ukey;not null;comment:team member id"`
				Name         string `gorm:"type:varchar(255);uniqueIndex:ukey;not null;comment:group name"`
				DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:display order"`
				CreatedAt    time.Time
				UpdatedAt    time.Time
			}

			if tx.Migrator().HasTable(&ProjectGroup{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&ProjectGroup{})
		},
	}

	MigrationHelper.Register(m)
}
