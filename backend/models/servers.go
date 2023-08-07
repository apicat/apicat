package models

import (
	"time"

	"github.com/apicat/apicat/backend/common/spec"
	"gorm.io/gorm"
)

type Servers struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectId    uint   `gorm:"type:bigint;index;not null;comment:项目id"`
	Description  string `gorm:"type:varchar(255);not null;comment:描述"`
	Url          string `gorm:"type:varchar(255);not null;comment:服务器地址"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewServers() *Servers {
	return &Servers{}
}

func (s *Servers) GetByProjectId(projectID uint) ([]Servers, error) {
	var servers []Servers
	return servers, Conn.Where("project_id = ?", projectID).Order("display_order asc").Find(&servers).Error
}

func (s *Servers) DeleteAndCreateServers(projectID uint, newServers []*Servers) error {
	// 开始事务
	Conn.Transaction(func(tx *gorm.DB) error {
		// 删除指定 ProjectId 的记录
		if err := tx.Where("project_id = ?", projectID).Delete(&Servers{}).Error; err != nil {
			return err
		}

		if len(newServers) == 0 {
			return nil
		}

		// 创建多个新记录
		if err := tx.Create(&newServers).Error; err != nil {
			return err
		}

		return nil
	})

	return nil
}

func ServersImport(projectID uint, servers []*spec.Server) {
	if len(servers) > 0 {
		for i, server := range servers {
			Conn.Create(&Servers{
				ProjectId:    projectID,
				Description:  server.Description,
				Url:          server.URL,
				DisplayOrder: i,
			})
		}
	}
}

func ServersExport(projectID uint) []*spec.Server {
	servers := []*Servers{}
	specServers := []*spec.Server{}
	if err := Conn.Where("project_id = ?", projectID).Find(&servers).Error; err == nil {
		for _, server := range servers {
			specServers = append(specServers, &spec.Server{
				Description: server.Description,
				URL:         server.Url,
			})
		}
	}

	return specServers
}
