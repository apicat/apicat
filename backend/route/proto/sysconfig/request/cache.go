package request

type RedisOption struct {
	Host     string `json:"host" binding:"required,gt=1"`
	Password string `json:"password" binding:"omitempty"`
	Database int    `json:"database" binding:"gte=0"`
}
