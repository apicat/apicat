package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID        uint   `gorm:"type:integer primary key autoincrement"`
	Email     string `gorm:"index;type:varchar(255);not null;comment:邮箱"`
	Username  string `gorm:"type:varchar(255);not null;comment:用户名"`
	Password  string `gorm:"type:varchar(255);not null;comment:密码"`
	Role      string `gorm:"type:varchar(255);not null;comment:角色:superadmin,admin,user"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func NewUsers(ids ...uint) (*Users, error) {
	users := &Users{}
	if len(ids) > 0 {
		if err := Conn.Take(users, ids[0]).Error; err != nil {
			return users, err
		}
		return users, nil
	}
	return users, nil
}

func (u *Users) GetByEmail(email string) error {
	return Conn.Where("email = ?", email).Take(u).Error
}

func (u *Users) List() ([]Users, error) {
	var users []Users
	return users, Conn.Order("created_at desc").Find(&users).Error
}

func (u *Users) Count() (int64, error) {
	var count int64
	return count, Conn.Model(&Users{}).Count(&count).Error
}

func (u *Users) Delete() error {
	return Conn.Delete(u).Error
}

func (u *Users) Save() error {
	return Conn.Save(u).Error
}
