package project

import (
	"context"
	"errors"
	"time"

	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/model/team"
)

type ProjectGroup struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	MemberID     uint   `gorm:"type:bigint;uniqueIndex:ukey;not null;comment:团队成员id"`
	Name         string `gorm:"type:varchar(255);uniqueIndex:ukey;not null;comment:分组名称"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func init() {
	model.RegMigrate(&ProjectGroup{})
}

func (pg *ProjectGroup) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if pg.ID != 0 {
		tx = tx.Take(pg, "id = ?", pg.ID)
	} else if pg.MemberID != 0 && pg.Name != "" {
		tx = tx.Take(pg, "member_id = ? AND name = ?", pg.MemberID, pg.Name)
	} else {
		return false, errors.New("query condition error")
	}
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

func (pg *ProjectGroup) Create(ctx context.Context, name string, member *team.TeamMember) (*ProjectGroup, error) {
	// 获取最大的display_order
	var maxDisplayOrder ProjectGroup
	if err := model.DB(ctx).Model(pg).Where("member_id = ?", member.ID).Order("display_order desc").First(&maxDisplayOrder).Error; err != nil {
		maxDisplayOrder = ProjectGroup{DisplayOrder: 0}
	}

	pg.MemberID = member.ID
	pg.Name = name
	pg.DisplayOrder = maxDisplayOrder.DisplayOrder + 1
	return pg, model.DB(ctx).Create(pg).Error
}

func (pg *ProjectGroup) Rename(ctx context.Context) error {
	return model.DB(ctx).Model(pg).Update("name", pg.Name).Error
}

func (pg *ProjectGroup) Delete(ctx context.Context) error {
	if err := model.DB(ctx).Model(&ProjectMember{}).Where("group_id = ?", pg.ID).Update("group_id", 0).Error; err != nil {
		return err
	}
	return model.DB(ctx).Delete(pg).Error
}
