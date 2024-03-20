package response

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	userbase "github.com/apicat/apicat/v2/backend/route/proto/user/base"
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
