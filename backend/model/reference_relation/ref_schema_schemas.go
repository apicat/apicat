package referencerelation

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
)

type RefSchemaSchemas struct {
	ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	RefSchemaID uint `gorm:"type:bigint;index;not null;comment:被引用的公共模型id"`
	SchemaID    uint `gorm:"type:bigint;not null;comment:引用ref_schema_id的公共模型id"`
}

func (r *RefSchemaSchemas) GetSchemas(ctx context.Context) ([]*RefSchemaSchemas, error) {
	var list []*RefSchemaSchemas
	tx := model.DB(ctx).Where("ref_schema_id = ?", r.RefSchemaID).Find(&list)
	return list, tx.Error
}

func (r *RefSchemaSchemas) GetSchemaIDs(ctx context.Context) ([]uint, error) {
	var list []uint
	tx := model.DB(ctx).Where("ref_schema_id = ?", r.RefSchemaID).Select("schema_id").Find(&list)
	return list, tx.Error
}

func BatchCreateRefSchemaSchemas(ctx context.Context, list []*RefSchemaSchemas) error {
	if len(list) == 0 {
		return nil
	}
	return model.DB(ctx).Create(list).Error
}

func BatchDelRefSchemaSchemas(ctx context.Context, ids ...uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Where("id in ?", ids).Delete(&RefSchemaSchemas{}).Error
}

// GetRefSchemaSchemas 获取引用指定公共模型的所有公共模型的引用关系
func GetRefSchemaSchemas(ctx context.Context, schemaIDs ...uint) ([]*RefSchemaSchemas, error) {
	var list []*RefSchemaSchemas
	tx := model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Find(&list)
	return list, tx.Error
}

// DelRefSchemaSchema 删除引用指定公共模型的指定公共模型的引用关系
// 用于删除文档时，删除该文档引用的公共模型
func DelRefSchemaSchema(ctx context.Context, schemaIDs []uint, responseID uint) error {
	return model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Where("schema_id = ?", responseID).Delete(&RefSchemaSchemas{}).Error
}

// DelRefSchemaSchemas 删除引用指定公共模型的所有公共模型的引用关系
// 用于删除公共模型时，删除所有引用了该公共模型的公共模型
func DelRefSchemaSchemas(ctx context.Context, schemaIDs uint) error {
	return model.DB(ctx).Where("ref_schema_id = ?", schemaIDs).Delete(&RefSchemaSchemas{}).Error
}
