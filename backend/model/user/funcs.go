package user

import (
	"context"

	"github.com/apicat/apicat/backend/model"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 通过oauth平台的用户id获取用户
func GetUserByOauth(ctx context.Context, oauthUsrID, typ string) (*User, error) {
	var (
		ob  Oauth2Bind
		usr User
		uid uint
	)

	tx := model.DB(ctx)
	ret := tx.Model(ob).Where("type = ? and oauth_uid = ?", typ, oauthUsrID).Pluck("user_id", &uid)
	if ret.Error != nil || uid < 1 {
		return nil, model.NotRecord(ret)
	}

	ret = tx.First(&usr, uid)
	if ret.Error != nil {
		return nil, model.NotRecord(ret)
	}

	return &usr, nil
}

// 通过oauth平台的用户id获取用户，如果被删除过，则恢复
func GetAndRecoverUserByOauth(ctx context.Context, oauthUsrID, typ string) (*User, error) {
	var (
		ob  Oauth2Bind
		usr User
	)

	ret := model.DB(ctx).Unscoped().Take(&ob, "type = ? and oauth_uid = ?", typ, oauthUsrID)
	if ret.Error != nil {
		return nil, model.NotRecord(ret)
	}

	ret = model.DB(ctx).First(&usr, ob.UserID)
	if ret.Error != nil {
		return nil, model.NotRecord(ret)
	}

	ret = model.DB(ctx).Unscoped().Model(&Oauth2Bind{}).Where("id = ?", ob.ID).Update("deleted_at", nil)
	if ret.Error != nil {
		return nil, model.NotRecord(ret)
	}
	return &usr, nil
}

func GetUsers(ctx context.Context, page, pageSize int, keywords string) ([]*User, error) {
	list := make([]*User, 0)

	tx := model.DB(ctx).Model(&User{})
	if keywords != "" {
		k := "%" + keywords + "%"
		tx.Where("name like ? OR email like ?", k, k)
	}
	if page > 0 && pageSize > 0 {
		tx = tx.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	if err := tx.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func GetUserCount(ctx context.Context, keywords string) (int64, error) {
	var count int64
	if keywords != "" {
		k := "%" + keywords + "%"
		return count, model.DB(ctx).Model(&User{}).Where("name like ? OR email like ?", k, k).Count(&count).Error
	}
	return count, model.DB(ctx).Model(&User{}).Count(&count).Error
}
