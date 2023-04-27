package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/apicat/apicat/commom/spec"
	"gorm.io/gorm"
)

type Definitions struct {
	ID           uint   `gorm:"type:integer primary key autoincrement"`
	ProjectId    uint   `gorm:"index;not null;comment:项目id"`
	ParentId     uint   `gorm:"not null;comment:父级id"`
	Name         string `gorm:"type:varchar(255);not null;comment:名称"`
	Description  string `gorm:"type:varchar(255);comment:描述"`
	Type         string `gorm:"type:varchar(255);not null;comment:类型:category,schema"`
	Schema       string `gorm:"type:mediumtext;comment:内容"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	CreatedBy    uint `gorm:"not null;default:0;comment:创建人id"`
	UpdatedAt    time.Time
	UpdatedBy    uint `gorm:"not null;default:0;comment:最后更新人id"`
	DeletedAt    *gorm.DeletedAt
	DeletedBy    uint `gorm:"not null;default:0;comment:删除人id"`
}

func NewDefinitions(ids ...uint) (*Definitions, error) {
	definition := &Definitions{}
	if len(ids) > 0 {
		if err := Conn.Take(definition, ids[0]).Error; err != nil {
			return definition, err
		}
		return definition, nil
	}
	return definition, nil
}

func (d *Definitions) List() ([]Definitions, error) {
	tx := Conn.Where("project_id = ?", d.ProjectId)
	if d.ParentId > 0 {
		tx = tx.Where("parent_id = ?", d.ParentId)
	}
	if d.Name != "" {
		tx = tx.Where("name = ?", d.Name)
	}
	if d.Type != "" {
		tx = tx.Where("type = ?", d.Type)
	}

	var definitions []Definitions
	return definitions, tx.Order("display_order asc").Order("id desc").Find(&definitions).Error
}

func (d *Definitions) Create() error {
	return Conn.Create(d).Error
}

func (d *Definitions) Save() error {
	return Conn.Save(d).Error
}

func (d *Definitions) Delete() error {
	if d.Type == "category" {
		Conn.Where("parent_id = ?", d.ID).Delete(&Definitions{})
	}
	return Conn.Delete(d).Error
}

func (d *Definitions) Creator() string {
	return ""
}

func (d *Definitions) Updater() string {
	return ""
}

func (d *Definitions) Deleter() string {
	return ""
}

func DefinitionsImport(projectID uint, schemas spec.Schemas) nameToIdMap {
	SchemasMap := make(nameToIdMap)

	if schemas == nil {
		return SchemasMap
	}

	for i, schema := range schemas {
		if schemaStr, err := json.Marshal(schema.Schema); err == nil {
			record := &Definitions{
				ProjectId:    projectID,
				ParentId:     0,
				Name:         schema.Name,
				Description:  schema.Description,
				Type:         "schema",
				Schema:       string(schemaStr),
				DisplayOrder: i,
			}

			if Conn.Create(record).Error == nil {
				SchemasMap[schema.Name] = record.ID
			}
		}
	}

	return SchemasMap
}

// func DefinitionsExport(projectID uint) spec.Schemas {
// 	var definitions []*Definitions
// 	specDefinitions := make(spec.Schemas, 0)

// 	if err := Conn.Where("project_id = ? AND type = ?", projectID, "schema").Find(&definitions).Error; err == nil {
// 		for _, definition := range definitions {
// 			schema := &spec.Schema{
// 				Schema: &jsonschema.Schema{},
// 			}
// 			if json.Unmarshal([]byte(definition.Schema), schema.Schema) == nil {
// 				schema.Name = definition.Name
// 				schema.Description = definition.Description
// 				specDefinitions = append(specDefinitions, schema)
// 			}
// 		}
// 	}
// 	return specDefinitions
// }

func DefinitionsDdereference(d *Definitions) error {
	ref := "{\"$ref\":\"#/definitions/schemas/" + strconv.FormatUint(uint64(d.ID), 10) + "\"}"

	definitions, _ := NewDefinitions()
	definitions.ProjectId = d.ProjectId
	definitionsList, err := definitions.List()
	if err != nil {
		return err
	}

	for _, definitions := range definitionsList {
		if strings.Contains(definitions.Schema, ref) {
			newContent := strings.Replace(definitions.Schema, ref, d.Schema, -1)
			definitions.Schema = newContent

			if err := definitions.Save(); err != nil {
				return err
			}
		}
	}

	return nil
}
