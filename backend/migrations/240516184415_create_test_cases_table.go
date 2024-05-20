package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184415",
		Migrate: func(tx *gorm.DB) error {

			type TestCase struct {
				ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				ProjectID    string `gorm:"type:varchar(24);index:idx_pid_cid;not null;comment:项目id"`
				CollectionID uint   `gorm:"type:bigint;index:idx_pid_cid;not null;comment:集合id"`
				Title        string `gorm:"type:varchar(255);not null;comment:测试用例标题"`
				Content      string `gorm:"type:mediumtext;comment:测试用例内容"`
				CreatedAt    time.Time
				UpdatedAt    time.Time
			}

			if tx.Migrator().HasTable(&TestCase{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&TestCase{})
		},
	}

	MigrationHelper.Register(m)
}
