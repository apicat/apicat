package models

import (
	"time"

	"gorm.io/gorm"
)

type Iterations struct {
	ID          uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	PublicID    string `gorm:"type:varchar(255);index;not null;comment:迭代公开id"`
	ProjectID   uint   `gorm:"type:bigint;not null;comment:项目id"`
	Title       string `gorm:"type:varchar(255);not null;comment:迭代标题"`
	Description string `gorm:"type:varchar(255);comment:迭代描述"`
	CreatedAt   time.Time
	CreatedBy   uint `gorm:"type:bigint;not null;default:0;comment:创建人id"`
	UpdatedAt   time.Time
	UpdatedBy   uint `gorm:"type:bigint;not null;default:0;comment:最后更新人id"`
	DeletedAt   gorm.DeletedAt
	DeletedBy   uint `gorm:"type:bigint;not null;default:0;comment:删除人id"`
}

func NewIterations(ids ...uint) (*Iterations, error) {
	if len(ids) > 0 {
		iteration := &Iterations{ID: ids[0]}
		if err := Conn.Take(iteration).Error; err != nil {
			return iteration, err
		}
		return iteration, nil
	}
	return &Iterations{}, nil
}

func (i *Iterations) List(page, pageSize int, pIDs ...uint) ([]*Iterations, error) {
	var (
		iterations []*Iterations
		query      *gorm.DB
	)

	if len(pIDs) > 0 {
		query = Conn.Where("project_id IN ?", pIDs).Order("created_at desc")
	} else {
		query = Conn.Order("created_at desc")
	}

	if page != 0 && pageSize != 0 {
		query = query.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	return iterations, query.Find(&iterations).Error
}

func (i *Iterations) Create() error {
	return Conn.Create(i).Error
}

func (i *Iterations) Update() error {
	return Conn.Save(i).Error
}

func (i *Iterations) Delete() error {
	return Conn.Delete(i).Error
}

func (i *Iterations) PlanningIterationApi(cIDs []uint) error {
	iterationApi, _ := NewIterationApis()
	iterationApis, err := iterationApi.List(i.ID)
	if err != nil {
		return err
	}

	if len(cIDs) == 0 && len(iterationApis) != 0 {
		if err := BatchDeleteIterationApi(iterationApis); err != nil {
			return err
		}
	}

	collection, _ := NewCollections()
	collection.ProjectId = i.ProjectID
	collections, err := collection.List()
	if err != nil {
		return err
	}
	collectionDict := map[uint]*Collections{}
	for _, v := range collections {
		collectionDict[v.ID] = v
	}

	iterationApiDict := map[uint]*IterationApis{}
	for _, v := range iterationApis {
		iterationApiDict[v.CollectionID] = v
	}

	wantPop := []*IterationApis{}
	wantPush := []*IterationApis{}

	// 找出iterationApis中存在但cIDs中不存在的元素
	for _, iterationApi := range iterationApis {
		found := false
		for _, cid := range cIDs {
			if iterationApi.CollectionID == cid {
				found = true
				break
			}
		}
		if !found {
			wantPop = append(wantPop, iterationApi)
		}
	}

	// 找出cIDs中存在但iterationApis中不存在的元素
	for _, cid := range cIDs {
		if _, ok := collectionDict[cid]; ok {
			found := false
			for _, iterationApi := range iterationApis {
				if cid == iterationApi.CollectionID {
					found = true
					break
				}
			}
			if !found {
				wantPush = append(wantPush, &IterationApis{
					IterationID:    i.ID,
					CollectionID:   collectionDict[cid].ID,
					CollectionType: collectionDict[cid].Type,
				})
			}
		}
	}

	if err := BatchDeleteIterationApi(wantPop); err != nil {
		return err
	}
	if err := BatchInsertIterationApi(wantPush); err != nil {
		return err
	}
	return nil
}

func (i *Iterations) Creator() string {
	user, err := NewUsers(i.CreatedBy)
	if err != nil {
		return ""
	}

	return user.Username
}

func (i *Iterations) Updater() string {
	user, err := NewUsers(i.UpdatedBy)
	if err != nil {
		return ""
	}

	return user.Username
}

func (i *Iterations) Deleter() string {
	user, err := NewUsers(i.DeletedBy)
	if err != nil {
		return ""
	}

	return user.Username
}

func (i *Iterations) IterationsCount(pIDs ...uint) (int64, error) {
	var count int64
	query := Conn.Model(&Iterations{})
	if len(pIDs) > 0 {
		query = query.Where("project_id IN ?", pIDs)
	}
	return count, query.Count(&count).Error
}

// func GetOneProjectIterations(pid, uid uint, page, pageSize int) ([]api.IterationSchemaData, error) {
// 	var res []api.IterationSchemaData
// 	pm, _ := NewProjectMembers()
// 	pm.UserID = uid
// 	pm.ProjectID = pid
// 	if err := pm.GetByUserIDAndProjectID(); err != nil {
// 		return nil, err
// 	}

// 	i, _ := NewIterations()
// 	iterations, err := i.List(page, pageSize, pid)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, v := range iterations {
// 		res = append(res, api.IterationSchemaData{})
// 	}
// 	return res, nil
// }
