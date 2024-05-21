package referencerelation

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
)

type ExceptParamCollection struct {
	ID            uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	ExceptParamID uint `gorm:"type:bigint;index;not null;comment:被排除的全局参数id"`
	CollectionID  uint `gorm:"type:bigint;not null;comment:排除except_param_id的文档id"`
}

func (e *ExceptParamCollection) GetCollections(ctx context.Context) ([]*ExceptParamCollection, error) {
	var list []*ExceptParamCollection
	tx := model.DB(ctx).Where("except_param_id = ?", e.ExceptParamID).Find(&list)
	return list, tx.Error
}

func (e *ExceptParamCollection) GetCollectionIDs(ctx context.Context) ([]uint, error) {
	var list []uint
	tx := model.DB(ctx).Model(&ExceptParamCollection{}).Where("except_param_id = ?", e.ExceptParamID).Select("collection_id").Scan(&list)
	return list, tx.Error
}

func BatchCreateExceptParamCollections(ctx context.Context, list []*ExceptParamCollection) error {
	if len(list) == 0 {
		return nil
	}
	return model.DB(ctx).Create(list).Error
}

func BatchDelExceptParamCollections(ctx context.Context, ids ...uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Where("id in ?", ids).Delete(&ExceptParamCollection{}).Error
}

// GetExceptParamCollection 获取指定collection排除指定parameters的排除关系
func GetExceptParamCollection(ctx context.Context, collectionID uint, paramIDs ...uint) ([]*ExceptParamCollection, error) {
	var list []*ExceptParamCollection
	tx := model.DB(ctx).Where("except_param_id in ?", paramIDs).Where("collection_id = ?", collectionID).Find(&list)
	return list, tx.Error
}

// GetExceptParamCollections 获取所有collections排除指定parameters的排除关系
func GetExceptParamCollections(ctx context.Context, paramIDs ...uint) ([]*ExceptParamCollection, error) {
	var list []*ExceptParamCollection
	tx := model.DB(ctx).Where("except_param_id in ?", paramIDs).Find(&list)
	return list, tx.Error
}

// DelExceptParamCollection 删除指定collection排除指定parameters的排除关系
// collectionID 排除全局参数的文档ID
// paramIDs 被排除的所有全局参数ID，因为except_param_id是索引字段，导致删除时需要传入所有被排除的全局参数ID
// 用于删除文档时，删除该文档排除的全局参数
func DelExceptParamCollection(ctx context.Context, collectionID uint, paramIDs ...uint) error {
	return model.DB(ctx).Where("except_param_id in ?", paramIDs).Where("collection_id = ?", collectionID).Delete(&ExceptParamCollection{}).Error
}

// DelExceptParamCollections 删除所有collections排除指定repsonse的排除关系
// 用于删除全局参数时，删除所有排除该全局参数的文档
func DelExceptParamCollections(ctx context.Context, paramID uint) error {
	return model.DB(ctx).Where("except_param_id = ?", paramID).Delete(&ExceptParamCollection{}).Error
}
