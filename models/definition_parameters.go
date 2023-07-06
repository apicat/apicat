package models

import (
	"encoding/json"
	"time"

	"github.com/apicat/apicat/common/spec"
	"github.com/apicat/apicat/common/spec/jsonschema"
)

type DefinitionParameters struct {
	ID        uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID uint   `gorm:"type:bigint;index;not null;comment:项目id"`
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

func DefinitionParametersImport(projectID uint, parameters spec.Schemas) virtualIDToIDMap {
	parametersMap := virtualIDToIDMap{}

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
				parametersMap[v.ID] = uint(dp.ID)
			}
		}
	}

	return parametersMap
}

func DefinitionParametersExport(projectID uint) spec.Schemas {
	parameters := []*DefinitionParameters{}
	specParameters := spec.Schemas{}

	if err := Conn.Where("project_id = ?", projectID).Find(&parameters).Error; err != nil {
		return specParameters
	}

	for _, v := range parameters {
		schema := &jsonschema.Schema{}
		if err := json.Unmarshal([]byte(v.Schema), &schema); err == nil {
			required := false
			if v.Required == 1 {
				required = true
			}

			specParameters = append(specParameters, &spec.Schema{
				ID:       int64(v.ID),
				Name:     v.Name,
				Required: required,
				Schema:   schema,
			})
		}
	}

	return specParameters
}
