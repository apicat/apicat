package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184407",
		Migrate: func(tx *gorm.DB) error {

			type Project struct {
				ID          string `gorm:"type:varchar(24);primarykey"`
				TeamID      string `gorm:"type:varchar(24);index;comment:team id"`
				MemberID    uint   `gorm:"type:bigint;comment:team member ID of the project manager"`
				Title       string `gorm:"type:varchar(255);not null;comment:project title"`
				Visibility  string `gorm:"type:varchar(32);not null;comment:project visibility:0-private,1-public"`
				ShareKey    string `gorm:"type:varchar(255);comment:project share key"`
				Description string `gorm:"type:varchar(1024);comment:project description"`
				Cover       string `gorm:"type:varchar(255);comment:project cover"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&Project{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&Project{})
		},
	}

	MigrationHelper.Register(m)
}
