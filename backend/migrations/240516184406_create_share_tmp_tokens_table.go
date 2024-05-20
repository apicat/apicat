package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184406",
		Migrate: func(tx *gorm.DB) error {

			type ShareTmpToken struct {
				ID           uint      `gorm:"type:bigint;primaryKey;autoIncrement"`
				ShareToken   string    `gorm:"type:varchar(255);index;not null;comment:md5的分享token"`
				Expiration   time.Time `gorm:"type:datetime;not null;comment:过期时间"`
				ProjectID    string    `gorm:"type:varchar(24);index;not null;comment:项目id"`
				CollectionID uint      `gorm:"type:bigint;index;comment:集合id"`
				CreatedAt    time.Time
				UpdatedAt    time.Time
			}

			if tx.Migrator().HasTable(&ShareTmpToken{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&ShareTmpToken{})
		},
	}

	MigrationHelper.Register(m)
}
