package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	type RefSchemaResponses struct {
		ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		RefSchemaID uint `gorm:"type:bigint;index;not null;comment:被引用的公共模型id"`
		ResponseID  uint `gorm:"type:bigint;not null;comment:引用ref_schema_id的公共响应id"`
	}

	type ResponseReference struct {
		ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		ResponseID  uint `gorm:"type:bigint;index;not null;comment:公共响应id"`
		RefSchemaID uint `gorm:"type:bigint;index;not null;comment:引用的公共模型id"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	m := &gormigrate.Migration{
		ID: "240516184801",
		Migrate: func(tx *gorm.DB) error {

			if tx.Migrator().HasTable(&RefSchemaResponses{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&RefSchemaResponses{})
		},
	}
	MigrationHelper.Register(m)

	md := &gormigrate.Migration{
		ID: "240516184802",
		Migrate: func(tx *gorm.DB) error {
			var list []*ResponseReference
			if err := tx.Find(&list).Error; err != nil {
				return err
			}

			var newList []*RefSchemaResponses
			for _, item := range list {
				newList = append(newList, &RefSchemaResponses{
					RefSchemaID: item.RefSchemaID,
					ResponseID:  item.ResponseID,
				})
			}

			if len(newList) == 0 {
				return nil
			}
			return tx.Create(&newList).Error
		},
	}
	MigrationHelper.Register(md)
}
