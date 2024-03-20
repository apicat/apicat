package project

import (
	"context"
	"errors"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/model/user"

	"gorm.io/gorm"
)

type Permission string

const (
	ProjectMemberManage Permission = "manage"
	ProjectMemberWrite  Permission = "write"
	ProjectMemberRead   Permission = "read"
	ProjectMemberNone   Permission = "none"
)

type ProjectMember struct {
	ID         uint       `gorm:"type:bigint;primaryKey"`
	ProjectID  string     `gorm:"type:varchar(24);uniqueIndex:ukey;not null;comment:项目id"`
	MemberID   uint       `gorm:"type:bigint;uniqueIndex:ukey;not null;comment:团队成员id"`
	GroupID    uint       `gorm:"type:bigint;not null;default:0;comment:分组id"`
	Permission Permission `gorm:"type:varchar(255);not null;comment:项目权限:manage,write,read"`
	FollowedAt *time.Time `gorm:"type:datetime;comment:关注项目时间"` // 不为空表示关注，字段类型为指针是为了在取消关注时，可以设置为null
	model.TimeModel
}

func init() {
	model.RegMigrate(&ProjectMember{})
}

var permissionRanking = map[Permission]int{
	ProjectMemberManage: 3,
	ProjectMemberWrite:  2,
	ProjectMemberRead:   1,
	ProjectMemberNone:   0,
}

func (p Permission) Greater(other Permission) bool {
	return permissionRanking[p] > permissionRanking[other]
}

func (p Permission) GreaterOrEqual(other Permission) bool {
	return permissionRanking[p] >= permissionRanking[other]
}

func (p Permission) Equal(other Permission) bool {
	return permissionRanking[p] == permissionRanking[other]
}

func (p Permission) Lower(other Permission) bool {
	return permissionRanking[p] < permissionRanking[other]
}

func (p Permission) LowerOrEqual(other Permission) bool {
	return permissionRanking[p] <= permissionRanking[other]
}

// Get 获取项目成员
func (pm *ProjectMember) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if pm.ID != 0 {
		tx = tx.Take(pm, "id = ?", pm.ID)
	} else if pm.ProjectID != "" && pm.MemberID != 0 {
		tx = tx.Take(pm, "project_id = ? AND member_id = ?", pm.ProjectID, pm.MemberID)
	} else {
		return false, errors.New("query condition error")
	}
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

func (pm *ProjectMember) MemberInfo(ctx context.Context, unscoped bool) (*team.TeamMember, error) {
	var (
		tm *team.TeamMember
		tx *gorm.DB
	)

	if unscoped {
		tx = model.DB(ctx).Unscoped()
	} else {
		tx = model.DB(ctx)
	}

	if err := tx.First(&tm, pm.MemberID).Error; err != nil {
		return nil, err
	} else {
		return tm, nil
	}
}

func (pm *ProjectMember) UserInfo(ctx context.Context, unscoped bool) (*user.User, error) {
	if tm, err := pm.MemberInfo(ctx, unscoped); err != nil {
		return nil, err
	} else {
		return tm.UserInfo(ctx, unscoped)
	}
}

// Create 创建项目成员
func (pm *ProjectMember) Create(ctx context.Context, tx *gorm.DB) (*ProjectMember, error) {
	if tx == nil {
		tx = model.DB(ctx)
	}
	// 先判断这个成员是否是被删除的成员
	var projectMember *ProjectMember
	err := tx.Unscoped().Where("project_id = ? AND  member_id = ?", pm.ProjectID, pm.MemberID).Take(&projectMember).Error
	if err == nil {
		// 如果存在但已经被删除则恢复
		return projectMember, tx.Unscoped().Model(&projectMember).Updates(map[string]interface{}{
			"permission": pm.Permission,
			"created_at": time.Now(),
			"updated_at": time.Now(),
			"deleted_at": nil,
		}).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果不存在，那么就创建
		return pm, tx.Create(pm).Error
	} else {
		return nil, err
	}
}

// Update 更新项目成员，目前只能更新permission和groupID
func (pm *ProjectMember) Update(ctx context.Context) error {
	if pm.ID == 0 {
		return nil
	}
	// 只能更新permission和groupID
	return model.DB(ctx).Model(&pm).Updates(map[string]interface{}{
		"permission": pm.Permission,
		"group_id":   pm.GroupID,
	}).Error
}

// Delete 删除项目成员
func (pm *ProjectMember) Delete(ctx context.Context) error {
	return model.DB(ctx).Delete(pm).Error
}

// UpdateFollow 更新关注项目时间
func (pm *ProjectMember) UpdateFollow(ctx context.Context, t *time.Time) error {
	return model.DB(ctx).Model(pm).Select("followed_at").UpdateColumn("followed_at", t).Error
}
