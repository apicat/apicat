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
	SchemaCategory = "category"
	SchemaSchema   = "schema"
)

type DefinitionSchema struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID    string `gorm:"type:varchar(24);index;not null;comment:project id"`
	ParentID     uint   `gorm:"type:bigint;not null;comment:parent schema id"`
	Name         string `gorm:"type:varchar(255);not null;comment:scheam name"`
	Description  string `gorm:"type:varchar(255);comment:schema description"`
	Type         string `gorm:"type:varchar(255);not null;comment:schema type:category,schema"`
	Schema       string `gorm:"type:mediumtext;comment:schema content"`
	DisplayOrder uint   `gorm:"type:int(11);not null;default:0;comment:display order"`
	VectorID     string `gorm:"type:varchar(255);comment:vector db record id"`
	CreatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:created by member id"`
	UpdatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:updated by member id"`
	DeletedBy    uint   `gorm:"type:bigint;default:null;comment:deleted by member id"`
	model.TimeModel
}

func (ds *DefinitionSchema) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if ds.ID != 0 && ds.ProjectID != "" {
		tx = tx.Take(ds, "id = ? AND project_id = ?", ds.ID, ds.ProjectID)
	} else if ds.ProjectID != "" && ds.Name != "" {
		tx = tx.Take(ds, "project_id = ? AND name = ? AND type = ?", ds.ProjectID, ds.Name, ds.Type)
	} else {
		return false, errors.New("query condition error")
	}
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

func (ds *DefinitionSchema) HasChildren(ctx context.Context) (bool, error) {
	tx := model.DB(ctx).Model(ds).Where("project_id = ? AND parent_id = ?", ds.ProjectID, ds.ID).Take(&DefinitionSchema{})
	return tx.Error == nil, model.NotRecord(tx)
}

func (ds *DefinitionSchema) Create(ctx context.Context, tm *team.TeamMember) error {
	if ds.Type == SchemaCategory {
		if err := model.DB(ctx).Model(ds).Where("project_id = ? AND parent_id = ?", ds.ProjectID, ds.ParentID).Update("display_order", gorm.Expr("display_order + ?", 1)).Error; err != nil {
			slog.ErrorContext(ctx, "DefinitionSchema.Create.UpdateOrder", "err", err)
		}
		ds.DisplayOrder = 1
	} else {
		if ds.DisplayOrder == 0 {
			// 获取最大的display_order
			var maxDisplayOrder DefinitionSchema
			if err := model.DB(ctx).Model(ds).Where("project_id = ? AND parent_id = ?", ds.ProjectID, ds.ParentID).Order("display_order desc").First(&maxDisplayOrder).Error; err != nil {
				maxDisplayOrder = DefinitionSchema{DisplayOrder: 0}
			}
			ds.DisplayOrder = maxDisplayOrder.DisplayOrder + 1
		}
	}

	ds.CreatedBy = tm.ID
	ds.UpdatedBy = tm.ID

	return model.DB(ctx).Create(ds).Error
}

func (ds *DefinitionSchema) Update(ctx context.Context, name, description, schema string, memberID uint) error {
	if ds.Type != SchemaCategory {
		h := &DefinitionSchemaHistory{
			SchemaID:    ds.ID,
			Name:        ds.Name,
			Description: ds.Description,
			Schema:      ds.Schema,
		}
		h.Create(ctx, memberID)
	}

	// 只能修改name、description、schema
	return model.DB(ctx).Model(ds).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
		"schema":      schema,
		"updated_by":  memberID,
	}).Error
}

func (ds *DefinitionSchema) Delete(ctx context.Context, tm *team.TeamMember) error {
	return model.DB(ctx).Model(ds).Updates(map[string]interface{}{
		"deleted_by": tm.ID,
		"deleted_at": time.Now(),
	}).Error
}

func (ds *DefinitionSchema) Sort(ctx context.Context, parentID, displayOrder uint) error {
	return model.DB(ctx).Model(ds).UpdateColumns(map[string]interface{}{
		"parent_id":     parentID,
		"display_order": displayOrder,
	}).Error
}

func (ds *DefinitionSchema) ToSpec() (*spec.DefinitionModel, error) {
	s := &spec.DefinitionModel{
		ID:          int64(ds.ID),
		ParentId:    uint64(ds.ParentID),
		Name:        ds.Name,
		Type:        ds.Type,
		Description: ds.Description,
	}

	if ds.Schema != "" {
		if err := json.Unmarshal([]byte(ds.Schema), &s.Schema); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (ds *DefinitionSchema) DelRef(ctx context.Context, refSchema *DefinitionSchema, deref bool) error {
	schemaSpec, err := ds.ToSpec()
	if err != nil {
		return err
	}

	refSchemaSpec, err := refSchema.ToSpec()
	if err != nil {
		return err
	}

	if deref {
		if err := schemaSpec.Deref(refSchemaSpec); err != nil {
			return err
		}
	} else {
		schemaSpec.DelRef(refSchemaSpec)
	}

	content, err := json.Marshal(schemaSpec.Schema)
	if err != nil {
		return err
	}

	return model.DB(ctx).Model(ds).Select("schema").UpdateColumn("schema", string(content)).Error
}
