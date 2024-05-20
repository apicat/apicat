package migrations

import (
	referencerelation "github.com/apicat/apicat/v2/backend/model/reference_relation"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240517145900",
		Migrate: func(tx *gorm.DB) error {
			var list []*referencerelationship.ParameterExcept
			if err := tx.Find(&list).Error; err != nil {
				return err
			}

			var newList []*referencerelation.ExceptParamCollection
			for _, item := range list {
				newList = append(newList, &referencerelation.ExceptParamCollection{
					ExceptParamID: item.ParameterID,
					CollectionID:  item.ExceptCollectionID,
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
