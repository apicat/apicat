package proto

type SetUserInfoData struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Username string `json:"username" binding:"required,lte=255"`
}

type ChangePasswordData struct {
	Password           string `json:"password" binding:"required,gte=6,lte=255"`
	NewPassword        string `json:"new_password" binding:"required,gte=6,lte=255"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,gte=6,lte=255,eqfield=NewPassword"`
}

type LoginEmail struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password" binding:"required,gte=6,lte=255"`
}

type RegisterEmail struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password" binding:"required,gte=6,lte=255"`
}
