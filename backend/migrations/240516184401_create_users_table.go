package migrations

import (
	"time"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184401",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				ID          uint      `gorm:"primarykey"`
				Name        string    `gorm:"type:varchar(255);comment:username"`
				Password    string    `gorm:"type:varchar(64);comment:password"`
				Email       string    `gorm:"type:varchar(255);uniqueIndex;comment:e-mail address"`
				Avatar      string    `gorm:"type:varchar(255);comment:user avatar"`
				Language    string    `gorm:"type:varchar(32);comment:language"` // zh-CN en-US
				Role        string    `gorm:"type:varchar(32);comment:role"`
				LastLoginIP string    `gorm:"type:varchar(15);comment:last login ip"`
				LastLoginAt time.Time `gorm:"type:datetime;not null;comment:last login time"`
				IsActive    bool      `gorm:"type:tinyint;not null;comment:is active"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&User{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&User{})
		},
	}

	MigrationHelper.Register(m)
}
