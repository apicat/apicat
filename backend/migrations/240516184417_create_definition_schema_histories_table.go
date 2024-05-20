package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184417",
		Migrate: func(tx *gorm.DB) error {

			type DefinitionSchemaHistory struct {
				ID          uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				SchemaID    uint   `gorm:"type:bigint;index;not null;comment:模型id"`
				Name        string `gorm:"type:varchar(255);not null;comment:名称"`
				Description string `gorm:"type:varchar(255);comment:描述"`
				Schema      string `gorm:"type:mediumtext;comment:内容"`
				CreatedBy   uint   `gorm:"type:bigint;not null;default:0;comment:创建人id"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&DefinitionSchemaHistory{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&DefinitionSchemaHistory{})
		},
	}

	MigrationHelper.Register(m)
}
