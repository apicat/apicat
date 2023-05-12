package models

import (
	"encoding/json"
	"time"

	"github.com/apicat/apicat/common/spec"
)

type DefinitionParameters struct {
	ID        uint   `gorm:"type:integer primary key autoincrement"`
	ProjectID uint   `gorm:"type:integer;index;not null;comment:项目id"`
	Name      string `gorm:"type:varchar(255);not null;comment:参数名称"`
	Required  int    `gorm:"type:tinyint(1);not null;comment:是否必传:0-否,1-是"`
	Schema    string `gorm:"type:mediumtext;comment:参数内容"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDefinitionParameters(ids ...uint) (*DefinitionParameters, error) {
	definitionParameters := &DefinitionParameters{}
	if len(ids) > 0 {
		if err := Conn.Take(definitionParameters, ids[0]).Error; err != nil {
			return definitionParameters, err
		}
		return definitionParameters, nil
	}
	return definitionParameters, nil
}

func (dp *DefinitionParameters) List() ([]*DefinitionParameters, error) {
	definitionParametersQuery := Conn.Where("project_id = ?", dp.ProjectID)

	var definitionParameters []*DefinitionParameters
	return definitionParameters, definitionParametersQuery.Find(&definitionParameters).Error
}

func (dp *DefinitionParameters) Create() error {
	return Conn.Create(dp).Error
}

func (dp *DefinitionParameters) Save() error {
	return Conn.Save(dp).Error
}

func (dp *DefinitionParameters) Delete() error {
	return Conn.Delete(dp).Error
}

func DefinitionParametersImport(projectID uint, parameters spec.Schemas) nameToIdMap {
	parametersMap := nameToIdMap{}

	if len(parameters) == 0 {
		return parametersMap
	}

	for _, v := range parameters {
		if schema, err := json.Marshal(v.Schema); err == nil {
			required := 0
			if v.Required {
				required = 1
			}

			dp := &DefinitionParameters{
				ProjectID: projectID,
				Name:      v.Name,
				Required:  required,
				Schema:    string(schema),
			}

			if dp.Create() == nil {
				parametersMap[v.Name] = uint(v.ID)
			}
		}
	}

	return parametersMap
}
