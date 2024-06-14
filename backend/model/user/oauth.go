package user

import (
	"github.com/apicat/apicat/v2/backend/model"
)

// Oauth2Bind oauth2绑定关系
type Oauth2Bind struct {
	ID       uint
	UserID   uint   `gorm:"type:bigint;uniqueIndex:ukey;comment:user id"`                  // github.com/apicat/apicat/v2 user.id
	Type     string `gorm:"type:varchar(32);uniqueIndex:ukey;not null;comment:oauth type"` // github
	OauthUID string `gorm:"type:varchar(255);comment:uid of the OAuth"`                    // github user.id
	model.TimeModel
}
