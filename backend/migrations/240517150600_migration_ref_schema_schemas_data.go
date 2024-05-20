package migrations

import (
	referencerelation "github.com/apicat/apicat/v2/backend/model/reference_relation"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240517150600",
		Migrate: func(tx *gorm.DB) error {
			var list []*referencerelationship.SchemaReference
			if err := tx.Find(&list).Error; err != nil {
				return err
			}

			var newList []*referencerelation.RefSchemaSchemas
			for _, item := range list {
				newList = append(newList, &referencerelation.RefSchemaSchemas{
					RefSchemaID: item.RefSchemaID,
					SchemaID:    item.SchemaID,
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
