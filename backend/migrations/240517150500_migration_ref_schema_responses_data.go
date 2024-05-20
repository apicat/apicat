package migrations

import (
	referencerelation "github.com/apicat/apicat/v2/backend/model/reference_relation"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240517150500",
		Migrate: func(tx *gorm.DB) error {
			var list []*referencerelationship.ResponseReference
			if err := tx.Find(&list).Error; err != nil {
				return err
			}

			var newList []*referencerelation.RefSchemaResponses
			for _, item := range list {
				newList = append(newList, &referencerelation.RefSchemaResponses{
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

	MigrationHelper.Register(m)
}
