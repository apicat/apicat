package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184405",
		Migrate: func(tx *gorm.DB) error {

			type Sysconfig struct {
				ID        uint   `gorm:"primarykey"`
				Type      string `gorm:"type:varchar(255);uniqueIndex:ukey;not null;comment:Configuration type"`
				Driver    string `gorm:"type:varchar(255);uniqueIndex:ukey;not null"`
				BeingUsed bool   `gorm:"type:tinyint;comment:is using"`
				Config    string `gorm:"type:varchar(512);"`
			}

			if tx.Migrator().HasTable(&Sysconfig{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&Sysconfig{})
		},
	}

	MigrationHelper.Register(m)
}
