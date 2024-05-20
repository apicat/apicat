package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	type RefResponseCollections struct {
		ID             uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		RefResponserID uint `gorm:"type:bigint;index;not null;comment:被引用的公共响应id"`
		CollectionID   uint `gorm:"type:bigint;not null;comment:引用ref_responser_id的文档id"`
	}

	type CollectionReference struct {
		ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
		CollectionID uint   `gorm:"type:bigint;index;not null;comment:集合id"`
		RefID        uint   `gorm:"type:bigint;index;not null;comment:引用节点id"`
		RefType      string `gorm:"type:varchar(255);not null;comment:引用节点类型:schema,response,parameter"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}

	m := &gormigrate.Migration{
		ID: "240516184601",
		Migrate: func(tx *gorm.DB) error {

			if tx.Migrator().HasTable(&RefResponseCollections{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&RefResponseCollections{})
		},
	}
	MigrationHelper.Register(m)

	md := &gormigrate.Migration{
		ID: "240516184602",
		Migrate: func(tx *gorm.DB) error {
			var list []*CollectionReference
			if err := tx.Where("ref_type = response").Find(&list).Error; err != nil {
				return err
			}

			var newList []*RefResponseCollections
			for _, item := range list {
				newList = append(newList, &RefResponseCollections{
					RefResponserID: item.RefID,
					CollectionID:   item.CollectionID,
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
