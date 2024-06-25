package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	type ExceptParamCollection struct {
		ID            uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		ExceptParamID uint `gorm:"type:bigint;index;not null;comment:excluded global parameter id"`
		CollectionID  uint `gorm:"type:bigint;not null;comment:collection id"`
	}

	type ParameterExcept struct {
		ID                 uint `gorm:"type:bigint;primaryKey;autoIncrement"`
		ParameterID        uint `gorm:"type:bigint;index;not null;comment:global parameter id"`
		ExceptCollectionID uint `gorm:"type:bigint;index;not null;comment:excluded collection id"`
		CreatedAt          time.Time
		UpdatedAt          time.Time
	}

	m := &gormigrate.Migration{
		ID: "240516184501",
		Migrate: func(tx *gorm.DB) error {
			if !tx.Migrator().HasTable(&ExceptParamCollection{}) {
				if err := tx.Migrator().CreateTable(&ExceptParamCollection{}); err != nil {
					return err
				}
			}

			if tx.Migrator().HasTable(&ParameterExcept{}) {
				var list []*ParameterExcept
				if err := tx.Find(&list).Error; err != nil {
					return err
				}

				var newList []*ExceptParamCollection
				for _, item := range list {
					newList = append(newList, &ExceptParamCollection{
						ExceptParamID: item.ParameterID,
						CollectionID:  item.ExceptCollectionID,
					})
				}

				if len(newList) == 0 {
					return nil
				}
				return tx.Create(&newList).Error
			}
			return nil
		},
	}
	MigrationHelper.Register(m)
}
