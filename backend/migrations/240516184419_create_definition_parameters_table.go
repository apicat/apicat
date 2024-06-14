package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184419",
		Migrate: func(tx *gorm.DB) error {

			type DefinitionParameter struct {
				ID        uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				ProjectID string `gorm:"type:varchar(24);index;not null;comment:project id"`
				In        string `gorm:"type:varchar(32);not null;comment:param in:header,cookie,query,path"`
				Name      string `gorm:"type:varchar(255);not null;comment:param name"`
				Required  bool   `gorm:"type:tinyint;not null;comment:is required"`
				Schema    string `gorm:"type:mediumtext;comment:param schema"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&DefinitionParameter{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&DefinitionParameter{})
		},
	}

	MigrationHelper.Register(m)
}
