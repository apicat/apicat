package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184411",
		Migrate: func(tx *gorm.DB) error {

			type Collection struct {
				ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				PublicID     string `gorm:"type:varchar(255);index;comment:collection public id"`
				ProjectID    string `gorm:"type:varchar(24);index;not null;comment:project id"`
				ParentID     uint   `gorm:"type:bigint;not null;comment:parent collection id"`
				Path         string `gorm:"type:varchar(255);not null;comment:request path"`
				Method       string `gorm:"type:varchar(255);not null;comment:request method"`
				Title        string `gorm:"type:varchar(255);not null;comment:collection title"`
				Type         string `gorm:"type:varchar(255);not null;comment:collection type:category,doc,http"`
				ShareKey     string `gorm:"type:varchar(255);comment:share key"`
				Content      string `gorm:"type:mediumtext;comment:doc content"`
				DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:display order"`
				CreatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:created by member id"`
				UpdatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:updated by member id"`
				DeletedBy    uint   `gorm:"type:bigint;default:null;comment:deleted by member id"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&Collection{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&Collection{})
		},
	}

	MigrationHelper.Register(m)
}
