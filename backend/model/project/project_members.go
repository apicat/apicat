package project

import (
	"github.com/apicat/apicat/backend/model"
	"time"

	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type ProjectMembers struct {
	ID         uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID  uint   `gorm:"type:bigint;index;not null;comment:项目id"`
	UserID     uint   `gorm:"type:bigint;index;not null;comment:用户id"`
	GroupID    uint   `gorm:"type:bigint;not null;default:0;comment:分组id"`
	Authority  string `gorm:"type:varchar(255);not null;comment:项目权限:manage,write,read"`
	FollowedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

func init() {
	model.RegMigrate(&ProjectMembers{})
}

var (
	ProjectMembersManage = "manage"
	ProjectMembersWrite  = "write"
	ProjectMembersRead   = "read"
)

func NewProjectMembers(ids ...uint) (*ProjectMembers, error) {
	members := &ProjectMembers{}
	if len(ids) > 0 {
		if err := model.Conn.Take(members, ids[0]).Error; err != nil {
			return members, err
		}
		return members, nil
	}
	return members, nil
}

func (pm *ProjectMembers) List(page, pageSize int) ([]ProjectMembers, error) {
	var projectMembers []ProjectMembers

	if page == 0 && pageSize == 0 {
		return projectMembers, model.Conn.Where("project_id = ?", pm.ProjectID).Order("created_at desc").Find(&projectMembers).Error
	}

	return projectMembers, model.Conn.Where("project_id = ?", pm.ProjectID).Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&projectMembers).Error
}

func (pm *ProjectMembers) GetByUserIDAndProjectID() error {
	return model.Conn.Where("user_id = ? and project_id = ?", pm.UserID, pm.ProjectID).Take(pm).Error
}

func (pm *ProjectMembers) Create() error {
	return model.Conn.Create(pm).Error
}

func (pm *ProjectMembers) Delete() error {
	return model.Conn.Delete(pm).Error
}

func (pm *ProjectMembers) Update() error {
	return model.Conn.Save(pm).Error
}

func (pm *ProjectMembers) Count() (int64, error) {
	var count int64
	return count, model.Conn.Model(&ProjectMembers{}).Where("project_id = ?", pm.ProjectID).Count(&count).Error
}

func DeleteAllMembersByProjectID(projectID uint) error {
	return model.Conn.Where("project_id = ?", projectID).Delete(&ProjectMembers{}).Error
}

func (pm *ProjectMembers) MemberIsManage() bool {
	return pm.Authority == ProjectMembersManage
}

func (pm *ProjectMembers) MemberHasWritePermission() bool {
	return slices.Contains([]string{ProjectMembersManage, ProjectMembersWrite}, pm.Authority)
}

func GetUserInvolvedProject(UserID uint, PMAuthorities ...string) ([]ProjectMembers, error) {
	var projectMembers []ProjectMembers
	query := model.Conn.Where("user_id = ?", UserID)
	if len(PMAuthorities) > 0 {
		query = query.Where("authority IN ?", PMAuthorities)
	}
	return projectMembers, query.Order("created_at desc").Find(&projectMembers).Error
}

func GetProjectGroupedByUser(UserID, GroupID uint) ([]ProjectMembers, error) {
	var projectMembers []ProjectMembers
	return projectMembers, model.Conn.Where("user_id = ? AND group_id = ?", UserID, GroupID).Find(&projectMembers).Error
}

func GetProjectFollowedByUser(UserID uint) ([]ProjectMembers, error) {
	var projectMembers []ProjectMembers
	return projectMembers, model.Conn.Where("user_id = ? AND followed_at is not null", UserID).Find(&projectMembers).Error
}
