package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// 创建一个全局的MigrationHelper实例
var MigrationHelper = &MigrationManager{}

// MigrationManager 管理所有的迁移
type MigrationManager struct {
	migrations []*gormigrate.Migration
}

// Register 注册一个迁移方法
func (m *MigrationManager) Register(migration *gormigrate.Migration) {
	m.migrations = append(m.migrations, migration)
}

// Sort 根据id对migrations进行排序
func (m *MigrationManager) Sort() {
	for i := 0; i < len(m.migrations); i++ {
		for j := i + 1; j < len(m.migrations); j++ {
			if m.migrations[i].ID > m.migrations[j].ID {
				m.migrations[i], m.migrations[j] = m.migrations[j], m.migrations[i]
			}
		}
	}
}

// Run 执行所有的迁移
func (m *MigrationManager) Run(db *gorm.DB) error {
	m.Sort()
	mg := gormigrate.New(db, gormigrate.DefaultOptions, m.migrations)
	return mg.Migrate()
}
