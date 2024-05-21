package user

import (
	"github.com/apicat/apicat/v2/backend/model"
)

// Oauth2Bind oauth2绑定关系
type Oauth2Bind struct {
	ID       uint
	UserID   uint   `gorm:"type:bigint;uniqueIndex:ukey;comment:用户id"`                    // github.com/apicat/apicat/v2 user.id
	Type     string `gorm:"type:varchar(32);uniqueIndex:ukey;not null;comment:对应oauth平台"` // github
	OauthUID string `gorm:"type:varchar(255);comment:对应oauth平台的用户id"`                     // github user.id
	model.TimeModel
}
