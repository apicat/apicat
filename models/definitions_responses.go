package models

import "time"

type DefinitionsResponses struct {
	ID           uint   `gorm:"type:integer primary key autoincrement"`
	ProjectID    uint   `gorm:"type:integer;index;not null;comment:项目id"`
	Name         string `gorm:"type:varchar(255);not null;comment:响应名称"`
	Code         int    `gorm:"type:int(11);not null;comment:Http status code"`
	Description  string `gorm:"type:varchar(255);not null;comment:状态描述"`
	Header       string `gorm:"type:mediumtext;comment:响应头"`
	Content      string `gorm:"type:mediumtext;comment:响应内容"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewDefinitionsResponses(ids ...uint) (*DefinitionsResponses, error) {
	definitionsResponses := &DefinitionsResponses{}
	if len(ids) > 0 {
		if err := Conn.Take(definitionsResponses, ids[0]).Error; err != nil {
			return definitionsResponses, err
		}
		return definitionsResponses, nil
	}
	return definitionsResponses, nil
}

func (dr *DefinitionsResponses) List() ([]*DefinitionsResponses, error) {
	definitionsResponsesQuery := Conn.Where("project_id = ?", dr.ProjectID)

	var definitionsResponses []*DefinitionsResponses
	return definitionsResponses, definitionsResponsesQuery.Find(&definitionsResponses).Error
}

func (dr *DefinitionsResponses) GetCountByName() (int64, error) {
	var count int64
	return count, Conn.Model(&DefinitionsResponses{}).Where("project_id = ? and name = ?", dr.ProjectID, dr.Name).Count(&count).Error
}

func (dr *DefinitionsResponses) GetCountExcludeTheID() (int64, error) {
	var count int64
	return count, Conn.Model(&DefinitionsResponses{}).Where("project_id = ? and name = ? and id != ?", dr.ProjectID, dr.Name, dr.ID).Count(&count).Error
}

func (dr *DefinitionsResponses) Create() error {
	return Conn.Create(dr).Error
}

func (dr *DefinitionsResponses) Update() error {
	return Conn.Save(dr).Error
}

func (dr *DefinitionsResponses) Delete() error {
	return Conn.Delete(dr).Error
}
