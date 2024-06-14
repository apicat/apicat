package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184423",
		Migrate: func(tx *gorm.DB) error {

			type CollectionReference struct {
				ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				CollectionID uint   `gorm:"type:bigint;index;not null;comment:集合id"`
				RefID        uint   `gorm:"type:bigint;index;not null;comment:引用节点id"`
				RefType      string `gorm:"type:varchar(255);not null;comment:引用节点类型:schema,response,parameter"`
				CreatedAt    time.Time
				UpdatedAt    time.Time
			}

			if tx.Migrator().HasTable(&CollectionReference{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&CollectionReference{})
		},
	}

	MigrationHelper.Register(m)
}
