package user

import (
	"context"
	"errors"
	"time"

	"github.com/apicat/apicat/v2/backend/model"

	"gorm.io/gorm"
)

const (
	LanguageZhCN = "zh-CN"
	LanguageEnUS = "en-US"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

var SupportedLanguages = map[string]bool{
	LanguageZhCN: true,
	LanguageEnUS: true,
}

type User struct {
	ID          uint      `gorm:"primarykey"`
	Name        string    `gorm:"type:varchar(255);comment:用户名"`
	Password    string    `gorm:"type:varchar(64);comment:密码"`
	Email       string    `gorm:"type:varchar(255);uniqueIndex;comment:邮箱"`
	Avatar      string    `gorm:"type:varchar(255);comment:头像"`
	Language    string    `gorm:"type:varchar(32);comment:语言"` // zh-CN en-US
	Role        string    `gorm:"type:varchar(32);comment:角色"`
	LastLoginIP string    `gorm:"type:varchar(15);comment:最后登录ip"`
	LastLoginAt time.Time `gorm:"type:datetime;not null;comment:最后登录时间"`
	IsActive    bool      `gorm:"type:tinyint;not null;comment:是否激活"`
	model.TimeModel
}

func init() {
	model.RegMigrate(&User{})
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password, err = hashPassword(u.Password)
	return
}

func (u *User) Create(ctx context.Context) error {
	return model.DB(ctx).Create(u).Error
}

func (u *User) CreateAndBindOauth(ctx context.Context, typ, oauthid string) error {
	return model.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(u).Error; err != nil {
			return err
		}
		return tx.Create(&Oauth2Bind{
			Type:     typ,
			UserID:   u.ID,
			OauthUID: oauthid,
		}).Error
	})
}

func (u *User) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if u.ID > 0 {
		tx = tx.Take(u, "id = ?", u.ID)
	} else if u.Email != "" {
		tx = tx.Take(u, "email = ?", u.Email)
	} else {
		return false, errors.New("xx")
	}
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

func (u *User) UpdateLastLogin(ctx context.Context, ip string) error {
	// 防止更新updated_at字段
	return model.DB(ctx).Model(u).Select("last_login_ip", "last_login_at").UpdateColumns(map[string]interface{}{
		"last_login_ip": ip,
		"last_login_at": time.Now(),
	}).Error
}

func (u *User) UpdateEmail(ctx context.Context) error {
	tx := model.DB(ctx).Model(u)
	return tx.UpdateColumn("email", u.Email).Error
}

func (u *User) SetActive(ctx context.Context) error {
	return model.DB(ctx).Model(u).Update("is_active", true).Error
}

func (u *User) CheckPassword(pwd string) bool {
	return checkPasswordHash(pwd, u.Password)
}

func (u *User) UpdatePassword(ctx context.Context) error {
	pwd, err := hashPassword(u.Password)
	if err != nil {
		return err
	}
	tx := model.DB(ctx).Model(u)
	return tx.UpdateColumn("password", pwd).Error
}

func (u *User) Update(ctx context.Context) error {
	if u.ID == 0 {
		return nil
	}
	// 只能更新name、avatar和language
	return model.DB(ctx).Model(u).Updates(map[string]interface{}{
		"name":     u.Name,
		"avatar":   u.Avatar,
		"language": u.Language,
	}).Error
}

// 获取用户的oauth绑定关系
func (u *User) Oauths(ctx context.Context, types ...string) ([]*Oauth2Bind, error) {
	var list []*Oauth2Bind
	if err := model.DB(ctx).Where("user_id = ? and type in (?)", u.ID, types).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// 创建用户的oauth绑定关系
func (u *User) BindOauth(ctx context.Context, typ, oauthid string) error {
	o := &Oauth2Bind{
		Type:     typ,
		UserID:   u.ID,
		OauthUID: oauthid,
	}
	return model.DB(ctx).Create(o).Error
}

// 创建用户的oauth绑定关系，如果被删除过，则修改被删除的记录
func (u *User) BindOrRecoverOauth(ctx context.Context, typ, oauthid string) error {
	var ob Oauth2Bind

	if err := model.DB(ctx).Unscoped().Take(&ob, "type = ? and oauth_uid = ?", typ, oauthid).Error; err == nil {
		return model.DB(ctx).Unscoped().Model(&Oauth2Bind{}).Where("id = ?", ob.ID).Updates(map[string]interface{}{
			"user_id":    u.ID,
			"deleted_at": nil,
		}).Error
	}

	o := &Oauth2Bind{
		Type:     typ,
		UserID:   u.ID,
		OauthUID: oauthid,
	}
	return model.DB(ctx).Create(o).Error
}

// 删除用户的oauth绑定关系
func (u *User) UnBindOauth(ctx context.Context, typ string) error {
	return model.DB(ctx).Delete(&Oauth2Bind{}, "user_id = ? and type = ?", u.ID, typ).Error
}

func (u *User) IsSysAdmin(ctx context.Context) bool {
	return u.Role == RoleAdmin
}

func (u *User) SetSysAdmin(ctx context.Context) error {
	return model.DB(ctx).Model(u).UpdateColumn("role", RoleAdmin).Error
}

func (u *User) Delete(ctx context.Context) error {
	return model.DB(ctx).Delete(u).Error
}
