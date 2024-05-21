package referencerelation

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
)

type RefSchemaResponses struct {
	ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	RefSchemaID uint `gorm:"type:bigint;index;not null;comment:被引用的公共模型id"`
	ResponseID  uint `gorm:"type:bigint;not null;comment:引用ref_schema_id的公共响应id"`
}

func (r *RefSchemaResponses) GetResponses(ctx context.Context) ([]*RefSchemaResponses, error) {
	var list []*RefSchemaResponses
	tx := model.DB(ctx).Where("ref_schema_id = ?", r.RefSchemaID).Find(&list)
	return list, tx.Error
}

func (r *RefSchemaResponses) GetResponseIDs(ctx context.Context) ([]uint, error) {
	var list []uint
	tx := model.DB(ctx).Model(&RefSchemaResponses{}).Where("ref_schema_id = ?", r.RefSchemaID).Select("response_id").Scan(&list)
	return list, tx.Error
}

func BatchCreateRefSchemaResponses(ctx context.Context, list []*RefSchemaResponses) error {
	if len(list) == 0 {
		return nil
	}
	return model.DB(ctx).Create(list).Error
}

func BatchDelRefSchemaResponses(ctx context.Context, ids ...uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Where("id in ?", ids).Delete(&RefSchemaResponses{}).Error
}

// GetRefSchemaResponse 获取指定response引用指定schemas的引用关系
func GetRefSchemaResponse(ctx context.Context, responseID uint, schemaID ...uint) ([]*RefSchemaResponses, error) {
	var list []*RefSchemaResponses
	tx := model.DB(ctx).Where("ref_schema_id in ?", schemaID).Where("response_id = ?", responseID).Find(&list)
	return list, tx.Error
}

// GetRefSchemaResponses 获取所有responses引用指定schemas的引用关系
func GetRefSchemaResponses(ctx context.Context, schemaIDs ...uint) ([]*RefSchemaResponses, error) {
	var list []*RefSchemaResponses
	tx := model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Find(&list)
	return list, tx.Error
}

// DelRefSchemaResponse 删除指定response引用指定schemas的引用关系
// responseID 引用公共模型的公共响应ID
// schemaIDs 被引用的所有公共模型ID，因为ref_schema_id是索引字段，导致删除时需要传入所有被引用的公共模型ID
// 用于删除公共响应时，删除该公共响应引用的公共模型
func DelRefSchemaResponse(ctx context.Context, responseID uint, schemaIDs ...uint) error {
	return model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Where("response_id = ?", responseID).Delete(&RefSchemaResponses{}).Error
}

// DelRefSchemaResponses 删除所有responses引用指定schema的引用关系
// 用于删除公共模型时，删除所有引用了该公共模型的公共响应
func DelRefSchemaResponses(ctx context.Context, schemaID uint) error {
	return model.DB(ctx).Where("ref_schema_id = ?", schemaID).Delete(&RefSchemaResponses{}).Error
}
