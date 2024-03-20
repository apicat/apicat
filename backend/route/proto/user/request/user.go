package request

import (
	"mime/multipart"

	userbase "github.com/apicat/apicat/backend/route/proto/user/base"
)

type UserIDOption struct {
	UserID uint `uri:"userID" json:"userID" query:"userID" binding:"required,gt=0"`
}

type SetUserSelfOption struct {
	userbase.NameOption
	userbase.LanguageOption
}

type ChangePwdOption struct {
	// 老密码
	PasswordOption
	// 新密码
	NewPassword string `json:"newPassword" binding:"required,gte=6,lte=64"`
	// 重复新密码
	ReNewPassword string `json:"reNewPassword" binding:"required,eqfield=NewPassword"`
}

type UploadAvatarOption struct {
	// 头像文件
	Avatar        *multipart.FileHeader `form:"avatar" binding:"required"`
	CroppedX      int                   `form:"croppedX" binding:"gte=0"`
	CroppedY      int                   `form:"croppedY" binding:"gte=0"`
	CroppedWidth  int                   `form:"croppedWidth" binding:"required,gte=1"`
	CroppedHeight int                   `form:"croppedHeight" binding:"required,eqfield=CroppedWidth"`
}

type UserListOption struct {
	Page     int    `query:"page"`
	PageSize int    `query:"pageSize"`
	Keywords string `query:"keywords" validate:"omitempty,lt=255"`
}

type ChangePasswordOption struct {
	UserIDOption
	PasswordOption
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}
