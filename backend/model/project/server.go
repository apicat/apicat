package project

import (
	"apicat-cloud/backend/model"
	"context"
	"errors"
	"time"
)

type Server struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID    string `gorm:"type:varchar(24);index;not null;comment:项目id"`
	Description  string `gorm:"type:varchar(255);not null;comment:描述"`
	URL          string `gorm:"type:varchar(255);not null;comment:服务器地址"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func init() {
	model.RegMigrate(&Server{})
}

// Get 获取项目URL
func (s *Server) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if s.ID != 0 {
		tx = tx.Take(s, "id = ?", s.ID)
	} else if s.ProjectID != "" && s.URL != "" {
		tx = tx.Take(s, "project_id = ? AND url = ?", s.ProjectID, s.URL)
	} else {
		return false, errors.New("query condition error")
	}
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

// Create 创建项目URL
func (s *Server) Create(ctx context.Context) (*Server, error) {
	// 获取最大的display_order
	var maxDisplayOrder Server
	if err := model.DB(ctx).Model(s).Where("project_id = ?", s.ProjectID).Order("display_order desc").First(&maxDisplayOrder).Error; err != nil {
		maxDisplayOrder = Server{DisplayOrder: 0}
	}

	s.DisplayOrder = maxDisplayOrder.DisplayOrder + 1
	return s, model.DB(ctx).Create(s).Error
}

// CheckRepeat 检查重复，gp.ID不为0时排除自身
func (s *Server) CheckRepeat(ctx context.Context) (bool, error) {
	var count int64
	tx := model.DB(ctx).Model(s).Where("project_id = ? AND url = ?", s.ProjectID, s.URL)
	if s.ID != 0 {
		tx = tx.Where("id != ?", s.ID)
	}
	err := tx.Count(&count).Error
	return count > 0, err
}

// Update 更新项目URL
func (s *Server) Update(ctx context.Context) error {
	return model.DB(ctx).Model(s).Updates(map[string]interface{}{
		"url":         s.URL,
		"description": s.Description,
	}).Error
}

// Delete 删除项目URL
func (s *Server) Delete(ctx context.Context) error {
	return model.DB(ctx).Delete(s).Error
}
