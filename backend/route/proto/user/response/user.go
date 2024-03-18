package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
	userbase "apicat-cloud/backend/route/proto/user/base"
)

type User struct {
	protobase.OnlyIdInfo
	UserData
	userbase.LanguageOption
	Github bool `json:"github"`
}

type UserList struct {
	protobase.PaginationInfo
	Items []User `json:"items"`
}
