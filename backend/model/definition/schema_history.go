package definition

import (
	"apicat-cloud/backend/model"
	"context"
	"errors"
	"log/slog"
	"time"
)

type DefinitionSchemaHistory struct {
	ID          uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	SchemaID    uint   `gorm:"type:bigint;index;not null;comment:模型id"`
	Name        string `gorm:"type:varchar(255);not null;comment:名称"`
	Description string `gorm:"type:varchar(255);comment:描述"`
	Schema      string `gorm:"type:mediumtext;comment:内容"`
	CreatedBy   uint   `gorm:"type:bigint;not null;default:0;comment:创建人id"`
	model.TimeModel
}

func init() {
	model.RegMigrate(&DefinitionSchemaHistory{})
}

func (dsh *DefinitionSchemaHistory) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if dsh.ID != 0 && dsh.SchemaID != 0 {
		tx = tx.Take(dsh, "id = ? AND schema_id = ?", dsh.ID, dsh.SchemaID)
	} else {
		return false, errors.New("query condition error")
	}
	return tx.Error == nil, model.NotRecord(tx)
}

func (dsh *DefinitionSchemaHistory) Create(ctx context.Context, memberID uint) error {
	var latestRecord DefinitionSchemaHistory
	err := model.DB(ctx).Last(&latestRecord, "schema_id = ?", dsh.SchemaID).Error
	if err == nil && latestRecord.CreatedBy == memberID && latestRecord.CreatedAt.Add(5*time.Minute).After(time.Now()) {
		return model.DB(ctx).Model(latestRecord).Updates(map[string]interface{}{
			"Name":        dsh.Name,
			"Description": dsh.Description,
			"Schema":      dsh.Schema,
		}).Error
	}

	dsh.CreatedBy = memberID
	return model.DB(ctx).Create(dsh).Error
}

func (dsh *DefinitionSchemaHistory) Restore(ctx context.Context, ds *DefinitionSchema, memberID uint) error {
	newDSH := &DefinitionSchemaHistory{
		SchemaID:    ds.ID,
		Name:        ds.Name,
		Description: ds.Description,
		Schema:      ds.Schema,
	}
	if err := newDSH.Create(ctx, memberID); err != nil {
		slog.ErrorContext(ctx, "DefinitionSchemaHistory.Restore.Create", "err", err)
	}

	return ds.Update(ctx, dsh.Name, dsh.Description, dsh.Schema, memberID)
}
