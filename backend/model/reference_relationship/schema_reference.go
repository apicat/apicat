package referencerelationship

import (
	"context"
	"errors"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
)

type SchemaReference struct {
	ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	SchemaID    uint `gorm:"type:bigint;index;not null;comment:公共模型id"`
	RefSchemaID uint `gorm:"type:bigint;index;not null;comment:引用的公共模型id"`
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

func (sr *SchemaReference) GetSchemaRefs(ctx context.Context) ([]*SchemaReference, error) {
	var list []*SchemaReference

	tx := model.DB(ctx)
	if sr.SchemaID != 0 {
		tx = tx.Where("schema_id = ?", sr.SchemaID)
	} else if sr.RefSchemaID != 0 {
		tx = tx.Where("ref_schema_id = ?", sr.RefSchemaID)
	} else {
		return nil, errors.New("query condition error")
	}

	return list, tx.Find(&list).Error
}

func (sr *SchemaReference) GetSchemaIDsByRef(ctx context.Context) ([]uint, error) {
	var list []uint
	tx := model.DB(ctx).Where("ref_schema_id = ?", sr.RefSchemaID).Select("schema_id").Find(&list)
	return list, tx.Error
}

func (sr *SchemaReference) DelBySchemaID(ctx context.Context) error {
	return model.DB(ctx).Where("schema_id = ?", sr.SchemaID).Delete(&SchemaReference{}).Error
}

func (sr *SchemaReference) DelByRefSchemaID(ctx context.Context) error {
	return model.DB(ctx).Where("ref_schema_id = ?", sr.RefSchemaID).Delete(&SchemaReference{}).Error
}
