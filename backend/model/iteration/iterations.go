package iteration

import (
	"errors"
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/model/user"
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

func init() {
	model.RegMigrate(&Iterations{})
}

func NewIterations(ids ...any) (*Iterations, error) {
	iteration := &Iterations{}

	if len(ids) > 0 {
		var err error

		switch ids[0].(type) {
		case string:
			err = model.Conn.Where("public_id = ?", ids[0]).Take(iteration).Error
		case uint:
			err = model.Conn.Take(iteration, ids[0]).Error
		default:
			err = errors.New("invalid id type")
		}

		if err != nil {
			return iteration, err
		}
		return iteration, nil
	}
	return iteration, nil
}

func (i *Iterations) List(page, pageSize int, pIDs ...uint) ([]*Iterations, error) {
	var (
		iterations []*Iterations
		query      *gorm.DB
	)

	if len(pIDs) > 0 {
		query = model.Conn.Where("project_id IN ?", pIDs).Order("created_at desc")
	} else {
		query = model.Conn.Order("created_at desc")
	}

	if page != 0 && pageSize != 0 {
		query = query.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	return iterations, query.Find(&iterations).Error
}

func (i *Iterations) Create() error {
	return model.Conn.Create(i).Error
}

func (i *Iterations) Update() error {
	return model.Conn.Save(i).Error
}

func (i *Iterations) Delete() error {
	return model.Conn.Delete(i).Error
}

func (i *Iterations) Creator() string {
	u, err := user.NewUsers(i.CreatedBy)
	if err != nil {
		return ""
	}

	return u.Username
}

func (i *Iterations) Updater() string {
	u, err := user.NewUsers(i.UpdatedBy)
	if err != nil {
		return ""
	}

	return u.Username
}

func (i *Iterations) Deleter() string {
	u, err := user.NewUsers(i.DeletedBy)
	if err != nil {
		return ""
	}

	return u.Username
}

func IterationsCount(pIDs ...uint) (int64, error) {
	var count int64
	query := model.Conn.Model(&Iterations{})
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
