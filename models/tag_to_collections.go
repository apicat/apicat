package models

import "time"

type TagToCollections struct {
	ID           uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	TagId        uint `gorm:"type:bigint;index;not null;comment:标签id"`
	CollectionId uint `gorm:"type:bigint;not null;comment:集合id"`
	DisplayOrder int  `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewTagToCollections() *TagToCollections {
	return &TagToCollections{}
}

func (ttc *TagToCollections) Create() error {
	return Conn.Create(ttc).Error
}

func CollectionToTagIds(collectionID uint) []uint {
	var (
		tagIds  []uint
		records []TagToCollections
	)
	if err := Conn.Where("collection_id = ?", collectionID).Find(&records).Error; err != nil {
		for _, v := range records {
			tagIds = append(tagIds, v.TagId)
		}
	}

	return tagIds
}
