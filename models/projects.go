package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Projects struct {
	ID          uint   `gorm:"type:integer primary key autoincrement"`
	PublicId    string `gorm:"type:varchar(255);uniqueIndex;not null;comment:项目公开id"`
	Title       string `gorm:"type:varchar(255);not null;comment:项目名称"`
	Visibility  int    `gorm:"type:tinyint(1);not null;comment:项目可见性:0私有,1公开"`
	Description string `gorm:"type:varchar(255);comment:项目描述"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func NewProjects(ids ...interface{}) (*Projects, error) {
	project := &Projects{}

	if len(ids) > 0 {
		var err error

		switch ids[0].(type) {
		case string:
			err = Conn.Where("public_id = ?", ids[0]).Take(project).Error
		case uint:
			err = Conn.Take(project, ids[0]).Error
		default:
			err = errors.New("invalid id type")
		}

		if err != nil {
			return project, err
		}
		return project, nil
	}
	return project, nil
}

func (p *Projects) Create() error {
	return Conn.Create(p).Error
}

func (p *Projects) Get(id string) error {
	return Conn.Where("public_id = ?", id).Take(p).Error
}

func (p *Projects) List() ([]Projects, error) {
	var projects []Projects
	return projects, Conn.Order("created_at desc").Find(&projects).Error
}

func (p *Projects) Delete() error {
	return Conn.Delete(p).Error
}

func (p *Projects) Save() error {
	return Conn.Save(p).Error
}
