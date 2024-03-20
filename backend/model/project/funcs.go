package project

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/module/spec"

	"gorm.io/gorm"
)

// GetProjects 获取成员项目
// permission: manage,write,read 为空时获取所有项目
func GetProjects(ctx context.Context, member *team.TeamMember, permission ...Permission) ([]*Project, error) {
	var list []*Project
	tx := model.DB(ctx)
	subQuery := tx.Model(&ProjectMember{}).Select("project_id").Where("member_id = ?", member.ID)
	if len(permission) > 0 {
		subQuery.Where("permission in (?)", permission)
	}

	return list, tx.Where("id in (?)", subQuery).Find(&list).Error
}

// 通过项目 ID 获取项目
func GetProjectsByIds(ctx context.Context, ids []string) ([]*Project, error) {
	var list []*Project
	tx := model.DB(ctx)
	return list, tx.Where("id in (?)", ids).Find(&list).Error
}

// GetFollowedProjectIDs 获取关注的项目id
func GetFollowedProjectIDs(ctx context.Context, member *team.TeamMember) ([]uint, error) {
	var list []uint
	tx := model.DB(ctx)
	tx = tx.Model(&ProjectMember{}).Select("project_id").Where("member_id = ? AND followed_at is not null", member.ID)
	if err := tx.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// TransferProject 传送项目
// p 需要移交的项目
// originMember 原项目管理员
// targetMember 目标项目管理员
func TransferProject(ctx context.Context, p *Project, originMember, targetMember *ProjectMember) error {
	return model.DB(ctx).Transaction(
		func(tx *gorm.DB) error {
			// 将项目拥有者修改为新的管理者
			if err := tx.Model(p).Update("member_id", targetMember.MemberID).Error; err != nil {
				return err
			}

			// 修改项目新拥有者为管理者
			targetMember.Permission = ProjectMemberManage
			if err := tx.Model(targetMember).Update("permission", ProjectMemberManage).Error; err != nil {
				return err
			}

			// 修改原拥有者角色为阅读者
			return tx.Model(originMember).Update("permission", ProjectMemberRead).Error
		},
	)
}

// BatchCreateMember 批量创建项目成员
func BatchCreateMember(ctx context.Context, projectID string, members []*team.TeamMember, permission Permission) {
	for _, member := range members {
		pm := &ProjectMember{
			ProjectID:  projectID,
			MemberID:   member.ID,
			Permission: permission,
		}
		pm.Create(ctx, nil)
	}
}

// GetProjectMembers 获取项目成员
func GetProjectMembers(ctx context.Context, pID string, page, pageSize int, permission ...Permission) ([]*ProjectMember, error) {
	var list = make([]*ProjectMember, 0)

	tx := model.DB(ctx).Where("project_id = ?", pID)

	if len(permission) > 0 {
		tx.Where("permission in (?)", permission)
	}
	if page > 0 && pageSize > 0 {
		tx = tx.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	err := tx.Find(&list).Error
	return list, err
}

// 通过项目成员 ID 获取项目成员记录
func GetProjectMemberRecordsByMemberID(ctx context.Context, mID uint) ([]*ProjectMember, error) {
	var list = make([]*ProjectMember, 0)
	err := model.DB(ctx).Model(&ProjectMember{}).Where("member_id = ?", mID).Find(&list).Error
	return list, err
}

// 通过项目分组 ID 获取项目成员记录
func GetProjectMemberRecordsByGroupID(ctx context.Context, mID uint, groupID uint) ([]*ProjectMember, error) {
	var list = make([]*ProjectMember, 0)
	err := model.DB(ctx).Model(&ProjectMember{}).Where("member_id = ? AND group_id = ?", mID, groupID).Find(&list).Error
	return list, err
}

// 通过项目权限获取项目成员记录
func GetProjectMemberRecordsByPermissions(ctx context.Context, mID uint, permission ...Permission) ([]*ProjectMember, error) {
	var list = make([]*ProjectMember, 0)
	err := model.DB(ctx).Model(&ProjectMember{}).Where("member_id = ? AND permission in (?)", mID, permission).Find(&list).Error
	return list, err
}

// 通过项目成员 ID 获取关注的项目成员记录
func GetProjectMemberRecordsWithFollowed(ctx context.Context, mID uint) ([]*ProjectMember, error) {
	var list = make([]*ProjectMember, 0)
	err := model.DB(ctx).Model(&ProjectMember{}).Where("member_id = ? AND followed_at is not null", mID).Find(&list).Error
	return list, err
}

// 获取项目成员数量
func GetProjectMembersCount(ctx context.Context, pID string, permission ...Permission) (int64, error) {
	var count int64
	tx := model.DB(ctx).Model(&ProjectMember{}).Where("project_id = ?", pID)
	if len(permission) > 0 {
		tx.Where("permission in (?)", permission)
	}
	err := tx.Count(&count).Error
	return count, err
}

// GetInvolvedProjectMembers 获取成员参与的项目成员信息
func GetInvolvedProjectMembers(ctx context.Context, member *team.TeamMember, permission ...Permission) ([]*ProjectMember, error) {
	var list = make([]*ProjectMember, 0)
	tx := model.DB(ctx)
	if len(permission) > 0 {
		tx = tx.Where("permission in (?)", permission)
	}
	err := tx.Find(&list, "member_id = ?", member.ID).Error
	return list, err
}

func GetProjectGroups(ctx context.Context, mID uint) ([]*ProjectGroup, error) {
	var list []*ProjectGroup
	err := model.DB(ctx).Where("member_id = ?", mID).Order("display_order asc").Find(&list).Error
	return list, err
}

func GroupSort(ctx context.Context, mID uint, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Transaction(
		func(tx *gorm.DB) error {
			for i, id := range ids {
				if err := tx.Model(&ProjectGroup{}).Where("id = ? AND member_id = ?", id, mID).Update("display_order", i+1).Error; err != nil {
					return err
				}
			}
			return nil
		},
	)
}

// GetServers 获取项目URL列表
func GetServers(ctx context.Context, projectID string) ([]*Server, error) {
	var list []*Server
	err := model.DB(ctx).Where("project_id = ?", projectID).Order("display_order asc").Find(&list).Error
	return list, err
}

// ServerSort 排序
func ServerSort(ctx context.Context, pID string, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Transaction(
		func(tx *gorm.DB) error {
			for i, id := range ids {
				if err := tx.Model(&Server{}).Where("id = ? AND project_id = ?", id, pID).Update("display_order", i+1).Error; err != nil {
					return err
				}
			}
			return nil
		},
	)
}

func ServersImport(ctx context.Context, projectID string, servers []*spec.Server) {
	if len(servers) > 0 {
		for i, server := range servers {
			model.DB(ctx).Create(&Server{
				ProjectID:    projectID,
				Description:  server.Description,
				URL:          server.URL,
				DisplayOrder: i + 1,
			})
		}
	}
}

func ExportServers(ctx context.Context, projectID string) []*spec.Server {
	specServers := make([]*spec.Server, 0)

	servers, err := GetServers(ctx, projectID)
	if err != nil {
		return specServers
	}

	for _, server := range servers {
		specServers = append(specServers, &spec.Server{
			Description: server.Description,
			URL:         server.URL,
		})
	}

	return specServers
}
