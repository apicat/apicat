package proto

type DBConfigItemData struct {
	Value string `json:"value"`
	Type  string `json:"type" binding:"required,oneof=value env"`
}

type DBConfigData struct {
	Host     DBConfigItemData `json:"host" binding:"required"`
	Port     DBConfigItemData `json:"port" binding:"required"`
	User     DBConfigItemData `json:"user" binding:"required"`
	Password DBConfigItemData `json:"password" binding:"required"`
	DBName   DBConfigItemData `json:"dbname" binding:"required"`
}
