package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184421",
		Migrate: func(tx *gorm.DB) error {

			type Iteration struct {
				ID          string `gorm:"type:varchar(24);primarykey"`
				TeamID      string `gorm:"type:varchar(24);not null;comment:team id"`
				ProjectID   string `gorm:"type:varchar(24);index;not null;comment:project id"`
				Title       string `gorm:"type:varchar(255);not null;comment:iteartion title"`
				Description string `gorm:"type:varchar(255);comment:iteration description"`
				CreatedBy   uint   `gorm:"type:bigint;not null;default:0;comment:created by member id"`
				UpdatedBy   uint   `gorm:"type:bigint;not null;default:0;comment:updated by member id"`
				DeletedBy   uint   `gorm:"type:bigint;default:null;comment:deleted by member id"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&Iteration{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&Iteration{})
		},
	}

	MigrationHelper.Register(m)
}
