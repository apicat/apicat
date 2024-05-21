package team

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/user"

	"gorm.io/gorm"
)

const (
	RoleOwner            = "owner"
	RoleAdmin            = "admin"
	RoleMember           = "member"
	MemberStatusActive   = "active"
	MemberStatusDeactive = "deactive"
)

type Role string

type TeamMember struct {
	ID              uint      `gorm:"primarykey"`
	TeamID          string    `gorm:"type:varchar(24);uniqueIndex:ukey;not null;comment:团队id"`
	UserID          uint      `gorm:"type:bigint;uniqueIndex:ukey;not null;comment:用户id"`
	Role            Role      `gorm:"type:varchar(32);comment:角色"`
	Status          string    `gorm:"type:varchar(32);default:active;comment:状态"`
	InvitationToken string    `gorm:"type:varchar(32);index;comment:邀请码"`
	InvitedBy       uint      `gorm:"type:bigint;default:0;comment:邀请人的TeamMemberID"`
	LastActiveAt    time.Time `gorm:"type:datetime;not null;comment:最后活跃时间"`
	model.TimeModel
}

var roleRanking = map[Role]int{
	RoleOwner:  3,
	RoleAdmin:  2,
	RoleMember: 1,
}

func (r Role) Greater(other Role) bool {
	return roleRanking[r] > roleRanking[other]
}

func (r Role) GreaterOrEqual(other Role) bool {
	return roleRanking[r] >= roleRanking[other]
}

func (r Role) Equal(other Role) bool {
	return roleRanking[r] == roleRanking[other]
}

func (r Role) Lower(other Role) bool {
	return roleRanking[r] < roleRanking[other]
}

func (r Role) LowerOrEqual(other Role) bool {
	return roleRanking[r] <= roleRanking[other]
}

func (t *TeamMember) UserInfo(ctx context.Context, unscoped bool) (*user.User, error) {
	var (
		usr user.User
		tx  *gorm.DB
	)

	if unscoped {
		tx = model.DB(ctx).Unscoped()
	} else {
		tx = model.DB(ctx)
	}

	if err := tx.First(&usr, t.UserID).Error; err != nil {
		return nil, err
	} else {
		return &usr, nil
	}
}

func (t *TeamMember) AfterCreate(tx *gorm.DB) error {
	t.InvitationToken = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v|%d|%d", t.TeamID, t.ID, time.Now().UnixNano()))))
	return tx.Model(t).Update("invitation_token", t.InvitationToken).Error
}

// UpdateRole 更新成员角色
func (t *TeamMember) Update(ctx context.Context) error {
	if t.ID == 0 {
		return nil
	}

	// 目前只能更新角色和状态
	return model.DB(ctx).Model(t).Updates(map[string]interface{}{
		"role":   t.Role,
		"status": t.Status,
	}).Error
}

// 修改成员最后活跃时间
func (t *TeamMember) UpdateActiveAt(ctx context.Context) error {
	return model.DB(ctx).Model(t).UpdateColumn("last_active_at", time.Now()).Error
}

// ResetInvitationToken 重置邀请码
func (t *TeamMember) ResetInvitationToken(ctx context.Context) error {
	t.InvitationToken = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v|%d|%d", t.TeamID, t.ID, time.Now().UnixNano()))))
	return model.DB(ctx).Model(t).Update("invitation_token", t.InvitationToken).Error
}

// Quit 退出团队
func (t *TeamMember) Quit(ctx context.Context) error {
	return model.DB(ctx).Delete(t).Error
}
