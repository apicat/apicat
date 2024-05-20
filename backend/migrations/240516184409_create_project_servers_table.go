package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184409",
		Migrate: func(tx *gorm.DB) error {

			type Server struct {
				ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				ProjectID    string `gorm:"type:varchar(24);index;not null;comment:项目id"`
				Description  string `gorm:"type:varchar(255);not null;comment:描述"`
				URL          string `gorm:"type:varchar(255);not null;comment:服务器地址"`
				DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
				CreatedAt    time.Time
				UpdatedAt    time.Time
			}

			if tx.Migrator().HasTable(&Server{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&Server{})
		},
	}

	MigrationHelper.Register(m)
}
