package team

import (
	"context"
	"time"

	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/model/user"

	"github.com/pkg-id/objectid"
	"gorm.io/gorm"
)

// Create 创建team，创建团队所有者
func Create(ctx context.Context, owner *user.User, name string) (*Team, error) {
	t := &Team{
		ID:      objectid.New().String(),
		Name:    name,
		OwnerID: owner.ID,
	}
	if err := model.DB(ctx).Transaction(
		func(tx *gorm.DB) error {
			ret := tx.Create(t)
			if ret.Error != nil {
				return ret.Error
			}
			// 添加默认管理员
			own := &TeamMember{
				TeamID:       t.ID,
				UserID:       owner.ID,
				Role:         RoleOwner,
				LastActiveAt: time.Now(),
			}
			return tx.Create(own).Error
		},
	); err != nil {
		return nil, err
	}
	return t, nil
}

// GetUserTeams 获取用户 team 列表
func GetUserTeams(ctx context.Context, uid uint, roles ...Role) ([]*Team, error) {
	var (
		list []*Team
		tid  []string
	)

	if len(roles) > 0 {
		model.DB(ctx).Model(&TeamMember{}).Where("user_id = ? AND role in (?)", uid, roles).Pluck("team_id", &tid)
	} else {
		model.DB(ctx).Model(&TeamMember{}).Where("user_id = ?", uid).Pluck("team_id", &tid)
	}

	tx := model.DB(ctx)
	if len(tid) < 1 {
		return list, nil
	} else if len(tid) == 1 {
		tx = tx.Where("id = ?", tid[0])
	} else {
		tx = tx.Where("id in (?)", tid)
	}

	if err := tx.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetTeam 通过 team ID 获取 team
func GetTeam(ctx context.Context, id string) (*Team, error) {
	var t Team
	tx := model.DB(ctx).Take(&t, "id = ?", id)
	if tx.Error != nil {
		return nil, model.NotRecord(tx)
	}
	return &t, nil
}

// GetMember 通过 member ID 获取成员
func GetMember(ctx context.Context, id uint) (*TeamMember, error) {
	var t TeamMember
	err := model.DB(ctx).Take(&t, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetMembers 获取团队成员列表
func GetMembers(ctx context.Context, tID string, page, pageSize int, status string, roles ...Role) ([]*TeamMember, error) {
	userIDs := make([]uint, 0)
	if err := model.DB(ctx).Model(&user.User{}).Where("is_active = ?", true).Pluck("id", &userIDs).Error; err != nil {
		return nil, err
	}

	list := make([]*TeamMember, 0)
	tx := model.DB(ctx).Model(&TeamMember{}).Where("team_id = ? AND user_id in (?)", tID, userIDs)

	if status == MemberStatusActive || status == MemberStatusDeactive {
		tx.Where("status = ?", status)
	}
	if len(roles) > 0 {
		tx.Where("role in (?)", roles)
	}
	if page > 0 && pageSize > 0 {
		tx = tx.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	if err := tx.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetMembersCount 获取团队成员数量
func GetMembersCount(ctx context.Context, tID string, roles ...Role) (int64, error) {
	var count int64
	tx := model.DB(ctx).Model(&TeamMember{}).Where("team_id = ?", tID)
	if len(roles) > 0 {
		tx.Where("role in (?)", roles)
	}
	err := tx.Count(&count).Error
	return count, err
}

// GetLastActiveTeam 获取用户上一次活跃的团队
func GetLastActiveTeam(ctx context.Context, uid uint) (*TeamMember, error) {
	var tm *TeamMember
	tx := model.DB(ctx).Where("user_id = ?", uid).Order("last_active_at desc").First(&tm)
	if tx.Error != nil {
		return nil, model.NotRecord(tx)
	}
	return tm, nil
}

// GetMemberByToken 根据邀请码获取成员
func GetMemberByToken(ctx context.Context, token string) (*TeamMember, error) {
	var t *TeamMember
	return t, model.DB(ctx).Where("invitation_token = ?", token).First(&t).Error
}
