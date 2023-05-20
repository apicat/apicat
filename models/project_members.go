package models

import (
	"time"

	"gorm.io/gorm"
)

type ProjectMembers struct {
	ID        uint   `gorm:"type:integer primary key autoincrement"`
	ProjectID uint   `gorm:"index;not null;comment:项目id"`
	UserID    uint   `gorm:"index;not null;comment:用户id"`
	Authority string `gorm:"type:varchar(255);not null;comment:项目权限:manage,write,read"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

var (
	ProjectMembersManage = "manage"
	ProjectMembersWrite  = "write"
	ProjectMembersRead   = "read"
)

func NewProjectMembers(ids ...uint) (*ProjectMembers, error) {
	members := &ProjectMembers{}
	if len(ids) > 0 {
		if err := Conn.Take(members, ids[0]).Error; err != nil {
			return members, err
		}
		return members, nil
	}
	return members, nil
}

func (pm *ProjectMembers) List() ([]ProjectMembers, error) {
	var projectMembers []ProjectMembers
	return projectMembers, Conn.Order("created_at desc").Find(&projectMembers).Error
}

func (pm *ProjectMembers) Get() error {
	return Conn.Take(pm).Error
}

func (pm *ProjectMembers) Create() error {
	return Conn.Create(pm).Error
}

func (pm *ProjectMembers) Delete() error {
	return Conn.Delete(pm).Error
}

func (pm *ProjectMembers) Update() error {
	return Conn.Save(pm).Error
}

func DeleteAllMembersByProjectID(projectID uint) error {
	return Conn.Where("project_id = ?", projectID).Delete(&ProjectMembers{}).Error
}
