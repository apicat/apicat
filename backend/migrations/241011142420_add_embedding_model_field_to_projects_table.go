package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	type project struct {
		ID             string `gorm:"type:varchar(24);primarykey"`
		TeamID         string `gorm:"type:varchar(24);index;comment:team id"`
		MemberID       uint   `gorm:"type:bigint;comment:team member ID of the project manager"`
		Title          string `gorm:"type:varchar(255);not null;comment:project title"`
		Visibility     string `gorm:"type:varchar(32);not null;comment:project visibility:0-private,1-public"`
		ShareKey       string `gorm:"type:varchar(255);comment:project share key"`
		Description    string `gorm:"type:varchar(1024);comment:project description"`
		Cover          string `gorm:"type:varchar(255);comment:project cover"`
		EmbeddingModel string `gorm:"type:varchar(255);comment:project embedding model"`
		model.TimeModel
	}

	m := &gormigrate.Migration{
		ID: "241011142420",
		Migrate: func(tx *gorm.DB) error {
			if tx.Migrator().HasTable(&project{}) {
				if !tx.Migrator().HasColumn(&project{}, "embedding_model") {
					if err := tx.Migrator().AddColumn(&project{}, "embedding_model"); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	MigrationHelper.Register(m)
}
