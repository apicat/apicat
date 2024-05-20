package global

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/module/spec"
)

const (
	ParameterInHeader = "header"
	ParameterInCookie = "cookie"
	ParameterInQuery  = "query"
	ParameterInPath   = "path"
)

type GlobalParameter struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID    string `gorm:"type:varchar(24);index;not null;comment:项目id"`
	In           string `gorm:"type:varchar(32);not null;comment:位置:header,cookie,query,path"`
	Name         string `gorm:"type:varchar(255);not null;comment:参数名称"`
	Required     bool   `gorm:"type:tinyint;not null;comment:是否必传"`
	Schema       string `gorm:"type:mediumtext;comment:参数内容"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	model.TimeModel
}

// Get 获取全局参数
func (gp *GlobalParameter) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if gp.ID != 0 {
		tx = tx.Take(gp, "id = ?", gp.ID)
	} else if gp.ProjectID != "" && gp.In != "" && gp.Name != "" {
		tx = tx.Take(gp, "project_id = ? AND `in` = ? AND name = ?", gp.ProjectID, gp.In, gp.Name)
	} else {
		return false, errors.New("query condition error")
	}
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

// Create 创建全局参数
func (gp *GlobalParameter) Create(ctx context.Context) error {
	// 获取最大的display_order
	var maxDisplayOrder GlobalParameter
	if err := model.DB(ctx).Model(gp).Where("project_id = ? AND `in` = ?", gp.ProjectID, gp.In).Order("display_order desc").First(&maxDisplayOrder).Error; err != nil {
		maxDisplayOrder = GlobalParameter{DisplayOrder: 0}
	}
	gp.DisplayOrder = maxDisplayOrder.DisplayOrder + 1
	return model.DB(ctx).Create(gp).Error
}

// Update 更新全局参数
func (gp *GlobalParameter) Update(ctx context.Context) error {
	// 只能修改name、required、schema
	return model.DB(ctx).Model(gp).Updates(map[string]interface{}{
		"name":     gp.Name,
		"required": gp.Required,
		"schema":   gp.Schema,
	}).Error
}

// Delete 删除全局参数
func (gp *GlobalParameter) Delete(ctx context.Context) error {
	return model.DB(ctx).Delete(gp).Error
}

// CheckRepeat 检查重复，gp.ID不为0时排除自身
func (gp *GlobalParameter) CheckRepeat(ctx context.Context) (bool, error) {
	var count int64
	tx := model.DB(ctx).Model(gp).Where("project_id = ? AND name = ? AND `in` = ?", gp.ProjectID, gp.Name, gp.In)
	if gp.ID != 0 {
		tx = tx.Where("id != ?", gp.ID)
	}
	err := tx.Count(&count).Error
	return count > 0, err
}

func (gp *GlobalParameter) ToSpec() (*spec.Parameter, error) {
	p := &spec.Parameter{
		ID:       int64(gp.ID),
		Name:     gp.Name,
		Required: gp.Required,
	}

	if gp.Schema != "" {
		if err := json.Unmarshal([]byte(gp.Schema), &p.Schema); err != nil {
			return p, err
		}
	}

	return p, nil
}
