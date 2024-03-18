package base

type ServiceOption struct {
	AppName        string `json:"appName" binding:"required,gt=1"`
	AppUrl         string `json:"appUrl" binding:"required,url"`
	AppServerBind  string `json:"appServerBind" binding:"required,gte=10"`
	MockUrl        string `json:"mockUrl" binding:"required,url"`
	MockServerBind string `json:"mockServerBind" binding:"required,gte=10"`
}

type GitHubClientID struct {
	ClientID string `json:"clientID" binding:"required,gt=1"`
}

type GitHubOption struct {
	GitHubClientID
	ClientSecret string `json:"clientSecret" binding:"required,gt=1"`
}

type ConfigList []*ConfigDetail

type ConfigDetail struct {
	Driver string                 `json:"driver" binding:"required,gt=1"`
	Use    bool                   `json:"use" binding:"required"`
	Config map[string]interface{} `json:"config" binding:"required"`
}

type MySQLDetail struct {
	Host     string `json:"host" binding:"required,gt=1"`
	Database string `json:"database" binding:"required,gt=1"`
	Username string `json:"username" binding:"required,gt=1"`
}
