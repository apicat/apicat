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
				Name        string    `gorm:"type:varchar(255);comment:用户名"`
				Password    string    `gorm:"type:varchar(64);comment:密码"`
				Email       string    `gorm:"type:varchar(255);uniqueIndex;comment:邮箱"`
				Avatar      string    `gorm:"type:varchar(255);comment:头像"`
				Language    string    `gorm:"type:varchar(32);comment:语言"` // zh-CN en-US
				Role        string    `gorm:"type:varchar(32);comment:角色"`
				LastLoginIP string    `gorm:"type:varchar(15);comment:最后登录ip"`
				LastLoginAt time.Time `gorm:"type:datetime;not null;comment:最后登录时间"`
				IsActive    bool      `gorm:"type:tinyint;not null;comment:是否激活"`
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
