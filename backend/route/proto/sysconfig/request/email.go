package request

type SMTPOption struct {
	Host     string `json:"host" binding:"required,gt=1"`
	User     string `json:"user" binding:"omitempty"`
	Address  string `json:"address" binding:"required,email"`
	Password string `json:"password" binding:"required,gt=1"`
}

type SendCloudOption struct {
	ApiUser   string `json:"apiUser" binding:"required,gt=1"`
	ApiKey    string `json:"apiKey" binding:"required,gt=1"`
	FromEmail string `json:"fromEmail" binding:"required,email"`
	FromName  string `json:"fromName" binding:"required,gt=1"`
}
