package referencerelationship

import (
	"context"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
)

type SchemaReference struct {
	ID          uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID   string `gorm:"type:varchar(24);index;not null;comment:项目id"`
	SchemaID    uint   `gorm:"type:bigint;index;not null;comment:公共模型id"`
	RefSchemaID uint   `gorm:"type:bigint;index;not null;comment:引用的公共模型id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func BatchCreateSchemaReference(ctx context.Context, list []*SchemaReference) error {
	if len(list) == 0 {
		return nil
	}
	return model.DB(ctx).Create(list).Error
}

func BatchDeleteSchemaReference(ctx context.Context, ids ...uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Where("id in ?", ids).Delete(&SchemaReference{}).Error
}

func GetSchemasReferences(ctx context.Context, projectID string) ([]*SchemaReference, error) {
	var list []*SchemaReference
	tx := model.DB(ctx).Where("project_id = ?", projectID).Find(&list)
	return list, tx.Error
}

func GetSchemaReferencesBySchema(ctx context.Context, projectID string, schemaID uint) ([]*SchemaReference, error) {
	var list []*SchemaReference
	tx := model.DB(ctx).Where("project_id = ? AND schema_id = ?", projectID, schemaID).Find(&list)
	return list, tx.Error
}

func GetSchemaReferencesByRefSchema(ctx context.Context, projectID string, refSchemaID uint) ([]*SchemaReference, error) {
	var list []*SchemaReference
	tx := model.DB(ctx).Where("project_id = ? AND ref_schema_id = ?", projectID, refSchemaID).Find(&list)
	return list, tx.Error
}
