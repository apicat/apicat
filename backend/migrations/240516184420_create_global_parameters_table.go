package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184420",
		Migrate: func(tx *gorm.DB) error {

			type GlobalParameter struct {
				ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				ProjectID    string `gorm:"type:varchar(24);index;not null;comment:项目id"`
				In           string `gorm:"type:varchar(32);not null;comment:位置:header,cookie,query,path"`
				Name         string `gorm:"type:varchar(255);not null;comment:参数名称"`
				Required     bool   `gorm:"type:tinyint;not null;comment:是否必传"`
				Schema       string `gorm:"type:mediumtext;comment:参数内容"`
				DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&GlobalParameter{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&GlobalParameter{})
		},
	}

	MigrationHelper.Register(m)
}
