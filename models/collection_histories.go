package models

import "time"

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
		return collectionHistories, Conn.Where("collection_id IN ?", collectionIDs).Find(&collectionHistories).Error
	}
	return collectionHistories, Conn.Find(&collectionHistories).Error
}

func (ch *CollectionHistories) Create() error {
	return Conn.Create(ch).Error
}

func (ch *CollectionHistories) Restore(uid uint) error {
	c, err := NewCollections(ch.CollectionId)
	if err != nil {
		return err
	}

	nch, _ := NewCollectionHistories()
	nch.CollectionId = c.ID
	nch.Title = c.Title
	nch.Type = c.Type
	nch.Content = c.Content
	nch.CreatedBy = uid
	if err := nch.Create(); err != nil {
		return err
	}

	c.Title = ch.Title
	c.Content = ch.Content
	c.UpdatedBy = uid
	if err := c.Update(); err != nil {
		return err
	}

	return nil
}
