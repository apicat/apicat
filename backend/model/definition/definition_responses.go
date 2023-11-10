package definition

import (
	"encoding/json"
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/module/spec"
	"time"
)

type DefinitionResponses struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID    uint   `gorm:"type:bigint;index;not null;comment:项目id"`
	Name         string `gorm:"type:varchar(255);not null;comment:响应名称"`
	Description  string `gorm:"type:varchar(255);not null;comment:状态描述"`
	Type         string `gorm:"type:varchar(255);not null;comment:响应类型:category,response"`
	Header       string `gorm:"type:mediumtext;comment:响应头"`
	Content      string `gorm:"type:mediumtext;comment:响应内容"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func init() {
	model.RegMigrate(&DefinitionResponses{})
}

func NewDefinitionResponses(ids ...uint) (*DefinitionResponses, error) {
	definitionResponses := &DefinitionResponses{}
	if len(ids) > 0 {
		if err := model.Conn.Take(definitionResponses, ids[0]).Error; err != nil {
			return definitionResponses, err
		}
		return definitionResponses, nil
	}
	return definitionResponses, nil
}

func (dr *DefinitionResponses) List() ([]*DefinitionResponses, error) {
	definitionResponsesQuery := model.Conn.Where("project_id = ?", dr.ProjectID)

	var definitionResponses []*DefinitionResponses
	return definitionResponses, definitionResponsesQuery.Order("display_order asc").Order("id desc").Find(&definitionResponses).Error
}

func (dr *DefinitionResponses) GetCountByName() (int64, error) {
	var count int64
	return count, model.Conn.Model(&DefinitionResponses{}).Where("project_id = ? and name = ?", dr.ProjectID, dr.Name).Count(&count).Error
}

func (dr *DefinitionResponses) GetCountExcludeTheID() (int64, error) {
	var count int64
	return count, model.Conn.Model(&DefinitionResponses{}).Where("project_id = ? and name = ? and id != ?", dr.ProjectID, dr.Name, dr.ID).Count(&count).Error
}

func (dr *DefinitionResponses) Create() error {
	var node *DefinitionResponses
	if err := model.Conn.Where("project_id = ?", dr.ProjectID).Order("display_order desc").First(&node).Error; err == nil {
		dr.DisplayOrder = node.DisplayOrder + 1
	}

	return model.Conn.Create(dr).Error
}

func (dr *DefinitionResponses) Update() error {
	return model.Conn.Save(dr).Error
}

func (dr *DefinitionResponses) Delete() error {
	return model.Conn.Delete(dr).Error
}

func DefinitionResponsesImport(projectID uint, responses spec.HTTPResponseDefines) model.VirtualIDToIDMap {
	responsesMap := model.VirtualIDToIDMap{}

	if len(responses) == 0 {
		return responsesMap
	}

	for i, response := range responses {
		header := ""
		if response.Header != nil {
			if headerByte, err := json.Marshal(response.Header); err == nil {
				header = string(headerByte)
			}
		}

		content := ""
		if response.Content != nil {
			if contentByte, err := json.Marshal(response.Content); err == nil {
				content = string(contentByte)
			}
		}

		dr := &DefinitionResponses{
			ProjectID:    projectID,
			Name:         response.Name,
			Description:  response.Description,
			Header:       header,
			Content:      content,
			DisplayOrder: i,
		}

		if dr.Create() == nil {
			responsesMap[response.ID] = dr.ID
		}
	}

	return responsesMap
}

func DefinitionResponsesExport(projectID uint) spec.HTTPResponseDefines {
	definitionResponses := []*DefinitionResponses{}
	specResponseDefines := spec.HTTPResponseDefines{}

	if err := model.Conn.Where("project_id = ?", projectID).Find(&definitionResponses).Error; err != nil {
		return specResponseDefines
	}

	for _, definitionResponse := range definitionResponses {
		header := spec.Schemas{}
		if err := json.Unmarshal([]byte(definitionResponse.Header), &header); err != nil {
			continue
		}

		content := spec.HTTPBody{}
		if err := json.Unmarshal([]byte(definitionResponse.Content), &content); err != nil {
			continue
		}

		specResponseDefines = append(specResponseDefines, spec.HTTPResponseDefine{
			ID:          int64(definitionResponse.ID),
			Name:        definitionResponse.Name,
			Description: definitionResponse.Description,
			Header:      header,
			Content:     content,
		})
	}

	return specResponseDefines
}
