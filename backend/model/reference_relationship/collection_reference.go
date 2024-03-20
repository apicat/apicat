package referencerelationship

import (
	"context"
	"time"

	"github.com/apicat/apicat/backend/model"
)

const (
	ReferenceSchema    = "schema"
	ReferenceResponse  = "response"
	ReferenceParameter = "parameter"
)

type CollectionReference struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID    string `gorm:"type:varchar(24);index;not null;comment:项目id"`
	CollectionID uint   `gorm:"type:bigint;index;not null;comment:集合id"`
	RefID        uint   `gorm:"type:bigint;index;not null;comment:引用节点id"`
	RefType      string `gorm:"type:varchar(255);not null;comment:引用节点类型:schema,response,parameter"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func init() {
	model.RegMigrate(&CollectionReference{}, &SchemaReference{}, &ResponseReference{}, &ParameterExcept{})
}

func BatchCreateCollectionReference(ctx context.Context, list []*CollectionReference) error {
	if len(list) == 0 {
		return nil
	}
	return model.DB(ctx).Create(list).Error
}

func BatchDeleteCollectionReference(ctx context.Context, ids ...uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Where("id in ?", ids).Delete(&CollectionReference{}).Error
}

func GetCollectionReferencesByCollection(ctx context.Context, projectID string, collectionID uint, refType ...string) ([]*CollectionReference, error) {
	var list []*CollectionReference
	tx := model.DB(ctx).Where("project_id = ? AND collection_id = ?", projectID, collectionID)
	if len(refType) > 0 {
		tx = tx.Where("ref_type in ?", refType)
	}
	tx = tx.Find(&list)
	return list, tx.Error
}

func GetCollectionReferencesByRef(ctx context.Context, projectID string, refID uint, refType string) ([]*CollectionReference, error) {
	var list []*CollectionReference
	tx := model.DB(ctx).Where("project_id = ? AND ref_id = ? AND ref_type = ?", projectID, refID, refType).Find(&list)
	return list, tx.Error
}

func GetCollectionRefByCIDs(ctx context.Context, projectID string, collectionIDs []uint) ([]*CollectionReference, error) {
	var list []*CollectionReference
	tx := model.DB(ctx).Where("project_id = ? AND collection_id in ?", projectID, collectionIDs).Find(&list)
	return list, tx.Error
}
