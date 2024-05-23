package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184402",
		Migrate: func(tx *gorm.DB) error {
			type Oauth2Bind struct {
				ID       uint
				UserID   uint   `gorm:"type:bigint;uniqueIndex:ukey;comment:user id"`                  // github.com/apicat/apicat/v2 user.id
				Type     string `gorm:"type:varchar(32);uniqueIndex:ukey;not null;comment:oauth type"` // github
				OauthUID string `gorm:"type:varchar(255);comment:uid of the OAuth"`                    // github user.id
				model.TimeModel
			}

			if tx.Migrator().HasTable(&Oauth2Bind{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&Oauth2Bind{})
		},
	}

	MigrationHelper.Register(m)
}
