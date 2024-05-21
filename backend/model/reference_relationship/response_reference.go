package referencerelationship

import (
	"context"
	"errors"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
)

type ResponseReference struct {
	ID          uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	ResponseID  uint `gorm:"type:bigint;index;not null;comment:公共响应id"`
	RefSchemaID uint `gorm:"type:bigint;index;not null;comment:引用的公共模型id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func BatchCreateResponseReference(ctx context.Context, list []*ResponseReference) error {
	if len(list) == 0 {
		return nil
	}
	return model.DB(ctx).Create(list).Error
}

func BatchDeleteResponseReference(ctx context.Context, ids ...uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Where("id in ?", ids).Delete(&ResponseReference{}).Error
}

func GetResponseReferencesByResponse(ctx context.Context, projectID string, responseID uint) ([]*ResponseReference, error) {
	var list []*ResponseReference
	tx := model.DB(ctx).Where("project_id = ? AND response_id = ?", projectID, responseID).Find(&list)
	return list, tx.Error
}

func GetResponseReferencesByRefSchema(ctx context.Context, projectID string, refSchemaID uint) ([]*ResponseReference, error) {
	var list []*ResponseReference
	tx := model.DB(ctx).Where("project_id = ? AND ref_schema_id = ?", projectID, refSchemaID).Find(&list)
	return list, tx.Error
}

func (rr *ResponseReference) GetResponseRefs(ctx context.Context) ([]*ResponseReference, error) {
	var list []*ResponseReference

	tx := model.DB(ctx)
	if rr.ResponseID != 0 {
		tx = tx.Where("response_id = ?", rr.ResponseID)
	} else if rr.RefSchemaID != 0 {
		tx = tx.Where("ref_schema_id = ?", rr.RefSchemaID)
	} else {
		return nil, errors.New("query condition error")
	}

	return list, tx.Find(&list).Error
}

func (rr *ResponseReference) GetResponseIDsByRef(ctx context.Context) ([]uint, error) {
	var list []uint
	tx := model.DB(ctx).Where("ref_schema_id = ?", rr.RefSchemaID).Select("response_id").Find(&list)
	return list, tx.Error
}

func (rr *ResponseReference) DelByResponseID(ctx context.Context) error {
	return model.DB(ctx).Where("response_id = ?", rr.ResponseID).Delete(&ResponseReference{}).Error
}

func (rr *ResponseReference) DelByRefSchemaID(ctx context.Context) error {
	return model.DB(ctx).Where("ref_schema_id = ?", rr.RefSchemaID).Delete(&ResponseReference{}).Error
}
