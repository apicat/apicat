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
				ProjectID string `gorm:"type:varchar(24);index;not null;comment:项目id"`
				In        string `gorm:"type:varchar(32);not null;comment:位置:header,cookie,query,path"`
				Name      string `gorm:"type:varchar(255);not null;comment:参数名称"`
				Required  bool   `gorm:"type:tinyint;not null;comment:是否必传"`
				Schema    string `gorm:"type:mediumtext;comment:参数内容"`
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
