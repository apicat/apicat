package proto

type CreateServer struct {
	Description string `json:"description" binding:"lte=255"`
	Url         string `json:"url" binding:"required,lte=255"`
}
