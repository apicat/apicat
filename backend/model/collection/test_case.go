package collection

import (
	"apicat-cloud/backend/model"
	"context"
	"errors"
	"time"
)

type TestCase struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID    string `gorm:"type:varchar(24);index:idx_pid_cid;not null;comment:项目id"`
	CollectionID uint   `gorm:"type:bigint;index:idx_pid_cid;not null;comment:集合id"`
	Title        string `gorm:"type:varchar(255);not null;comment:测试用例标题"`
	Content      string `gorm:"type:mediumtext;comment:测试用例内容"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func init() {
	model.RegMigrate(&TestCase{})
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
