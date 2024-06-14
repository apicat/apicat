package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184416",
		Migrate: func(tx *gorm.DB) error {

			type DefinitionSchema struct {
				ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				ProjectID    string `gorm:"type:varchar(24);index;not null;comment:project id"`
				ParentID     uint   `gorm:"type:bigint;not null;comment:parent schema id"`
				Name         string `gorm:"type:varchar(255);not null;comment:scheam name"`
				Description  string `gorm:"type:varchar(255);comment:schema description"`
				Type         string `gorm:"type:varchar(255);not null;comment:schema type:category,schema"`
				Schema       string `gorm:"type:mediumtext;comment:schema content"`
				DisplayOrder uint   `gorm:"type:int(11);not null;default:0;comment:display order"`
				CreatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:created by member id"`
				UpdatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:updated by member id"`
				DeletedBy    uint   `gorm:"type:bigint;default:null;comment:deleted by member id"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&DefinitionSchema{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&DefinitionSchema{})
		},
	}

	MigrationHelper.Register(m)
}
