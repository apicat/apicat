package iteration

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/model/team"

	"github.com/pkg-id/objectid"

	"gorm.io/gorm"
)

type Iteration struct {
	ID          string `gorm:"type:varchar(24);primarykey"`
	TeamID      string `gorm:"type:varchar(24);not null;comment:团队id"`
	ProjectID   string `gorm:"type:varchar(24);index;not null;comment:项目id"`
	Title       string `gorm:"type:varchar(255);not null;comment:迭代标题"`
	Description string `gorm:"type:varchar(255);comment:迭代描述"`
	CreatedBy   uint   `gorm:"type:bigint;not null;default:0;comment:创建人id"`
	UpdatedBy   uint   `gorm:"type:bigint;not null;default:0;comment:最后更新人id"`
	DeletedBy   uint   `gorm:"type:bigint;default:null;comment:删除人id"`
	model.TimeModel
}

func init() {
	model.RegMigrate(&Iteration{})
}

func (i *Iteration) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx).Take(i, "id = ?", i.ID)
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

func (i *Iteration) Create(ctx context.Context, member *team.TeamMember) error {
	i.ID = objectid.New().String()
	i.TeamID = member.TeamID
	i.CreatedBy = member.ID
	i.UpdatedBy = member.ID
	return model.DB(ctx).Create(i).Error
}

func (i *Iteration) Update(ctx context.Context, member *team.TeamMember) error {
	// 只能更新Title和Description
	return model.DB(ctx).Model(i).Updates(map[string]interface{}{
		"title":       i.Title,
		"description": i.Description,
		"updated_by":  member.ID,
	}).Error
}

func (i *Iteration) Delete(ctx context.Context, member *team.TeamMember) error {
	return model.DB(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Model(&i).Updates(map[string]interface{}{
				"deleted_by": &member.ID,
			}).Error; err != nil {
				return err
			}

			if err := tx.Delete(&IterationApi{}, "iteration_id = ?", i.ID).Error; err != nil {
				return err
			}

			return tx.Delete(i).Error
		},
	)
}

func (i *Iteration) PlanningIterationApi(ctx context.Context, collections []*collection.Collection) error {
	iterationApis, err := GetIterationApi(ctx, i)
	if err != nil {
		return err
	}

	iterationApiDict := map[uint]*IterationApi{}
	for _, v := range iterationApis {
		iterationApiDict[v.CollectionID] = v
	}

	collectionDict := map[uint]*collection.Collection{}
	for _, v := range collections {
		collectionDict[v.ID] = v
	}

	wantPop := make([]uint, 0)
	wantPush := make([]*collection.Collection, 0)

	// 找出 iterationApis 中存在但 collections 中不存在的元素
	for _, v := range iterationApis {
		if _, ok := collectionDict[v.CollectionID]; !ok {
			wantPop = append(wantPop, v.CollectionID)
		}
	}

	// 找出 collections 中存在但 iterationApis 中不存在的元素
	for _, c := range collections {
		if _, ok := iterationApiDict[c.ID]; !ok {
			wantPush = append(wantPush, c)
		}
	}

	if err := i.BatchCreateCollection(ctx, wantPush); err != nil {
		return err
	}
	return i.BatchDeleteCollection(ctx, wantPop)
}

// GetIterationApiCount 获取迭代涉及的api数量
func (i *Iteration) GetIterationApiCount(ctx context.Context) (int64, error) {
	var count int64
	err := model.DB(ctx).Model(&IterationApi{}).Where("iteration_id = ? AND collection_type != ?", i.ID, collection.CategoryType).Count(&count).Error
	return count, err
}

func (i *Iteration) BatchCreateCollection(ctx context.Context, collections []*collection.Collection) error {
	if len(collections) == 0 {
		return nil
	}

	ias := make([]*IterationApi, 0)
	for _, c := range collections {
		ias = append(ias, &IterationApi{
			IterationID:    i.ID,
			CollectionID:   c.ID,
			CollectionType: c.Type,
		})
	}
	return model.DB(ctx).Create(ias).Error
}

func (i *Iteration) BatchDeleteCollection(ctx context.Context, collectionIDs []uint) error {
	if len(collectionIDs) == 0 {
		return nil
	}

	return model.DB(ctx).Where("iteration_id = ? AND collection_id in (?)", i.ID, collectionIDs).Delete(&IterationApi{}).Error
}

func (i *Iteration) GetCollectionIDs(ctx context.Context) ([]uint, error) {
	var cIDs []uint
	return cIDs, model.DB(ctx).Model(&IterationApi{}).Where("iteration_id = ?", i.ID).Pluck("collection_id", &cIDs).Error
}

func (i *Iteration) ProjectInfo(ctx context.Context) (*project.Project, error) {
	p := &project.Project{ID: i.ProjectID}
	if _, err := p.Get(ctx); err != nil {
		return nil, err
	}
	return p, nil
}
