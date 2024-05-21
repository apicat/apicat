package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184407",
		Migrate: func(tx *gorm.DB) error {

			type Project struct {
				ID          string `gorm:"type:varchar(24);primarykey"`
				TeamID      string `gorm:"type:varchar(24);index;comment:团队id"`
				MemberID    uint   `gorm:"type:bigint;comment:项目管理者的团队成员id"`
				Title       string `gorm:"type:varchar(255);not null;comment:项目名称"`
				Visibility  string `gorm:"type:varchar(32);not null;comment:项目可见性:0私有,1公开"`
				ShareKey    string `gorm:"type:varchar(255);comment:项目分享密码"`
				Description string `gorm:"type:varchar(255);comment:项目描述"`
				Cover       string `gorm:"type:varchar(255);comment:项目封面"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&Project{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&Project{})
		},
	}

	MigrationHelper.Register(m)
}
