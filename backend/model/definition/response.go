package definition

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/module/spec"

	"gorm.io/gorm"
)

const (
	ResponseCategory = "category"
	ResponseResponse = "response"
)

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

func init() {
	model.RegMigrate(&DefinitionResponse{})
}

func (dr *DefinitionResponse) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if dr.ID != 0 && dr.ProjectID != "" {
		tx = tx.Take(dr, "id = ? AND project_id = ?", dr.ID, dr.ProjectID)
	} else if dr.ProjectID != "" && dr.Name != "" {
		tx = tx.Take(dr, "project_id = ? AND name = ? AND type = ?", dr.ProjectID, dr.Name, dr.Type)
	} else {
		return false, errors.New("query condition error")
	}
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

func (dr *DefinitionResponse) HasChildren(ctx context.Context) (bool, error) {
	tx := model.DB(ctx).Model(dr).Where("project_id = ? AND parent_id = ?", dr.ProjectID, dr.ID).Take(&DefinitionResponse{})
	return tx.Error == nil, model.NotRecord(tx)
}

func (dr *DefinitionResponse) Create(ctx context.Context, tm *team.TeamMember) error {
	if dr.Type == ResponseCategory {
		if err := model.DB(ctx).Model(dr).Where("project_id = ? AND parent_id = ?", dr.ProjectID, dr.ParentID).Update("display_order", gorm.Expr("display_order + ?", 1)).Error; err != nil {
			slog.ErrorContext(ctx, "DefinitionResponse.Create.UpdateOrder", "err", err)
		}
		dr.DisplayOrder = 1
	} else {
		if dr.DisplayOrder == 0 {
			// 获取最大的display_order
			var maxDisplayOrder DefinitionResponse
			if err := model.DB(ctx).Model(dr).Where("project_id = ? AND parent_id = ?", dr.ProjectID, dr.ParentID).Order("display_order desc").First(&maxDisplayOrder).Error; err != nil {
				maxDisplayOrder = DefinitionResponse{DisplayOrder: 0}
			}
			dr.DisplayOrder = maxDisplayOrder.DisplayOrder + 1
		}
	}

	dr.CreatedBy = tm.ID
	dr.UpdatedBy = tm.ID

	return model.DB(ctx).Create(dr).Error
}

func (dr *DefinitionResponse) Update(ctx context.Context, memberID uint) error {
	// 只能修改name、description、header、content
	return model.DB(ctx).Model(dr).Updates(map[string]interface{}{
		"name":        dr.Name,
		"description": dr.Description,
		"header":      dr.Header,
		"content":     dr.Content,
		"updated_by":  memberID,
	}).Error
}

func (dr *DefinitionResponse) Delete(ctx context.Context, tm *team.TeamMember) error {
	return model.DB(ctx).Model(dr).Updates(map[string]interface{}{
		"deleted_by": tm.ID,
		"deleted_at": time.Now(),
	}).Error
}

func (dr *DefinitionResponse) Sort(ctx context.Context, parentID, displayOrder uint) error {
	return model.DB(ctx).Model(dr).UpdateColumns(map[string]interface{}{
		"parent_id":     parentID,
		"display_order": displayOrder,
	}).Error
}

func (dr *DefinitionResponse) ToSpec() (*spec.DefinitionResponse, error) {
	r := &spec.DefinitionResponse{
		BasicResponse: spec.BasicResponse{
			ID:          int64(dr.ID),
			Name:        dr.Name,
			Description: dr.Description,
		},
		ParentId: int64(dr.ParentID),
		Type:     dr.Type,
	}

	if dr.Header != "" {
		if err := json.Unmarshal([]byte(dr.Header), &r.Header); err != nil {
			return nil, err
		}
	}
	if dr.Content != "" {
		if err := json.Unmarshal([]byte(dr.Content), &r.Content); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (dr *DefinitionResponse) DelRef(ctx context.Context, refSchema *DefinitionSchema, deref bool) error {
	responseSpec, err := dr.ToSpec()
	if err != nil {
		return err
	}

	refSchemaSpec, err := refSchema.ToSpec()
	if err != nil {
		return err
	}

	if deref {
		if err := responseSpec.Deref(refSchemaSpec); err != nil {
			return err
		}
	} else {
		responseSpec.DelRef(refSchemaSpec)
	}

	content, err := json.Marshal(responseSpec.Content)
	if err != nil {
		return err
	}

	return model.DB(ctx).Model(dr).Select("content").UpdateColumn("content", string(content)).Error
}
