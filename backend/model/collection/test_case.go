package collection

import (
	"context"
	"errors"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
)

type TestCase struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID    string `gorm:"type:varchar(24);index:idx_pid_cid;not null;comment:project id"`
	CollectionID uint   `gorm:"type:bigint;index:idx_pid_cid;not null;comment:collection id"`
	Title        string `gorm:"type:varchar(255);not null;comment:test case title"`
	Content      string `gorm:"type:mediumtext;comment:test case content"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (t *TestCase) Create(ctx context.Context) error {
	return model.DB(ctx).Create(t).Error
}

func (t *TestCase) CreateWithoutCtx() error {
	return model.DBWithoutCtx().Create(t).Error
}

func (t *TestCase) Get(ctx context.Context) (bool, error) {
	if t.ID == 0 {
		return false, errors.New("id cannot be 0")
	}
	tx := model.DB(ctx).Take(t, "id = ?", t.ID)
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

func (t *TestCase) Update(ctx context.Context, title, content string) error {
	return model.DB(ctx).Model(t).Updates(map[string]interface{}{
		"title":   title,
		"content": content,
	}).Error
}

func (t *TestCase) Delete(ctx context.Context) error {
	return model.DB(ctx).Delete(t).Error
}
