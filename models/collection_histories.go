package models

import "time"

type CollectionHistories struct {
	ID           uint   `gorm:"type:integer primary key autoincrement"`
	CollectionId uint   `gorm:"index;not null;comment:标签id"`
	Title        string `gorm:"type:varchar(255);not null;comment:名称"`
	Type         string `gorm:"type:varchar(255);not null;comment:类型:category,doc,http"`
	Content      string `gorm:"type:mediumtext;comment:内容"`
	CreatedAt    time.Time
	CreatedBy    uint `gorm:"not null;default:0;comment:创建人id"`
}
