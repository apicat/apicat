package response

import userbase "github.com/apicat/apicat/backend/route/proto/user/base"

type UserData struct {
	userbase.EmailOption
	userbase.NameOption
	userbase.AvatarOption
	Role string `json:"role" binding:"required"`
}
