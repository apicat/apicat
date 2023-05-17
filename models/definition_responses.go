package models

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apicat/apicat/common/apicat_struct"
	"github.com/apicat/apicat/common/spec"
)

type DefinitionResponses struct {
	ID           uint   `gorm:"type:integer primary key autoincrement"`
	ProjectID    uint   `gorm:"type:integer;index;not null;comment:项目id"`
	Name         string `gorm:"type:varchar(255);not null;comment:响应名称"`
	Description  string `gorm:"type:varchar(255);not null;comment:状态描述"`
	Type         string `gorm:"type:varchar(255);not null;comment:响应类型:category,response"`
	Header       string `gorm:"type:mediumtext;comment:响应头"`
	Content      string `gorm:"type:mediumtext;comment:响应内容"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewDefinitionResponses(ids ...uint) (*DefinitionResponses, error) {
	definitionResponses := &DefinitionResponses{}
	if len(ids) > 0 {
		if err := Conn.Take(definitionResponses, ids[0]).Error; err != nil {
			return definitionResponses, err
		}
		return definitionResponses, nil
	}
	return definitionResponses, nil
}

func (dr *DefinitionResponses) List() ([]*DefinitionResponses, error) {
	definitionResponsesQuery := Conn.Where("project_id = ?", dr.ProjectID)

	var definitionResponses []*DefinitionResponses
	return definitionResponses, definitionResponsesQuery.Find(&definitionResponses).Error
}

func (dr *DefinitionResponses) GetCountByName() (int64, error) {
	var count int64
	return count, Conn.Model(&DefinitionResponses{}).Where("project_id = ? and name = ?", dr.ProjectID, dr.Name).Count(&count).Error
}

func (dr *DefinitionResponses) GetCountExcludeTheID() (int64, error) {
	var count int64
	return count, Conn.Model(&DefinitionResponses{}).Where("project_id = ? and name = ? and id != ?", dr.ProjectID, dr.Name, dr.ID).Count(&count).Error
}

func (dr *DefinitionResponses) Create() error {
	return Conn.Create(dr).Error
}

func (dr *DefinitionResponses) Update() error {
	return Conn.Save(dr).Error
}

func (dr *DefinitionResponses) Delete() error {
	return Conn.Delete(dr).Error
}

func DefinitionResponsesImport(projectID uint, responses spec.HTTPResponseDefines) virtualIDToIDMap {
	responsesMap := virtualIDToIDMap{}

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

	if err := Conn.Where("project_id = ?", projectID).Find(&definitionResponses).Error; err != nil {
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

func DefinitionsResponseUnRef(dr *DefinitionResponses) error {
	ref := "\"$ref\":\"#/definitions/responses/" + strconv.Itoa(int(dr.ID)) + "\""

	collections, _ := NewCollections()
	collections.ProjectId = dr.ProjectID
	collectionList, err := collections.List()
	if err != nil {
		return err
	}

	header := []interface{}{}
	if err := json.Unmarshal([]byte(dr.Header), &header); err != nil {
		return err
	}

	content := map[string]interface{}{}
	if err := json.Unmarshal([]byte(dr.Content), &content); err != nil {
		return err
	}

	data := map[string]interface{}{
		"name":        dr.Name,
		"description": dr.Description,
		"header":      header,
		"content":     content,
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	newStr := string(dataJson)[1 : len(string(dataJson))-1]

	for _, collection := range collectionList {
		if strings.Contains(collection.Content, ref) {
			newContent := strings.Replace(collection.Content, ref, newStr, -1)
			collection.Content = newContent

			if err := collection.Update(); err != nil {
				return err
			}
		}
	}

	return nil
}

func DefinitionsResponseDelRef(dr *DefinitionResponses) error {
	re1 := regexp.MustCompile(`,{"code":\d+,"\$ref":"#/definitions/responses/` + strconv.Itoa(int(dr.ID)) + `"}`)
	re2 := regexp.MustCompile(`{"code":\d+,"\$ref":"#/definitions/responses/` + strconv.Itoa(int(dr.ID)) + `"}`)

	collections, _ := NewCollections()
	collections.ProjectId = dr.ProjectID
	collectionList, err := collections.List()
	if err != nil {
		return err
	}

	emptyResponse := apicat_struct.TypeEmptyStructure()["response"]

	for _, collection := range collectionList {
		matchRe1 := re1.FindString(collection.Content)
		if matchRe1 != "" {
			newContent := strings.Replace(collection.Content, matchRe1, "", -1)
			collection.Content = newContent
		} else {
			matchRe2 := re2.FindString(collection.Content)
			if matchRe2 != "" {
				newContent := strings.Replace(collection.Content, matchRe2, emptyResponse, -1)
				collection.Content = newContent
			}
		}

		if err := collection.Update(); err != nil {
			return err
		}
	}
	return nil
}
