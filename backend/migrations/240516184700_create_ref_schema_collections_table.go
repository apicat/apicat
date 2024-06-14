package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	type RefSchemaCollections struct {
		ID           uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		RefSchemaID  uint `gorm:"type:bigint;index;not null;comment:referenced definition schema id"`
		CollectionID uint `gorm:"type:bigint;not null;comment:collection id"`
	}

	type CollectionReference struct {
		ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
		CollectionID uint   `gorm:"type:bigint;index;not null;comment:collection id"`
		RefID        uint   `gorm:"type:bigint;index;not null;comment:ref node id"`
		RefType      string `gorm:"type:varchar(255);not null;comment:ref node type:schema,response,parameter"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}

	m := &gormigrate.Migration{
		ID: "240516184701",
		Migrate: func(tx *gorm.DB) error {
			if !tx.Migrator().HasTable(&RefSchemaCollections{}) {
				if err := tx.Migrator().CreateTable(&RefSchemaCollections{}); err != nil {
					return err
				}
			}

			var list []*CollectionReference
			if err := tx.Where("ref_type = ?", "schema").Find(&list).Error; err != nil {
				return err
			}

			var newList []*RefSchemaCollections
			for _, item := range list {
				newList = append(newList, &RefSchemaCollections{
					RefSchemaID:  item.RefID,
					CollectionID: item.CollectionID,
				})
			}

			if len(newList) == 0 {
				return nil
			}
			return tx.Create(&newList).Error
		},
	}
	MigrationHelper.Register(m)
}
