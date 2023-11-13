package proto

type TrashsRecoverQuery struct {
	CollectionID []uint `form:"collection-id" binding:"required,dive,gte=0"`
}

type TrashsRecoverBody struct {
	Category uint `json:"category" binding:"gte=0"`
}
