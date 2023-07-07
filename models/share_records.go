package models

import (
	"time"
)

type ShareRecords struct {
	ID           uint      `gorm:"type:bigint;primaryKey;autoIncrement"`
	ShareToken   string    `gorm:"type:varchar(255);index;not null;comment:md5的分享token"`
	Expiration   time.Time `gorm:"type:datetime;not null;comment:过期时间"`
	ProjectID    uint      `gorm:"type:bigint;index;not null;comment:项目id"`
	CollectionID uint      `gorm:"type:bigint;index;comment:集合id"`
	MemberID     uint      `gorm:"type:bigint;index;not null;comment:发起分享的成员id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewShareRecords() *ShareRecords {
	return &ShareRecords{}
}

func (sr *ShareRecords) GetByShareToken() error {
	return Conn.Where("share_token = ?", sr.ShareToken).Take(sr).Error
}

func (sr *ShareRecords) Create() error {
	return Conn.Create(sr).Error
}

func (sr *ShareRecords) DeleteByProjectID() error {
	return Conn.Where("project_id = ?", sr.ProjectID).Delete(sr).Error
}

func (sr *ShareRecords) DeleteByCollectionID() error {
	return Conn.Where("collection_id = ?", sr.CollectionID).Delete(sr).Error
}
