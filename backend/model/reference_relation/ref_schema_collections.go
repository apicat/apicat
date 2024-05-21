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

// GetRefSchemaCollection 获取指定collection引用指定schemas的引用关系
func GetRefSchemaCollection(ctx context.Context, collectionID uint, schemaIDs ...uint) ([]*RefSchemaCollections, error) {
	var list []*RefSchemaCollections
	tx := model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Where("collection_id = ?", collectionID).Find(&list)
	return list, tx.Error
}

// GetRefSchemaCollections 获取所有collections引用指定schemas的引用关系
func GetRefSchemaCollections(ctx context.Context, schemaIDs ...uint) ([]*RefSchemaCollections, error) {
	var list []*RefSchemaCollections
	tx := model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Find(&list)
	return list, tx.Error
}

// DelRefSchemaCollection 删除指定collection引用指定schema的引用关系
// collectionID 引用公共模型的文档ID
// schemaIDs 被引用的所有公共模型ID，因为ref_schema_id是索引字段，导致删除时需要传入所有被引用的公共模型ID
// 用于删除文档时，删除该文档引用的公共模型
func DelRefSchemaCollection(ctx context.Context, collectionID uint, schemaIDs ...uint) error {
	return model.DB(ctx).Where("ref_schema_id in ?", schemaIDs).Where("collection_id = ?", collectionID).Delete(&RefSchemaCollections{}).Error
}

// DelRefSchemaCollections 删除所有collections引用指定schema的引用关系
// 用于删除公共模型时，删除所有引用了该公共模型的文档
func DelRefSchemaCollections(ctx context.Context, schemaID uint) error {
	return model.DB(ctx).Where("ref_schema_id = ?", schemaID).Delete(&RefSchemaCollections{}).Error
}
