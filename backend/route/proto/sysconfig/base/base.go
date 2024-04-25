package base

type ServiceOption struct {
	AppName string `json:"appName" binding:"required,gt=1"`
	AppUrl  string `json:"appUrl" binding:"required,url"`
	MockUrl string `json:"mockUrl" binding:"required,url"`
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
