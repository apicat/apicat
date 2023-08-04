package models

import "time"

type DefinitionSchemaHistories struct {
	ID          uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	SchemaID    uint   `gorm:"type:bigint;index;not null;comment:模型id"`
	Name        string `gorm:"type:varchar(255);not null;comment:名称"`
	Description string `gorm:"type:varchar(255);comment:描述"`
	Type        string `gorm:"type:varchar(255);not null;comment:类型:category,schema"`
	Schema      string `gorm:"type:mediumtext;comment:内容"`
	CreatedAt   time.Time
	CreatedBy   uint `gorm:"type:bigint;not null;default:0;comment:创建人id"`
}

func NewDefinitionSchemaHistories(ids ...uint) (*DefinitionSchemaHistories, error) {
	dsh := &DefinitionSchemaHistories{}
	if len(ids) > 0 {
		if err := Conn.Take(dsh, ids[0]).Error; err != nil {
			return dsh, err
		}
		return dsh, nil
	}
	return dsh, nil
}

func (dsh *DefinitionSchemaHistories) List(schemsIDs ...uint) ([]*DefinitionSchemaHistories, error) {
	var definitionSchemaHistories []*DefinitionSchemaHistories
	if len(schemsIDs) > 0 {
		return definitionSchemaHistories, Conn.Where("schema_id IN ?", schemsIDs).Order("created_at desc").Find(&definitionSchemaHistories).Error
	}
	return definitionSchemaHistories, Conn.Order("created_at desc").Find(&definitionSchemaHistories).Error
}

func (dsh *DefinitionSchemaHistories) Create() error {
	return Conn.Create(dsh).Error
}

func (dsh *DefinitionSchemaHistories) Restore(ds *DefinitionSchemas, uid uint) error {
	ndsh, _ := NewDefinitionSchemaHistories()
	ndsh.SchemaID = ds.ID
	ndsh.Name = ds.Name
	ndsh.Description = ds.Description
	ndsh.Type = ds.Type
	ndsh.Schema = ds.Schema
	ndsh.CreatedBy = uid
	if err := ndsh.Create(); err != nil {
		return err
	}

	ds.Name = ndsh.Name
	ds.Description = ndsh.Description
	ds.Schema = ndsh.Schema
	ds.UpdatedBy = uid
	if err := ds.Save(); err != nil {
		return err
	}

	return nil
}
