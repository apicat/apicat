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
	tx := model.DB(ctx).Where("ref_schema_id = ?", r.RefSchemaID).Select("response_id").Find(&list)
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

// GetRefSchemaResponses 获取引用指定公共模型的所有公共响应的引用关系
func GetRefSchemaResponses(ctx context.Context, schemaIDs ...uint) ([]*RefSchemaResponses, error) {
	var list []*RefSchemaResponses
	tx := model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Find(&list)
	return list, tx.Error
}

// DelRefSchemaResponse 删除引用指定公共模型s的指定公共响应的引用关系
// 用于删除公共响应时，删除该公共响应引用的公共模型
func DelRefSchemaResponse(ctx context.Context, schemaIDs []uint, responseID uint) error {
	return model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Where("response_id = ?", responseID).Delete(&RefSchemaResponses{}).Error
}

// DelRefSchemaResponses 删除引用指定公共模型的所有公共响应的引用关系
// 用于删除公共模型时，删除所有引用了该公共模型的公共响应
func DelRefSchemaResponses(ctx context.Context, schemaIDs uint) error {
	return model.DB(ctx).Where("ref_schema_id = ?", schemaIDs).Delete(&RefSchemaResponses{}).Error
}
