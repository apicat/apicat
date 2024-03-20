package referencerelationship

import (
	"context"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
)

type ResponseReference struct {
	ID          uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID   string `gorm:"type:varchar(24);index;not null;comment:项目id"`
	ResponseID  uint   `gorm:"type:bigint;index;not null;comment:公共响应id"`
	RefSchemaID uint   `gorm:"type:bigint;index;not null;comment:引用的公共模型id"`
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
