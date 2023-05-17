package models

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apicat/apicat/common/apicat_struct"
	"github.com/apicat/apicat/common/spec"
	"github.com/apicat/apicat/common/spec/jsonschema"
	"gorm.io/gorm"
)

type DefinitionSchemas struct {
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

func NewDefinitionSchemas(ids ...uint) (*DefinitionSchemas, error) {
	definition := &DefinitionSchemas{}
	if len(ids) > 0 {
		if err := Conn.Take(definition, ids[0]).Error; err != nil {
			return definition, err
		}
		return definition, nil
	}
	return definition, nil
}

func (d *DefinitionSchemas) List() ([]DefinitionSchemas, error) {
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

	var definitions []DefinitionSchemas
	return definitions, tx.Order("display_order asc").Order("id desc").Find(&definitions).Error
}

func (d *DefinitionSchemas) Get() error {
	tx := Conn.Where("project_id = ?", d.ProjectId)
	if d.Name != "" {
		tx = tx.Where("name = ?", d.Name)
	}
	if d.Type != "" {
		tx = tx.Where("type = ?", d.Type)
	}
	return tx.Take(d).Error
}

func (d *DefinitionSchemas) Create() error {
	return Conn.Create(d).Error
}

func (d *DefinitionSchemas) Save() error {
	return Conn.Save(d).Error
}

func (d *DefinitionSchemas) Delete() error {
	if d.Type == "category" {
		Conn.Where("parent_id = ?", d.ID).Delete(&DefinitionSchemas{})
	}
	return Conn.Delete(d).Error
}

func (d *DefinitionSchemas) Creator() string {
	return ""
}

func (d *DefinitionSchemas) Updater() string {
	return ""
}

func (d *DefinitionSchemas) Deleter() string {
	return ""
}

func DefinitionSchemasImport(projectID uint, schemas spec.Schemas) virtualIDToIDMap {
	schemasMap := virtualIDToIDMap{}

	if len(schemas) == 0 {
		return schemasMap
	}

	for i, schema := range schemas {
		if schemaStr, err := json.Marshal(schema.Schema); err == nil {
			record := &DefinitionSchemas{
				ProjectId: projectID,
				Name:      schema.Name,
				Type:      "schema",
			}
			if record.Get() == nil {
				schemasMap[schema.ID] = record.ID
			} else {
				ds := &DefinitionSchemas{
					ProjectId:    projectID,
					ParentId:     0,
					Name:         schema.Name,
					Description:  schema.Description,
					Type:         "schema",
					Schema:       string(schemaStr),
					DisplayOrder: i,
				}

				if ds.Create() == nil {
					schemasMap[schema.ID] = ds.ID
				}
			}
		}
	}

	definitionschemas := []*DefinitionSchemas{}
	if err := Conn.Where("project_id = ? AND type = ?", projectID, "schema").Find(&definitionschemas).Error; err != nil {
		return schemasMap
	}

	for _, v := range definitionschemas {
		schema := replaceVirtualIDToID(v.Schema, schemasMap, "#/definitions/schemas/")
		v.Schema = schema
		v.Save()
	}

	return schemasMap
}

func DefinitionSchemasExport(projectID uint) spec.Schemas {
	definitions := []*DefinitionSchemas{}
	specDefinitionSchemas := spec.Schemas{}

	if err := Conn.Where("project_id = ? AND type = ?", projectID, "schema").Find(&definitions).Error; err != nil {
		return specDefinitionSchemas
	}

	for _, v := range definitions {
		schema := &jsonschema.Schema{}
		if err := json.Unmarshal([]byte(v.Schema), &schema); err == nil {
			specDefinitionSchemas = append(specDefinitionSchemas, &spec.Schema{
				ID:          int64(v.ID),
				Name:        v.Name,
				Description: v.Description,
				Schema:      schema,
			})
		}
	}

	return specDefinitionSchemas
}

func DefinitionIdToName(content string, idToNameMap IdToNameMap) string {
	re := regexp.MustCompile(`#/definitions/schemas/\d+`)
	reID := regexp.MustCompile(`\d+`)

	for {
		match := re.FindString(content)
		if match == "" {
			break
		}

		schemasIDStr := reID.FindString(match)
		if schemasIDInt, err := strconv.Atoi(schemasIDStr); err == nil {
			schemasID := uint(schemasIDInt)
			content = strings.Replace(content, match, "#/definitions/schemas/"+idToNameMap[schemasID], -1)
		}
	}

	return content
}

func DefinitionsSchemaUnRefByDefinitionsSchema(d *DefinitionSchemas, isUnRef int) error {
	ref := "{\"$ref\":\"#/definitions/schemas/" + strconv.FormatUint(uint64(d.ID), 10) + "\"}"

	definitions, _ := NewDefinitionSchemas()
	definitions.ProjectId = d.ProjectId
	definitionsList, err := definitions.List()
	if err != nil {
		return err
	}

	sourceJson := map[string]interface{}{}
	if err := json.Unmarshal([]byte(d.Schema), &sourceJson); err != nil {
		return err
	}
	typeEmptyStructure := apicat_struct.TypeEmptyStructure()

	for _, definitions := range definitionsList {
		if strings.Contains(definitions.Schema, ref) {
			newStr := typeEmptyStructure[sourceJson["type"].(string)]
			if isUnRef == 1 {
				newStr = d.Schema
			}

			newContent := strings.Replace(definitions.Schema, ref, newStr, -1)
			definitions.Schema = newContent

			if err := definitions.Save(); err != nil {
				return err
			}
		}
	}

	return nil
}

func DefinitionsSchemaUnRefByDefinitionsResponse(d *DefinitionSchemas, isUnRef int) error {
	ref := "{\"$ref\":\"#/definitions/schemas/" + strconv.FormatUint(uint64(d.ID), 10) + "\"}"

	definitionResponses, _ := NewDefinitionResponses()
	definitionResponses.ProjectID = d.ProjectId
	definitionResponsesList, err := definitionResponses.List()
	if err != nil {
		return err
	}

	sourceJson := map[string]interface{}{}
	if err := json.Unmarshal([]byte(d.Schema), &sourceJson); err != nil {
		return err
	}
	typeEmptyStructure := apicat_struct.TypeEmptyStructure()

	for _, definitionResponse := range definitionResponsesList {
		if strings.Contains(definitionResponse.Content, ref) {
			newStr := typeEmptyStructure[sourceJson["type"].(string)]
			if isUnRef == 1 {
				newStr = d.Schema
			}

			newContent := strings.Replace(definitionResponse.Content, ref, newStr, -1)
			definitionResponse.Content = newContent

			if err := definitionResponse.Update(); err != nil {
				return err
			}
		}
	}

	return nil
}

func DefinitionsSchemaUnRefByCollections(d *DefinitionSchemas, isUnRef int) error {
	ref := "{\"$ref\":\"#/definitions/schemas/" + strconv.FormatUint(uint64(d.ID), 10) + "\"}"

	collections, _ := NewCollections()
	collections.ProjectId = d.ProjectId
	collectionList, err := collections.List()
	if err != nil {
		return err
	}

	sourceJson := map[string]interface{}{}
	if err := json.Unmarshal([]byte(d.Schema), &sourceJson); err != nil {
		return err
	}
	typeEmptyStructure := apicat_struct.TypeEmptyStructure()

	for _, collection := range collectionList {
		if strings.Contains(collection.Content, ref) {
			newStr := typeEmptyStructure[sourceJson["type"].(string)]
			if isUnRef == 1 {
				newStr = d.Schema
			}

			newContent := strings.Replace(collection.Content, ref, newStr, -1)
			collection.Content = newContent

			if err := collection.Update(); err != nil {
				return err
			}
		}
	}

	return nil
}
