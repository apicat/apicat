package migrations

import (
	"time"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184408",
		Migrate: func(tx *gorm.DB) error {

			type ProjectMember struct {
				ID         uint               `gorm:"type:bigint;primaryKey"`
				ProjectID  string             `gorm:"type:varchar(24);uniqueIndex:ukey;not null;comment:项目id"`
				MemberID   uint               `gorm:"type:bigint;uniqueIndex:ukey;not null;comment:团队成员id"`
				GroupID    uint               `gorm:"type:bigint;not null;default:0;comment:分组id"`
				Permission project.Permission `gorm:"type:varchar(255);not null;comment:项目权限:manage,write,read"`
				FollowedAt *time.Time         `gorm:"type:datetime;comment:关注项目时间"` // 不为空表示关注，字段类型为指针是为了在取消关注时，可以设置为null
				model.TimeModel
			}

			if tx.Migrator().HasTable(&ProjectMember{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&ProjectMember{})
		},
	}

	MigrationHelper.Register(m)
}
