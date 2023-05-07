package models

import (
	"encoding/json"
	"time"

	"github.com/apicat/apicat/common/spec"
)

type Commons struct {
	ID           uint   `gorm:"type:integer primary key autoincrement"`
	ProjectId    uint   `gorm:"index;not null;comment:项目id"`
	Type         string `gorm:"type:varchar(255);not null;comment:类型:parameter,response"`
	Content      string `gorm:"type:mediumtext;comment:内容"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewCommons(ids ...uint) (*Commons, error) {
	commons := &Commons{}
	if len(ids) > 0 {
		if err := Conn.Take(commons, ids[0]).Error; err != nil {
			return commons, err
		}
		return commons, nil
	}
	return commons, nil
}

func (c *Commons) List() ([]Commons, error) {
	tx := Conn.Where("project_id = ?", c.ProjectId)
	if c.ID > 0 {
		tx = tx.Where("id = ?", c.ID)
	}
	if c.Type != "" {
		tx = tx.Where("type = ?", c.Type)
	}
	commons := []Commons{}
	return commons, tx.Find(&commons).Error
}

func (c *Commons) Get() error {
	tx := Conn.Where("project_id = ?", c.ProjectId)
	if c.ID > 0 {
		tx = tx.Where("id = ?", c.ID)
	}
	if c.Type != "" {
		tx = tx.Where("type = ?", c.Type)
	}
	return tx.First(c).Error
}

func (c *Commons) Create() error {
	return Conn.Create(c).Error
}

func (c *Commons) Update() error {
	return Conn.Save(c).Error
}

func (c *Commons) Delete() error {
	return Conn.Delete(c).Error
}

// 判断记录是否存在
func (c *Commons) IsExist() bool {
	db := Conn.Model(&Commons{}).Where("project_id = ?", c.ProjectId)
	if c.ID > 0 {
		db = db.Where("id = ?", c.ID)
	}
	if c.Type != "" {
		db = db.Where("type = ?", c.Type)
	}

	var count int64
	db.Count(&count)
	return count > 0
}

func CommonsImport(projectID uint, commons *spec.Common) map[string]nameToIdMap {
	var commonsMap = map[string]nameToIdMap{
		"parameter": make(nameToIdMap),
		"response":  make(nameToIdMap),
	}

	if commons.Parameters == nil && commons.Responses == nil {
		return commonsMap
	}

	for i, parameter := range commons.Parameters {
		if parameterStr, err := json.Marshal(parameter); err == nil {
			record := &Commons{
				ProjectId:    projectID,
				Type:         "parameter",
				Content:      string(parameterStr),
				DisplayOrder: i,
			}

			if Conn.Create(record).Error == nil {
				commonsMap["parameter"][parameter.Name] = record.ID
			}
		}
	}

	for i, response := range commons.Responses {
		if responseStr, err := json.Marshal(response); err == nil {
			record := &Commons{
				ProjectId:    projectID,
				Type:         "response",
				Content:      string(responseStr),
				DisplayOrder: i,
			}

			if Conn.Create(record).Error == nil {
				commonsMap["response"][response.Name] = record.ID
			}
		}
	}

	return commonsMap
}

// func CommonsExport(projectID uint) *spec.Common {
// 	var commons []*Commons

// 	specCommon := &spec.Common{
// 		Parameters: &spec.HTTPParameters{},
// 		Responses:  []spec.HTTPResponse{},
// 	}

// 	if err := Conn.Where("project_id = ?", projectID).Find(&commons).Error; err == nil {
// 		for _, c := range commons {
// 			if c.Type == "parameter" {
// 				_ = json.Unmarshal([]byte(c.Content), &specCommon.Parameters)
// 			} else {
// 				var response spec.HTTPResponse
// 				if json.Unmarshal([]byte(c.Content), &response) == nil {
// 					specCommon.Responses = append(specCommon.Responses, response)
// 				}
// 			}
// 		}
// 	}
// 	return specCommon
// }
