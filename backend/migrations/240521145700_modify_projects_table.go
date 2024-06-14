package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240521145700",
		Migrate: func(tx *gorm.DB) error {
			type Project struct {
				Description string `gorm:"type:varchar(1024);comment:project description"`
			}
			if tx.Migrator().HasTable(&Project{}) {
				if tx.Migrator().HasColumn(&Project{}, "description") {
					return tx.Migrator().AlterColumn(&Project{}, "description")
				}
			}
			return nil
		},
	}
	MigrationHelper.Register(m)
}
