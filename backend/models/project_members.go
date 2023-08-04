package models

import (
	"time"

	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type ProjectMembers struct {
	ID        uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID uint   `gorm:"type:bigint;index;not null;comment:项目id"`
	UserID    uint   `gorm:"type:bigint;index;not null;comment:用户id"`
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

func (pm *ProjectMembers) List(page, pageSize int) ([]ProjectMembers, error) {
	var projectMembers []ProjectMembers

	if page == 0 && pageSize == 0 {
		return projectMembers, Conn.Where("project_id = ?", pm.ProjectID).Order("created_at desc").Find(&projectMembers).Error
	}

	return projectMembers, Conn.Where("project_id = ?", pm.ProjectID).Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&projectMembers).Error
}

func (pm *ProjectMembers) GetByUserIDAndProjectID() error {
	return Conn.Where("user_id = ? and project_id = ?", pm.UserID, pm.ProjectID).Take(pm).Error
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

func (pm *ProjectMembers) Count() (int64, error) {
	var count int64
	return count, Conn.Model(&ProjectMembers{}).Where("project_id = ?", pm.ProjectID).Count(&count).Error
}

func DeleteAllMembersByProjectID(projectID uint) error {
	return Conn.Where("project_id = ?", projectID).Delete(&ProjectMembers{}).Error
}

func (pm *ProjectMembers) MemberIsManage() bool {
	return pm.Authority == ProjectMembersManage
}

func (pm *ProjectMembers) MemberHasWritePermission() bool {
	return slices.Contains([]string{ProjectMembersManage, ProjectMembersWrite}, pm.Authority)
}

func GetUserInvolvedProject(UserID uint) ([]ProjectMembers, error) {
	var projectMembers []ProjectMembers
	return projectMembers, Conn.Where("user_id = ?", UserID).Order("created_at desc").Find(&projectMembers).Error
}
