package models

import (
	"fmt"
	"time"
)

type CollectionHistories struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	CollectionId uint   `gorm:"type:bigint;index;not null;comment:集合id"`
	Title        string `gorm:"type:varchar(255);not null;comment:名称"`
	Type         string `gorm:"type:varchar(255);not null;comment:类型:category,doc,http"`
	Content      string `gorm:"type:mediumtext;comment:内容"`
	CreatedAt    time.Time
	CreatedBy    uint `gorm:"type:bigint;not null;default:0;comment:创建人id"`
}

func NewCollectionHistories(ids ...uint) (*CollectionHistories, error) {
	if len(ids) > 0 {
		ch := &CollectionHistories{ID: ids[0]}
		if err := Conn.Take(ch).Error; err != nil {
			return ch, err
		}
		return ch, nil
	}
	return &CollectionHistories{}, nil
}

func (ch *CollectionHistories) List(collectionIDs ...uint) ([]*CollectionHistories, error) {
	var collectionHistories []*CollectionHistories
	if len(collectionIDs) > 0 {
		return collectionHistories, Conn.Where("collection_id IN ?", collectionIDs).Order("created_at desc").Find(&collectionHistories).Error
	}
	return collectionHistories, Conn.Order("created_at desc").Find(&collectionHistories).Error
}

func (ch *CollectionHistories) Create() error {
	return Conn.Create(ch).Error
}

func (ch *CollectionHistories) Restore(collection *Collections, uid uint) error {
	fmt.Printf("ch: %+v\n", *ch)
	fmt.Printf("collection: %+v\n", *collection)
	return collection.UpdateContent(true, ch.Title, ch.Content, uid)
}
