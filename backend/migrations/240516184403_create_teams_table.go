package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184403",
		Migrate: func(tx *gorm.DB) error {
			type Team struct {
				ID      string `gorm:"type:varchar(24);primarykey"`
				Name    string `gorm:"type:varchar(255);not null;comment:team name"`
				Avatar  string `gorm:"type:varchar(255);comment:team avatar"`
				OwnerID uint   `gorm:"type:bigint;comment:team owner id"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&Team{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&Team{})
		},
	}

	MigrationHelper.Register(m)
}
