package user

import (
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/model/project"
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID        uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	Email     string `gorm:"type:varchar(255);index;not null;comment:邮箱"`
	Username  string `gorm:"type:varchar(255);not null;comment:用户名"`
	Password  string `gorm:"type:varchar(255);not null;comment:密码"`
	Role      string `gorm:"type:varchar(255);not null;comment:角色:superadmin,admin,user"`
	IsEnabled int    `gorm:"type:tinyint(1);not null;default:1;comment:是否启用:0停用,1启用"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func init() {
	model.RegMigrate(&Users{})
}

func NewUsers(ids ...uint) (*Users, error) {
	users := &Users{}
	if len(ids) > 0 {
		if err := model.Conn.Take(users, ids[0]).Error; err != nil {
			return users, err
		}
		return users, nil
	}
	return users, nil
}

func (u *Users) GetByEmail(email string) error {
	return model.Conn.Where("email = ?", email).Take(u).Error
}

func (u *Users) List(page, pageSize int) ([]Users, error) {
	var users []Users

	if page == 0 && pageSize == 0 {
		return users, model.Conn.Order("created_at desc").Find(&users).Error
	}

	return users, model.Conn.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&users).Error
}

func (u *Users) Count() (int64, error) {
	var count int64
	return count, model.Conn.Model(&Users{}).Count(&count).Error
}

func (u *Users) Delete() error {
	pms, err := project.GetUserInvolvedProject(u.ID)
	if err != nil {
		return err
	}

	for _, pm := range pms {
		if err := pm.Delete(); err != nil {
			return err
		}
	}

	return model.Conn.Delete(u).Error
}

func (u *Users) Save() error {
	return model.Conn.Save(u).Error
}

func UserListByIDs(ids []uint) ([]*Users, error) {
	var users []*Users
	return users, model.Conn.Where("id in (?)", ids).Find(&users).Error
}
