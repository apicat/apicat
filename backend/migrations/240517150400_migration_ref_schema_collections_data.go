package migrations

import (
	"fmt"

	referencerelation "github.com/apicat/apicat/v2/backend/model/reference_relation"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240517150400",
		Migrate: func(tx *gorm.DB) error {
			var list []*referencerelationship.CollectionReference
			if err := tx.Where("ref_type = ?", referencerelationship.ReferenceSchema).Find(&list).Error; err != nil {
				return err
			}

			var newList []*referencerelation.RefSchemaCollections
			for _, item := range list {
				newList = append(newList, &referencerelation.RefSchemaCollections{
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
	fmt.Println("240517150400_migration_ref_schema_collections_data.go")
	MigrationHelper.Register(m)
}
