package migrations

import (
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184418",
		Migrate: func(tx *gorm.DB) error {

			type DefinitionResponse struct {
				ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				ProjectID    string `gorm:"type:varchar(24);index;not null;comment:项目id"`
				ParentID     uint   `gorm:"type:bigint;not null;comment:父级id"`
				Name         string `gorm:"type:varchar(255);not null;comment:响应名称"`
				Description  string `gorm:"type:varchar(255);not null;comment:状态描述"`
				Type         string `gorm:"type:varchar(255);not null;comment:响应类型:category,response"`
				Header       string `gorm:"type:mediumtext;comment:响应头"`
				Content      string `gorm:"type:mediumtext;comment:响应内容"`
				DisplayOrder uint   `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
				CreatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:创建成员id"`
				UpdatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:最后更新成员id"`
				DeletedBy    uint   `gorm:"type:bigint;default:null;comment:删除成员id"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&DefinitionResponse{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&DefinitionResponse{})
		},
	}

	MigrationHelper.Register(m)
}
