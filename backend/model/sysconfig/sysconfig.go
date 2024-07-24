package sysconfig

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
)

type Sysconfig struct {
	ID        uint   `gorm:"primarykey"`
	Type      string `gorm:"type:varchar(255);uniqueIndex:ukey;not null;comment:Configuration type"`
	Driver    string `gorm:"type:varchar(255);uniqueIndex:ukey;not null"`
	BeingUsed bool   `gorm:"type:tinyint;comment:is using"`
	Config    string `gorm:"type:varchar(512);"`
	Extra     string `gorm:"type:varchar(512);"`
}

func (o *Sysconfig) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx).Take(o, "type = ? and driver = ?", o.Type, o.Driver)
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

func (o *Sysconfig) GetByUse(ctx context.Context) (bool, error) {
	tx := model.DB(ctx).Take(o, "type = ? and being_used = ?", o.Type, 1)
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

func (o *Sysconfig) Update(ctx context.Context) error {
	return model.DB(ctx).Save(o).Error
}

func (o *Sysconfig) Create(ctx context.Context) error {
	return model.DB(ctx).Create(o).Error
}
