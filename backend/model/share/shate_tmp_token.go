package share

import (
	"context"
	"errors"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/project"
)

type ShareTmpToken struct {
	ID           uint      `gorm:"type:bigint;primaryKey;autoIncrement"`
	ShareToken   string    `gorm:"type:varchar(255);index;not null;comment:share token"`
	Expiration   time.Time `gorm:"type:datetime;not null;comment:expiration time"`
	ProjectID    string    `gorm:"type:varchar(24);index;not null;comment:project id"`
	CollectionID uint      `gorm:"type:bigint;index;comment:collection id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (stt *ShareTmpToken) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if stt.ID != 0 {
		tx = tx.Take(stt, "id = ?", stt.ID)
	} else if stt.ShareToken != "" {
		tx = tx.Take(stt, "share_token = ?", stt.ShareToken)
	} else {
		return false, errors.New("query condition error")
	}
	return tx.Error == nil, model.NotRecord(tx)
}

func (stt *ShareTmpToken) Create(ctx context.Context) error {
	return model.DB(ctx).Create(stt).Error
}

func (stt *ShareTmpToken) Delete(ctx context.Context) error {
	return model.DB(ctx).Delete(stt).Error
}

func DeleteShareTmpTokenByProject(ctx context.Context, p *project.Project) error {
	return model.DB(ctx).Delete(&ShareTmpToken{}, "project_id = ?", p.ID).Error
}

func DeleteShareTmpTokenByCollection(ctx context.Context, c *collection.Collection) error {
	return model.DB(ctx).Delete(&ShareTmpToken{}, "collection_id = ?", c.ID).Error
}

// DeleteProjectShareTmpTokens 项目关闭分享清除分享token
func DeleteProjectShareTmpTokens(ctx context.Context, p *project.Project) error {
	return model.DB(ctx).Delete(&ShareTmpToken{}, "project_id = ? AND collection_id = 0", p.ID).Error
}

// DeleteCollectionShareTmpTokens 集合关闭分享清除分享token
func DeleteCollectionShareTmpTokens(ctx context.Context, cIDs ...uint) error {
	return model.DB(ctx).Delete(&ShareTmpToken{}, "collection_id IN (?)", cIDs).Error
}
