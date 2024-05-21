package referencerelation

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
)

type RefResponseCollections struct {
	ID             uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	RefResponserID uint `gorm:"type:bigint;index;not null;comment:被引用的公共响应id"`
	CollectionID   uint `gorm:"type:bigint;not null;comment:引用ref_responser_id的文档id"`
}

func (r *RefResponseCollections) GetCollections(ctx context.Context) ([]*RefResponseCollections, error) {
	var list []*RefResponseCollections
	tx := model.DB(ctx).Where("ref_responser_id = ?", r.RefResponserID).Find(&list)
	return list, tx.Error
}

func (r *RefResponseCollections) GetCollectionIDs(ctx context.Context) ([]uint, error) {
	var list []uint
	tx := model.DB(ctx).Model(&RefResponseCollections{}).Where("ref_responser_id = ?", r.RefResponserID).Select("collection_id").Scan(&list)
	return list, tx.Error
}

func BatchCreateRefResponseCollections(ctx context.Context, list []*RefResponseCollections) error {
	if len(list) == 0 {
		return nil
	}
	return model.DB(ctx).Create(list).Error
}

func BatchDelRefResponseCollections(ctx context.Context, ids ...uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Where("id in ?", ids).Delete(&RefResponseCollections{}).Error
}

// GetRefResponseCollection 获取指定collection引用指定responses的引用关系
func GetRefResponseCollection(ctx context.Context, collectionID uint, responseIDs ...uint) ([]*RefResponseCollections, error) {
	var list []*RefResponseCollections
	tx := model.DB(ctx).Where("ref_responser_id in ?", responseIDs).Where("collection_id = ?", collectionID).Find(&list)
	return list, tx.Error
}

// GetRefResponseCollections 获取所有collections引用指定repsonses的引用关系
func GetRefResponseCollections(ctx context.Context, responseIDs ...uint) ([]*RefResponseCollections, error) {
	var list []*RefResponseCollections
	tx := model.DB(ctx).Where("ref_responser_id in ?", responseIDs).Find(&list)
	return list, tx.Error
}

// DelRefResponseCollection 删除指定collection引用指定responses的引用关系
// collectionID 引用公共响应的文档ID
// responseIDs 被引用的所有公共响应ID，因为ref_responser_id是索引字段，导致删除时需要传入所有被引用的公共响应ID
// 用于删除文档时，删除该文档引用的公共响应
func DelRefResponseCollection(ctx context.Context, collectionID uint, responseIDs ...uint) error {
	return model.DB(ctx).Where("ref_responser_id in ?", responseIDs).Where("collection_id = ?", collectionID).Delete(&RefResponseCollections{}).Error
}

// DelRefResponseCollections 删除所有collections引用指定repsonse的引用关系
// 用于删除公共响应时，删除所有引用了该公共响应的文档
func DelRefResponseCollections(ctx context.Context, responseID uint) error {
	return model.DB(ctx).Where("ref_responser_id = ?", responseID).Delete(&RefResponseCollections{}).Error
}
