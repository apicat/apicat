package project

import (
	"errors"
	"github.com/apicat/apicat/backend/model"
	"time"

	"gorm.io/gorm"
)

type Projects struct {
	ID            uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	PublicId      string `gorm:"type:varchar(255);uniqueIndex;not null;comment:项目公开id"`
	Title         string `gorm:"type:varchar(255);not null;comment:项目名称"`
	Visibility    int    `gorm:"type:tinyint(1);not null;comment:项目可见性:0私有,1公开"`
	SharePassword string `gorm:"type:varchar(255);comment:项目分享密码"`
	Description   string `gorm:"type:varchar(255);comment:项目描述"`
	Cover         string `gorm:"type:varchar(255);comment:项目封面"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func init() {
	model.RegMigrate(&Projects{})
}

func NewProjects(ids ...interface{}) (*Projects, error) {
	project := &Projects{}

	if len(ids) > 0 {
		var err error

		switch ids[0].(type) {
		case string:
			err = model.Conn.Where("public_id = ?", ids[0]).Take(project).Error
		case uint:
			err = model.Conn.Take(project, ids[0]).Error
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
	return model.Conn.Create(p).Error
}

func (p *Projects) Get(id string) error {
	return model.Conn.Where("public_id = ?", id).Take(p).Error
}

func (p *Projects) List(ids ...uint) ([]Projects, error) {
	var projects []Projects
	if len(ids) > 0 {
		return projects, model.Conn.Where("id IN ?", ids).Order("created_at desc").Find(&projects).Error
	}
	return projects, model.Conn.Order("created_at desc").Find(&projects).Error
}

func (p *Projects) Delete() error {
	if err := DeleteAllMembersByProjectID(p.ID); err != nil {
		return err
	}

	return model.Conn.Delete(p).Error
}

func (p *Projects) Save() error {
	return model.Conn.Save(p).Error
}
