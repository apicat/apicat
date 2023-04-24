package models

import (
	"encoding/json"
	"time"

	"github.com/apicat/apicat/commom/spec"
)

type CommonResponses struct {
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

func NewCommonResponses(ids ...uint) (*CommonResponses, error) {
	commonResponses := &CommonResponses{}
	if len(ids) > 0 {
		if err := Conn.Take(commonResponses, ids[0]).Error; err != nil {
			return commonResponses, err
		}
		return commonResponses, nil
	}
	return commonResponses, nil
}

func (cr *CommonResponses) List() ([]*CommonResponses, error) {
	commonResponsesQuery := Conn.Where("project_id = ?", cr.ProjectID)

	var commonResponses []*CommonResponses
	return commonResponses, commonResponsesQuery.Find(&commonResponses).Error
}

func (cr *CommonResponses) GetCountByName() (int64, error) {
	var count int64
	return count, Conn.Model(&CommonResponses{}).Where("project_id = ? and name = ?", cr.ProjectID, cr.Name).Count(&count).Error
}

func (cr *CommonResponses) GetCountExcludeTheID() (int64, error) {
	var count int64
	return count, Conn.Model(&CommonResponses{}).Where("project_id = ? and name = ? and id != ?", cr.ProjectID, cr.Name, cr.ID).Count(&count).Error
}

func (cr *CommonResponses) Create() error {
	return Conn.Create(cr).Error
}

func (cr *CommonResponses) Update() error {
	return Conn.Save(cr).Error
}

func (cr *CommonResponses) Delete() error {
	return Conn.Delete(cr).Error
}

func CommonResponsesImport(projectID uint, responses spec.HTTPResponses) nameToIdMap {
	var ResponsesMap nameToIdMap

	if responses == nil {
		return ResponsesMap
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

		record := &CommonResponses{
			ProjectID:    projectID,
			Name:         response.Name,
			Code:         response.Code,
			Description:  response.Description,
			Header:       header,
			Content:      content,
			DisplayOrder: i,
		}

		if Conn.Create(record).Error == nil {
			ResponsesMap[response.Name] = record.ID
		}
	}

	return ResponsesMap
}
