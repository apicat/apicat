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
				Name    string `gorm:"type:varchar(255);not null;comment:团队名"`
				Avatar  string `gorm:"type:varchar(255);comment:头像"`
				OwnerID uint   `gorm:"type:bigint;comment:团队所有者id"`
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
