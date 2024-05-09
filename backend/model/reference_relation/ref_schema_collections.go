package referencerelation

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
)

type RefSchemaCollections struct {
	ID           uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	RefSchemaID  uint `gorm:"type:bigint;index;not null;comment:被引用的公共模型id"`
	CollectionID uint `gorm:"type:bigint;not null;comment:引用ref_schema_id的文档id"`
}

func (r *RefSchemaCollections) GetCollections(ctx context.Context) ([]*RefSchemaCollections, error) {
	var list []*RefSchemaCollections
	tx := model.DB(ctx).Where("ref_schema_id = ?", r.RefSchemaID).Find(&list)
	return list, tx.Error
}

func (r *RefSchemaCollections) GetCollectionIDs(ctx context.Context) ([]uint, error) {
	var list []uint
	tx := model.DB(ctx).Where("ref_schema_id = ?", r.RefSchemaID).Select("collection_id").Find(&list)
	return list, tx.Error
}

func BatchCreateRefSchemaCollections(ctx context.Context, list []*RefSchemaCollections) error {
	if len(list) == 0 {
		return nil
	}
	return model.DB(ctx).Create(list).Error
}

func BatchDelRefSchemaCollections(ctx context.Context, ids ...uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Where("id in ?", ids).Delete(&RefSchemaCollections{}).Error
}

// GetRefSchemaCollections 获取引用指定公共模型的所有文档的引用关系
func GetRefSchemaCollections(ctx context.Context, schemaIDs ...uint) ([]*RefSchemaCollections, error) {
	var list []*RefSchemaCollections
	tx := model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Find(&list)
	return list, tx.Error
}

// DelRefSchemaCollection 删除引用指定公共模型的指定文档的引用关系
// 用于删除文档时，删除该文档引用的公共模型
func DelRefSchemaCollection(ctx context.Context, schemaIDs []uint, collectionID uint) error {
	return model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Where("collection_id = ?", collectionID).Delete(&RefSchemaCollections{}).Error
}

// DelRefSchemaCollections 删除引用指定公共模型的所有文档的引用关系
// 用于删除公共模型时，删除所有引用了该公共模型的文档
func DelRefSchemaCollections(ctx context.Context, schemaIDs uint) error {
	return model.DB(ctx).Where("ref_schema_id = ?", schemaIDs).Delete(&RefSchemaCollections{}).Error
}
