package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184428",
		Migrate: func(tx *gorm.DB) error {

			type RefResponseCollections struct {
				ID             uint `gorm:"type:bigint;primaryKey;autoIncrement"`
				RefResponserID uint `gorm:"type:bigint;index;not null;comment:被引用的公共响应id"`
				CollectionID   uint `gorm:"type:bigint;not null;comment:引用ref_responser_id的文档id"`
			}

			if tx.Migrator().HasTable(&RefResponseCollections{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&RefResponseCollections{})
		},
	}

	MigrationHelper.Register(m)
}
