package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184421",
		Migrate: func(tx *gorm.DB) error {

			type Iteration struct {
				ID          string `gorm:"type:varchar(24);primarykey"`
				TeamID      string `gorm:"type:varchar(24);not null;comment:团队id"`
				ProjectID   string `gorm:"type:varchar(24);index;not null;comment:项目id"`
				Title       string `gorm:"type:varchar(255);not null;comment:迭代标题"`
				Description string `gorm:"type:varchar(255);comment:迭代描述"`
				CreatedBy   uint   `gorm:"type:bigint;not null;default:0;comment:创建人id"`
				UpdatedBy   uint   `gorm:"type:bigint;not null;default:0;comment:最后更新人id"`
				DeletedBy   uint   `gorm:"type:bigint;default:null;comment:删除人id"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&Iteration{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&Iteration{})
		},
	}

	MigrationHelper.Register(m)
}
