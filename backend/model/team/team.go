package team

import (
	"context"
	"errors"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/user"

	"gorm.io/gorm"
)

type Team struct {
	ID      string `gorm:"type:varchar(24);primarykey"`
	Name    string `gorm:"type:varchar(255);not null;comment:团队名"`
	Avatar  string `gorm:"type:varchar(255);comment:头像"`
	OwnerID uint   `gorm:"type:bigint;comment:团队所有者id"`
	model.TimeModel
}

func init() {
	model.RegMigrate(&Team{})
}

// Get 获取team
func (t *Team) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx).Take(t, "id = ?", t.ID)
	return tx.Error == nil, model.NotRecord(tx)
}

// Update 更新team
func (t *Team) Update(ctx context.Context) error {
	if t.ID == "" {
		return nil
	}
	// 只能更新name和avatar
	return model.DB(ctx).Model(t).Updates(map[string]interface{}{
		"name":   t.Name,
		"avatar": t.Avatar,
	}).Error
}

// Delete 删除team
func (t *Team) Delete(ctx context.Context) error {
	if t.ID == "" {
		return nil
	}
	return model.DB(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Where("team_id = ?", t.ID).Delete(&TeamMember{}).Error; err != nil {
				return err
			}

			return tx.Delete(t).Error
		},
	)
}

// Transfer 转让team
func (t *Team) Transfer(ctx context.Context, currentMember, targetMember *TeamMember) error {
	if t.ID == "" {
		return nil
	}
	return model.DB(ctx).Transaction(
		func(tx *gorm.DB) error {
			currentMember.Role = RoleMember
			if err := currentMember.Update(ctx); err != nil {
				return err
			}
			targetMember.Role = RoleOwner
			if err := targetMember.Update(ctx); err != nil {
				return err
			}

			if err := tx.Model(t).Update("owner_id", targetMember.UserID).Error; err != nil {
				return err
			}
			return nil
		},
	)
}

// HasMember 判断团队是否有此成员
func (t *Team) HasMember(ctx context.Context, tm *TeamMember) (bool, error) {
	tx := model.DB(ctx).Where("team_id = ? and user_id = ?", t.ID, tm.UserID).Take(&tm)
	return tx.Error == nil, model.NotRecord(tx)
}

// AddMember 团队添加成员
func (t *Team) AddMember(ctx context.Context, inviteBy uint, u *user.User) (*TeamMember, error) {
	var tm *TeamMember
	err := model.DB(ctx).Unscoped().Where("team_id = ? AND deleted_at is not null AND user_id = ?", t.ID, u.ID).Take(&tm).Error
	if err == nil {
		// 如果存在但已经被删除则恢复
		return tm, model.DB(ctx).Unscoped().Model(tm).Updates(map[string]interface{}{
			"status":         MemberStatusActive,
			"role":           RoleMember,
			"last_active_at": time.Now(),
			"created_at":     time.Now(),
			"updated_at":     time.Now(),
			"deleted_at":     nil,
		}).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果不存在，那么就创建
		return tm, model.DB(ctx).Transaction(
			func(tx *gorm.DB) error {
				tm.Role = RoleMember
				tm.TeamID = t.ID
				tm.UserID = u.ID
				tm.InvitedBy = inviteBy
				tm.LastActiveAt = time.Now()
				if err := tx.Create(tm).Error; err != nil {
					return err
				}
				return nil
			},
		)
	} else {
		return nil, err
	}
}
