package collection

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/team"
)

type CollectionHistory struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	CollectionID uint   `gorm:"type:bigint;index;not null;comment:集合id"`
	Title        string `gorm:"type:varchar(255);not null;comment:名称"`
	Content      string `gorm:"type:mediumtext;comment:内容"`
	CreatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:创建人id"`
	model.TimeModel
}

func init() {
	model.RegMigrate(&CollectionHistory{})
}

func (ch *CollectionHistory) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if ch.ID != 0 && ch.CollectionID != 0 {
		tx = tx.Take(ch, "id = ? AND collection_id = ?", ch.ID, ch.CollectionID)
	} else {
		return false, errors.New("query condition error")
	}
	return tx.Error == nil, model.NotRecord(tx)
}

func (ch *CollectionHistory) Create(ctx context.Context, memberID uint) error {
	var latestRecord CollectionHistory
	err := model.DB(ctx).Last(&latestRecord, "collection_id = ?", ch.CollectionID).Error
	if err == nil && latestRecord.CreatedBy == memberID && latestRecord.CreatedAt.Add(5*time.Minute).After(time.Now()) {
		return model.DB(ctx).Model(latestRecord).Updates(map[string]interface{}{
			"Title":   ch.Title,
			"Content": ch.Content,
		}).Error
	}

	ch.CreatedBy = memberID
	return model.DB(ctx).Create(ch).Error
}

func (ch *CollectionHistory) Restore(ctx context.Context, c *Collection, tm *team.TeamMember) error {
	newCH := &CollectionHistory{
		CollectionID: c.ID,
		Title:        c.Title,
		Content:      c.Content,
	}
	if err := newCH.Create(ctx, tm.ID); err != nil {
		slog.ErrorContext(ctx, "CollectionHistory.Restore.Create", "err", err)
	}

	return c.Update(ctx, ch.Title, ch.Content, tm.ID)
}
