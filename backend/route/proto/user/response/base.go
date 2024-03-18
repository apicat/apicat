package response

import userbase "apicat-cloud/backend/route/proto/user/base"

type UserData struct {
	userbase.EmailOption
	userbase.NameOption
	userbase.AvatarOption
	Role string `json:"role" binding:"required"`
}
